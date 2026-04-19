package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	_, err := os.Open("missing.txt")

	// Check specific error type
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File does not exist")
		return
	}

	if err != nil {
		fmt.Println("Other error:", err)
	}
}
