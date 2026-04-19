# Go Basics Cheatsheet

## Variables and Types
```go
var name string = "Go"
age := 10 // Short declaration
const Pi = 3.14
```

## Basic Types
- `bool`: true, false
- `string`: "hello"
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`
- `byte` (alias for uint8)
- `rune` (alias for int32, represents a Unicode code point)
- `float32`, `float64`
- `complex64`, `complex128`

## Control Flow
### If Statement
```go
if age >= 18 {
    fmt.Println("Adult")
} else {
    fmt.Println("Minor")
}
```

### Switch Statement
```go
switch os := runtime.GOOS; os {
case "darwin":
    fmt.Println("OS X")
case "linux":
    fmt.Println("Linux")
default:
    fmt.Println("Other")
}
```

### For Loop
```go
// Standard loop
for i := 0; i < 10; i++ { ... }

// "While" loop
for i < 10 { ... }

// Infinite loop
for { ... }
```
