package main

import "fmt"

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Sender 1
	go func() {
		ch1 <- "Message from channel 1"
	}()

	// Sender 2
	go func() {
		ch2 <- "Message from channel 2"
	}()

	// Select waits for any channel to respond first
	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	}
}
