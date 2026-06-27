package broker

import (
	"fmt"

	"github.com/zhitoo/golang-web-api/config"

	natsbroker "github.com/zhitoo/golang-web-api/broker/nats"
)

func NewBroker(cfg *config.Config) (Broker, error) {
	switch cfg.App.Broker {
	case "nats":
		return natsbroker.New(cfg), nil
	default:
		return nil, fmt.Errorf("unknown broker: %s", cfg.App.Broker)
	}
}

func NewDurableBroker(cfg *config.Config) (DurableBroker, error) {
	b, err := NewBroker(cfg)
	if err != nil {
		return nil, err
	}
	db, ok := b.(DurableBroker)
	if !ok {
		return nil, fmt.Errorf("broker %s does not support durable operations", cfg.App.Broker)
	}
	return db, nil
}
