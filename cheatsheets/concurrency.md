# Concurrency Cheatsheet

## Goroutines
```go
go funcName()
```

## Channels
```go
ch := make(chan int)
ch <- 42      // Send
val := <-ch   // Receive
```

## Buffered Channels
```go
ch := make(chan int, 10)
```

## Select
```go
select {
case msg1 := <-c1:
    fmt.Println(msg1)
case c2 <- "hello":
    fmt.Println("Sent")
default:
    fmt.Println("Nothing")
}
```

## WaitGroups
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

## Mutex
```go
var mu sync.Mutex
mu.Lock()
// critical section
mu.Unlock()
```
