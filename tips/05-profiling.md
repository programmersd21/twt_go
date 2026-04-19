# 05 — Profiling: Measure, Locate, Fix, Verify

---

## The Principle

> Optimize nothing until pprof tells you where to look.

Intuition-driven optimization in Go wastes time and introduces bugs. The profiler costs minutes; guessing costs hours. Every optimization must be followed by a benchmark showing improvement.

---

## pprof Integration

### HTTP Endpoint (Always On in Services)

```go
import (
    _ "net/http/pprof" // registers /debug/pprof/* handlers
    "net/http"
)

func main() {
    // Expose pprof on a separate port — never the public port
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // ... start your actual server
}
```

Available endpoints:
- `/debug/pprof/heap` — heap allocations
- `/debug/pprof/goroutine` — all goroutine stacks
- `/debug/pprof/profile?seconds=30` — 30s CPU profile
- `/debug/pprof/block` — goroutine blocking events
- `/debug/pprof/mutex` — mutex contention
- `/debug/pprof/trace?seconds=5` — execution trace

### Programmatic Profiling

```go
import (
    "os"
    "runtime/pprof"
)

// CPU profile
f, _ := os.Create("cpu.prof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()

// Heap profile (snapshot at this moment)
f2, _ := os.Create("heap.prof")
defer func() {
    runtime.GC() // get up-to-date statistics
    pprof.WriteHeapProfile(f2)
    f2.Close()
}()
```

---

## CPU Profiling Workflow

### Step 1: Capture

```bash
# From a running service
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# From a test/benchmark
go test -bench=BenchmarkFoo -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Step 2: Analyze in pprof Interactive Mode

```
(pprof) top10           # top 10 functions by cumulative CPU
(pprof) top10 -cum      # sort by cumulative (includes callees)
(pprof) list FuncName   # annotated source for a function
(pprof) web             # flame graph in browser (requires graphviz)
(pprof) peek FuncName   # callers/callees of a function
```

### Step 3: Read the Flame Graph

```
Wide bars = function consumes significant CPU
Tall stacks = deep call chains
Flat top = CPU spent in that exact function (no callees)
```

### Step 4: Interpret `top` Output

```
(pprof) top5
Showing top 5 nodes out of 87
      flat  flat%   sum%        cum   cum%
   1.23s   41.0%  41.0%     1.23s  41.0%  runtime.mallocgc
   0.85s   28.3%  69.3%     2.10s  70.0%  encoding/json.Marshal
   0.42s   14.0%  83.3%     0.42s  14.0%  runtime.scanobject
   0.21s    7.0%  90.3%     0.21s   7.0%  sync.(*Mutex).Lock
   0.15s    5.0%  95.3%     0.15s   5.0%  strings.(*Builder).WriteString
```

- `flat`: CPU in this function only
- `cum`: CPU in this function + all it calls
- `mallocgc` at top = allocation problem, not logic problem — go to heap profile

---

## Heap Profiling Workflow

### Capture

```bash
# From service
go tool pprof http://localhost:6060/debug/pprof/heap

# From benchmark
go test -bench=. -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Heap Profile Modes

```bash
# Default: inuse_space (currently allocated, not freed)
go tool pprof -inuse_space mem.prof

# Total allocated (including freed — shows allocation rate)
go tool pprof -alloc_space mem.prof    # most useful for GC tuning

# Object counts
go tool pprof -alloc_objects mem.prof
go tool pprof -inuse_objects mem.prof
```

### What to Look For

```
(pprof) top10 -alloc_space
1024MB  52.1%  encoding/json.(*encodeState).marshal
 512MB  26.1%  bytes.(*Buffer).WriteString
 256MB  13.0%  fmt.Sprintf
```

`encoding/json` and `fmt.Sprintf` in top allocators = switch JSON libraries, use `strconv`.

---

## Benchmark Writing: Production-Grade

### Correct Benchmark Structure

```go
func BenchmarkOrderProcessing(b *testing.B) {
    // 1. Build fixtures once — outside measurement
    orders := make([]*Order, 1000)
    for i := range orders {
        orders[i] = &Order{
            ID:     uuid.New(),
            Amount: decimal.NewFromFloat(float64(i) * 9.99),
            Items:  generateItems(5),
        }
    }

    svc := NewOrderService(
        repository.NewMemoryStore(),
        NewPricingEngine(),
    )
    ctx := context.Background()

    // 2. Reset timer after setup
    b.ResetTimer()
    b.ReportAllocs()

    // 3. Loop with b.Loop() (Go 1.24+) or i < b.N (1.23 and below)
    for b.Loop() {
        for _, o := range orders {
            if _, err := svc.Process(ctx, o); err != nil {
                b.Fatal(err)
            }
        }
    }
}
```

