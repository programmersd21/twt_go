package main

import "fmt"

func main() {

	for i := 0; i < 10; i++ {

		// Immediately exits loop when condition matches
		if i == 5 {
			break
		}

		fmt.Println("i =", i)
	}
}
