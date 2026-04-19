package main

import "fmt"

func main() {

	// Slice (flexible size)
	numbers := []int{10, 20, 30}

	fmt.Println(numbers)

	// Append (VERY important in Go)
	numbers = append(numbers, 40)

	fmt.Println(numbers)

	// Slice operations
	fmt.Println(numbers[1:3]) // from index 1 to 2
}
