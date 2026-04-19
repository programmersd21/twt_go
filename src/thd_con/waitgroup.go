package main

import (
	"fmt"
	"sync"
)

// Worker function
func task(wg *sync.WaitGroup) {
	defer wg.Done() // signal completion
	fmt.Println("Task running")
}

func main() {
	var wg sync.WaitGroup

	// Add 2 goroutines to wait for
	wg.Add(2)

	go task(&wg)
	go task(&wg)

	// Wait until both finish
	wg.Wait()

	fmt.Println("All tasks complete")
}
