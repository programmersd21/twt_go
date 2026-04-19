package main

import "fmt"

// Identity returns the same value it receives.
// Useful for demonstrating generic pass-through behavior.
func Identity[T any](v T) T {
	return v
}

func main() {
	a := Identity
	b := Identity("Go Generics")

	fmt.Println(a)
	fmt.Println(b)
}
