package main

import "fmt"

func main() {

	// Array: fixed size, same type
	var arr [3]int = [3]int{1, 2, 3}

	fmt.Println(arr)

	// Access element
	fmt.Println(arr[0])

	// Modify element
	arr[1] = 99

	fmt.Println(arr)
}
