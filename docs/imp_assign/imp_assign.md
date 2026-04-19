## ⚙️ Implicit Assignment in Go

Go is **strictly typed**, so it does *not* do full “automatic type magic” like JavaScript or Python. But it *does* allow **implicit assignment in limited, safe cases**.

---

# 🧠 What “implicit assignment” means

It means:

> You assign a value **without explicitly stating its type**, and Go still figures it out.

---

# 🔹 1. Type inference (most common form)

Go can infer the type on the right side.

```go
a := 10
b := "hello"
c := 3.14
d := true
```

### What’s happening:

* `a` becomes `int`
* `b` becomes `string`
* `c` becomes `float64`
* `d` becomes `bool`

👉 This is *not dynamic typing*. The type is locked at compile time.

---

# 🔹 2. Implicit assignment in expressions (safe widening)

Go allows assignment when types are **compatible and safe**, especially with constants.

```go
var x int = 10
var y float64 = 10
```

✔ This works because `10` is a **constant** and can be converted safely.

---

# ❌ 3. What Go does NOT allow (important)

Go avoids hidden conversions.

```go
var a int = 10
var b int64 = a // ❌ error
```

You must explicitly convert:

```go
var b int64 = int64(a) // ✅ correct
```

---

# 🔹 4. Implicit assignment in function returns

When a function returns values, Go assigns them automatically.

```go
func add(a int, b int) int {
	return a + b
}

result := add(2, 3)
```

Here:

* `result` is implicitly assigned `int`

---

# 🔹 5. Short variable declaration (`:=`)

This is Go’s main implicit assignment tool.

```go
name := "Go"
age := 20
```

Equivalent to:

```go
var name string = "Go"
var age int = 20
```

---

# ⚡ Key Rule

Go allows implicit assignment only when:

* Type is **clear from context**
* No **unsafe conversion**
* No **hidden precision loss**

---

# 🧠 Mental Model

* `:=` → “figure it out and lock the type”
* `=` → “assign explicitly, types must match”
* Go → “no surprises, no silent conversions”

---

# 🚀 Summary

Go supports **limited implicit assignment via type inference**, but avoids hidden type conversions by design. This makes code more predictable, safer, and easier to reason about at scale.
