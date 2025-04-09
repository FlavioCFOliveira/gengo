package gengo

import (
	"math"
	"math/rand/v2"
	"time"
)

func Int8Between(min, max int8) int8 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	diff := int16(max) - int16(min) + 1
	if diff <= 0 {
		return 0
	}

	val := int16(rand.Int64N(int64(diff))) + int16(min)
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

	val := int32(rand.Int64N(int64(diff))) + int32(min)

	return int16(val)
}
func Int16() int16 {
	return Int16Between(math.MinInt16, math.MaxInt16)
}

func Int32Between(min, max int32) int32 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}

	// Handle edge cases to prevent overflow
	if min == math.MinInt32 && max == math.MaxInt32 {
		return int32(int64(rand.Int32()))
	}

	// Safe calculation
	diff := uint32(max) - uint32(min) + 1
	return min + int32(rand.Uint32()%diff)
}
func Int32() int32 {
	return Int32Between(math.MinInt32, math.MaxInt32)
}

func IntBetween(min, max int) int {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	if min < math.MinInt32 {
		min = math.MinInt32
	}
	if max > math.MaxInt32 {
		max = math.MaxInt32
	}

	diff := int64(max) - int64(min) + 1
	if diff <= 0 {
		return 0
	}

	r := rand.Int64N(diff)
	return int(r) + min
}
func Int() int {
	return rand.Int()
}

func Int64Between(min, max int64) int64 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}

	// Handle the case where subtracting min from max overflows
	if uint64(max)-uint64(min)+1 == 0 {
		return rand.Int64() // Full int64 range
	}

	diff := max - min + 1
	if diff <= 0 {
		return 0
	}

	return rand.Int64N(diff) + min
}
func Int64() int64 {
	return rand.Int64()
}

// ----------

func Byte() byte {
	return UInt8()
}

func UInt8Between(min, max uint8) uint8 {
	if min > max {
		min, max = max, min
	}

	return min + uint8(rand.IntN(int(max-min+1)))
}
func UInt8() uint8 {
	return uint8(rand.UintN(math.MaxUint8))
}

func UInt16Between(min, max uint16) uint16 {
	if min > max {
		min, max = max, min
	}

	return min + uint16(rand.IntN(int(max-min+1)))
}
func UInt16() uint16 {
	return uint16(rand.UintN(math.MaxUint16))
}

func UInt32Between(min, max uint32) uint32 {
	if min > max {
		min, max = max, min
	}

	// Usa o timestamp atual como a semente
	seed := uint32(time.Now().UnixNano())
	// Aplica uma operação bit a bit para distribuir a aleatoriedade
	random := seed ^ (seed << 21) ^ (seed >> 35) ^ (seed << 4)

	// Calcula o valor aleatório dentro do intervalo [min, max]
	random = min + (random % (max - min + 1))

	// Ajuste para garantir a variabilidade no número de dígitos
	numDigits := 1 + random%10
	factor := uint32(1)
	for i := uint32(1); i < numDigits; i++ {
		factor *= 10
	}

	// Ajusta o número aleatório para ter o número de dígitos desejado
	if factor > max {
		factor = max
	}
	random = min + (random % (factor - min + 1))
	if random > max {
		random = max
	}

	return random
}
func UInt32() uint32 {
	return uint32(rand.UintN(math.MaxUint32))
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
	return uint64(rand.UintN(math.MaxUint64))
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
