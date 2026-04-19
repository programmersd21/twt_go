package main

import "fmt"

func main() {

	a := 20

	// Outer condition
	if a > 0 {

		// Inner condition only checked if outer is true
		if a < 50 {
			fmt.Println("Between 1 and 49")
		}
	}
}
