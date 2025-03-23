package GoGen

import "testing"

func TestGenerateUInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t := GenerateUInt()
		println(t)
	}
}

func TestGenerateUInt8(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t := GenerateUInt8()
		println(t)
	}
}

func TestGenerateUInt16(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t := GenerateUInt16()
		println(t)
	}
}

func TestGenerateUInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t := GenerateUInt64()
		println(t)
	}
}
