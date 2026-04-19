package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("data.txt")

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Guaranteed cleanup
	defer file.Close()

	fmt.Println("Working with file safely")
}
