package main

import "fmt"

// Print is a generic function that works with ANY type.
// T is a type parameter (placeholder type).
func Print[T any](value T) {
	fmt.Println(value)
}

func main() {
	// Type inference (Go figures out T automatically)
	Print(10)
	Print("hello")
	Print(3.14)
}
