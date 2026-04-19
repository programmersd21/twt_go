package main

import (
	"errors"
	"fmt"
)

func main() {
	// Creating a simple error
	err1 := errors.New("something went wrong")
	fmt.Println(err1)

	// Creating formatted error
	value := 10
	err2 := fmt.Errorf("invalid value: %d", value)
	fmt.Println(err2)
}
