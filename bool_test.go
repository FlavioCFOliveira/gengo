package gengo

import "testing"

func TestBool(t *testing.T) {
	var trueCount, falseCount int
	for i := 0; i < 1000; i++ {
		if Bool() {
			trueCount++
		} else {
			falseCount++
		}
	}
	if trueCount == 0 {
		t.Error("Bool() never returned true in 1000 iterations")
	}
	if falseCount == 0 {
		t.Error("Bool() never returned false in 1000 iterations")
	}
}

func BenchmarkBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bool()
	}
}
