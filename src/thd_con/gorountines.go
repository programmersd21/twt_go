package main

import (
	"fmt"
	"time"
)

// This function runs as a goroutine (lightweight thread)
func sayHello() {
	fmt.Println("Hello from goroutine")
}

func main() {
	// Start goroutine (runs concurrently with main)
	go sayHello()

	// Prevent program from exiting immediately
	time.Sleep(time.Second)

	fmt.Println("Main function finished")
}
