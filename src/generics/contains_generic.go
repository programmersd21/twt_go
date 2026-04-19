package main

import "fmt"

// Contains checks if a value exists in a slice.
// "comparable" constraint allows using == and != safely.
func Contains[T comparable](arr []T, val T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func main() {
	nums := []int{1, 2, 3, 4}

	fmt.Println(Contains(nums, 3)) // true
	fmt.Println(Contains(nums, 9)) // false
}
