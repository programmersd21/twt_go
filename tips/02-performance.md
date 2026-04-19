# 02 — Performance: Go Memory, Allocations, and GC Pressure

---

## The Allocator Is Your Enemy

Every heap allocation costs: GC scan time, potential STW pause, pointer indirection. The goal is not zero allocations — it's *intentional* allocations only.

```go
// Measure first. Always.
// go test -bench=. -benchmem -count=5 ./...
// -benchmem shows allocs/op and B/op — the two metrics that matter
```

---

## Slice Preallocation

### The Pattern: Know Your Size

```go
// BAD: O(log n) allocations due to growth
func collectIDs(orders []*Order) []string {
    var ids []string
    for _, o := range orders {
        ids = append(ids, o.ID)
    }
    return ids
}

// GOOD: 1 allocation, exact capacity
func collectIDs(orders []*Order) []string {
    ids := make([]string, 0, len(orders))
    for _, o := range orders {
        ids = append(ids, o.ID)
    }
    return ids
}
```

### When You Don't Know the Size

```go
// Use a growth estimate — even a rough one beats zero
func filterActive(orders []*Order) []*Order {
    // assume ~50% active: better than starting at 0
    result := make([]*Order, 0, len(orders)/2)
    for _, o := range orders {
        if o.Status == StatusActive {
            result = append(result, o)
        }
    }
    return result
}
```

### Slice Reuse with Reset

```go
// Reuse the backing array between calls
type Processor struct {
    buf []*Event
}

func (p *Processor) Process(events []*Event) {
    p.buf = p.buf[:0]            // reset length, keep capacity
    for _, e := range events {
        if e.Type == TypeImportant {
            p.buf = append(p.buf, e)
        }
    }
    p.flush(p.buf)
}
```

---

## Memory Allocation Minimization

### Stack-Allocated Values Stay Fast

```go
// This stays on stack if it doesn't escape
type Point struct{ X, Y float64 }

func distance(a, b Point) float64 {
    dx := a.X - b.X
    dy := a.Y - b.Y
    return math.Sqrt(dx*dx + dy*dy)
}
// Point is passed by value — no heap allocation, no GC involvement
```

### Avoid Pointer Returns for Small, Frequent Types

```go
// BAD: causes heap allocation, GC tracking
func newPoint(x, y float64) *Point {
    return &Point{x, y}
}

// GOOD: returned by value, stays on caller's stack
func newPoint(x, y float64) Point {
    return Point{x, y}
}
```

### Maps: Pre-size Aggressively

```go
// BAD: repeated rehashing
m := make(map[string]int)
for _, item := range largeSlice {
    m[item.Key]++
}

// GOOD: single allocation
m := make(map[string]int, len(largeSlice))
for _, item := range largeSlice {
    m[item.Key]++
}
```

### Struct Layout: Minimize Padding

```go
// BAD: 24 bytes due to padding
type Bad struct {
    A bool    // 1 byte + 7 pad
    B float64 // 8 bytes
    C bool    // 1 byte + 7 pad
}

// GOOD: 10 bytes (+ 6 pad to align = 16 bytes, still better)
type Good struct {
    B float64 // 8 bytes
    A bool    // 1 byte
    C bool    // 1 byte + 6 pad
}

// Use: go vet -composites or fieldalignment linter
// golang.org/x/tools/go/analysis/passes/fieldalignment
```

---

## String Optimization

### strings.Builder for Concatenation Loops

```go
// BAD: O(n²) allocations — new string every iteration
func joinParts(parts []string) string {
    var result string
    for _, p := range parts {
        result += p + ", "
    }
    return result
}

// GOOD: single allocation
func joinParts(parts []string) string {
    var b strings.Builder
    b.Grow(estimateSize(parts)) // pre-size if you can estimate
    for i, p := range parts {
        b.WriteString(p)
        if i < len(parts)-1 {
            b.WriteString(", ")
        }
    }
    return b.String()
}

// BETTER for this specific case: stdlib already does it
strings.Join(parts, ", ")
```

### []byte ↔ string Zero-Copy Conversions

```go
// Standard conversion allocates. For hot paths, use unsafe.
import "unsafe"

// string → []byte without allocation (READ ONLY — do not mutate the slice)
func stringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

// []byte → string without allocation
func bytesToString(b []byte) string {
    return unsafe.String(unsafe.SliceData(b), len(b))
}
// Use these only in hot paths where you can guarantee no mutation.
// As of Go 1.20, unsafe.StringData and unsafe.SliceData are the safe way.
```

### Avoid string([]byte) in Log/Format Hot Paths

```go
// BAD: allocates new string on every call
func (h *Handler) log(data []byte) {
    slog.Info("received", "body", string(data))
}

// GOOD: pass []byte directly if logger supports it, or cache the conversion
```

