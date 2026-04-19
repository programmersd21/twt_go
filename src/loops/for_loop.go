package main

import "fmt"

func main() {

	// Classic loop has 3 parts:
	// 1. initialization
	// 2. condition
	// 3. post statement

	for i := 0; i < 5; i++ {

		// Runs while i < 5
		fmt.Println("i =", i)
	}
}
