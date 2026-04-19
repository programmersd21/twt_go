# Functions and Interfaces Cheatsheet

## Functions
```go
func add(x int, y int) int {
    return x + y
}

func swap(x, y string) (string, string) {
    return y, x
}
```

## Methods
```go
type Rect struct {
    Width, Height int
}

func (r Rect) Area() int {
    return r.Width * r.Height
}
```

## Interfaces
```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```
