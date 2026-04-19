package main

import "fmt"

func main() {

	x := 1

	// fallthrough forces next case to run
	switch x {

	case 1:
		fmt.Println("One")
		fallthrough

	case 2:
		fmt.Println("Two (forced execution)")
	}
}
