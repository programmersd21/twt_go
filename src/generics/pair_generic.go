package main

import "fmt"

// Pair returns two values of different types.
// A and B are independent type parameters.
func Pair[A any, B any](a A, b B) (A, B) {
	return a, b
}

func main() {
	x, y := Pair(10, "Go")

	fmt.Println(x) // 10
	fmt.Println(y) // Go
}
