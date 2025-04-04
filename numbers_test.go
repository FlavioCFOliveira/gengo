package gengo

import (
	"fmt"
	"math"
	"testing"
)

func TestInt8(t *testing.T) {
	for i := 0; i < 1000; i++ {
		println(Int8())
	}
}

func TestInt16(t *testing.T) {
	for i := 0; i < 1000; i++ {
		println(Int16())
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		Int()
	}
}
func TestIntBetween(t *testing.T) {
	for i := 0; i < 1000; i++ {
		println(IntBetween(math.MinInt, math.MaxInt))
	}
}

func TestInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(Int64())
	}
}
func TestInt64Between(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(Int64Between(math.MinInt64, math.MaxInt64))
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
