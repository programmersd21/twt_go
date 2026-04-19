package main

import "fmt"

func main() {

	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	// Iteration order is NOT guaranteed
	for k, v := range m {

		fmt.Println("Key:", k, "Value:", v)
	}
}
