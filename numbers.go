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

func GenerateUInt8() uint8 {
	return GenerateUInt8Between(0, math.MaxUint8)
}

func GenerateUInt16Between(min, max uint16) uint16 {
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

func GenerateUInt16() uint16 {
	return GenerateUInt16Between(0, math.MaxUint16)
}

func GenerateUIntBetween(min, max uint32) uint32 {
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

func GenerateUInt() uint32 {
	return GenerateUIntBetween(0, math.MaxUint32)
}

func GenerateUInt64Between(min, max uint64) uint64 {
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
