package main

import "fmt"

type Student struct {
	Name  string
	Marks int
}

func main() {

	students := []Student{
		{"A", 90},
		{"B", 85},
		{"C", 95},
	}

	for _, s := range students {
		fmt.Println(s.Name, s.Marks)
	}
}
