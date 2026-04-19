# Go Generics — Type Parameters Explained

Generics in Go allow you to write **type-safe reusable code** without losing Go’s static typing guarantees.

They were introduced in Go 1.18 to solve one major problem:

> Writing the same logic for multiple types without duplication.

---

# 🧠 Core Idea

Before generics, Go forced duplication:

```go
func AddInt(a, b int) int { return a + b }
func AddFloat(a, b float64) float64 { return a + b }
```

With generics:

```go
func Add[T any](a, b T) T {
    return a + b // (needs constraints in real usage)
}
```

---

# 🔧 1. Generic Functions

## Basic syntax

```go
func Print[T any](value T) {
    fmt.Println(value)
}
```

## Usage

```go
Print[int](10)
Print[string]("hello")
```

👉 Type is inferred in most cases:

```go
Print(10)
Print("hello")
```

---

# ⚙️ 2. Type Constraints

Go uses constraints to restrict allowed types.

## any (no restriction)

```go
func Identity[T any](v T) T {
    return v
}
```

---

## comparable (supports == and !=)

```go
func Contains[T comparable](arr []T, val T) bool {
    for _, v := range arr {
        if v == val {
            return true
        }
    }
    return false
}
```

---

# 📦 3. Generic Data Structures

## Generic Slice Wrapper

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() T {
    last := len(s.items) - 1
    val := s.items[last]
    s.items = s.items[:last]
    return val
}
```

## Usage

```go
var s Stack[int]
s.Push(10)
s.Push(20)
```

---

# 🧱 4. Multiple Type Parameters

```go
func Pair[A any, B any](a A, b B) (A, B) {
    return a, b
}
```

Usage:

```go
x, y := Pair(10, "Go")
```

---

# 🔬 5. Type Constraints with Interfaces

You can define custom constraints.

```go
type Number interface {
    int | int64 | float64
}
```

```go
func Sum[T Number](a, b T) T {
    return a + b
}
```

---

# ⚡ 6. Key Rules

* Generics are **compile-time only**
* No runtime overhead (in most cases)
* Must satisfy constraints
* Cannot use operators unless allowed by constraint

---

# 🧠 Mental Model

* Non-generic → duplicate logic per type
* Generic → one logic, many types
* Constraint → rulebook for allowed types

---

# 🚀 Summary

Generics in Go provide:

* Type-safe reuse
* Cleaner abstraction
* Reduced duplication
* Controlled flexibility via constraints

They keep Go simple while enabling modern, scalable code patterns used in real-world systems.
