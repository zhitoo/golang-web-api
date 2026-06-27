# Broker & Job System

## Quick Examples

### Example 1: Simple Job (Default Queue)

```go
// 1. Define job
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

// 2. Dispatch
job.Dispatch(&SendWelcomeEmailJob{To: "ali@test.com", Name: "Ali"})

// 3. Worker (main.go)
job.NewWorker(broker, "default").Start()
```

### Example 2: Job with Custom Queue

```go
type SendSmsJob struct {
    job.Base
    Mobile string `json:"mobile"`
    Body   string `json:"body"`
}

func (j *SendSmsJob) Handle(ctx context.Context) error {
    return smsProvider.Send(j.Mobile, j.Body)
}

func (j *SendSmsJob) Queue() string { return "sms" }  // custom queue

func init() {
    job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}

// Dispatch
job.Dispatch(&SendSmsJob{Mobile: "09121234567", Body: "Your code: 123456"})

// Worker (main.go)
job.NewWorker(broker, "sms").Start()
```

### Example 3: Multiple Workers (Load Balancing)

```go
// 3 workers on same queue - each job goes to ONE worker
for i := 0; i < 3; i++ {
    job.NewWorker(broker, "sms").Start()
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
// Different jobs, different queues
func (j *SendSmsJob) Queue() string { return "sms" }
func (j *SendEmailJob) Queue() string { return "email" }
func (j *ProcessOrderJob) Queue() string { return "order" }

// Start workers for each queue
job.NewWorker(broker, "sms").Start()
job.NewWorker(broker, "email").Start()
job.NewWorker(broker, "order").Start()
```

### Example 5: Dispatch to Specific Queue

```go
// Job has default queue "email"
job.Dispatch(&SendEmailJob{To: "user@test.com"})

// Override queue at dispatch time
job.DispatchOnQueue(&SendEmailJob{To: "user@test.com"}, "priority")
```

### Example 6: Durable Jobs (Survive Restart)

```go
// Startup: create stream
d.EnsureStream("JOBS", []string{"jobs.>"})

// Worker with durable subscription
w := job.NewWorker(broker, "sms")
w.StartDurable("sms-worker-1")
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
// broker/broker.go
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

**Subscribe vs QueueSubscribe:**

| Method | Behavior | Use Case |
|--------|----------|----------|
| `Subscribe` | Every subscriber gets every message (fan-out) | Broadcasting, events |
| `QueueSubscribe` | One subscriber per message (load-balanced) | Job workers |

**Durable (JetStream):**
- Messages survive broker restarts
- Automatic retry on failure
- Stream must be created before publishing

### Job Package (`pkg/job/`)

#### Job Struct

```go
type Job struct {
    ID        string
    Type      string    // e.g., "SendSmsJob"
    Payload   []byte    // JSON-serialized job data
    Status    Status    // pending, processing, success, failed
    Attempts  int
    MaxRetry  int
    CreatedAt time.Time
    RunAt     time.Time
}
```

#### JobHandler Interface

```go
type JobHandler interface {
    Handle(ctx context.Context) error
    Queue() string
}

type Base struct{}
func (b Base) Queue() string { return "default" }
```

#### Registry

Maps job type names to factory functions for deserialization:

```go
var registry = map[string]Factory{}

func RegisterJob(name string, factory Factory) {
    registry[name] = factory
}
```

#### Dispatch

```go
func Dispatch(j JobHandler) error              // uses j.Queue()
func DispatchOnQueue(j JobHandler, queue string) // explicit queue
```

Flow:
1. Marshal job to JSON
2. Create `Job` struct with ID, type name, payload
3. Publish to `jobs.<queue>` subject

#### Worker

```go
func NewWorker(broker Broker, queue string) *Worker
func (w *Worker) Start() error                    // QueueSubscribe
func (w *Worker) StartDurable(durableString) error  // JetStream
```

Flow:
1. Subscribe to `jobs.<queue>` with QueueSubscribe
2. On message: unmarshal `Job` struct
3. Look up factory in registry by `Type`
4. Create handler instance, unmarshal payload
5. Call `handler.Handle(ctx)`

### Complete Example: OTP SMS

**Step 1: Define the Job**

```go
// app/modules/otp/sms_job.go
package otp

import (
    "context"
    "github.com/zhitoo/golang-web-api/pkg/job"
)

type SendSmsJob struct {
    job.Base
    Mobile string `json:"mobile"`
    Body   string `json:"body"`
}

func (j *SendSmsJob) Handle(ctx context.Context) error {
    return smsProvider.Send(j.Mobile, j.Body)
}

func (j *SendSmsJob) Queue() string { return "sms" }

func init() {
    job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}
```

**Step 2: Dispatch the Job**

```go
// app/modules/otp/service.go
func SendOTP(mobile string) (string, error) {
    otp := generateOTP()

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
// cmd/main.go
func main() {
    d := job.NewDispatcher(broker)
    job.SetDispatcher(d)

    smsWorker := job.NewWorker(broker, "sms")
    smsWorker.Start()
}
```

**Step 4: Run**

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
├── job.go          # Job struct, Broker, Dispatcher, Consumer
├── handler.go      # JobHandler interface, Base struct
├── registry.go     # RegisterJob, factory registry
├── dispatch.go     # Dispatch, DispatchOnQueue
└── worker.go       # Worker with QueueSubscribe

broker/
├── broker.go       # Broker interface
├── factory.go      # NewBroker factory
└── nats/
    └── nats.go     # NATS implementation

app/modules/otp/
├── sms_job.go      # SendSmsJob definition
└── service.go      # Dispatches SendSmsJob in SendOTP()
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
