# 09 — Interfaces, Duck Typing & The Infamous nil Trap

## 🎯 Learning Objectives
- Explain how interfaces work internally (type-value pair)
- Understand why `(*Type)(nil) != nil` in interfaces
- Use empty interface (`any`) correctly
- Know when to use value vs pointer receivers with interfaces

---

## 🧠 Core Concepts

### Interface Internal Structure
```
[Interface Header]
| Type Pointer (8B) | Data Pointer (8B) |
```
An interface is nil **only when both pointers are nil**.

### The nil Trap
```go
var p *MyError = nil
var err error = p
fmt.Println(err == nil) // FALSE!
```
Why? Type Pointer = `*MyError` (non-nil!), Data Pointer = nil.

---

## ⚠️ Common Traps
- Returning typed nil: interface is NOT nil
- Type assertion panic: use `v, ok := x.(int)`
- Value receiver on pointer type: not in method set of `*Type`

---

## 📝 Interview Questions
Q: Why does `return err` where `err` is `*MyError = nil` fail the `!= nil` check?
A: The interface has type `*MyError` and value `nil`. Since type is non-nil, interface is non-nil.

Q: How do you fix it?
A: Return `nil` directly (untyped), or return `error(nil)`.
