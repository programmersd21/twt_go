# Error Handling Cheatsheet

## Basic Error Handling
```go
func doSomething() error {
    return errors.New("something went wrong")
}

err := doSomething()
if err != nil {
    // handle error
}
```

## Custom Errors
```go
type MyError struct {
    Code int
    Msg  string
}

func (e *MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Msg)
}
```

## Defer
```go
func main() {
    defer fmt.Println("Last")
    fmt.Println("First")
}
```

## Recover (Panic Handling)
```go
func safe() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from", r)
        }
    }()
    panic("Ouch!")
}
```
