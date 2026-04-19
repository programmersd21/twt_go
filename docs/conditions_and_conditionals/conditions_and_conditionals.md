# Go Conditions & Conditionals — Fundamentals

Conditions and conditionals are the **control system of program flow** in Go. They decide what code runs, when it runs, and under what logic.

---

# 🧠 Core Idea

A condition is any expression that evaluates to a **boolean value**:

* `true` → execute block
* `false` → skip block

Go uses strict typing, so conditions must always resolve cleanly to `bool`.

---

# 🔀 1. If Statements

The most basic control structure in Go.

## Syntax

```go
if condition {
    // runs if condition is true
}
```

## Example

```go
a := 10

if a > 5 {
    fmt.Println("Greater than 5")
}
```

---

## ⚡ Key Rules

* No parentheses required around condition
* Curly braces `{}` are mandatory
* Condition must evaluate to `bool`

```go
if a == 10 {
    fmt.Println("Match")
}
```

---

## 🔗 If–Else Chain

```go
if a > 10 {
    fmt.Println("Big")
} else if a == 10 {
    fmt.Println("Equal")
} else {
    fmt.Println("Small")
}
```

---

## 🧱 Nested If

```go
if a > 0 {
    if a < 100 {
        fmt.Println("Valid range")
    }
}
```

---

# ⚖️ 2. Comparison Operators

Used to build conditions.

| Operator | Meaning          |
| -------- | ---------------- |
| `==`     | equal            |
| `!=`     | not equal        |
| `<`      | less than        |
| `>`      | greater than     |
| `<=`     | less or equal    |
| `>=`     | greater or equal |

---

## Example

```go
if a >= 10 {
    fmt.Println("Valid")
}
```

---

## ⚠️ Type Safety Rule

Go does NOT allow comparing incompatible types.

```go
// ❌ invalid
// 10 == "10"
```

No need for `===` like in JavaScript.

---

# 🔁 3. Switch Statements

Used for multi-way branching.

---

## Standard Switch

```go
switch value {
case 1:
    fmt.Println("One")
case 2:
    fmt.Println("Two")
default:
    fmt.Println("Other")
}
```

---

## ⚡ Key Properties

* No need for `break` (implicit break)
* Cleaner than long if–else chains
* Only matching case runs

---

## 🧠 Example

```go
x := 2

switch x {
case 1:
    fmt.Println("One")
case 2:
    fmt.Println("Two")
}
```

---

# 🧩 4. Naked Switch

Switch without a value — evaluates conditions directly.

Behaves like an **if–else chain**.

---

## Example

```go
x := 10

switch {
case x < 5:
    fmt.Println("Small")
case x < 20:
    fmt.Println("Medium")
default:
    fmt.Println("Large")
}
```

---

# 🔥 5. Fallthrough

Forces execution into the next case.

Normally, Go exits after a match.

---

## Example

```go
x := 1

switch x {
case 1:
    fmt.Println("One")
    fallthrough
case 2:
    fmt.Println("Two")
}
```

---

## ⚠️ Behavior Rule

* `fallthrough` ignores the next case condition
* It simply continues execution
* Must be used carefully

---

# 🧠 Mental Model

* `if` → single decision gate
* `if–else` → branching logic tree
* `switch` → multi-path router
* `naked switch` → condition-driven router
* `fallthrough` → forced chain execution

---

# 🚀 Summary

Go conditionals are built for **clarity and predictability**:

* Strict boolean logic
* Minimal syntax overhead
* No hidden branching behavior
* Easy-to-trace execution flow
