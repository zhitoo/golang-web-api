# Broker & Job System

## Quick Examples

### Example 1: Simple Job (Default Queue)

```go
// 📁 app/modules/notification/email_job.go
package notification

import (
    "context"
    "fmt"

    "github.com/zhitoo/golang-web-api/pkg/job"
)

type SendWelcomeEmailJob struct {
    job.Base
    To   string `json:"to"`
    Name string `json:"name"`
}

func (j *SendWelcomeEmailJob) Handle(ctx context.Context) error {
    fmt.Printf("Welcome %s!\n", j.Name)
    return nil
}

func init() {
    job.RegisterJob("SendWelcomeEmailJob", func() job.JobHandler {
        return &SendWelcomeEmailJob{}
    })
}
```

```go
// 📁 app/modules/notification/service.go
package notification

import "github.com/zhitoo/golang-web-api/pkg/job"

func RegisterUser(name, email string) error {
    // user registration logic...

    job.Dispatch(&SendWelcomeEmailJob{To: email, Name: name})
    return nil
}
```

```go
// 📁 cmd/main.go
package main

import "github.com/zhitoo/golang-web-api/pkg/job"

func main() {
    // ... broker setup ...

    job.NewWorker(broker, "default").Start()
}
```

### Example 2: Job with Custom Queue

```go
// 📁 app/modules/otp/sms_job.go
package otp

import (
    "context"

    "github.com/zhitoo/golang-web-api/pkg/job"
    "github.com/zhitoo/golang-web-api/pkg/logging"
)

var log logging.ScopedLogger

func init() {
    log = logging.NewLogger(config.GetConfig()).With(logging.General, logging.Api)
    job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}

type SendSmsJob struct {
    job.Base
    Mobile string `json:"mobile"`
    Body   string `json:"body"`
}

func (j *SendSmsJob) Handle(ctx context.Context) error {
    log.Info("sending SMS", map[logging.ExtraKey]any{"mobile": j.Mobile})
    return smsProvider.Send(j.Mobile, j.Body)
}

func (j *SendSmsJob) Queue() string { return "sms" }
```

```go
// 📁 app/modules/otp/service.go
package otp

import (
    "fmt"
    "time"

    "github.com/zhitoo/golang-web-api/database/cache"
    "github.com/zhitoo/golang-web-api/pkg/job"
)

func SendOTP(mobile string) (string, error) {
    otp := fmt.Sprintf("%06d", rand.Intn(1000000))

    job.Dispatch(&SendSmsJob{
        Mobile: mobile,
        Body:   fmt.Sprintf("Your code: %s", otp),
    })

    cache.SetValue("otp:"+mobile, otp, 2*time.Minute)
    return otp, nil
}
```

```go
// 📁 cmd/main.go
package main

import "github.com/zhitoo/golang-web-api/pkg/job"

func main() {
    // ... broker setup ...

    job.NewWorker(broker, "sms").Start()
}
```

### Example 3: Multiple Workers (Load Balancing)

```go
// 📁 cmd/main.go
package main

import "github.com/zhitoo/golang-web-api/pkg/job"

func main() {
    // ... broker setup ...

    // 3 workers - each job goes to ONE worker
    for i := 0; i < 3; i++ {
        job.NewWorker(broker, "sms").Start()
    }
}
```

```
                    ┌──────────────┐
                    │   Worker 1   │
POST /otp/send ──▶ ├──────────────┤ ──▶ Worker 2 processes it
                    │   Worker 2   │
                    ├──────────────┤
                    │   Worker 3   │
                    └──────────────┘
```

### Example 4: Multiple Queues

```go
// 📁 app/modules/otp/sms_job.go
func (j *SendSmsJob) Queue() string { return "sms" }

// 📁 app/modules/notification/email_job.go
func (j *SendEmailJob) Queue() string { return "email" }

// 📁 app/modules/order/order_job.go
func (j *ProcessOrderJob) Queue() string { return "order" }
```

```go
// 📁 cmd/main.go
package main

import "github.com/zhitoo/golang-web-api/pkg/job"

func main() {
    // ... broker setup ...

    job.NewWorker(broker, "sms").Start()
    job.NewWorker(broker, "email").Start()
    job.NewWorker(broker, "order").Start()
}
```

### Example 5: Dispatch to Specific Queue

