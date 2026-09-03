package main

import "testing"

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"two", 2, true},
		{"three", 3, true},
		{"four", 4, false},
		{"seventeen", 17, true},
		{"one", 1, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrime(tt.n); got != tt.want {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	if Reverse("hello") != "olleh" { t.Error() }
	if Reverse("Hello, 世界") != "界世 ,olleH" { t.Errorf("got %q", Reverse("Hello, 世界")) }
}

func TestFibonacci(t *testing.T) {
	if Fibonacci(0) != 0 { t.Error() }
	if Fibonacci(1) != 1 { t.Error() }
	if Fibonacci(10) != 55 { t.Errorf("got %d", Fibonacci(10)) }
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}
