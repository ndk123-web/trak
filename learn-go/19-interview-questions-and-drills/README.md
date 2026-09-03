# 19 — High-Yield Go Interview Questions & Drills

## 🎯 Learning Objectives
- Master the top 10 Go interview gotchas
- Explain GMP scheduler, escape analysis, and nil interfaces
- Write correct concurrent code under pressure
- Handle edge cases in slices, maps, and channels

---

## 🧠 Top 10 Gotchas

1. **GMP Scheduler**: G (goroutine), M (OS thread), P (logical processor)
2. **`(*T)(nil) != nil`**: Interface with typed nil is not nil
3. **Map Concurrency**: Maps panic on concurrent writes without mutexes
4. **Defer Evaluation**: Arguments evaluated immediately, execution deferred in LIFO
5. **Sub-slice Memory Leaks**: Large backing arrays pinned by small windows
6. **Slice Growth**: Capacity doubles below 256, then grows by ~25%
7. **Channel Closing**: Panic on send to closed channel; panic on close of closed channel
8. **Context Cancellation**: Flows downward only; not calling cancel() leaks goroutines
9. **JSON Numbers**: Decode to float64 by default; use `json.Number` or custom types
10. **Race Detector**: `go test -race` adds 5-10x overhead; detects data races at runtime

---

## 📝 Final Drill
Write a program that:
1. Creates a worker pool of 3 goroutines
2. Each worker reads integers from a channel and squares them
3. Results are collected in another channel
4. The program times out after 100ms using context
5. Returns all successfully computed squares
