# 06 — Compiler Optimization: Inlining, Escape, and Hot Path Control

---

## How the Go Compiler Optimizes

The Go compiler (gc) applies these passes in order:
1. Parsing and type checking
2. Escape analysis
3. Inlining
4. Dead code elimination
5. Register allocation (SSA backend)
6. Machine code generation

You interact with 1, 2, and 3. The rest is automatic.

---

## Inlining

### What Gets Inlined

The compiler inlines small functions to eliminate call overhead (stack frame setup, argument passing, return). The budget is measured in "AST nodes" — approximately 80 nodes for direct inlining.

```bash
# See what was inlined and what wasn't
go build -gcflags='-m=1' ./... 2>&1 | grep "inlining call"
go build -gcflags='-m=2' ./... 2>&1 | grep "cannot inline"
```

### Functions That Cannot Be Inlined

```go
// 1. Too complex (> 80 AST nodes)
func complexFunction(...) { ... } // "function too complex for inlining"

// 2. Recursive functions
func factorial(n int) int {
    if n <= 1 { return 1 }
    return n * factorial(n-1) // "recursive: cannot inline"
}

// 3. Functions with closure captures (in some cases)
func withClosure(x int) func() int {
    return func() int { return x } // may not be inlined
}

// 4. Functions containing defer (pre Go 1.14)
// As of Go 1.14+, defer no longer prevents inlining in most cases
```

### Designing for Inlining

```go
// GOOD: short, no side effects, no closures — will inline
func (o *Order) IsExpired() bool {
    return time.Now().After(o.ExpiresAt)
}

// BAD: too many operations — won't inline
func (o *Order) ValidateAndProcess() error {
    if o.ID == uuid.Nil {
        return ErrInvalidID
    }
    if o.Amount.LessThanOrEqual(decimal.Zero) {
        return ErrInvalidAmount
    }
    if time.Now().After(o.ExpiresAt) {
        return ErrExpired
    }
    // 10 more checks...
    return o.process()
}
```

### Forcing Inlining (Go 1.22+)

```go
//go:inline  // hint to the compiler — not a guarantee
func fastPath(x int) int {
    return x * 2
}
```

### Preventing Inlining

```go
//go:noinline
func neverInline(x int) int {
    return x * 2
}
// Use when: benchmarking (to isolate a specific function),
// or debugging (to see it in pprof traces by name)
```

---

## Escape Analysis Flags

```bash
# Level 1: what escapes and where
go build -gcflags='-m' ./...

# Level 2: why it escapes (with reasons)
go build -gcflags='-m=2' ./...

# Example output:
# ./service.go:42:9: &Order literal escapes to heap
# ./service.go:67:14: moved to heap: result
# ./handler.go:23:6: parameter req does not escape
```

### Escape Analysis in Tests

```go
// Write a test that asserts no heap allocations
func TestNoAlloc(t *testing.T) {
    s := NewFastService()
    allocs := testing.AllocsPerRun(100, func() {
        _ = s.HotPath(42)
    })
    if allocs > 0 {
        t.Fatalf("expected 0 allocations, got %v", allocs)
    }
}
```

---

## Abstraction Cost

### Interface Dispatch vs Direct Call

Interface method calls involve an indirect function call through the itab pointer. In tight loops, this matters.

```go
// Direct call — compiler can inline, no vtable lookup
type ConcreteProcessor struct{}
func (p *ConcreteProcessor) Process(x int) int { return x * 2 }

// Interface call — indirect, cannot inline through interface
type Processor interface {
    Process(x int) int
}

func processAll(p Processor, vals []int) []int {
    out := make([]int, len(vals))
    for i, v := range vals {
        out[i] = p.Process(v) // indirect call every iteration
    }
    return out
}

// In benchmarks: interface dispatch adds ~1-3ns per call
// For millions of iterations, this compounds
```

### Generic Functions vs Interface Functions

```go
// Go 1.18+ generics: compiler can specialize per concrete type
// Dictionary-based dispatch (default) = similar to interface
// Stenciling (with GCShape) = closer to direct call

// The compiler uses GCShape stenciling: one instantiation per GC shape.
// *T and T share a shape group if they have the same memory layout.
// For performance-critical generic code, benchmark against interface version.

func ProcessGeneric[T interface{ Process(int) int }](p T, vals []int) []int {
    out := make([]int, len(vals))
    for i, v := range vals {
        out[i] = p.Process(v) // may be devirtualized if type is known
    }
    return out
}
```

### Function Values Are Closures — They Escape

```go
// Storing a func value in a struct causes a heap allocation
type Handler struct {
    onEvent func(Event) // this func value will escape
}

// For hot paths, prefer interface or concrete method over func value
type EventHandler interface {
    OnEvent(Event)
}
```

---

## Compiler Directives Reference

```go
//go:noescape
// Tells compiler the function doesn't cause its pointer args to escape.
// Only valid on external (assembly) functions.

//go:nosplit
// Prevents stack split check — use only in very hot runtime-level code.
// Dangerous: can cause stack overflow if not carefully bounded.

//go:noinline
// Prevents inlining. Use for benchmarks, debugging, or stack trace clarity.

//go:inline (Go 1.22+)
// Hint to inline even if slightly over budget. Not a hard guarantee.

//go:linkname localname importpath.name
// Access unexported symbols from other packages. Fragile; avoid in userland.

//go:generate command
// Code generation marker. Executed by `go generate`.
```

---

## Hot Path Optimization Rules

### Rule 1: Eliminate Allocations Before Reducing CPU

