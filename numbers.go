package gengo

import (
	"math"
	"math/rand"
)

func Int8Between(min, max int8) int8 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	diff := int16(max) - int16(min) + 1
	if diff <= 0 {
		return 0
	}
	val := int16(rand.Int63n(int64(diff))) + int16(min)
	return int8(val)
}

func Int8() int8 {
	return Int8Between(math.MinInt8, math.MaxInt8)
}

func Int16Between(min, max int16) int16 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	diff := int32(max) - int32(min) + 1
	if diff <= 0 {
		return 0
	}

	val := int32(rand.Int63n(int64(diff))) + int32(min)
	return int16(val)
}

func Int16() int16 {
	return Int16Between(math.MinInt16, math.MaxInt16)
}

func IntBetween(min, max int) int {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}

	// Convert to int64 to avoid overflow on subtraction
	diff := int64(max) - int64(min) + 1
	if diff <= 0 {
		return 0
	}

	// Generate a value in [0, diff)
	r := rand.Int63n(diff)
	return int(r) + min
}

func Int() int {
	return IntBetween(math.MinInt, math.MaxInt)
}

func Int64Between(min, max int64) int64 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	diff := max - min + 1
	if diff <= 0 {
		return 0
	}
	return rand.Int63n(diff) + min
}

func Int64() int64 {
	return Int64Between(math.MinInt64, math.MaxInt64)
}

// ----------
func Byte() byte {
	return UInt8()
}

func UInt8Between(min, max uint8) uint8 {
	if min > max {
		min, max = max, min
	}

	diff := uint16(max) - uint16(min) + 1
	if diff <= 0 {
		return 0
	}
	val := rand.Int63n(int64(diff)) + int64(min)
	return uint8(val)
}

func UInt8() uint8 {
	return UInt8Between(0, math.MaxUint8)
}

func UInt16Between(min, max uint16) uint16 {
	if min > max {
		min, max = max, min
	}

	diff := uint32(max) - uint32(min) + 1
	if diff == 0 {
		return 0
	}
	val := rand.Int63n(int64(diff)) + int64(min)
	return uint16(val)
}

func UInt16() uint16 {
	return UInt16Between(0, math.MaxUint16)
}

func UIntBetween(min, max uint32) uint32 {
	if min > max {
		min, max = max, min
	}

	diff := uint64(max) - uint64(min) + 1
	if diff == 0 {
		return 0
	}

	val := rand.Int63n(int64(diff)) + int64(min)
	return uint32(val)
}

func UInt() uint32 {
	return UIntBetween(0, math.MaxUint32)
}

func UInt64Between(min, max uint64) uint64 {
	if min > max {
		min, max = max, min
	}
	diff := max - min + 1
	if diff == 0 {
		// Full uint64 range
		return rand.Uint64()
	}
	return min + (rand.Uint64() % diff)
}

func UInt64() uint64 {
	return UInt64Between(0, math.MaxUint64)
}

// --------------

func Float32() float32 {
	return Float32Between(math.SmallestNonzeroFloat32, math.MaxFloat32)
}
func Float32Between(min, max float32) float32 {
	if min > max {
		min, max = max, min
	}
	return rand.Float32()*(max-min) + min
}

func Float64() float64 {
	return Float64Between(math.SmallestNonzeroFloat64, math.MaxFloat64)
}
func Float64Between(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}

	return rand.Float64()*(max-min) + min
}

func Complex64() complex64 {
	realPartMin := Float32()
	realPartMax := Float32()
	imagPartMin := Float32()
	imagPartMax := Float32()

	if realPartMin > realPartMax {
		realPartMin, realPartMax = realPartMax, realPartMin
	}

	if imagPartMin > imagPartMax {
		imagPartMin, imagPartMax = imagPartMax, imagPartMin
	}

	return Complex64Between(realPartMin, realPartMax, imagPartMin, imagPartMax)
}
func Complex64Between(minReal, maxReal, minImag, maxImag float32) complex64 {
	realPart := Float32Between(minReal, maxReal)
	imagPart := Float32Between(minImag, maxImag)
	return complex(realPart, imagPart)
}

func Complex128() complex128 {
	realPartMin := Float64()
	realPartMax := Float64()
	imagPartMin := Float64()
	imagPartMax := Float64()

	if realPartMin > realPartMax {
		realPartMin, realPartMax = realPartMax, realPartMin
	}

	if imagPartMin > imagPartMax {
		imagPartMin, imagPartMax = imagPartMax, imagPartMin
	}

	return Complex128Between(realPartMin, realPartMax, imagPartMin, imagPartMax)
}
func Complex128Between(minReal, maxReal, minImag, maxImag float64) complex128 {
	realPart := Float64Between(minReal, maxReal)
	imagPart := Float64Between(minImag, maxImag)
	return complex(realPart, imagPart)
}
