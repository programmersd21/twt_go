package main

import "fmt"

func main() {

	x := 2

	// Switch evaluates one value and matches cases
	switch x {

	case 1:
		fmt.Println("One")

	case 2:
		fmt.Println("Two")

	case 3:
		fmt.Println("Three")

	default:
		fmt.Println("Unknown")
	}
}