### Parallel Benchmarks

```go
func BenchmarkCacheGet(b *testing.B) {
    cache := NewLRUCache(10000)
    populateCache(cache, 10000)

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        keys := generateKeys(100)
        i := 0
        for pb.Next() {
            cache.Get(keys[i%len(keys)])
            i++
        }
    })
}
```

### Table-Driven Benchmarks

```go
func BenchmarkJSONMarshal(b *testing.B) {
    cases := []struct {
        name string
        size int
    }{
        {"small/16B", 16},
        {"medium/1KB", 1024},
        {"large/100KB", 100 * 1024},
    }

    for _, tc := range cases {
        payload := generatePayload(tc.size)
        b.Run(tc.name, func(b *testing.B) {
            b.SetBytes(int64(tc.size))
            b.ReportAllocs()
            for b.Loop() {
                out, err := json.Marshal(payload)
                if err != nil {
                    b.Fatal(err)
                }
                _ = out
            }
        })
    }
}
```

### benchstat for Statistical Comparison

```bash
# Collect baselines
go test -bench=. -benchmem -count=10 ./... > old.txt

# Make your change, then:
go test -bench=. -benchmem -count=10 ./... > new.txt

# Compare with statistical significance
benchstat old.txt new.txt
```

Output:

```
name              old time/op    new time/op    delta
OrderProcess-8      1.23ms ± 2%    0.87ms ± 1%  -29.27%  (p=0.000 n=10+10)

name              old alloc/op   new alloc/op   delta
OrderProcess-8      48.2kB ± 0%    12.1kB ± 0%  -74.90%  (p=0.000 n=10+10)
```

---

## Goroutine Profile: Finding Leaks

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
(pprof) top
(pprof) traces  # full stack traces for every goroutine state
```

Goroutine count should be bounded. If it grows without limit:

```go
// Check goroutine count in metrics
metrics.Gauge("runtime.goroutines", float64(runtime.NumGoroutine()))

// Common cause: goroutines blocked on channel send/receive
// diagnosis:
(pprof) traces | grep "chan send"
(pprof) traces | grep "chan receive"
```

---

## Execution Tracer

For microsecond-level timing, scheduler behavior, and GC events:

```bash
# Capture 5-second trace
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out
go tool trace trace.out
```

The trace viewer shows:
- Goroutine scheduling latency
- GC stop-the-world events
- Heap growth over time
- Network/syscall blocking

---

## Block and Mutex Profiling

These are disabled by default — enable explicitly.

```go
import "runtime"

func init() {
    runtime.SetBlockProfileRate(1)       // profile every blocking event
    runtime.SetMutexProfileFraction(1)   // profile every mutex contention
    // In production: use 100 or 1000 to sample 1% or 0.1%
}
```

```bash
# After enabling:
go tool pprof http://localhost:6060/debug/pprof/block
go tool pprof http://localhost:6060/debug/pprof/mutex

(pprof) top
# Shows which mutex/channel operations are contended
```

---

## Performance Debugging Workflow

```
1. OBSERVE:    Metrics show latency spike / high CPU / memory growth
2. CAPTURE:    pprof profile during the symptom (not before, not after)
3. IDENTIFY:   top -cum to find the hot function; list FuncName for source
4. HYPOTHESIZE: Allocations? Lock contention? Algorithm? I/O?
5. CHANGE:     One change at a time — no shotgun optimization
6. VERIFY:     benchstat shows statistically significant improvement
7. DEPLOY:     Monitor metrics; confirm improvement holds under real load
8. DOCUMENT:   Record what changed, why, and what the numbers were
```

### Profiling Checklist

```
□ Is the profile captured during actual load? (cold profiling is useless)
□ Did you use -count=10 and benchstat? (single benchmark runs have noise)
□ Did you check alloc_space, not just inuse_space?
□ Did you enable block/mutex profiling if latency is the problem?
□ Did you check goroutine count for leaks?
□ Did you run with GODEBUG=gctrace=1 to see GC frequency?
```

```bash
# GC trace output:
GODEBUG=gctrace=1 go run .
# gc 14 @2.345s 3%: 0.12+1.2+0.04 ms clock, 0.48+0.22/1.1/0+0.16 ms cpu,
#    4->5->2 MB, 5 MB goal, 8 P
# 4->5->2 MB: heap before GC -> heap after marking -> heap after sweep
# 3%: % of runtime spent in GC — should be < 5% for most services
```
