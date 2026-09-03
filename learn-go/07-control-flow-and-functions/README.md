# 07 — Control Flow, Defer Stack & Panic/Recover

## 🎯 Learning Objectives
- Predict defer execution order (LIFO)
- Understand when defer arguments are evaluated
- Use panic/recover for exception-like control flow (sparingly!)
- Know all variants of `for` loop

---

## 🧠 Core Concepts

### Defer Rules
1. **LIFO order**: Last defer executes first
2. **Arguments evaluated immediately**: `defer fmt.Println(i)` captures `i` at defer time
3. **Deferred functions run on return**: Even if panic occurs

### Panic/Recover
- `panic` unwinds the stack, running deferred functions
- `recover` only works inside a deferred function
- **Best practice**: Don't use panic for normal errors

---

## ⚠️ Common Traps
- Defer in a loop creates many deferred calls
- Recover outside defer does nothing
- Panic in a goroutine only unwinds that goroutine

---

## 📝 Interview Questions
Q: What does this print?
```go
for i := 0; i < 3; i++ { defer fmt.Println(i) }
```
A: `2, 1, 0` — LIFO order, `i` evaluated at defer time.

Q: Can you recover from a panic in a different goroutine?
A: **No.** Each goroutine has its own stack.
