# 10 — Generics (Go 1.18+)

## 🎯 Learning Objectives
- Write generic functions with type parameters
- Define custom type constraints using interfaces
- Understand `any` vs `comparable`
- Know when NOT to use generics

---

## 🧠 Core Concepts

### Type Parameters
```go
func Max[T comparable](a, b T) T { /* ERROR: comparable doesn't support > */ }
```

### Ordered Constraint
```go
type Number interface {
    ~int | ~float64
}
```
`~int` means "any type whose underlying type is int".

---

## ⚠️ Common Traps
- `comparable` only supports `==` and `!=`, not `<` or `>`
- Generics increase compile time
- Type inference: `Max(1, 2)` infers `T=int`

---

## 📝 Interview Questions
Q: What's the difference between `any` and `comparable`?
A: `any` = `interface{}` (no constraints). `comparable` = types that support `==` and `!=`.

Q: Can methods have their own type parameters?
A: No. But the receiver type can be generic.
