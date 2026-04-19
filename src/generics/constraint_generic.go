package main

import "fmt"

// Number is a type constraint (only allows numeric types listed).
type Number interface {
	int | int64 | float64
}

// Sum works only for numeric types defined in Number.
func Sum[T Number](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(Sum(10, 20))
	fmt.Println(Sum(2.5, 3.5))
}
