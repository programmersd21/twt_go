package main

import "fmt"

func main() {

	a := 10
	b := 20

	// Equality check
	if a == b {
		fmt.Println("Equal")
	}

	// Not equal
	if a != b {
		fmt.Println("Not equal")
	}

	// Greater than
	if a < b {
		fmt.Println("a is smaller")
	}

	// Greater or equal
	if b >= 20 {
		fmt.Println("b is at least 20")
	}
}
