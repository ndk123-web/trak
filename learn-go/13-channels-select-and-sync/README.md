# 13 — Channels, Select & Synchronization Primitives

## 🎯 Learning Objectives
- Choose between buffered and unbuffered channels
- Use `select` for timeouts and non-blocking operations
- Use `sync.WaitGroup`, `Mutex`, `RWMutex`, and `sync.Pool`
- Understand channel closing semantics

---

## 🧠 Core Concepts

### Buffered vs Unbuffered
| | Unbuffered | Buffered |
|---|---|---|
| Size | `make(chan T)` | `make(chan T, N)` |
| Send | Blocks until receiver ready | Blocks only when buffer full |
| Receive | Blocks until sender ready | Blocks only when buffer empty |

### Select
```go
select {
case v := <-ch1:
    // handle v
case ch2 <- val:
    // sent
case <-time.After(100 * time.Millisecond):
    // timeout
default:
    // non-blocking
}
```

### sync.Pool
- Object pool for reuse (reduces GC pressure)
- Not for long-lived objects (GC can clear pool)
- Great for temporary buffers

---

## ⚠️ Common Traps
- **Closing a closed channel panics**
- **Sending to a closed channel panics**
- **Closing a nil channel blocks forever**
- **RWMutex**: Writers starve readers if constantly arriving

---

## 📝 Interview Questions
Q: What's the difference between `make(chan int)` and `make(chan int, 1)`?
A: Unbuffered requires both sender and receiver to be ready simultaneously. Buffered allows sender to continue if buffer has space.

Q: Can you range over a channel that was never closed?
A: Yes, but it will block forever after the last value unless there's a timeout or another case.
