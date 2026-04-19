package main

import "fmt"

func main() {

	// Map declaration
	studentMarks := map[string]int{
		"Alice": 90,
		"Bob":   85,
	}

	fmt.Println(studentMarks)

	// Access value
	fmt.Println(studentMarks["Alice"])

	// Add new key
	studentMarks["Charlie"] = 95

	fmt.Println(studentMarks)

	// Delete key
	delete(studentMarks, "Bob")

	fmt.Println(studentMarks)
}
