package gengo

import (
	"math"
	"math/rand/v2"
)

// Int8Between returns a random int8 in [min, max]. Arguments are swapped if min > max.
func Int8Between(min, max int8) int8 {
	if min > max {
		min, max = max, min
	}
	return int8(int32(min) + rand.N(int32(max)-int32(min)+1)) //nolint:gosec // result is in [min,max], guaranteed to fit in int8
}

// Int8 returns a random int8 across the full [math.MinInt8, math.MaxInt8] range.
func Int8() int8 {
	return int8(rand.Uint32()) //nolint:gosec // low 8 bits of a random uint32; every bit pattern is a valid int8
}

// Int16Between returns a random int16 in [min, max]. Arguments are swapped if min > max.
func Int16Between(min, max int16) int16 {
	if min > max {
		min, max = max, min
	}
	return int16(int32(min) + rand.N(int32(max)-int32(min)+1)) //nolint:gosec // result is in [min,max], guaranteed to fit in int16
}

// Int16 returns a random int16 across the full [math.MinInt16, math.MaxInt16] range.
func Int16() int16 {
	return int16(rand.Uint32()) //nolint:gosec // low 16 bits of a random uint32; every bit pattern is a valid int16
}

// Int32Between returns a random int32 in [min, max]. Arguments are swapped if min > max.
func Int32Between(min, max int32) int32 {
	if min > max {
		min, max = max, min
	}
	return int32(int64(min) + rand.N(int64(max)-int64(min)+1)) //nolint:gosec // result is in [min,max], guaranteed to fit in int32
}

// Int32 returns a random int32 across the full [math.MinInt32, math.MaxInt32] range.
func Int32() int32 {
	return int32(rand.Uint32()) //nolint:gosec // same bit width; every bit pattern is a valid int32
}

// IntBetween returns a random int in [min, max]. Arguments are swapped if min > max.
func IntBetween(min, max int) int {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1 //nolint:gosec // two's-complement widening; the difference is always correct in uint64
	if rangeSize == 0 {                          // overflows to 0 only for the full int range
		return int(rand.Uint64()) //nolint:gosec // same bit width on 64-bit; every bit pattern is a valid int
	}
	return int(uint64(min) + rand.Uint64N(rangeSize)) //nolint:gosec // result fits within the original [min,max] range
}

// Int returns a random int across the full [math.MinInt, math.MaxInt] range.
func Int() int {
	return int(rand.Uint64()) //nolint:gosec // same bit width on 64-bit; every bit pattern is a valid int
}

// Int64Between returns a random int64 in [min, max]. Arguments are swapped if min > max.
func Int64Between(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1 //nolint:gosec // two's-complement widening; the difference is always correct in uint64
	if rangeSize == 0 {                          // overflows to 0 only for the full int64 range
		return int64(rand.Uint64()) //nolint:gosec // same bit width; every bit pattern is a valid int64
	}
	return int64(uint64(min) + rand.Uint64N(rangeSize)) //nolint:gosec // result fits within the original [min,max] range
}

// Int64 returns a random int64 across the full [math.MinInt64, math.MaxInt64] range.
func Int64() int64 {
	return int64(rand.Uint64()) //nolint:gosec // same bit width; every bit pattern is a valid int64
}

// ----------

// Byte returns a random byte (alias for UInt8).
func Byte() byte {
	return UInt8()
}

// UInt8Between returns a random uint8 in [min, max]. Arguments are swapped if min > max.
func UInt8Between(min, max uint8) uint8 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint32(max) - uint32(min) + 1 // uint32 prevents uint8 overflow
	return uint8(uint32(min) + rand.Uint32N(rangeSize)) //nolint:gosec // result is in [min,max], guaranteed to fit in uint8
}

// UInt8 returns a random uint8 across the full [0, math.MaxUint8] range.
func UInt8() uint8 {
	return uint8(rand.Uint32()) //nolint:gosec // low 8 bits of a random uint32; every bit pattern is a valid uint8
}

// UInt16Between returns a random uint16 in [min, max]. Arguments are swapped if min > max.
func UInt16Between(min, max uint16) uint16 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint32(max) - uint32(min) + 1 // uint32 prevents uint16 overflow
	return uint16(uint32(min) + rand.Uint32N(rangeSize)) //nolint:gosec // result is in [min,max], guaranteed to fit in uint16
}

// UInt16 returns a random uint16 across the full [0, math.MaxUint16] range.
func UInt16() uint16 {
	return uint16(rand.Uint32()) //nolint:gosec // low 16 bits of a random uint32; every bit pattern is a valid uint16
}

// UInt32Between returns a random uint32 in [min, max]. Arguments are swapped if min > max.
func UInt32Between(min, max uint32) uint32 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1 // uint64 prevents uint32 overflow
	if rangeSize > math.MaxUint32 {
		return rand.Uint32()
	}
	return min + rand.Uint32N(uint32(rangeSize)) //nolint:gosec // rangeSize is verified to fit in uint32
}

// UInt32 returns a random uint32 across the full [0, math.MaxUint32] range.
func UInt32() uint32 {
	return rand.Uint32()
}

// UInt64Between returns a random uint64 in [min, max]. Arguments are swapped if min > max.
func UInt64Between(min, max uint64) uint64 {
	if min > max {
		min, max = max, min
	}
	rangeSize := max - min + 1
	if rangeSize == 0 { // full uint64 range
		return rand.Uint64()
	}
	return min + rand.Uint64N(rangeSize)
}

// UInt64 returns a random uint64 across the full [0, math.MaxUint64] range.
func UInt64() uint64 {
	return rand.Uint64()
}

// --------------

// Float32 returns a random float32 in (0, math.MaxFloat32].
func Float32() float32 {
	return Float32Between(math.SmallestNonzeroFloat32, math.MaxFloat32)
}

// Float32Between returns a random float32 in [min, max]. Arguments are swapped if min > max.
func Float32Between(min, max float32) float32 {
	if min > max {
		min, max = max, min
	}
	return rand.Float32()*(max-min) + min
}

// Float64 returns a random float64 in (0, math.MaxFloat64].
func Float64() float64 {
	return Float64Between(math.SmallestNonzeroFloat64, math.MaxFloat64)
}

// Float64Between returns a random float64 in [min, max]. Arguments are swapped if min > max.
func Float64Between(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	return rand.Float64()*(max-min) + min
}

// Complex64 returns a random complex64 with real and imaginary parts in [0, 1).
func Complex64() complex64 {
	return complex(rand.Float32(), rand.Float32())
}

// Complex64Between returns a random complex64 with real part in [minReal, maxReal]
// and imaginary part in [minImag, maxImag].
func Complex64Between(minReal, maxReal, minImag, maxImag float32) complex64 {
	return complex(Float32Between(minReal, maxReal), Float32Between(minImag, maxImag))
}

// Complex128 returns a random complex128 with real and imaginary parts in [0, 1).
func Complex128() complex128 {
	return complex(rand.Float64(), rand.Float64())
}

// Complex128Between returns a random complex128 with real part in [minReal, maxReal]
// and imaginary part in [minImag, maxImag].
func Complex128Between(minReal, maxReal, minImag, maxImag float64) complex128 {
	return complex(Float64Between(minReal, maxReal), Float64Between(minImag, maxImag))
}
