package main

import "fmt"

func main() {
	// Create channel for int values
	ch := make(chan int)

	// Sender goroutine
	go func() {
		ch <- 42 // send value into channel
	}()

	// Receiver
	value := <-ch

	fmt.Println("Received:", value)
}
