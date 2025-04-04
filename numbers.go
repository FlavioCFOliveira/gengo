package gengo

import (
	"math"
	"math/rand"
)

func Int8Between(min, max int8) int8 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return int8(rand.Intn(int(max-min+1)) + int(min))
}

func Int8() int8 {
	return Int8Between(math.MinInt8, math.MaxInt8)
}

func Int16Between(min, max int16) int16 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return int16(rand.Intn(int(max-min+1)) + int(min))
}

func Int16() int16 {
	return Int16Between(math.MinInt16, math.MaxInt16)
}

func IntBetween(min, max int) int {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return rand.Intn(max-min+1) + min
}

func Int() int {
	return IntBetween(math.MinInt, math.MaxInt)
}

func Int64Between(min, max int64) int64 {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return rand.Int63n(max-min+1) + min
}

func Int64() int64 {
	return Int64Between(math.MinInt64, math.MaxInt64)
}

// ----------

func UInt8Between(min, max uint8) uint8 {
	if min > max {
		min, max = max, min
	}

	rangeSize := uint16(max) - uint16(min) + 1 // usa uint16 para evitar overflow

	var randNum uint8
	limit := uint16(256) - (uint16(256) % rangeSize) // maior múltiplo de rangeSize dentro do uint8

	for {
		r := uint16(rand.Uint32() & 0xFF) // pega só os 8 bits menos significativos
		if r < limit {
			randNum = uint8(r % rangeSize)
			break
		}
	}

	return min + randNum
}

func UInt8() uint8 {
	return UInt8Between(0, math.MaxUint8)
}

func UInt16Between(min, max uint16) uint16 {
	if min > max {
		min, max = max, min
	}

	// Calcula o tamanho do intervalo de forma segura (usando uint32 para evitar overflow)
	rangeSize := uint32(max) - uint32(min) + 1

	var randNum uint16
	limit := uint32(65536) - (uint32(65536) % rangeSize) // maior múltiplo de rangeSize dentro de uint16

	for {
		// Gera número aleatório com uint32
		r := uint32(rand.Uint32() & 0xFFFF) // Pega os 16 bits menos significativos
		if r < limit {
			randNum = uint16(r % rangeSize)
			break
		}
	}

	return min + randNum
}

func UInt16() uint16 {
	return UInt16Between(0, math.MaxUint16)
}

func UIntBetween(min, max uint32) uint32 {
	if min > max {
		min, max = max, min
	}

	// rangeSize como uint64 para evitar overflow (ex: max = math.MaxUint32, min = 0)
	rangeSize := uint64(max) - uint64(min) + 1

	// Gera número aleatório com uint64 para pegar todo o range
	randNum := uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
	randNum = randNum % rangeSize

	return min + uint32(randNum)
}

func UInt() uint32 {
	return UIntBetween(0, math.MaxUint32)
}

func UInt64Between(min, max uint64) uint64 {
	if min > max {
		min, max = max, min
	}

	// Calcula tamanho do intervalo como uint128 (conceitualmente)
	rangeSize := max - min + 1

	var randNum uint64

	if rangeSize == 0 {
		// Caso especial: intervalo completo, pois max - min + 1 == 0 no overflow uint64
		// Então retorna número aleatório completo
		randNum = rand.Uint64()
	} else {
		// Garante valor uniforme dentro do intervalo com mod
		// Usa rejection sampling para evitar viés (bias)
		limit := ^uint64(0) - (^uint64(0) % rangeSize) // maior múltiplo de rangeSize que cabe em uint64

		for {
			r := rand.Uint64()
			if r <= limit {
				randNum = r % rangeSize
				break
			}
		}
	}

	return min + randNum
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
	realPart := Float32()
	imagPart := Float32()
	return complex(realPart, imagPart)
}
func Complex64Between(minReal, maxReal, minImag, maxImag float32) complex64 {
	realPart := Float32Between(minReal, maxReal)
	imagPart := Float32Between(minImag, maxImag)
	return complex(realPart, imagPart)
}

func Complex128() complex128 {
	realPart := Float64()
	imagPart := Float64()
	return complex(realPart, imagPart)
}
func Complex128Between(minReal, maxReal, minImag, maxImag float64) complex128 {
	realPart := Float64Between(minReal, maxReal)
	imagPart := Float64Between(minImag, maxImag)
	return complex(realPart, imagPart)
}
