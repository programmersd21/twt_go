# 03. In-Memory Key-Value Store with TTL (Intermediate/Advanced)

**Concepts:** `Generics`, `Maps`, `Mutexes`, `Goroutines`, `Structs`.

### Description
Build a thread-safe in-memory cache where you can store data with an optional "Time-To-Live" (TTL).
- **Features:** 
  - `Set(key, value, ttl)`: Store a value.
  - `Get(key)`: Retrieve a value.
  - `Delete(key)`: Remove a value.
- **Goal:** Use `sync.Mutex` to prevent race conditions during concurrent access. Use `Generics` to allow storing any data type. Use a background goroutine to automatically clean up expired items.

---

## How to Get Started
1. Create a new directory for your project (e.g., `projects/kv-store`).
2. Initialize a go module: `go mod init kv-store`.
3. Start coding! Refer back to the `src/` and `docs/` folders for help.
