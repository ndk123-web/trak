# 08 — Pointers, Value Semantics & Mutation

## 🎯 Learning Objectives
- Decide when to use pointer vs value receivers
- Understand pass-by-value for all function arguments
- Predict mutation behavior with pointers
- Know when copying is expensive

---

## 🧠 Core Concepts

### Go is Always Pass-by-Value
- **Value**: Entire struct/array copied
- **Pointer**: Pointer value (8 bytes) copied, but points to shared data

### Pointer vs Value Receivers
| | Value Receiver | Pointer Receiver |
|---|---|---|
| Can mutate? | No | Yes |
| Nil safe? | Yes | Must check nil |
| Method set | Both value and pointer vars | Only pointer vars |

---

## ⚠️ Common Traps
- Pointer to loop variable: all goroutines see same `v`
- Map values are not addressable: `&m["key"]` invalid
- Slice elements ARE addressable: `&slice[i]` valid

---

## 📝 Interview Questions
Q: What's the output of `modify(x)` where `modify` takes `int`?
A: Original `x` unchanged — Go is pass-by-value.

Q: Should `Config` methods use pointer or value receivers?
A: Usually pointer, because configs are often large and may need mutation.
