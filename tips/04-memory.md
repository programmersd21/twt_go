# 04 — Memory: Stack, Heap, Escape Analysis, and Object Reuse

---

## Stack vs Heap: What Actually Happens

### Stack

- Allocated and freed at function call/return — O(1), no GC
- Size grows dynamically (starting at 8KB per goroutine)
- Variables that don't escape the function live here
- LIFO allocation means pointer arithmetic is trivial

### Heap

- Managed by the Go GC (tricolor mark-and-sweep)
- Every allocation adds to GC scan work
- GC pauses (STW for root scanning, concurrent marking)
- Fragmentation over time increases GC scan time

```go
// Stack allocation — compiler can prove it doesn't escape
func add(a, b int) int {
    result := a + b // lives on stack
    return result   // value copied out, result freed immediately
}

// Heap allocation — escapes because it's returned by pointer
func newBuffer() *bytes.Buffer {
    buf := bytes.Buffer{} // escapes to heap
    return &buf
}
```

---

## Escape Analysis In Depth

### What Triggers Escape

The compiler runs escape analysis to determine whether a variable can remain on the stack or must be promoted to the heap. Key triggers:

**1. Address taken and returned:**

```go
func escape1() *int {
    x := 42
    return &x // x escapes — outlives the function frame
}
```

**2. Assigned to an interface:**

```go
func escape2() {
    x := 42
    var i any = x // x is boxed — escapes to heap
    _ = i
}
```

**3. Captured by a goroutine closure:**

```go
func escape3() {
    x := 42
    go func() {
        fmt.Println(x) // x escapes — goroutine may outlive the current frame
    }()
}
```

**4. Size not known at compile time:**

```go
func escape4(n int) []byte {
    return make([]byte, n) // n unknown at compile time — heap allocated
}
```

**5. Too large for stack:**

```go
func escape5() {
    var arr [1 << 16]byte // 64KB — exceeds stack budget, escapes to heap
    _ = arr
}
```

**6. Passed to a function that stores it:**

```go
type Server struct{ handler http.Handler }

func (s *Server) setHandler(h http.Handler) {
    s.handler = h // h escapes — stored past the function call
}
```

### Reading Escape Analysis Output

```bash
go build -gcflags='-m=2' ./... 2>&1 | grep escape
```

```
./handler.go:15:6: moved to heap: req
./cache.go:42:14: &entry escapes to heap
./pool.go:8:12: make([]byte, n) escapes to heap
```

`-m=1` gives high-level escapes, `-m=2` gives reasons. Always use `-m=2` when debugging.

### Non-Escaping Interface Trick

```go
// This escapes (interface causes boxing):
func logValue(v any) {
    fmt.Println(v)
}
logValue(42) // 42 escapes

// This doesn't (concrete type, inlinable):
func logInt(v int) {
    fmt.Println(v)
}
logInt(42) // stays on stack
```

---

## sync.Pool: Object Reuse Without GC Pressure

### Correct Pool Usage

```go
var encoderPool = sync.Pool{
    New: func() any {
        return &json.Encoder{}
    },
}

// More realistic: pool the buffer the encoder writes into
var bufPool = sync.Pool{
    New: func() any {
        return bytes.NewBuffer(make([]byte, 0, 512))
    },
}

func encodeToBytes(v any) ([]byte, error) {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset() // ALWAYS reset before use
    defer bufPool.Put(buf)

    if err := json.NewEncoder(buf).Encode(v); err != nil {
        return nil, err
    }

    // Must copy out — buf will be returned to pool
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result, nil
}
```

### Pool Semantics You Must Internalize

```go
// Pool items are NOT guaranteed to survive between GC runs.
// GC can clear the pool at any time.
// DO NOT use Pool as a cache with guaranteed persistence.

// Pool is per-P (processor), so access is mostly lock-free.
// Get/Put are fast — suitable for allocations in hot paths.

// If New is nil and pool is empty, Get returns nil.
var badPool sync.Pool // no New — Get() can return nil
obj := badPool.Get() // obj is nil
_ = obj.(*MyType)    // PANIC
```

### Pool Anti-Patterns

```go
// ANTI-PATTERN 1: storing objects with pointers to outside state
type Worker struct {
    conn *net.Conn // pool can't safely reset pointer state
}

// ANTI-PATTERN 2: not resetting before use
buf := pool.Get().(*bytes.Buffer)
// forgot buf.Reset() — previous contents leak into this call

// ANTI-PATTERN 3: returning a Pool object after storing it elsewhere
buf := pool.Get().(*bytes.Buffer)
s.buffers = append(s.buffers, buf) // now buf is in two places
pool.Put(buf)                      // pool recycles it, s.buffers has stale ref
```

---

## String Handling Optimization

### Strings Are Immutable Byte Slices

```go
// string header in Go: { ptr *byte, len int } — 16 bytes on 64-bit
// Strings are immutable. Any modification requires a new allocation.

s := "hello"
b := []byte(s) // COPY — new backing array
b[0] = 'H'
s2 := string(b) // COPY — new string
```

