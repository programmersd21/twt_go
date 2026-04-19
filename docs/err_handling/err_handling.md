# Go Error Handling — Fundamentals

Go does not use exceptions like many languages. Instead, it uses **explicit error values** returned from functions.

This makes error handling **predictable, visible, and mandatory**.

---

# 🧠 Core Idea

In Go:

* Errors are just values (`error` type)
* Functions return errors explicitly
* You must handle them manually

```go
value, err := someFunction()
```

---

# 🔁 1. Basic Error Pattern

Most Go functions return:

* result
* error

## Example

```go
result, err := strconv.Atoi("123")

if err != nil {
    fmt.Println("Error occurred:", err)
    return
}

fmt.Println(result)
```

---

# ⚡ 2. The `error` Type

Go has a built-in interface for errors:

```go
type error interface {
    Error() string
}
```

Any type implementing `Error()` is an error.

---

# 🧱 3. Creating Errors

## Using `errors.New`

```go
import "errors"

err := errors.New("something went wrong")
fmt.Println(err)
```

---

## Using `fmt.Errorf`

Adds formatting support:

```go
import "fmt"

err := fmt.Errorf("invalid value: %d", 10)
fmt.Println(err)
```

---

# 🔀 4. If Error Pattern (Core Go Idiom)

This is the most common pattern in Go:

```go
file, err := os.Open("data.txt")

if err != nil {
    fmt.Println("failed to open file:", err)
    return
}

defer file.Close()
```

---

# 🧠 5. `defer` + Error Handling

Used for cleanup even if errors happen.

```go
file, err := os.Open("data.txt")
if err != nil {
    return
}

defer file.Close()
```

---

# 🔥 6. Wrapping Errors (Modern Go)

Go allows error chaining:

```go
import "fmt"

err := fmt.Errorf("operation failed: %w", originalErr)
```

* `%w` wraps the original error
* Preserves error history

---

# 🧪 7. Checking Specific Errors

```go
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("File does not exist")
}
```

---

# 🧠 Mental Model

* Errors are **not exceptions**
* Errors are **returned values you must check**
* Control flow is explicit, not hidden

---

# ⚡ Key Patterns

### Standard pattern

```go
if err != nil {
    return err
}
```

### Clean chaining style

```go
if err := doWork(); err != nil {
    return err
}
```

---

# 🚀 Summary

Go error handling is designed for:

* Predictability
* Explicit control flow
* No hidden exceptions
* Easy debugging in production systems

You don’t “catch” errors — you **handle them directly where they occur**.
