## 🧾 Go `%` Formatters (fmt package) — Clean Breakdown

In Go, `%` format verbs are used with `fmt.Printf`, `fmt.Sprintf`, etc.
They control **how values are printed**, not how they are stored.

Think of them as:

> “Rendering instructions for output”

---

# 🔹 Core idea

```go
fmt.Printf("format", value)
```

Example:

```go
fmt.Printf("%d", 10)
```

---

# 🔢 1. Integer formatters

### `%d` — Decimal (base 10)

```go
fmt.Printf("%d\n", 42)
```

---

### `%b` — Binary

```go
fmt.Printf("%b\n", 10) // 1010
```

---

### `%o` — Octal

```go
fmt.Printf("%o\n", 10)
```

---

### `%x` / `%X` — Hexadecimal

```go
fmt.Printf("%x\n", 255) // ff
fmt.Printf("%X\n", 255) // FF
```

---

# 🔣 2. Floating-point formatters

### `%f` — Default float

```go
fmt.Printf("%f\n", 3.14159)
```

---

### `%.2f` — Precision control

```go
fmt.Printf("%.2f\n", 3.14159) // 3.14
```

---

### `%e` — Scientific notation

```go
fmt.Printf("%e\n", 12345.0)
```

---

# ✍️ 3. Strings & characters

### `%s` — String

```go
fmt.Printf("%s\n", "Go")
```

---

### `%q` — Quoted string

```go
fmt.Printf("%q\n", "Go") // "Go"
```

---

### `%c` — Character (Unicode rune)

```go
fmt.Printf("%c\n", 65) // A
```

---

# 🔘 4. Booleans

### `%t` — true / false

```go
fmt.Printf("%t\n", true)
```

---

# 🧠 5. Type & debugging formatters

### `%T` — Type of value

```go
fmt.Printf("%T\n", 10) // int
```

---

### `%v` — Default value (generic)

```go
fmt.Printf("%v\n", 10)
fmt.Printf("%v\n", "hello")
```

---

### `%+v` — Struct with field names

```go
fmt.Printf("%+v\n", struct{A int}{A: 10})
```

---

### `%#v` — Go-syntax representation

```go
fmt.Printf("%#v\n", 10)
```

---

# 📦 6. Pointer formatter

### `%p` — Memory address

```go
x := 10
fmt.Printf("%p\n", &x)
```

---

# ⚡ 7. Width & alignment

### Right-align (default)

```go
fmt.Printf("%5d\n", 10)
```

### Left-align

```go
fmt.Printf("%-5d\n", 10)
```

---

# 🧠 Mental Model

* `%d` → numbers
* `%f` → decimals
* `%s` → text
* `%t` → logic
* `%v` → “just show it”
* `%T` → “what are you?”

---

# 🚀 Summary

Go formatters are a **controlled output system**, not decoration.
They let you precisely define how data appears in logs, CLI tools, APIs, and debugging output.
