package main

import "fmt"

func main() {

	i := 0

	// Only condition is used
	// Works like a "while loop"
	for i < 5 {

		fmt.Println("i =", i)

		// manual increment is required
		i++
	}
}
