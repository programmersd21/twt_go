package main

import "fmt"

func main() {

	str := "Go"

	// range iterates over runes (Unicode safe)
	for i, ch := range str {

		fmt.Println("Index:", i, "Char:", ch)
	}
}
