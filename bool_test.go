package gengo

import "testing"

func TestBool(t *testing.T) {
	for i := 0; i < 1000; i++ {
		Bool()
	}
}
