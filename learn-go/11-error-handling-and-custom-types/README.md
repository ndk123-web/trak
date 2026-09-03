# 11 — Error Handling, Wrapping & Sentinel Errors

## 🎯 Learning Objectives
- Create custom error types
- Wrap errors with `%w` for unwrapping
- Use `errors.Is()` and `errors.As()` correctly
- Design sentinel errors vs structured errors

---

## 🧠 Core Concepts

### Error Wrapping
```go
fmt.Errorf("db query failed: %w", err)
```
The `%w` verb wraps the error, allowing `errors.Is()` to match.

### Sentinel Errors
```go
var ErrNotFound = errors.New("not found")
```
Package-level error values used for identity checking.

### errors.Is vs errors.As
- `errors.Is(err, target)`: Checks if `err` or any wrapped error equals `target`
- `errors.As(err, &target)`: Checks if `err` or any wrapped error is assignable to `target` type

---

## ⚠️ Common Traps
- `%v` instead of `%w`: `errors.Is` won't work
- `errors.As` requires a pointer: `errors.As(err, &myErr)`
- Comparing errors with `==` instead of `errors.Is` breaks wrapped error chains

---

## 📝 Interview Questions
Q: What's the difference between `fmt.Errorf("%v", err)` and `fmt.Errorf("%w", err)`?
A: `%w` wraps the error, preserving it in the chain for `errors.Is()`. `%v` just formats the string.

Q: When should you use sentinel errors vs custom error types?
A: Sentinel errors for simple identity checks (`ErrNotFound`). Custom types when you need structured data (status codes, retry info).
