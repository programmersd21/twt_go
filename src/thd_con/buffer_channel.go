package main

import "fmt"

func main() {
	// Buffered channel with capacity 2
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	// Reading values
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
