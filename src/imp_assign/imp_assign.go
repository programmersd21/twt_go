package main

import "fmt"

func main() {

	// =========================
	// 🔹 Basic implicit assignment
	// =========================

	// Go infers type as int because 42 is an integer literal
	a := 42

	// Go infers type as string
	b := "Hello Go"

	// Go infers type as bool
	c := true

	fmt.Println(a, b, c)

	// =========================
	// 🔹 Float inference
	// =========================

	// Default inferred type is float64 (NOT float32)
	x := 3.14

	fmt.Println(x)

	// =========================
	// 🔹 Type is fixed after inference
	// =========================

	y := 100 // inferred as int

	// y = "text" ❌ invalid: type mismatch (int vs string)

	y = 200 // valid: same type

	fmt.Println(y)

	// =========================
	// 🔹 Multiple variable implicit assignment
	// =========================

	p, q := 10, 20

	// Go infers both as int
	fmt.Println(p + q)

	// =========================
	// 🔹 Mixed type inference (important rule)
	// =========================

	m := 10
	n := 2.5

	// fmt.Println(m + n) ❌ invalid (int + float64 mismatch)

	// Fix: explicit conversion required
	result := float64(m) + n

	fmt.Println(result)
}
