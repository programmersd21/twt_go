# Data Structures Cheatsheet

## Arrays
```go
var a [5]int
a[0] = 10
b := [2]string{"hello", "world"}
```

## Slices
```go
s := []int{1, 2, 3}
s = append(s, 4)
subset := s[1:3] // [2, 3]
```

## Maps
```go
m := make(map[string]int)
m["key"] = 42
delete(m, "key")

val, ok := m["key"] // Check if exists
```

## Structs
```go
type User struct {
    Name string
    Age  int
}

u := User{Name: "Alice", Age: 30}
```

## Pointers
```go
x := 10
p := &x         // p is a pointer to x
fmt.Println(*p) // dereference p to get x's value
```
