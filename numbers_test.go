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
func BenchmarkInt8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int8()
	}
}
func BenchmarkInt8Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int8Between(10, 127)
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
func BenchmarkInt16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int16()
	}
}
func BenchmarkInt16Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32()
	}
}
func BenchmarkInt32Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int()
	}
}
func BenchmarkIntBetween(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int64()
	}
}
func BenchmarkInt64Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkUInt8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UInt8()
	}
}
func BenchmarkUInt8Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UInt8Between(1, 255)
	}
}
func BenchmarkByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkUInt16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UInt16()
	}
}
func BenchmarkUInt16Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkUInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UInt32()
	}
}
func BenchmarkUInt32Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
func BenchmarkUInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UInt64()
	}
}
func BenchmarkUInt64Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UInt64Between(1, 1000000)
	}
}

func TestFloat32(t *testing.T) {
	for i := 0; i < loop; i++ {
		Float32()
	}
}
func BenchmarkFloat32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32()
	}
}

func TestFloat64(t *testing.T) {
	for i := 0; i < loop; i++ {
		Float64()
	}
}
func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64()
	}
}

func TestComplex64(t *testing.T) {
	for i := 0; i < loop; i++ {
		Complex64()
	}
}
func BenchmarkComplex64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Complex64()
	}
}

func TestComplex128(t *testing.T) {
	for i := 0; i < loop; i++ {
		Complex128()
	}
}
func BenchmarkComplex128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Complex128()
	}
}
