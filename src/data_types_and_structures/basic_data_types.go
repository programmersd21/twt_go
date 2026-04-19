package main

import "fmt"

func main() {

	// =========================
	// 🔢 Integers
	// =========================

	var a int = 10   // default signed integer
	var b int8 = 127 // 8-bit signed integer (-128 to 127)
	var c uint = 50  // unsigned integer (no negatives)

	fmt.Println(a, b, c)

	// =========================
	// 🔢 Floating point numbers
	// =========================

	var x float32 = 3.14
	var y float64 = 99.999999

	fmt.Println(x, y)

	// =========================
	// 🔤 Strings
	// =========================

	var name string = "Go Language"

	fmt.Println(name)

	// Strings are immutable
	// name[0] = 'X' ❌ NOT allowed

	// =========================
	// 🔘 Boolean
	// =========================

	var isActive bool = true

	fmt.Println(isActive)

	// =========================
	// ⚡ Type inference (Go style)
	// =========================

	autoInt := 42
	autoString := "fast inference"

	fmt.Println(autoInt, autoString)
}
