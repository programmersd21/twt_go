# 03 — Concurrency: Production Go Patterns

---

## The Goroutine Model

Goroutines are cheap (initial stack ~2KB) but not free. Each goroutine in existence is tracked by the scheduler and GC. Unbounded goroutine creation is a resource leak.

**Rules:**
- Every goroutine must have a defined lifetime
- Every goroutine must have a way to be stopped
- Never launch a goroutine without knowing who owns its cleanup

```go
// BAD: fire-and-forget with no lifecycle control
func HandleEvent(e Event) {
    go process(e) // leaked if process panics or hangs
}

// GOOD: goroutine lifecycle is owned by the caller
func (w *Worker) HandleEvent(ctx context.Context, e Event) {
    w.wg.Add(1)
    go func() {
        defer w.wg.Done()
        defer func() {
            if r := recover(); r != nil {
                w.logger.Error("panic in worker", "recovered", r)
            }
        }()
        w.process(ctx, e)
    }()
}

func (w *Worker) Shutdown() {
    w.wg.Wait()
}
```

---

## Channels: Buffered vs Unbuffered

### Unbuffered Channel = Synchronization Point

Both sender and receiver must be ready simultaneously. Use for hand-offs and signaling.

```go
// Signal completion
done := make(chan struct{})
go func() {
    doWork()
    close(done) // signal, not send
}()
<-done // blocks until doWork completes
```

### Buffered Channel = Decoupling + Throughput

Producer doesn't block until buffer is full. Buffer size is a capacity contract.

```go
// Producer/consumer with controlled buffering
jobs := make(chan Job, 100) // producer can be 100 ahead before blocking

// Producer
go func() {
    defer close(jobs)
    for _, j := range jobList {
        jobs <- j
    }
}()

// Consumer
for j := range jobs {
    process(j)
}
```

### Channel Direction in Signatures

Lock down channel direction in function signatures — compiler enforces it.

```go
func producer(out chan<- int) {  // send-only
    for i := 0; i < 10; i++ {
        out <- i
    }
    close(out)
}

func consumer(in <-chan int) {  // receive-only
    for v := range in {
        fmt.Println(v)
    }
}
```

### Closing Channels: Rules

- Only the **sender** closes a channel
- Closing a nil channel panics
- Sending to a closed channel panics
- Multiple close() calls panic — use `sync.Once` if multiple closers

```go
type SafeChannel[T any] struct {
    ch     chan T
    once   sync.Once
}

func (sc *SafeChannel[T]) Close() {
    sc.once.Do(func() { close(sc.ch) })
}
```

---

## Worker Pool: Production Pattern

### Fixed Worker Pool with Context Cancellation

```go
type WorkerPool struct {
    jobs    chan Job
    results chan Result
    wg      sync.WaitGroup
}

func NewWorkerPool(workers, queueSize int) *WorkerPool {
    return &WorkerPool{
        jobs:    make(chan Job, queueSize),
        results: make(chan Result, queueSize),
    }
}

func (p *WorkerPool) Start(ctx context.Context, workers int) {
    for i := 0; i < workers; i++ {
        p.wg.Add(1)
        go func() {
            defer p.wg.Done()
            for {
                select {
                case job, ok := <-p.jobs:
                    if !ok {
                        return
                    }
                    result := job.Execute()
                    select {
                    case p.results <- result:
                    case <-ctx.Done():
                        return
                    }
                case <-ctx.Done():
                    return
                }
            }
        }()
    }
}

func (p *WorkerPool) Submit(j Job) {
    p.jobs <- j
}

func (p *WorkerPool) Shutdown() {
    close(p.jobs)
    p.wg.Wait()
    close(p.results)
}
```

### Semaphore Pattern for Bounded Concurrency

```go
// Limit concurrent goroutines without a full worker pool
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore {
    return make(chan struct{}, n)
}

func (s Semaphore) Acquire() { s <- struct{}{} }
func (s Semaphore) Release() { <-s }

// Usage: limit to 10 concurrent HTTP requests
sem := NewSemaphore(10)
var wg sync.WaitGroup

for _, url := range urls {
    wg.Add(1)
    go func(u string) {
        defer wg.Done()
        sem.Acquire()
        defer sem.Release()
        fetch(u)
    }(url)
}
wg.Wait()
```

