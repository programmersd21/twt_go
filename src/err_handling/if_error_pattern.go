package main

import (
	"fmt"
	"os"
)

func main() {
	// Try opening a file
	file, err := os.Open("data.txt")

	// Handle error immediately
	if err != nil {
		fmt.Println("failed to open file:", err)
		return
	}

	// Ensure cleanup happens
	defer file.Close()

	fmt.Println("File opened successfully")
}