---

## GC Pressure Reduction

### sync.Pool for Temporary Objects

```go
var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func processRequest(data []byte) ([]byte, error) {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf)

    // use buf freely
    if err := json.NewEncoder(buf).Encode(data); err != nil {
        return nil, err
    }
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result, nil
}
```

**Pool rules:**
- Objects must be safe to reset and reuse
- Never store pointers to Pool objects outside Pool lifecycle
- Pool contents are cleared on GC — don't use for persistent state
- Always `buf.Reset()` before use, not after (so you don't miss cleanup on panic)

### Avoid Interface Boxing in Hot Paths

```go
// BAD: every int stored in interface{} allocates
func sumInterface(vals []any) int {
    sum := 0
    for _, v := range vals {
        sum += v.(int)
    }
    return sum
}

// GOOD: typed slice, zero boxing
func sumTyped(vals []int) int {
    sum := 0
    for _, v := range vals {
        sum += v
    }
    return sum
}
```

### Value Types in Channels Over Pointers (Small Structs)

```go
// For small structs (≤ 64 bytes): value is cheaper — avoids heap escape
type Event struct {
    Type      uint8
    Timestamp int64
    UserID    uint32
} // 13 bytes

ch := make(chan Event, 1024) // values stored inline in channel buffer
```

---

## Benchmarking Mindset

### Benchmark Structure That Actually Measures Things

```go
func BenchmarkProcessOrder(b *testing.B) {
    // Setup outside the loop
    orders := generateTestOrders(1000)
    svc := NewOrderService(newMockRepo())

    b.ResetTimer()   // exclude setup from measurement
    b.ReportAllocs() // always report allocations

    for b.Loop() { // Go 1.24+: b.Loop() replaces i < b.N pattern
        for _, o := range orders {
            _, _ = svc.Process(context.Background(), o)
        }
    }
}

// Run with:
// go test -bench=BenchmarkProcessOrder -benchmem -count=5 -benchtime=3s
// Use benchstat for statistical analysis across runs
```

### Sub-benchmarks for Comparative Analysis

```go
func BenchmarkJSON(b *testing.B) {
    payload := generatePayload(1024)

    b.Run("stdlib/marshal", func(b *testing.B) {
        b.ReportAllocs()
        for b.Loop() {
            _, _ = json.Marshal(payload)
        }
    })

    b.Run("sonic/marshal", func(b *testing.B) {
        b.ReportAllocs()
        for b.Loop() {
            _, _ = sonic.Marshal(payload)
        }
    })
}
```

### Prevent Compiler Elimination

```go
var sink any // package-level sink

func BenchmarkCompute(b *testing.B) {
    for b.Loop() {
        result := expensiveComputation(42)
        sink = result // prevent dead-code elimination
    }
}
```

---

## Escape Analysis Awareness

Force the compiler to tell you what escapes:

```bash
go build -gcflags='-m=2' ./...
```

### Common Escape Triggers

```go
// 1. Returning pointer to local variable
func newFoo() *Foo {
    f := Foo{} // escapes to heap
    return &f
}

// 2. Storing in interface
var i any = Foo{} // Foo escapes to heap (boxing)

// 3. Closures capturing variables
x := 42
go func() { fmt.Println(x) }() // x escapes

// 4. Slice too large for stack
buf := make([]byte, 1<<20) // large slices always escape

// 5. Passed to function that stores it
func store(s *Server, h *Handler) {
    s.handler = h // h escapes — stored beyond function lifetime
}
```

### The Stack Budget

Go's default goroutine stack starts at 8KB and grows. Variables up to ~64KB can remain on stack. Beyond that, heap allocation is forced regardless of escape analysis.

```go
// This will always heap-allocate
bigBuf := make([]byte, 65536+1)

// This may stack-allocate (compiler decides)
var smallBuf [512]byte
```

---

## Hot Path Optimization Rules

1. **Measure first.** Profiler before eyeballing.
2. **Reduce allocations before reducing CPU.** GC pauses are non-deterministic.
3. **Batch I/O operations.** 1 syscall with 1000 items beats 1000 syscalls.
4. **Inline small functions.** Keep hot-path functions < 80 AST nodes.
5. **Avoid fmt.Sprintf in loops.** Use strconv for number conversions.

```go
// BAD: fmt.Sprintf allocates for string formatting
for _, id := range ids {
    key := fmt.Sprintf("user:%d", id)
    cache.Get(key)
}

// GOOD: strconv.AppendInt reuses buffer
var key []byte
for _, id := range ids {
    key = append(key[:0], "user:"...)
    key = strconv.AppendInt(key, id, 10)
    cache.Get(string(key))
}
```