### errgroup: Canonical Goroutine Fan-Out

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, ids []string) ([]*User, error) {
    g, ctx := errgroup.WithContext(ctx)
    users := make([]*User, len(ids))

    for i, id := range ids {
        i, id := i, id // capture loop vars (pre Go 1.22)
        g.Go(func() error {
            u, err := fetchUser(ctx, id)
            if err != nil {
                return fmt.Errorf("fetch user %s: %w", id, err)
            }
            users[i] = u
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return users, nil
}
```

With a concurrency limit:

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(10) // max 10 goroutines at a time

for _, id := range ids {
    id := id
    g.Go(func() error {
        return processID(ctx, id)
    })
}
```

---

## Mutex vs Channel: The Decision

| Scenario | Use |
|---|---|
| Protecting shared state | `sync.Mutex` or `sync.RWMutex` |
| Ownership transfer | Channel |
| Fan-out/fan-in pipelines | Channel |
| Simple counters/flags | `sync/atomic` |
| Rate limiting, throttling | Buffered channel (semaphore) |
| Publish-subscribe | Channel + goroutine per subscriber |
| Communicating sequential processes | Channel |
| Cache/map with reads >> writes | `sync.RWMutex` |

### Mutex Best Practices

```go
type OrderCache struct {
    mu    sync.RWMutex
    cache map[string]*Order
}

func (c *OrderCache) Get(id string) (*Order, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    o, ok := c.cache[id]
    return o, ok
}

func (c *OrderCache) Set(id string, o *Order) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[id] = o
}

// Lock contention pattern: minimize lock scope
func (c *OrderCache) Refresh(id string, fetch func() (*Order, error)) error {
    // Check with read lock first
    c.mu.RLock()
    _, exists := c.cache[id]
    c.mu.RUnlock()

    if exists {
        return nil
    }

    // Fetch outside the lock
    o, err := fetch()
    if err != nil {
        return err
    }

    // Write with write lock
    c.mu.Lock()
    c.cache[id] = o
    c.mu.Unlock()
    return nil
}
```

### sync.RWMutex: Only Worth It When Read-Heavy

RWMutex has higher overhead than Mutex for write-heavy workloads. Benchmark before assuming it helps.

---

## Select Statement Patterns

### Priority Select: Drain High-Priority First

```go
// Go's select is random among ready cases.
// For priority, use nested selects:
func (w *Worker) run(ctx context.Context) {
    for {
        // Check urgent channel first (non-blocking)
        select {
        case job := <-w.urgentJobs:
            w.process(job)
            continue
        default:
        }

        // Then block on any
        select {
        case job := <-w.urgentJobs:
            w.process(job)
        case job := <-w.normalJobs:
            w.process(job)
        case <-ctx.Done():
            return
        }
    }
}
```

### Timeout Pattern

```go
func fetchWithTimeout(ctx context.Context, id string) (*Result, error) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    resultCh := make(chan *Result, 1)
    errCh := make(chan error, 1)

    go func() {
        r, err := expensiveFetch(id)
        if err != nil {
            errCh <- err
            return
        }
        resultCh <- r
    }()

    select {
    case r := <-resultCh:
        return r, nil
    case err := <-errCh:
        return nil, err
    case <-ctx.Done():
        return nil, fmt.Errorf("fetchWithTimeout: %w", ctx.Err())
    }
}
```

### Tick + Done Pattern

```go
func (s *Syncer) run(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := s.sync(ctx); err != nil {
                s.logger.Error("sync failed", "err", err)
            }
        case <-ctx.Done():
            s.logger.Info("syncer stopping", "reason", ctx.Err())
            return
        }
    }
}
```

---

## Race Condition Avoidance

### Always Run the Race Detector

```bash
go test -race ./...
go run -race main.go
```

The race detector adds ~5–20x overhead — use in CI, not production.

### Common Race Patterns and Fixes

**Loop variable capture (pre Go 1.22):**

```go
// BUG: all goroutines share the same `item` variable
for _, item := range items {
    go func() {
        process(item) // item is the loop variable, data race
    }()
}

// FIX: Go 1.22+ loop variables are per-iteration — no capture needed
// For Go 1.21 and below:
for _, item := range items {
    item := item // rebind
    go func() {
        process(item)
    }()
}
```

**Sharing slices across goroutines:**

```go
// BAD: concurrent append on shared slice
var results []int
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        results = append(results, n*n) // DATA RACE
    }(i)
}

// GOOD: pre-allocate, write to known index
results := make([]int, 100)
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        results[n] = n * n // safe: each goroutine owns its index
    }(i)
}
```

**Map concurrent access:**

```go
// sync.Map for high-concurrency, mostly-read scenarios
var m sync.Map

go func() { m.Store("key", "value") }()
go func() {
    if v, ok := m.Load("key"); ok {
        fmt.Println(v)
    }
}()

// For write-heavy: prefer map + sync.RWMutex
// sync.Map has higher overhead than a guarded map in write-heavy scenarios
```

### Context Propagation: The Contract

```go
// Context MUST be the first param and MUST be passed through
func (s *Service) GetOrder(ctx context.Context, id string) (*Order, error) {
    // Never store context in struct fields
    // Never pass nil context — use context.Background() or context.TODO()
    return s.repo.FindByID(ctx, id)
}

// Cancellation propagates automatically — check it in long loops
func processLargeDataset(ctx context.Context, items []Item) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        if err := processItem(ctx, item); err != nil {
            return err
        }
    }
    return nil
}
```
