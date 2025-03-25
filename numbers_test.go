package gengo

import "testing"

func TestUInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = UInt()
	}
}

func TestUInt8(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = UInt8()
	}
}

func TestUInt16(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = UInt16()
	}
}

func TestUInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = UInt64()
	}
}
