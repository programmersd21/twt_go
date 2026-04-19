package main

import (
	"fmt"
	"strconv"
)

func main() {
	// strconv.Atoi returns (int, error)
	value, err := strconv.Atoi("123")

	// Always check error first
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}

	fmt.Println("Parsed value:", value)
}
