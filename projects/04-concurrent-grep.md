# 04. Concurrent File Search Tool (Grep-lite) (Advanced)

**Concepts:** `File I/O`, `Concurrency`, `Channels`, `String Processing`.

### Description
A tool that searches for a specific string within all files in a directory (and its subdirectories).
- **Features:** 
  - Spawn a goroutine for each file (or a pool of workers).
  - Use channels to report back the line numbers and file names where the string was found.
  - Handle large directories efficiently.
- **Goal:** Combine file system operations with high-performance concurrency.

---

## How to Get Started
1. Create a new directory for your project (e.g., `projects/grep-lite`).
2. Initialize a go module: `go mod init grep-lite`.
3. Start coding! Refer back to the `src/` and `docs/` folders for help.
