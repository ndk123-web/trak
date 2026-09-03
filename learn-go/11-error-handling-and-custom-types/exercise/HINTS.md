# Hints for 11

## TODO 1
```go
return fmt.Errorf("validation failed: %w", ErrInvalidInput)
```

## TODO 2
```go
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
```

## TODO 3
```go
var ve *ValidationError
if errors.As(err, &ve) {
    return ve, true
}
return nil, false
```
