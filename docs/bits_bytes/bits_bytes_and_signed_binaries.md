# Bits, Bytes & Signed Binary — Fundamentals

---

## 🧩 Bits and Bytes
- A **bit** is the smallest unit of data in computing and can be either **0 or 1**
- A **byte** consists of **8 bits**
- Bytes are the standard unit used to represent and store data in computers

👉 In short:
- 1 bit = binary digit (0 or 1)
- 1 byte = 8 bits

---

## ➖ Signed Binary Representation
Binary can represent both positive and negative numbers using a **sign bit**.

### How it works:
- The **leftmost bit** is used as the sign indicator
- If the sign bit is:
  - `0` → the number is **positive**
  - `1` → the number is **negative**

This allows binary systems to store signed integers efficiently.

---

## 📊 Value Ranges and Bit Growth
- Increasing the number of bits increases the **range of values** that can be represented
- However, introducing a **sign bit reduces the usable range for magnitude**, since one bit is reserved for sign information

### Key insight:
- More bits → larger range overall
- Signed representation → splits range into positive and negative halves

---

## ⚙️ Summary
- Bits are the foundation of all digital data
- Bytes group bits into meaningful storage units
- Signed binary enables representation of negative numbers using a dedicated sign bit
- Bit-width directly determines the numeric range available in a system
