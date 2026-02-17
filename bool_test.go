package gengo

import "testing"

func TestBool(t *testing.T) {
	for i := 0; i < 1000; i++ {
		Bool()
	}
}

func BenchmarkBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bool()
	}
}
