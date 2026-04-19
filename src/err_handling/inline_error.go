package main

import (
	"fmt"
	"os"
)

func main() {
	// Inline error handling (very common Go style)
	if file, err := os.Open("data.txt"); err != nil {
		fmt.Println("error:", err)
	} else {
		defer file.Close()
		fmt.Println("File opened")
	}
}
