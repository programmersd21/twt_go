package main

import "fmt"

func main() {

	// =========================
	// 🔹 Basic arithmetic setup
	// =========================

	a := 10
	b := 3

	// =========================
	// ➕ Addition
	// =========================

	sum := a + b
	fmt.Println("Addition:", sum)

	// =========================
	// ➖ Subtraction
	// =========================

	diff := a - b
	fmt.Println("Subtraction:", diff)

	// =========================
	// ✖ Multiplication
	// =========================

	product := a * b
	fmt.Println("Multiplication:", product)

	// =========================
	// ➗ Division (IMPORTANT BEHAVIOR)
	// =========================

	// Integer division → decimal part is dropped
	division := a / b
	fmt.Println("Integer Division:", division)

	// Correct float division
	floatDivision := float64(a) / float64(b)
	fmt.Println("Float Division:", floatDivision)

	// =========================
	// 🔁 Modulus (remainder)
	// =========================

	mod := a % b
	fmt.Println("Remainder (Modulus):", mod)

	// =========================
	// ⚠️ Operator constraints in Go
	// =========================

	x := 5
	y := 2.5

	// ❌ x + y is invalid (int + float mismatch)

	// Fix using explicit conversion
	result := float64(x) + y
	fmt.Println("Mixed addition:", result)

	// =========================
	// 🔥 Increment / Decrement
	// =========================

	count := 0

	count++ // increment by 1
	fmt.Println("Increment:", count)

	count-- // decrement by 1
	fmt.Println("Decrement:", count)

	// =========================
	// ⚡ Compound assignments
	// =========================

	n := 10

	n += 5 // same as n = n + 5
	fmt.Println("n += 5:", n)

	n -= 3 // same as n = n - 3
	fmt.Println("n -= 3:", n)

	n *= 2 // same as n = n * 2
	fmt.Println("n *= 2:", n)

	n /= 4 // same as n = n / 4
	fmt.Println("n /= 4:", n)
}