```go
// 📁 app/modules/notification/service.go
package notification

import "github.com/zhitoo/golang-web-api/pkg/job"

func SendNotification(user *User, message string) {
    // Normal: uses Queue() which returns "email"
    job.Dispatch(&SendEmailJob{To: user.Email, Body: message})

    // Override: force "priority" queue
    job.DispatchOnQueue(&SendEmailJob{To: user.Email, Body: message}, "priority")
}
```

### Example 6: Durable Jobs (Survive Restart)

```go
// 📁 cmd/main.go
package main

import "github.com/zhitoo/golang-web-api/pkg/job"

func main() {
    d := job.NewDispatcher(broker)
    job.SetDispatcher(d)

    // Create stream on startup
    d.EnsureStream("JOBS", []string{"jobs.>"})

    // Durable worker
    w := job.NewWorker(broker, "sms")
    w.StartDurable("sms-worker-1")
}
```

---

## How It Works

### Architecture

```
┌─────────────┐     Publish      ┌─────────────┐     Consume      ┌─────────────┐
│   Handler   │ ───────────────▶ │     NATS    │ ───────────────▶ │   Worker    │
│  (Dispatch) │                  │   (Broker)  │                  │  (Process)  │
└─────────────┘                  └─────────────┘                  └─────────────┘
```

Three layers:

1. **Broker** (`broker/`) — NATS connection and pub/sub primitives
2. **Job Package** (`pkg/job/`) — dispatcher, consumer, worker abstractions
3. **Business Jobs** (`app/modules/*/`) — concrete job definitions

### Broker Layer

```go
// 📁 broker/broker.go
package broker

type Broker interface {
    Publish(subject string, data []byte) error
    Subscribe(subject string, handler func(msg []byte)) error
    QueueSubscribe(subject string, queue string, handler func(msg []byte)) error
}

type DurableBroker interface {
    Broker
    PublishDurable(subject string, data []byte) error
    SubscribeDurable(subject string, durable string, handler func(msg []byte)) error
    CreateStream(name string, subjects []string) error
}
```

```go
// 📁 broker/nats/nats.go
package nats

import (
    "github.com/nats-io/nats.go"
    "github.com/zhitoo/golang-web-api/config"
)

type NatsBroker struct {
    conn *nats.Conn
    js   nats.JetStreamContext
    cfg  *config.Config
}

func (b *NatsBroker) QueueSubscribe(subject string, queue string, handler func(msg []byte)) error {
    _, err := b.conn.QueueSubscribe(subject, queue, func(m *nats.Msg) {
        handler(m.Data)
    })
    return err
}
```

**Subscribe vs QueueSubscribe:**

| Method | Behavior | Use Case |
|--------|----------|----------|
| `Subscribe` | Every subscriber gets every message (fan-out) | Broadcasting, events |
| `QueueSubscribe` | One subscriber per message (load-balanced) | Job workers |

### Job Package (`pkg/job/`)

#### Job Struct

```go
// 📁 pkg/job/job.go
package job

import "time"

type Status string

const (
    Pending    Status = "pending"
    Processing Status = "processing"
    Success    Status = "success"
    Failed     Status = "failed"
)

type Job struct {
    ID        string
    Type      string    // e.g., "SendSmsJob"
    Payload   []byte    // JSON-serialized job data
    Status    Status
    Attempts  int
    MaxRetry  int
    CreatedAt time.Time
    RunAt     time.Time
}
```

#### JobHandler Interface

```go
// 📁 pkg/job/handler.go
package job

import "context"

type JobHandler interface {
    Handle(ctx context.Context) error
    Queue() string
}

type Base struct{}

func (b Base) Queue() string { return "default" }
```

#### Registry

```go
// 📁 pkg/job/registry.go
package job

import "reflect"

type Factory func() JobHandler

var registry = map[string]Factory{}

func RegisterJob(name string, factory Factory) {
    registry[name] = factory
}

func getFactory(name string) (Factory, bool) {
    f, ok := registry[name]
    return f, ok
}

func jobTypeName(j JobHandler) string {
    t := reflect.TypeOf(j)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    return t.Name()
}
```

#### Dispatch

```go
// 📁 pkg/job/dispatch.go
package job

import (
    "encoding/json"
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
    payload, _ := json.Marshal(j)

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
```

#### Worker

```go
// 📁 pkg/job/worker.go
package job

import (
    "context"
    "encoding/json"
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

func (w *Worker) process(j Job) error {
    factory, _ := getFactory(j.Type)
    handler := factory()
    json.Unmarshal(j.Payload, handler)
    return handler.Handle(context.Background())
}
```

### Complete Example: OTP SMS

