## ➗ Arithmetic Operators in Go — Clean Breakdown

Go keeps arithmetic **simple, predictable, and strict**. No hidden behavior, no surprises.

---

# 🔢 Core Arithmetic Operators

## ➕ Addition (`+`)

Adds two values.

```go id="g1a2b3"
a := 10
b := 5
c := a + b // 15
```

Also used for string concatenation:

```go id="g4h8k2"
s := "Go" + "Lang"
```

---

## ➖ Subtraction (`-`)

Subtracts right from left.

```go id="k9m2p1"
a := 10
b := 5
c := a - b // 5
```

---

## ✖️ Multiplication (`*`)

Multiplies values.

```go id="x8n3q7"
a := 10
b := 5
c := a * b // 50
```

---

## ➗ Division (`/`)

Divides left by right.

```go id="d2f9t6"
a := 10
b := 5
c := a / b // 2
```

### ⚠️ Important behavior:

* Integer division **truncates decimals**

```go id="r7s1u3"
a := 10
b := 3
c := a / b // 3 (not 3.33)
```

---

## 🧮 Modulus (`%`)

Returns remainder.

```go id="m4p8v2"
a := 10
b := 3
c := a % b // 1
```

---

# ⚡ Compound Assignment Operators

These combine operation + assignment.

## `+=`

```go id="c1d3e5"
a := 10
a += 5 // 15
```

## `-=`

```go id="f6g8h1"
a := 10
a -= 5 // 5
```

## `*=`

```go id="j2k4l6"
a := 10
a *= 2 // 20
```

## `/=`

```go id="n7p9q2"
a := 10
a /= 2 // 5
```

## `%=`

```go id="s3t5u7"
a := 10
a %= 3 // 1
```

---

# 🧠 Type Rules (VERY IMPORTANT)

Go is strict:

```go id="v1w2x3"
var a int = 10
var b float64 = 3.5

// ❌ invalid
// c := a + b
```

### ✔ Correct way:

```go id="y4z5a6"
c := float64(a) + b
```

---

# ⚡ Key Behavior Summary

* All arithmetic is **type-safe**
* No automatic mixing of `int`, `float64`, etc.
* Integer division **drops decimals**
* `%` only works on integers
* Strings can only use `+`

---

# 🧠 Mental Model

* `+` → combine / grow
* `-` → reduce / offset
* `*` → scale
* `/` → split
* `%` → leftovers

---

# 🚀 Summary

Go arithmetic is intentionally **minimal and strict**.
You get predictable math behavior without hidden conversions or runtime surprises — which is exactly why it scales well in systems and backend engineering.
