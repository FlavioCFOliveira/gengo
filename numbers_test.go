package gengo

import "testing"

func TestInt8(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Int8()
	}
}
func TestInt16(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Int16()
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Int()
	}
}

func TestInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Int64()
	}
}

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
func TestByte(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Byte()
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

func TestFloat32(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Float32()
	}
}
func TestFloat64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Float64()
	}
}

func TestComplex64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Complex64()
	}
}
func TestComplex128(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = Complex128()
	}
}
