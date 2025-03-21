package GoGen

import (
	"math"
	"math/rand"
)

func GenerateInt8Between(min, max int8) int8 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return int8(rand.Intn(int(max-min+1)) + int(min))
}

func GenerateInt8() int8 {
	return GenerateInt8Between(math.MinInt8, math.MaxInt8)
}

func GenerateInt16Between(min, max int16) int16 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return int16(rand.Intn(int(max-min+1)) + int(min))
}

func GenerateInt16() int16 {
	return GenerateInt16Between(math.MinInt16, math.MaxInt16)
}

func GenerateIntBetween(min, max int) int {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return rand.Intn(max-min+1) + min
}

func GenerateInt() int {
	return GenerateIntBetween(math.MinInt, math.MaxInt)
}

func GenerateInt64Between(min, max int64) int64 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return rand.Int63n(max-min+1) + min
}

func GenerateInt64() int64 {
	return GenerateInt64Between(math.MinInt64, math.MaxInt64)
}

// ----------

func GenerateUInt8Between(min, max uint8) uint8 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return uint8(rand.Intn(int(max-min+1)) + int(min))
}

func GenerateUInt8() uint8 {
	return GenerateUInt8Between(0, math.MaxUint8)
}

func GenerateUInt16Between(min, max uint16) uint16 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return uint16(rand.Intn(int(max-min+1)) + int(min))
}

func GenerateUInt16() uint16 {
	return GenerateUInt16Between(0, math.MaxUint16)
}

func GenerateUIntBetween(min, max uint32) uint32 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return rand.Uint32()%(max-min+1) + min
}

func GenerateUInt() uint32 {
	return GenerateUIntBetween(0, math.MaxUint32)
}

func GenerateUInt64Between(min, max uint64) uint64 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return rand.Uint64()%(max-min+1) + min
}

func GenerateUInt64() uint64 {
	return GenerateUInt64Between(0, math.MaxUint64)
}

// --------------

func GenerateFloat32(min, max float32) float32 {
	if min > max {
		min, max = max, min
	}
	return rand.Float32()*(max-min) + min
}
func GenerateFloat64(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}

	return rand.Float64()*(max-min) + min
}
func GenerateComplex64(minReal, maxReal, minImag, maxImag float32) complex64 {
	realPart := GenerateFloat32(minReal, maxReal)
	imagPart := GenerateFloat32(minImag, maxImag)
	return complex(realPart, imagPart)
}
func GenerateComplex128(minReal, maxReal, minImag, maxImag float64) complex128 {
	realPart := GenerateFloat64(minReal, maxReal)
	imagPart := GenerateFloat64(minImag, maxImag)
	return complex(realPart, imagPart)
}
