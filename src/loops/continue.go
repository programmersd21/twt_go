package main

import "fmt"

func main() {

	for i := 0; i < 5; i++ {

		// Skips current iteration
		if i == 2 {
			continue
		}

		fmt.Println("i =", i)
	}
}
