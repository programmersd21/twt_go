package main

import (
	"fmt"
)

func main() {
	baseErr := fmt.Errorf("low-level failure")

	// Wrap original error
	err := fmt.Errorf("operation failed: %w", baseErr)

	fmt.Println(err)
}
