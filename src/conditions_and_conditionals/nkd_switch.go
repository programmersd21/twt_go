package main

import "fmt"

func main() {

	x := 15

	// Each case is a full condition (like if-else ladder)
	switch {

	case x < 10:
		fmt.Println("Small")

	case x < 20:
		fmt.Println("Medium")

	default:
		fmt.Println("Large")
	}
}
