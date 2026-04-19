package main

import "fmt"

func main() {
	var a int = -10
	var b int8 = 100
	var c int16 = 20000
	var d int32 = 100000
	var e int64 = 1000000000

	var u1 uint = 10
	var u2 uint8 = 255
	var u3 byte = 100
	var u4 uint16 = 50000
	var u5 uint32 = 4000000000
	var u6 uint64 = 1000000000000

	fmt.Println("Signed Integers:")
	fmt.Println(a, b, c, d, e)

	fmt.Println("Unsigned Integers:")
	fmt.Println(u1, u2, u3, u4, u5, u6)
}
