# 14 — Context Package: Cancellation Trees

## 🎯 Learning Objectives
- Use `context.WithCancel`, `WithTimeout`, `WithDeadline`
- Propagate cancellation through the call tree
- Use `context.WithValue` sparingly (it's not a general store!)
- Handle context cancellation in HTTP handlers and goroutines

---

## 🧠 Core Concepts

### Context Types
- `context.Background()`: Root context, never cancelled
- `context.WithCancel(parent)`: Returns ctx + cancel func
- `context.WithTimeout(parent, duration)`: Auto-cancels after duration
- `context.WithDeadline(parent, time)`: Auto-cancels at specific time
- `context.WithValue(parent, key, val)`: Carries request-scoped values

### Cancellation Propagation
When a parent context is cancelled, ALL children are cancelled automatically.
```go
ctx, cancel := context.WithCancel(context.Background())
childCtx, _ := context.WithCancel(ctx) // child inherits parent cancellation
cancel() // both ctx and childCtx are cancelled
```

### WithValue Best Practices
- Use private key types to avoid collisions
- Only for request-scoped data (request IDs, auth tokens)
- NEVER for optional parameters or configuration

---

## ⚠️ Common Traps
- **Not calling cancel()**: Causes goroutine leak
- **Using string keys in WithValue**: Collisions between packages
- **Storing large objects in context**: Increases memory pressure
- **Checking ctx.Err() before starting work**: Should check during work

---

## 📝 Interview Questions
Q: What happens if you don't call `cancel()` on a `WithTimeout` context?
A: The goroutine created by `WithTimeout` leaks until the deadline passes.

Q: Can a child context cancel its parent?
A: **No.** Cancellation flows downward only.
