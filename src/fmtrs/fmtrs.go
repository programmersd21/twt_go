package main

import "fmt"

func main() {

	// =========================
	// 🔹 1. Basic printing
	// =========================

	fmt.Println("Hello", "Go") // adds space + newline automatically

	fmt.Print("Hello ") // no newline
	fmt.Print("Go\n")   // manual newline

	// =========================
	// 🔹 2. Formatted output (fmt.Printf)
	// =========================

	name := "Soumalya"
	age := 18

	// %s = string, %d = integer
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	// =========================
	// 🔹 3. Common format verbs
	// =========================

	x := 42
	y := 3.14159
	z := true

	// %v → default format (most commonly used)
	fmt.Printf("Default: %v %v %v\n", x, y, z)

	// %#v → Go-syntax representation (useful for debugging)
	fmt.Printf("Go syntax: %#v\n", name)

	// %T → type of variable
	fmt.Printf("Type of x: %T\n", x)

	// =========================
	// 🔹 4. Numeric formatting
	// =========================

	num := 255

	// %d → decimal
	fmt.Printf("Decimal: %d\n", num)

	// %b → binary
	fmt.Printf("Binary: %b\n", num)

	// %x → hexadecimal
	fmt.Printf("Hex: %x\n", num)

	// =========================
	// 🔹 5. Floating-point formatting
	// =========================

	f := 3.14159265

	// %.2f → 2 decimal places
	fmt.Printf("Rounded: %.2f\n", f)

	// =========================
	// 🔹 6. Width & alignment
	// =========================

	// %5d → width 5 (right aligned)
	fmt.Printf("|%5d|\n", 42)

	// %-5d → left aligned
	fmt.Printf("|%-5d|\n", 42)

	// =========================
	// 🔹 7. Strings formatting
	// =========================

	s := "GoLang"

	// %s → normal string
	fmt.Printf("%s\n", s)

	// %q → quoted string
	fmt.Printf("%q\n", s)

	// =========================
	// 🔹 8. Struct-style debug output
	// =========================

	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "Alex", Age: 22}

	// %v → simple
	fmt.Printf("%v\n", p)

	// %+v → includes field names
	fmt.Printf("%+v\n", p)

	// %#v → Go syntax representation
	fmt.Printf("%#v\n", p)
}
