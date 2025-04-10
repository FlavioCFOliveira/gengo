package gengo

import (
	"testing"
)

var loop = 1000

func TestInt8(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int8()
	}
}
func TestInt8Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int8Between(10, 100)
	}
}

func TestInt16(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int16()
	}
}
func TestInt16Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int16Between(100, 1000)
	}
}
func TestInt32(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int32()
	}
}
func TestInt32Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int32Between(100, 1000)
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int()
	}
}

func TestIntBetween(t *testing.T) {
	for i := 0; i < loop; i++ {
		IntBetween(1, 3000)
	}
}

func TestInt64(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int64()
	}
}
func TestInt64Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int64Between(10000, 1000000)
	}
}

func TestUInt8(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt8()
	}
}
func TestUInt8Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt8Between(1, 255)
	}
}
func TestByte(t *testing.T) {
	for i := 0; i < loop; i++ {
		Byte()
	}
}

func TestUInt16(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt16()
	}
}
func TestUInt16Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt16Between(1, 1000)
	}
}
func TestUInt32(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt32()
	}
}
func TestUInt32Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt32Between(1, 100000)
	}
}

func TestUInt64(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt64()
	}
}

func TestUInt64Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt64Between(1, 1000000)
	}
}

func TestFloat32(t *testing.T) {
	for i := 0; i < loop; i++ {
		Float32()
	}
}
func TestFloat64(t *testing.T) {
	for i := 0; i < loop; i++ {
		Float64()
	}
}

func TestComplex64(t *testing.T) {
	for i := 0; i < loop; i++ {
		Complex64()
	}
}
func TestComplex128(t *testing.T) {
	for i := 0; i < loop; i++ {
		Complex128()
	}
}
