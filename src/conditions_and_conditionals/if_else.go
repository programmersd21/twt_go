package main

import "fmt"

func main() {

	a := 3

	// If condition is true → first block runs
	// Else → fallback block runs
	if a > 5 {
		fmt.Println("Greater than 5")
	} else {
		fmt.Println("5 or less")
	}
}
