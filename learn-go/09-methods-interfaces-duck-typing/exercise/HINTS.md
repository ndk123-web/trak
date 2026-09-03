# Hints for 09

## TODO 1
Return `nil` directly, not a typed nil pointer.

## TODO 2
```go
switch v.(type) {
case int: return "int"
case string: return "string"
case bool: return "bool"
default: return fmt.Sprintf("%T", v)
}
```

## TODO 3
`return err == nil` — Go's `==` on interfaces checks both pointers.
