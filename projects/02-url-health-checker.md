# 02. Concurrent URL Health Checker (Intermediate)

**Concepts:** `Goroutines`, `Channels`, `WaitGroups`, `Error Handling`.

### Description
Create a tool that takes a list of URLs and checks if they are online (return status 200).
- **Features:** Use goroutines to check multiple URLs simultaneously. Use a channel to collect the results and a `Waitgroup` to wait for all checks to finish.
- **Goal:** Understand Go's lightweight concurrency model and how to safely communicate between goroutines.
- **Bonus:** Add a timeout for each request using a `context` (Advanced concept).

---

## How to Get Started
1. Create a new directory for your project (e.g., `projects/url-checker`).
2. Initialize a go module: `go mod init url-checker`.
3. Start coding! Refer back to the `src/` and `docs/` folders for help.
