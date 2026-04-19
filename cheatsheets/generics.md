# Generics Cheatsheet (Go 1.18+)

## Generic Function
```go
func Sum[T int | float64](nums []T) T {
    var total T
    for _, v := range nums {
        total += v
    }
    return total
}
```

## Generic Struct
```go
type Box[T any] struct {
    Value T
}
```

## Type Constraints
```go
type Number interface {
    int | int64 | float64
}

func Add[T Number](a, b T) T {
    return a + b
}
```
