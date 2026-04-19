# Go Concurrency & Threading — Complete Guide

Go is designed for **high-concurrency systems**. Instead of traditional heavy OS threads, Go uses a lightweight model built around **goroutines and channels**.

---

# 🧠 Core Idea

* **Concurrency** → dealing with many tasks at once
* **Parallelism** → executing tasks at the same time (on multiple CPU cores)

Go focuses on **concurrency first**, and enables parallelism when possible.

---

# 🧵 1. Goroutines (Lightweight Threads)

A goroutine is a **lightweight function execution unit** managed by the Go runtime.

## 🔥 Syntax

```go
go functionName()
```

---

## Example

```go
package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine")
}

func main() {
	go sayHello() // runs concurrently

	time.Sleep(time.Second) // wait for goroutine
	fmt.Println("Main function")
}
```

---

## ⚡ Key Points

* Extremely lightweight (thousands can run easily)
* Managed by Go runtime, not OS directly
* Starts with `go` keyword

---

# 🔄 2. Concurrency vs Parallelism

## Concurrency

* Multiple tasks in progress
* Switching between tasks

## Parallelism

* Multiple tasks executing at the same time
* Requires multiple CPU cores

---

## Mental Model

* Concurrency = multitasking
* Parallelism = multitasking + multiple workers

---

# 📡 3. Channels (Communication System)

Channels allow goroutines to **communicate safely**.

---

## Creating a Channel

```go
ch := make(chan int)
```

---

## Sending Data

```go
ch <- 10
```

---

## Receiving Data

```go
value := <-ch
```

---

## Full Example

```go
package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 42
	}()

	value := <-ch
	fmt.Println(value)
}
```

---

# 🔒 4. Buffered Channels

Channels with capacity.

```go
ch := make(chan int, 2)
```

---

## Example

```go
ch := make(chan int, 2)

ch <- 1
ch <- 2

fmt.Println(<-ch)
fmt.Println(<-ch)
```

---

# 🔁 5. WaitGroups (Synchronization)

Used to wait for multiple goroutines.

```go
package main

import (
	"fmt"
	"sync"
)

func task(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Task running")
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go task(&wg)

	wg.Wait()
	fmt.Println("All tasks complete")
}
```

---

# ⚠️ 6. Race Conditions

Happens when multiple goroutines access shared data.

```go
counter := 0

// unsafe concurrent access
counter++
```

---

## Fix: Mutex

```go
import "sync"

var mu sync.Mutex

mu.Lock()
counter++
mu.Unlock()
```

---

# 🧠 7. Select Statement

Used to handle multiple channels.

```go
select {
case msg1 := <-ch1:
	fmt.Println(msg1)
case msg2 := <-ch2:
	fmt.Println(msg2)
}
```

---

# ⚡ Key Concepts Summary

* Goroutines → lightweight tasks
* Channels → communication between goroutines
* WaitGroups → synchronization
* Mutex → shared data protection
* Select → multi-channel control

---

# 🧠 Mental Model

* Goroutine → worker
* Channel → pipeline
* WaitGroup → task counter
* Mutex → lock on shared resource
* Select → traffic controller

---

# 🚀 Final Summary

Go concurrency is built for:

* Massive scalability
* Minimal overhead
* Safe communication
* Predictable synchronization

Instead of managing threads manually, Go gives you a **runtime-managed concurrency system that scales automatically** across CPUs.
