package nats

import (
	"fmt"
	"log"

	"github.com/zhitoo/golang-web-api/config"

	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	cfg  *config.Config
}

func New(cfg *config.Config) *NatsBroker {
	return &NatsBroker{cfg: cfg}
}

func (b *NatsBroker) Connect() error {
	var err error
	b.conn, err = nats.Connect(
		fmt.Sprintf("nats://%s:%s", b.cfg.Nats.Host, b.cfg.Nats.Port),
		nats.UserInfo("nats", b.cfg.Nats.Password),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	if b.cfg.Nats.JetStream.Enabled {
		b.js, err = b.conn.JetStream()
		if err != nil {
			return fmt.Errorf("failed to create jetstream: %w", err)
		}
	}

	return nil
}

func (b *NatsBroker) Publish(subject string, data []byte) error {
	return b.conn.Publish(subject, data)
}

func (b *NatsBroker) Subscribe(subject string, handler func(msg []byte)) error {
	_, err := b.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	})
	return err
}

func (b *NatsBroker) PublishDurable(subject string, data []byte) error {
	if b.js == nil {
		return fmt.Errorf("jetstream not enabled")
	}
	_, err := b.js.Publish(subject, data)
	return err
}

func (b *NatsBroker) SubscribeDurable(subject string, durable string, handler func(msg []byte)) error {
	if b.js == nil {
		return fmt.Errorf("jetstream not enabled")
	}
	_, err := b.js.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	}, nats.Durable(durable))
	return err
}

func (b *NatsBroker) CreateStream(name string, subjects []string) error {
	if b.js == nil {
		return fmt.Errorf("jetstream not enabled")
	}
	_, err := b.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
	})
	return err
}

func (b *NatsBroker) DeleteStream(name string) error {
	if b.js == nil {
		return fmt.Errorf("jetstream not enabled")
	}
	return b.js.DeleteStream(name)
}

func (b *NatsBroker) Close() {
	if b.conn != nil {
		b.conn.Close()
	}
}
