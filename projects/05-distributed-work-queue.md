# 05. Basic Distributed Work Queue (Advanced)

**Concepts:** `Interfaces`, `Channels`, `Select`, `WaitGroups`, `Structs`.

### Description
Implement a simple system where "Producer" goroutines send "Jobs" to a central queue, and multiple "Worker" goroutines process them.
- **Features:**
  - Define a `Job` interface.
  - Different types of jobs (e.g., EmailJob, ImageProcessingJob).
  - A central dispatcher that uses `select` to manage job flow.
- **Goal:** Master complex orchestration and the power of Go interfaces.

---

## How to Get Started
1. Create a new directory for your project (e.g., `projects/work-queue`).
2. Initialize a go module: `go mod init work-queue`.
3. Start coding! Refer back to the `src/` and `docs/` folders for help.
