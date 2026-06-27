# Broker & Job System

This document explains how the NATS broker and job queue system work in this project.

## Overview

```
┌─────────────┐     Publish      ┌─────────────┐     Consume      ┌─────────────┐
│   Handler   │ ───────────────▶ │     NATS    │ ───────────────▶ │   Worker    │
│  (Dispatch) │                  │   (Broker)  │                  │  (Process)  │
└─────────────┘                  └─────────────┘                  └─────────────┘
```

The system has three layers:

1. **Broker** (`broker/`) — NATS connection and pub/sub primitives
2. **Job Package** (`pkg/job/`) — low-level dispatcher, consumer, and high-level job abstractions
3. **Business Jobs** (`app/modules/*/`) — concrete job definitions (e.g., SendSmsJob)

## How It Works

### 1. Broker Layer

The broker connects to NATS and provides pub/sub methods:

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

### 2. Job Package (`pkg/job/`)

#### Low-Level Types

```go
// pkg/job/job.go

// Job struct - serialized and sent over NATS
type Job struct {
    ID        string
    Type      string    // job type name (e.g., "SendSmsJob")
    Payload   []byte    // JSON-serialized job data
    Status    Status    // pending, processing, success, failed
    Attempts  int
    MaxRetry  int
    CreatedAt time.Time
    RunAt     time.Time
}

// Dispatcher - publishes jobs to NATS subjects
type Dispatcher struct {
    broker     Broker
    durable    DurableBroker
    useDurable bool
}

// Consumer - subscribes and processes jobs
type Consumer struct {
    broker     Broker
    handler    MessageHandler  // func(Job) error
}
```

#### High-Level Job Abstraction

```go
// pkg/job/handler.go

// JobHandler - what your jobs implement
type JobHandler interface {
    Handle(ctx context.Context) error
    Queue() string  // which queue to use
}

// Base - provides default Queue() returning "default"
type Base struct{}
func (b Base) Queue() string { return "default" }
```

#### Registry

```go
// pkg/job/registry.go

// Maps job type names to factory functions
var registry = map[string]Factory{}

func RegisterJob(name string, factory Factory) {
    registry[name] = factory
}
```

The registry is used for deserialization. When a worker receives a job, it looks up the factory by `Type` name, creates a new instance, unmarshals the payload, and calls `Handle()`.

#### Dispatch

```go
// pkg/job/dispatch.go

func Dispatch(j JobHandler) error              // uses j.Queue()
func DispatchOnQueue(j JobHandler, queue string) // explicit queue
```

Flow:
1. Marshal job to JSON
2. Create `Job` struct with ID, type name, payload
3. Publish to `jobs.<queue>` subject

#### Worker

```go
// pkg/job/worker.go

func NewWorker(broker Broker, queue string) *Worker
func (w *Worker) Start() error                    // QueueSubscribe
func (w *Worker) StartDurable(durableName string) error  // JetStream
```

Flow:
1. Subscribe to `jobs.<queue>` with QueueSubscribe
2. On message: unmarshal `Job` struct
3. Look up factory in registry by `Type`
4. Create handler instance, unmarshal payload
5. Call `handler.Handle(ctx)`

### 3. Business Jobs

Jobs are defined in their respective modules:

```go
// app/modules/otp/sms_job.go

type SendSmsJob struct {
    job.Base                          // provides default Queue()
    Mobile string `json:"mobile"`
    Body   string `json:"body"`
}

func (j *SendSmsJob) Handle(ctx context.Context) error {
    // send SMS logic
    return nil
}

func (j *SendSmsJob) Queue() string { return "sms" }

func init() {
    job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}
```

## Complete Example: OTP SMS

### Step 1: Define the Job

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
    // Call SMS provider API
    return smsProvider.Send(j.Mobile, j.Body)
}

func (j *SendSmsJob) Queue() string { return "sms" }

