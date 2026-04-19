package main

import "fmt"

// Speaker defines a behavior contract.
// Any type that implements Speak() string automatically satisfies this interface.
type Speaker interface {
	Speak() string
}

// User is a concrete type (data model)
type User struct {
	Name string
}

// Speak implements the Speaker interface for User.
// This is implicit — no "implements" keyword needed in Go.
func (u User) Speak() string {
	return "Hello, " + u.Name
}

// Robot is another type implementing the same interface
type Robot struct {
	ID int
}

// Speak implementation for Robot
func (r Robot) Speak() string {
	return fmt.Sprintf("Beep boop. I am robot #%d", r.ID)
}

// say accepts ANY type that satisfies the Speaker interface
func say(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	u := User{Name: "Soumalya"}
	r := Robot{ID: 7}

	say(u)
	say(r)
}
