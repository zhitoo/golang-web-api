package job

import "context"

type JobHandler interface {
	Handle(ctx context.Context) error
	Queue() string
}

type Base struct{}

func (b Base) Queue() string { return "default" }
