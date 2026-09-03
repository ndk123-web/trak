package main

import (
	"context"
	"testing"
	"time"
)

func TestCheckTimeout(t *testing.T) {
	if !CheckTimeout() { t.Error("expected timeout") }
}

func TestParentCancelsChild(t *testing.T) {
	if !ParentCancelsChild() { t.Error("expected child cancelled") }
}

func TestGetRequestID(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestIDKey{}, "abc-123")
	if GetRequestID(ctx) != "abc-123" { t.Error() }
	empty := context.Background()
	if GetRequestID(empty) != "" { t.Error() }
}
