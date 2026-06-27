package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var globalDispatcher *Dispatcher

func SetDispatcher(d *Dispatcher) {
	globalDispatcher = d
}

func Dispatch(j JobHandler) error {
	return DispatchOnQueue(j, j.Queue())
}

func DispatchOnQueue(j JobHandler, queue string) error {
	if globalDispatcher == nil {
		return fmt.Errorf("dispatcher not initialized")
	}

	payload, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	jb := Job{
		ID:        uuid.New().String(),
		Type:      jobTypeName(j),
		Payload:   payload,
		Status:    Pending,
		CreatedAt: time.Now(),
	}

	subject := "jobs." + queue
	return globalDispatcher.Dispatch(subject, jb)
}
