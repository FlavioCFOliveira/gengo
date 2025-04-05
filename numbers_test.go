package gengo

import (
	"math"
	"testing"
)

var loop = 1000

func TestInt8(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int8()
	}
}

func TestInt16(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int16()
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int()
	}
}
func TestInt32(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int32()
	}
}
func TestIntBetween(t *testing.T) {
	for i := 0; i < loop; i++ {
		IntBetween(math.MinInt, math.MaxInt)
	}
}

func TestInt64(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int64()
	}
}
func TestInt64Between(t *testing.T) {
	for i := 0; i < loop; i++ {
		Int64Between(math.MinInt64, math.MaxInt64)
	}
}

func TestUInt(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt()
	}
}

func TestUInt8(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt8()
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

func TestUInt64(t *testing.T) {
	for i := 0; i < loop; i++ {
		UInt64()
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
