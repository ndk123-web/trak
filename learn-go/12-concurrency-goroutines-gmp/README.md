# 12 — Goroutines & The GMP Scheduler

## 🎯 Learning Objectives
- Explain the GMP model (Goroutine, Machine, Processor)
- Understand why goroutines are lightweight (2KB vs 2MB stacks)
- Use `runtime.GOMAXPROCS` to control parallelism
- Know the difference between concurrency and parallelism

---

## 🧠 Core Concepts

### GMP Model
- **G** = Goroutine: User-level lightweight thread (starts at 2KB stack)
- **M** = Machine: OS thread (maps 1:1 to kernel thread)
- **P** = Processor: Logical processor (schedules Gs onto Ms)

### Concurrency vs Parallelism
- **Concurrency**: Dealing with multiple things at once (structure)
- **Parallelism**: Doing multiple things at once (execution)
You can have concurrency without parallelism (single core).

### Stack Growth
- Goroutine stacks start at 2KB and grow dynamically
- OS threads start at 2MB (fixed)
- This is why you can spawn millions of goroutines but not threads

---

## ⚠️ Common Traps
- **Goroutine leak**: Starting a goroutine that never exits
- **Closure variable capture**: Loop variables are shared
- **Main exits before goroutines finish**: Use `sync.WaitGroup`

---

## 📝 Interview Questions
Q: What's the difference between a goroutine and an OS thread?
A: Goroutines are user-space, start at 2KB, grow dynamically, and are multiplexed onto OS threads by the Go scheduler. OS threads are kernel-managed, start at 2MB, and are expensive.

Q: What does `runtime.GOMAXPROCS(0)` return?
A: The number of logical CPUs available. Setting it to N limits parallelism to N OS threads.
