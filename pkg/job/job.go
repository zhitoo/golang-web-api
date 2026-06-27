package job

import (
	"encoding/json"
	"fmt"
	"time"
)

type Status string

const (
	Pending    Status = "pending"
	Processing Status = "processing"
	Success    Status = "success"
	Failed     Status = "failed"
)

type Job struct {
	ID        string
	Type      string
	Payload   []byte
	Status    Status
	Attempts  int
	MaxRetry  int
	CreatedAt time.Time
	RunAt     time.Time
}

type Broker interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, handler func(msg []byte)) error
}

type DurableBroker interface {
	Broker
	PublishDurable(subject string, data []byte) error
	SubscribeDurable(subject string, durable string, handler func(msg []byte)) error
	CreateStream(name string, subjects []string) error
}

type Dispatcher struct {
	broker     Broker
	durable    DurableBroker
	useDurable bool
}

func NewDispatcher(broker Broker) *Dispatcher {
	d := &Dispatcher{broker: broker}
	if db, ok := broker.(DurableBroker); ok {
		d.durable = db
		d.useDurable = true
	}
	return d
}

func (d *Dispatcher) Dispatch(subject string, job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if d.useDurable && d.durable != nil {
		return d.durable.PublishDurable(subject, data)
	}
	return d.broker.Publish(subject, data)
}

func (d *Dispatcher) DispatchDurable(subject string, job Job) error {
	if !d.useDurable || d.durable == nil {
		return fmt.Errorf("durable broker not available")
	}
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return d.durable.PublishDurable(subject, data)
}

func (d *Dispatcher) EnsureStream(name string, subjects []string) error {
	if !d.useDurable || d.durable == nil {
		return fmt.Errorf("durable broker not available")
	}
	return d.durable.CreateStream(name, subjects)
}

type Handler func(job Job) error

type Consumer struct {
	broker     Broker
	durable    DurableBroker
	useDurable bool
	handler    Handler
}

func NewConsumer(broker Broker, handler Handler) *Consumer {
	c := &Consumer{broker: broker, handler: handler}
	if db, ok := broker.(DurableBroker); ok {
		c.durable = db
		c.useDurable = true
	}
	return c
}

func (c *Consumer) Consume(subject string) error {
	return c.broker.Subscribe(subject, c.processMessage)
}

func (c *Consumer) ConsumeDurable(subject string, durable string) error {
	if !c.useDurable || c.durable == nil {
		return fmt.Errorf("durable broker not available")
	}
	return c.durable.SubscribeDurable(subject, durable, c.processMessage)
}

func (c *Consumer) processMessage(msg []byte) {
	var job Job
	if err := json.Unmarshal(msg, &job); err != nil {
		return
	}
	_ = c.handler(job)
}
