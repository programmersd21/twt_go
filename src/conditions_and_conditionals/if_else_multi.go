package main

import "fmt"

func main() {

	a := 10

	// Multiple conditions evaluated top to bottom
	// First TRUE block executes, rest are skipped

	if a > 10 {
		fmt.Println("Greater than 10")
	} else if a == 10 {
		fmt.Println("Exactly 10")
	} else {
		fmt.Println("Less than 10")
	}
}