### Interning Frequently Used Strings

```go
// For high-cardinality repeated strings (e.g., user agents, status codes):
type StringInterner struct {
    mu sync.RWMutex
    m  map[string]string
}

func (si *StringInterner) Intern(s string) string {
    si.mu.RLock()
    if interned, ok := si.m[s]; ok {
        si.mu.RUnlock()
        return interned
    }
    si.mu.RUnlock()

    si.mu.Lock()
    defer si.mu.Unlock()
    if interned, ok := si.m[s]; ok {
        return interned
    }
    si.m[s] = s
    return s
}
```

### strings.Builder vs bytes.Buffer

```go
// strings.Builder: optimized for building strings
// - Never reallocates for Grow() pre-sized builds
// - WriteString is the fast path
var sb strings.Builder
sb.Grow(256) // pre-size to avoid reallocations
for _, part := range parts {
    sb.WriteString(part)
}
result := sb.String() // single allocation to produce final string

// bytes.Buffer: when you need []byte output or io.Writer interface
var buf bytes.Buffer
buf.Grow(256)
for _, part := range parts {
    buf.WriteString(part)
}
result := buf.Bytes() // zero-copy if you don't need a string
```

### String Conversion in Hot Loops

```go
// Benchmark: string(key) in map lookup allocates if key is []byte
// Workaround for Go 1.20+: use string(key) in map index — compiler optimizes
// this specific pattern to avoid allocation in map lookups

var m map[string]int
key := []byte("some-key")

// This specific form is optimized by the compiler — no allocation:
v := m[string(key)]

// But this is NOT optimized — allocates:
k := string(key)
v = m[k]
```

---

## Object Reuse Strategies

### Pre-Allocated Ring Buffers

```go
type RingBuffer[T any] struct {
    buf  []T
    head int
    tail int
    size int
    cap  int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
    return &RingBuffer[T]{
        buf: make([]T, capacity),
        cap: capacity,
    }
}

func (r *RingBuffer[T]) Push(v T) bool {
    if r.size == r.cap {
        return false // full
    }
    r.buf[r.tail] = v
    r.tail = (r.tail + 1) % r.cap
    r.size++
    return true
}

func (r *RingBuffer[T]) Pop() (T, bool) {
    var zero T
    if r.size == 0 {
        return zero, false
    }
    v := r.buf[r.head]
    r.buf[r.head] = zero // clear to allow GC of pointed-to values
    r.head = (r.head + 1) % r.cap
    r.size--
    return v, true
}
```

### Clearing Slices Safely

```go
// Go 1.21+: clear() builtin zeros a slice in one operation
items := make([]*Order, 100)
// ... fill items ...
clear(items) // all pointers zeroed — GC can collect the Orders

// Pre-1.21:
for i := range items {
    items[i] = nil
}
items = items[:0]
```

### Arena-Style Allocation for Request Lifetimes

```go
// For objects that all die together (e.g., per-request), use a single
// large allocation and subdivide manually:

type Arena struct {
    buf  []byte
    used int
}

func NewArena(size int) *Arena {
    return &Arena{buf: make([]byte, size)}
}

func (a *Arena) Alloc(n int) []byte {
    if a.used+n > len(a.buf) {
        return make([]byte, n) // fall back to heap
    }
    slice := a.buf[a.used : a.used+n]
    a.used += n
    return slice
}

// Reset at end of request — all sub-allocations freed in one step
func (a *Arena) Reset() {
    a.used = 0
    // buf backing array is reused; no GC work needed
}
```

---

## GC Tuning

### GOGC and GOMEMLIMIT

```go
// GOGC: controls heap growth trigger (default 100 = GC when heap doubles)
// GOMEMLIMIT: hard memory limit (Go 1.19+)

// In container environments, always set GOMEMLIMIT to ~90% of container limit:
// GOMEMLIMIT=900MiB go run .

// Or programmatically:
import "runtime/debug"

func init() {
    debug.SetMemoryLimit(900 * 1024 * 1024) // 900 MB
    debug.SetGCPercent(50) // GC sooner, smaller pauses
}
```

### Triggering GC Manually (Rarely Justified)

```go
// After a known large allocation spike — e.g., loading a large dataset
// that's been processed and can now be freed:
data := loadHugeDataset()
process(data)
data = nil
runtime.GC() // force immediate collection
```

### Reading GC Stats

```go
var stats runtime.MemStats
runtime.ReadMemStats(&stats)

fmt.Printf("Alloc: %v MiB\n", stats.Alloc/1024/1024)
fmt.Printf("TotalAlloc: %v MiB\n", stats.TotalAlloc/1024/1024)
fmt.Printf("NumGC: %v\n", stats.NumGC)
fmt.Printf("PauseNs last: %v\n", stats.PauseNs[(stats.NumGC+255)%256])
```

**Key metrics to watch:**
- `HeapInuse` vs `HeapSys`: fragmentation indicator
- `PauseNs` histogram: GC pause latency
- `NumGC` over time: GC frequency
- `NextGC`: next collection trigger point