func init() {
    job.RegisterJob("SendSmsJob", func() job.JobHandler { return &SendSmsJob{} })
}
```

### Step 2: Dispatch the Job

```go
// app/modules/otp/service.go
func SendOTP(mobile string) (string, error) {
    otp := generateOTP()

    // Dispatch SMS job (async, non-blocking)
    job.Dispatch(&SendSmsJob{
        Mobile: mobile,
        Body:   fmt.Sprintf("Your code: %s", otp),
    })

    // Store OTP in Redis
    cache.SetValue("otp:"+mobile, otp, 2*time.Minute)
    return otp, nil
}
```

### Step 3: Start the Worker

```go
// cmd/main.go
func main() {
    // ... broker setup ...

    d := job.NewDispatcher(broker)
    job.SetDispatcher(d)

    // Start SMS worker
    smsWorker := job.NewWorker(broker, "sms")
    smsWorker.Start()

    // ... rest of app ...
}
```

### Step 4: Run

1. User calls `POST /api/v1/otp/send` with `{"mobile": "09121234567"}`
2. Handler calls `SendOTP()`
3. `SendOTP()` dispatches `SendSmsJob` to queue `sms`
4. Worker receives job, calls `Handle()`
5. SMS is sent to user

## Scenarios

### Scenario 1: Single Worker

```
                    ┌──────────────┐
                    │   Worker 1   │
                    │   (sms)      │
POST /otp/send ──▶ │              │ ──▶ SMS Sent
                    └──────────────┘
```

One worker processes all SMS jobs.

### Scenario 2: Multiple Workers (Load Balancing)

```
                    ┌──────────────┐
                    │   Worker 1   │
POST /otp/send ──▶ ├──────────────┤ ──▶ SMS Sent
                    │   Worker 2   │ ──▶ SMS Sent
                    ├──────────────┤
                    │   Worker 3   │ ──▶ SMS Sent
                    └──────────────┘
```

With `QueueSubscribe`, each job goes to only ONE worker. Good for scaling.

```go
// Start 3 workers on same queue
for i := 0; i < 3; i++ {
    w := job.NewWorker(broker, "sms")
    w.Start()
}
```

### Scenario 3: Multiple Queues

```
                    ┌──────────────┐
                    │  SMS Worker  │  queue: sms
POST /otp/send ──▶ │              │ ──▶ SMS Sent
                    └──────────────┘

                    ┌──────────────┐
                    │ Email Worker │  queue: email
POST /register ───▶ │              │ ──▶ Email Sent
                    └──────────────┘
```

Different job types go to different queues. Each queue has its own worker(s).

```go
// Define jobs with different queues
func (j *SendSmsJob) Queue() string { return "sms" }
func (j *SendEmailJob) Queue() string { return "email" }

// Start workers for each queue
job.NewWorker(broker, "sms").Start()
job.NewWorker(broker, "email").Start()
```

### Scenario 4: Durable Jobs (JetStream)

```go
// Create stream on startup
d.EnsureStream("JOBS", []string{"jobs.>"})

// Start durable worker
w := job.NewWorker(broker, "sms")
w.StartDurable("sms-worker-1")
```

- Messages survive broker restart
- Delivery guaranteed
- Good for important jobs (payments, orders)

### Scenario 5: Fan-Out (Broadcast)

```go
// Use Subscribe instead of QueueSubscribe
consumer := job.NewConsumer(broker, func(j job.Job) error {
    log.Printf("Event received: %s", j.Type)
    return nil
})
consumer.Consume("events.user.registered")  // all subscribers get it
```

Use `Subscribe` when you want ALL subscribers to receive every message (e.g., notifications, analytics).

### Scenario 6: Dispatch to Specific Queue

```go
// Job has default queue "email"
job.Dispatch(&SendEmailJob{To: "user@test.com"})

// But you want it processed by "priority" queue instead
job.DispatchOnQueue(&SendEmailJob{To: "user@test.com"}, "priority")
```

## Subject Naming

All job subjects follow the pattern: `jobs.<queue>`

| Queue | Subject | Description |
|-------|---------|-------------|
| `default` | `jobs.default` | Default queue |
| `sms` | `jobs.sms` | SMS jobs |
| `email` | `jobs.email` | Email jobs |
| `order` | `jobs.order` | Order processing |

## File Structure

```
pkg/job/
├── job.go          # Job struct, Broker, Dispatcher, Consumer (low-level)
├── handler.go      # JobHandler interface, Base struct
├── registry.go     # RegisterJob, factory registry
├── dispatch.go     # Dispatch, DispatchOnQueue (global dispatcher)
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

### Check NATS streams

```bash
# List streams
nats stream ls

# Check stream info
nats stream info JOBS

# Monitor messages
nats sub "jobs.sms"
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `nats: no response from stream` | Stream doesn't exist | Call `EnsureStream()` on startup |
| `jetstream not enabled` | JetStream disabled in NATS config | Add `jetstream {}` to nats.conf |
| `unknown job type` | Job not registered | Call `RegisterJob()` in `init()` |
| `dispatcher not initialized` | `SetDispatcher()` not called | Call it in `main.go` after broker connect |
