package gengo

import (
	"math"
	"testing"
)

var loop = 1000

// ---------- Int8 ----------

func TestInt8(t *testing.T) {
	seen := make(map[int8]bool)
	for i := 0; i < loop; i++ {
		seen[Int8()] = true
	}
	if len(seen) < 10 {
		t.Errorf("Int8() low variance: %d unique values in %d calls", len(seen), loop)
	}
}
func TestInt8Between(t *testing.T) {
	const min, max = int8(10), int8(100)
	for i := 0; i < loop; i++ {
		v := Int8Between(min, max)
		if v < min || v > max {
			t.Fatalf("Int8Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	// swapped args must produce same valid range
	for i := 0; i < loop; i++ {
		v := Int8Between(max, min)
		if v < min || v > max {
			t.Fatalf("Int8Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	// min == max must return that exact value
	if v := Int8Between(42, 42); v != 42 {
		t.Fatalf("Int8Between(42, 42) = %d, want 42", v)
	}
	// full range must not panic
	for i := 0; i < loop; i++ {
		Int8Between(math.MinInt8, math.MaxInt8)
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

// ---------- Int16 ----------

func TestInt16(t *testing.T) {
	seen := make(map[int16]bool)
	for i := 0; i < loop; i++ {
		seen[Int16()] = true
	}
	if len(seen) < 10 {
		t.Errorf("Int16() low variance: %d unique values in %d calls", len(seen), loop)
	}
}
func TestInt16Between(t *testing.T) {
	const min, max = int16(100), int16(1000)
	for i := 0; i < loop; i++ {
		v := Int16Between(min, max)
		if v < min || v > max {
			t.Fatalf("Int16Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := Int16Between(max, min)
		if v < min || v > max {
			t.Fatalf("Int16Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	if v := Int16Between(500, 500); v != 500 {
		t.Fatalf("Int16Between(500, 500) = %d, want 500", v)
	}
	for i := 0; i < loop; i++ {
		Int16Between(math.MinInt16, math.MaxInt16)
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

// ---------- Int32 ----------

func TestInt32(t *testing.T) {
	seen := make(map[int32]bool)
	for i := 0; i < loop; i++ {
		seen[Int32()] = true
	}
	if len(seen) < 10 {
		t.Errorf("Int32() low variance: %d unique values in %d calls", len(seen), loop)
	}
}
func TestInt32Between(t *testing.T) {
	const min, max = int32(100), int32(1000)
	for i := 0; i < loop; i++ {
		v := Int32Between(min, max)
		if v < min || v > max {
			t.Fatalf("Int32Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := Int32Between(max, min)
		if v < min || v > max {
			t.Fatalf("Int32Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	if v := Int32Between(77, 77); v != 77 {
		t.Fatalf("Int32Between(77, 77) = %d, want 77", v)
	}
	for i := 0; i < loop; i++ {
		Int32Between(math.MinInt32, math.MaxInt32)
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

// ---------- Int ----------

func TestInt(t *testing.T) {
	hasNeg := false
	for i := 0; i < loop; i++ {
		if Int() < 0 {
			hasNeg = true
			break
		}
	}
	if !hasNeg {
		t.Error("Int() never returned a negative value in 1000 iterations")
	}
}
func TestIntBetween(t *testing.T) {
	const min, max = 1, 3000
	for i := 0; i < loop; i++ {
		v := IntBetween(min, max)
		if v < min || v > max {
			t.Fatalf("IntBetween(%d, %d) = %d: out of range", min, max, v)
		}
	}
	// negative range
	const negMin, negMax = -500, -10
	for i := 0; i < loop; i++ {
		v := IntBetween(negMin, negMax)
		if v < negMin || v > negMax {
			t.Fatalf("IntBetween(%d, %d) = %d: out of range", negMin, negMax, v)
		}
	}
	// swapped
	for i := 0; i < loop; i++ {
		v := IntBetween(max, min)
		if v < min || v > max {
			t.Fatalf("IntBetween(%d, %d) = %d: out of range", max, min, v)
		}
	}
	if v := IntBetween(42, 42); v != 42 {
		t.Fatalf("IntBetween(42, 42) = %d, want 42", v)
	}
	// full int range triggers rangeSize==0 overflow path
	for i := 0; i < loop; i++ {
		IntBetween(math.MinInt, math.MaxInt)
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

// ---------- Int64 ----------

func TestInt64(t *testing.T) {
	hasNeg := false
	for i := 0; i < loop; i++ {
		if Int64() < 0 {
			hasNeg = true
			break
		}
	}
	if !hasNeg {
		t.Error("Int64() never returned a negative value in 1000 iterations")
	}
}
func TestInt64Between(t *testing.T) {
	const min, max = int64(10000), int64(1000000)
	for i := 0; i < loop; i++ {
		v := Int64Between(min, max)
		if v < min || v > max {
			t.Fatalf("Int64Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := Int64Between(max, min)
		if v < min || v > max {
			t.Fatalf("Int64Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	if v := Int64Between(99, 99); v != 99 {
		t.Fatalf("Int64Between(99, 99) = %d, want 99", v)
	}
	// full int64 range must not panic
	for i := 0; i < loop; i++ {
		Int64Between(math.MinInt64, math.MaxInt64)
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

// ---------- UInt8 / Byte ----------

func TestUInt8(t *testing.T) {
	seen := make(map[uint8]bool)
	for i := 0; i < loop; i++ {
		seen[UInt8()] = true
	}
	if len(seen) < 10 {
		t.Errorf("UInt8() low variance: %d unique values in %d calls", len(seen), loop)
	}
}
func TestUInt8Between(t *testing.T) {
	const min, max = uint8(10), uint8(200)
	for i := 0; i < loop; i++ {
		v := UInt8Between(min, max)
		if v < min || v > max {
			t.Fatalf("UInt8Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := UInt8Between(max, min)
		if v < min || v > max {
			t.Fatalf("UInt8Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	if v := UInt8Between(5, 5); v != 5 {
		t.Fatalf("UInt8Between(5, 5) = %d, want 5", v)
	}
	// full uint8 range must not panic
	for i := 0; i < loop; i++ {
		UInt8Between(0, math.MaxUint8)
	}
}
func TestByte(_ *testing.T) {
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

// ---------- UInt16 ----------

func TestUInt16(t *testing.T) {
	seen := make(map[uint16]bool)
	for i := 0; i < loop; i++ {
		seen[UInt16()] = true
	}
	if len(seen) < 10 {
		t.Errorf("UInt16() low variance: %d unique values in %d calls", len(seen), loop)
	}
}
func TestUInt16Between(t *testing.T) {
	const min, max = uint16(1), uint16(1000)
	for i := 0; i < loop; i++ {
		v := UInt16Between(min, max)
		if v < min || v > max {
			t.Fatalf("UInt16Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := UInt16Between(max, min)
		if v < min || v > max {
			t.Fatalf("UInt16Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	if v := UInt16Between(500, 500); v != 500 {
		t.Fatalf("UInt16Between(500, 500) = %d, want 500", v)
	}
	// full uint16 range must not panic
	for i := 0; i < loop; i++ {
		UInt16Between(0, math.MaxUint16)
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

// ---------- UInt32 ----------

func TestUInt32(_ *testing.T) {
	for i := 0; i < loop; i++ {
		UInt32()
	}
}
func TestUInt32Between(t *testing.T) {
	const min, max = uint32(1), uint32(100000)
	for i := 0; i < loop; i++ {
		v := UInt32Between(min, max)
		if v < min || v > max {
			t.Fatalf("UInt32Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := UInt32Between(max, min)
		if v < min || v > max {
			t.Fatalf("UInt32Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	// full uint32 range must not panic
	for i := 0; i < loop; i++ {
		UInt32Between(0, math.MaxUint32)
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

// ---------- UInt64 ----------

func TestUInt64(_ *testing.T) {
	for i := 0; i < loop; i++ {
		UInt64()
	}
}
func TestUInt64Between(t *testing.T) {
	const min, max = uint64(1), uint64(1000000)
	for i := 0; i < loop; i++ {
		v := UInt64Between(min, max)
		if v < min || v > max {
			t.Fatalf("UInt64Between(%d, %d) = %d: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := UInt64Between(max, min)
		if v < min || v > max {
			t.Fatalf("UInt64Between(%d, %d) = %d: out of range", max, min, v)
		}
	}
	// full uint64 range must not panic
	for i := 0; i < loop; i++ {
		UInt64Between(0, math.MaxUint64)
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

// ---------- Float32 ----------

func TestFloat32(t *testing.T) {
	for i := 0; i < loop; i++ {
		v := Float32()
		if v <= 0 {
			t.Fatalf("Float32() = %v: expected > 0", v)
		}
	}
}
func TestFloat32Between(t *testing.T) {
	const min, max = float32(-100.0), float32(100.0)
	for i := 0; i < loop; i++ {
		v := Float32Between(min, max)
		if v < min || v > max {
			t.Fatalf("Float32Between(%f, %f) = %f: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := Float32Between(max, min)
		if v < min || v > max {
			t.Fatalf("Float32Between(%f, %f) = %f: out of range", max, min, v)
		}
	}
	// equal bounds must return that exact value
	const eq32 = float32(3.14)
	if v := Float32Between(eq32, eq32); v != eq32 {
		t.Fatalf("Float32Between(%v, %v) = %v, want %v", eq32, eq32, v, eq32)
	}
	// NaN input propagates: Float32Between(NaN, x) returns NaN (documented behavior)
	nan32 := float32(math.NaN())
	if result := Float32Between(nan32, 1.0); !math.IsNaN(float64(result)) {
		t.Fatalf("Float32Between(NaN, 1.0) = %v, want NaN", result)
	}
}
func BenchmarkFloat32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32()
	}
}
func BenchmarkFloat32Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32Between(-100, 100)
	}
}

// ---------- Float64 ----------

func TestFloat64(t *testing.T) {
	for i := 0; i < loop; i++ {
		v := Float64()
		if v <= 0 {
			t.Fatalf("Float64() = %v: expected > 0", v)
		}
	}
}
func TestFloat64Between(t *testing.T) {
	const min, max = float64(-100.0), float64(100.0)
	for i := 0; i < loop; i++ {
		v := Float64Between(min, max)
		if v < min || v > max {
			t.Fatalf("Float64Between(%f, %f) = %f: out of range", min, max, v)
		}
	}
	for i := 0; i < loop; i++ {
		v := Float64Between(max, min)
		if v < min || v > max {
			t.Fatalf("Float64Between(%f, %f) = %f: out of range", max, min, v)
		}
	}
	// equal bounds must return that exact value
	const eq64 = float64(2.718281828)
	if v := Float64Between(eq64, eq64); v != eq64 {
		t.Fatalf("Float64Between(%v, %v) = %v, want %v", eq64, eq64, v, eq64)
	}
	// NaN input propagates: Float64Between(NaN, x) returns NaN (documented behavior)
	nan64 := math.NaN()
	if result := Float64Between(nan64, 1.0); !math.IsNaN(result) {
		t.Fatalf("Float64Between(NaN, 1.0) = %v, want NaN", result)
	}
}
func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64()
	}
}
func BenchmarkFloat64Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64Between(-100, 100)
	}
}

// ---------- Complex64 ----------

func TestComplex64(t *testing.T) {
	for i := 0; i < loop; i++ {
		v := Complex64()
		if real(v) < 0 || real(v) >= 1 {
			t.Fatalf("Complex64() real = %v: out of [0, 1)", real(v))
		}
		if imag(v) < 0 || imag(v) >= 1 {
			t.Fatalf("Complex64() imag = %v: out of [0, 1)", imag(v))
		}
	}
}
func TestComplex64Between(t *testing.T) {
	const (
		minReal, maxReal = float32(-10.0), float32(10.0)
		minImag, maxImag = float32(0.0), float32(5.0)
	)
	for i := 0; i < loop; i++ {
		v := Complex64Between(minReal, maxReal, minImag, maxImag)
		if real(v) < minReal || real(v) > maxReal {
			t.Fatalf("Complex64Between real = %f: out of [%f, %f]", real(v), minReal, maxReal)
		}
		if imag(v) < minImag || imag(v) > maxImag {
			t.Fatalf("Complex64Between imag = %f: out of [%f, %f]", imag(v), minImag, maxImag)
		}
	}
	// swapped args
	for i := 0; i < loop; i++ {
		v := Complex64Between(maxReal, minReal, maxImag, minImag)
		if real(v) < minReal || real(v) > maxReal {
			t.Fatalf("Complex64Between (swapped) real = %f: out of [%f, %f]", real(v), minReal, maxReal)
		}
		if imag(v) < minImag || imag(v) > maxImag {
			t.Fatalf("Complex64Between (swapped) imag = %f: out of [%f, %f]", imag(v), minImag, maxImag)
		}
	}
}
func BenchmarkComplex64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Complex64()
	}
}
func BenchmarkComplex64Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Complex64Between(-10, 10, 0, 5)
	}
}

// ---------- Complex128 ----------

func TestComplex128(t *testing.T) {
	for i := 0; i < loop; i++ {
		v := Complex128()
		if real(v) < 0 || real(v) >= 1 {
			t.Fatalf("Complex128() real = %v: out of [0, 1)", real(v))
		}
		if imag(v) < 0 || imag(v) >= 1 {
			t.Fatalf("Complex128() imag = %v: out of [0, 1)", imag(v))
		}
	}
}
func TestComplex128Between(t *testing.T) {
	const (
		minReal, maxReal = float64(-10.0), float64(10.0)
		minImag, maxImag = float64(0.0), float64(5.0)
	)
	for i := 0; i < loop; i++ {
		v := Complex128Between(minReal, maxReal, minImag, maxImag)
		if real(v) < minReal || real(v) > maxReal {
			t.Fatalf("Complex128Between real = %f: out of [%f, %f]", real(v), minReal, maxReal)
		}
		if imag(v) < minImag || imag(v) > maxImag {
			t.Fatalf("Complex128Between imag = %f: out of [%f, %f]", imag(v), minImag, maxImag)
		}
	}
	for i := 0; i < loop; i++ {
		v := Complex128Between(maxReal, minReal, maxImag, minImag)
		if real(v) < minReal || real(v) > maxReal {
			t.Fatalf("Complex128Between (swapped) real = %f: out of [%f, %f]", real(v), minReal, maxReal)
		}
		if imag(v) < minImag || imag(v) > maxImag {
			t.Fatalf("Complex128Between (swapped) imag = %f: out of [%f, %f]", imag(v), minImag, maxImag)
		}
	}
}
func BenchmarkComplex128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Complex128()
	}
}
func BenchmarkComplex128Between(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Complex128Between(-10, 10, 0, 5)
	}
}
