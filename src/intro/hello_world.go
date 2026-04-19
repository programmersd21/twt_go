// Every Go program starts with a package declaration.
// 'main' is special — it tells Go this is an executable program.
package main

// Importing the fmt package.
// fmt = formatted I/O (used for printing text to console)
import "fmt"

// main function is the entry point of the program.
// Execution ALWAYS starts here.
func main() {

	// Println prints text followed by a newline.
	// This is the simplest way to output to the terminal.
	fmt.Println("Hello, World!")

	// You can print multiple values too:
	fmt.Println("Go is running successfully 🚀")

	// You can also format output:
	name := "Soumalya"
	fmt.Printf("Welcome, %s! 👋\n", name)
}
