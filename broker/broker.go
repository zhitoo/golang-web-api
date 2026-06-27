package broker

type Broker interface {
	Connect() error
	Publish(subject string, data []byte) error
	Subscribe(subject string, handler func(msg []byte)) error
	Close()
}

type DurableBroker interface {
	Broker
	PublishDurable(subject string, data []byte) error
	SubscribeDurable(subject string, durable string, handler func(msg []byte)) error
	CreateStream(name string, subjects []string) error
	DeleteStream(name string) error
}