A function with 0 allocations running 3x slower than one with 10 allocations will win in GC-heavy workloads. GC pauses are unpredictable; allocation-free code is deterministic.

### Rule 2: Move Type Assertions Out of Loops

```go
// BAD: type assertion on every iteration
func processItems(items []any) {
    for _, item := range items {
        order := item.(*Order) // repeated assertion overhead
        process(order)
    }
}

// GOOD: type the input correctly; if you must use any, assert once outside
func processItems(items []*Order) { // use concrete type
    for _, item := range items {
        process(item)
    }
}
```

### Rule 3: Batch Small Writes

```go
// BAD: one syscall per write
for _, record := range records {
    conn.Write(record.Marshal())
}

// GOOD: buffer writes, one syscall for the batch
w := bufio.NewWriterSize(conn, 64*1024)
for _, record := range records {
    w.Write(record.Marshal())
}
w.Flush()
```

### Rule 4: Avoid fmt in Hot Paths

```go
// fmt.Sprintf allocates. For key construction, use strconv + []byte tricks.

// BAD
key := fmt.Sprintf("cache:%s:%d", userID, version)

// GOOD
var key [64]byte
n := copy(key[:], "cache:")
n += copy(key[n:], userID)
key[n] = ':'
n++
n += len(strconv.AppendInt(key[n:n], int64(version), 10))
cacheKey := string(key[:n]) // single allocation at the end
```

### Rule 5: Range Over Slice of Values vs Pointers

```go
// Slice of values: elements are laid out contiguously in memory.
// Cache-friendly traversal.
type Point struct{ X, Y float64 }
points := []Point{{1, 2}, {3, 4}, {5, 6}}
for _, p := range points { // p copied by value — that's fine for small structs
    _ = p.X + p.Y
}

// Slice of pointers: each element is a pointer to a separate heap allocation.
// Each access is a potential cache miss.
ptrs := []*Point{...}
for _, p := range ptrs { // dereference = potential cache miss
    _ = p.X + p.Y
}

// Benchmark difference grows with slice size.
// Rule: use value slices for hot iteration; pointer slices for mutation.
```

---

## Type Performance Hierarchy

From fastest to slowest in terms of operation cost:

```
1. bool, int8-64, uint8-64, float32/64      — register ops, zero allocation
2. [N]T (fixed array, small)                — stack-allocated, cache-friendly
3. struct (value, small, no pointers)       — same as fixed array
4. string                                   — header on stack, data may be static
5. []T (slice, pre-allocated)               — 3-word header, contiguous backing
6. map[K]V                                  — hash table, GC-tracked
7. *T (pointer to small struct)             — indirect access, may miss cache
8. interface{ M() }                         — 2-word (itab+data), indirect call
9. func(T) T                               — closure, potential heap allocation
10. chan T                                  — runtime object, full synchronization
```

---

## Bounds Check Elimination (BCE)

The compiler proves bounds at compile time and elides runtime checks.

```go
// Hint the compiler by asserting length up front
func sumFirst4(s []int) int {
    _ = s[3] // hint: proves len(s) >= 4, eliminates bounds checks below
    return s[0] + s[1] + s[2] + s[3]
}

// Verify BCE with:
go build -gcflags='-d=ssa/check_bce/debug=1' ./...
```

---

## Compiler Directives for Unsafe Fast Paths

```go
// memmove-backed copy — compiler optimizes when T has no pointers
copy(dst, src) // becomes memmove if no GC write barriers needed

// For types with no pointers, copy is extremely fast.
// The compiler detects pointer-free types and skips write barriers.
type PacketHeader struct {
    Seq      uint32
    Checksum uint32
    Length   uint16
    Flags    uint8
    Reserved uint8
} // no pointers — copy is a raw memmove
```

---

## Go Optimization Mental Model (2026)

**The 10 Laws of Go Performance Engineering:**

1. **Measure with pprof before touching code.** Every optimization decision must be backed by a profiler output. Intuition is wrong ~60% of the time in Go.

2. **Allocations kill latency predictability.** GC pauses are non-deterministic. A function with 0 allocations will have more predictable tail latency than a faster function that allocates heavily.

3. **Escape analysis is the first filter.** Before any other optimization, understand what escapes to the heap and why. Use `-gcflags='-m=2'`. Fix unnecessary escapes before touching algorithms.

4. **Inlining is free performance.** Keep hot-path functions under the inline budget (~80 AST nodes). The compiler turns call overhead into zero. Validate with `-gcflags='-m=1'`.

5. **Interface dispatch is not free.** Every interface method call is an indirect branch. In tight loops over millions of items, this is measurable. Use generics, concrete types, or function specialization in hot paths.

6. **Cache coherence beats algorithmic cleverness.** A cache-friendly O(n) over a contiguous `[]T` beats an O(log n) with pointer chasing through scattered heap objects at small n. Know your working set size.

7. **sync.Pool eliminates GC pressure at the cost of code complexity.** Pool is the right tool for objects with clear reset semantics that are allocated frequently and discarded quickly. Not a general-purpose cache.

8. **goroutines are cheap but not free.** Unbounded goroutine creation is a resource leak. Worker pools with bounded concurrency are the production pattern. errgroup.SetLimit() is your friend.

9. **benchstat over single benchmark runs.** Noise in benchmarks is real. Always use `-count=10` and benchstat with p-values. A 10% improvement with p=0.2 is noise; a 30% improvement with p=0.000 is real.

10. **GOMEMLIMIT is mandatory in containers.** Set `GOMEMLIMIT` to ~90% of your container's memory limit. Without it, the Go GC does not know your container boundary and will OOM. This is the single most impactful deployment change for memory-constrained Go services.