**Step 1: Define the Job**

```go
// 📁 app/modules/otp/sms_job.go
package otp

import (
    "context"

    "github.com/zhitoo/golang-web-api/config"
    "github.com/zhitoo/golang-web-api/pkg/job"
    "github.com/zhitoo/golang-web-api/pkg/logging"
)

var log logging.ScopedLogger

func init() {
    cfg := config.GetConfig()
    log = logging.NewLogger(cfg).With(logging.General, logging.Api)
    job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}

type SendSmsJob struct {
    job.Base
    Mobile string `json:"mobile"`
    Body   string `json:"body"`
}

func (j *SendSmsJob) Handle(ctx context.Context) error {
    log.Info("sending SMS", map[logging.ExtraKey]any{"mobile": j.Mobile})
    return smsProvider.Send(j.Mobile, j.Body)
}

func (j *SendSmsJob) Queue() string { return "sms" }
```

**Step 2: Dispatch the Job**

```go
// 📁 app/modules/otp/service.go
package otp

import (
    "fmt"
    "time"

    "github.com/zhitoo/golang-web-api/database/cache"
    "github.com/zhitoo/golang-web-api/pkg/job"
)

func SendOTP(mobile string) (string, error) {
    otp := fmt.Sprintf("%06d", rand.Intn(1000000))

    job.Dispatch(&SendSmsJob{
        Mobile: mobile,
        Body:   fmt.Sprintf("Your code: %s", otp),
    })

    cache.SetValue("otp:"+mobile, otp, 2*time.Minute)
    return otp, nil
}
```

**Step 3: Start the Worker**

```go
// 📁 cmd/main.go
package main

import (
    "github.com/zhitoo/golang-web-api/app"
    "github.com/zhitoo/golang-web-api/broker"
    "github.com/zhitoo/golang-web-api/config"
    "github.com/zhitoo/golang-web-api/database/cache"
    "github.com/zhitoo/golang-web-api/database/db"
    "github.com/zhitoo/golang-web-api/pkg/job"
    "github.com/zhitoo/golang-web-api/pkg/logging"
)

func main() {
    cfg := config.GetConfig()
    logger := logging.NewLogger(cfg)

    cache.InitRedis(cfg)
    defer cache.CloseRedis()

    b, _ := broker.NewBroker(cfg)
    b.Connect()
    defer b.Close()

    d := job.NewDispatcher(b)
    job.SetDispatcher(d)

    // Start SMS worker
    smsWorker := job.NewWorker(b, "sms")
    smsWorker.Start()

    db.InitDb(cfg)
    defer db.CloseDb()

    app.InitServer(cfg)
}
```

**Step 4: Flow**

1. User calls `POST /api/v1/otp/send`
2. Handler calls `SendOTP()`
3. `SendOTP()` dispatches `SendSmsJob` to queue `sms`
4. Worker receives job, calls `Handle()`
5. SMS is sent

---

## Subject Naming

All job subjects follow: `jobs.<queue>`

| Queue | Subject | Description |
|-------|---------|-------------|
| `default` | `jobs.default` | Default queue |
| `sms` | `jobs.sms` | SMS jobs |
| `email` | `jobs.email` | Email jobs |
| `order` | `jobs.order` | Order processing |

## File Structure

```
pkg/job/
├── job.go          # Job struct, Status, Broker, Dispatcher, Consumer
├── handler.go      # JobHandler interface, Base struct
├── registry.go     # RegisterJob, factory registry
├── dispatch.go     # Dispatch, DispatchOnQueue, SetDispatcher
└── worker.go       # Worker with QueueSubscribe

broker/
├── broker.go       # Broker interface, DurableBroker interface
├── factory.go      # NewBroker factory
└── nats/
    └── nats.go     # NATS implementation

app/modules/otp/
├── sms_job.go      # SendSmsJob definition + init() registration
└── service.go      # SendOTP() dispatches SendSmsJob
```

## Debugging

```bash
# List streams
nats stream ls

# Check stream info
nats stream info JOBS

# Monitor messages
nats sub "jobs.sms"
```

| Error | Cause | Solution |
|-------|-------|----------|
| `nats: no response from stream` | Stream doesn't exist | Call `EnsureStream()` on startup |
| `jetstream not enabled` | JetStream disabled | Add `jetstream {}` to nats.conf |
| `unknown job type` | Job not registered | Call `RegisterJob()` in `init()` |
| `dispatcher not initialized` | `SetDispatcher()` not called | Call it in `main.go` |
