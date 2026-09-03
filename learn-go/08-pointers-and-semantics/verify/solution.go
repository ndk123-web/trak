package main

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func Swap(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

func ToPointers(items []Rectangle) []*Rectangle {
	out := make([]*Rectangle, len(items))
	for i := range items {
		out[i] = &items[i]
	}
	return out
}

func main() {}
