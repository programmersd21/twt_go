package main

import (
	"fmt"
	"time"
)

// Simulates a task
func task(name string) {
	for i := 0; i < 3; i++ {
		fmt.Println(name, "->", i)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	// Concurrency: tasks overlap in execution
	go task("Task A")
	go task("Task B")

	time.Sleep(2 * time.Second)

	fmt.Println("Main done")
}
