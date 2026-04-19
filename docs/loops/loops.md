# Go Loops — Complete Guide

Go has a **single looping construct**: `for`. Everything else is built from it. There is no `while` or `do-while` keyword — Go simplifies loops into one flexible system.

---

# 🧠 Core Idea

In Go, all looping patterns are variations of:

```go
for initialization; condition; post {
    // code
}
```

Or simplified forms of it.

---

# 🔁 1. Classic For Loop (C-style)

This is the most explicit loop form.

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

### Structure:

* Initialization → `i := 0`
* Condition → `i < 5`
* Post statement → `i++`

---

# 🔄 2. While-style Loop

Go does not have `while`, but this form behaves like one.

```go
i := 0

for i < 5 {
    fmt.Println(i)
    i++
}
```

### Key idea:

* Only condition is used
* Acts exactly like a `while` loop

---

# ♾️ 3. Infinite Loop

A loop with no condition.

```go
for {
    fmt.Println("running")
}
```

### Control mechanisms:

* `break` → exit loop
* `continue` → skip iteration

---

# 📦 4. Range Loop

Used for iterating over collections.

## Arrays / Slices

```go
nums := []int{10, 20, 30}

for i, v := range nums {
    fmt.Println(i, v)
}
```

* `i` → index
* `v` → value

---

## Strings

```go
for i, ch := range "Go" {
    fmt.Println(i, ch)
}
```

* Iterates over Unicode characters (runes)

---

## Maps

```go
m := map[string]int{"a": 1, "b": 2}

for k, v := range m {
    fmt.Println(k, v)
}
```

* Order is NOT guaranteed

---

# ⛔ 5. Break Statement

Exits the loop immediately.

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break
    }
    fmt.Println(i)
}
```

---

# ⏭️ 6. Continue Statement

Skips current iteration.

```go
for i := 0; i < 5; i++ {
    if i == 2 {
        continue
    }
    fmt.Println(i)
}
```

---

# 🏷️ 7. Labeled Loops

Used to break or continue outer loops.

```go
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            break outer
        }
        fmt.Println(i, j)
    }
}
```

---

# 🧠 Mental Model

* `for` → universal loop engine
* condition-only → while behavior
* empty `for` → infinite loop
* `range` → structured iteration tool
* `break` → exit loop
* `continue` → skip step
* label → control nested loops

---

# 🚀 Summary

Go intentionally avoids multiple loop keywords. Instead:

* One loop construct (`for`)
* Multiple patterns via syntax variation
* Explicit control flow (`break`, `continue`, labels)

This design keeps looping behavior consistent, predictable, and easy to reason about in large systems.
