package job

import (
	"context"
	"encoding/json"
	"fmt"
)

type Worker struct {
	consumer *Consumer
	queue    string
}

func NewWorker(broker Broker, queue string) *Worker {
	w := &Worker{queue: queue}
	w.consumer = NewConsumer(broker, w.process)
	return w
}

func (w *Worker) Start() error {
	subject := "jobs." + w.queue
	return w.consumer.ConsumeQueue(subject, w.queue)
}

func (w *Worker) StartDurable(durableName string) error {
	subject := "jobs." + w.queue
	return w.consumer.ConsumeDurable(subject, durableName)
}

func (w *Worker) process(j Job) error {
	factory, ok := getFactory(j.Type)
	if !ok {
		return fmt.Errorf("unknown job type: %s", j.Type)
	}

	handler := factory()
	if err := json.Unmarshal(j.Payload, handler); err != nil {
		return fmt.Errorf("failed to unmarshal job %s: %w", j.ID, err)
	}

	ctx := context.Background()
	return handler.Handle(ctx)
}
