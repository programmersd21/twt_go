package main

import "fmt"

// Define a structure
type Person struct {
	Name string
	Age  int
}

func main() {

	// Create struct instance
	p1 := Person{
		Name: "Soumalya",
		Age:  18,
	}

	fmt.Println(p1)

	// Access fields
	fmt.Println(p1.Name)

	// Modify fields
	p1.Age = 19

	fmt.Println(p1)
}
