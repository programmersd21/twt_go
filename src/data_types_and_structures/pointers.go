package main

import "fmt"

func main() {

	a := 10

	// Pointer stores memory address
	var p *int = &a

	fmt.Println("Value:", a)
	fmt.Println("Address:", p)
	fmt.Println("Dereferenced:", *p)

	// Modify via pointer
	*p = 50

	fmt.Println("Updated a:", a)
}
