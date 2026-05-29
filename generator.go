package gengo

import (
	"math"
	"math/bits"
	"math/rand/v2"
	"time"
)

// Generator produces random values from an explicit, seeded source, enabling
// reproducible output: two generators created with the same seed yield the same
// sequence of values. Its methods mirror the package-level functions exactly,
// drawing from the generator's own source instead of the global one.
//
// A Generator is NOT safe for concurrent use by multiple goroutines (it wraps a
// math/rand/v2 *Rand). Use one Generator per goroutine, or use the package-level
// functions (which are safe for concurrent use) when reproducibility is not
// required.
type Generator struct {
	r *rand.Rand
}

// New returns a Generator seeded with the given value. The same seed always
// produces the same sequence of values, which is useful for reproducing a
// specific run (for example, a failing test case).
func New(seed uint64) *Generator {
	return &Generator{r: rand.New(rand.NewPCG(seed, seed))}
}

// NewSource returns a Generator that draws from the given source, allowing any
// math/rand/v2 Source (for example rand.NewPCG or rand.NewChaCha8) to be used.
// The source must not be nil.
func NewSource(src rand.Source) *Generator {
	return &Generator{r: rand.New(src)}
}

// ---------- Integers ----------

// Int8Between is the seeded-generator equivalent of [Int8Between].
func (g *Generator) Int8Between(min, max int8) int8 {
	if min > max {
		min, max = max, min
	}
	return int8(int32(min) + g.r.Int32N(int32(max)-int32(min)+1))
}

// Int8 is the seeded-generator equivalent of [Int8].
func (g *Generator) Int8() int8 {
	return int8(g.r.Uint32())
}

// Int16Between is the seeded-generator equivalent of [Int16Between].
func (g *Generator) Int16Between(min, max int16) int16 {
	if min > max {
		min, max = max, min
	}
	return int16(int32(min) + g.r.Int32N(int32(max)-int32(min)+1))
}

// Int16 is the seeded-generator equivalent of [Int16].
func (g *Generator) Int16() int16 {
	return int16(g.r.Uint32())
}

// Int32Between is the seeded-generator equivalent of [Int32Between].
func (g *Generator) Int32Between(min, max int32) int32 {
	if min > max {
		min, max = max, min
	}
	return int32(int64(min) + g.r.Int64N(int64(max)-int64(min)+1))
}

// Int32 is the seeded-generator equivalent of [Int32].
func (g *Generator) Int32() int32 {
	return int32(g.r.Uint32())
}

// IntBetween is the seeded-generator equivalent of [IntBetween].
func (g *Generator) IntBetween(min, max int) int {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1
	if rangeSize == 0 {
		return int(g.r.Uint64())
	}
	return int(uint64(min) + g.r.Uint64N(rangeSize))
}

// Int is the seeded-generator equivalent of [Int].
func (g *Generator) Int() int {
	return int(g.r.Uint64())
}

// Int64Between is the seeded-generator equivalent of [Int64Between].
func (g *Generator) Int64Between(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1
	if rangeSize == 0 {
		return int64(g.r.Uint64())
	}
	return int64(uint64(min) + g.r.Uint64N(rangeSize))
}

// Int64 is the seeded-generator equivalent of [Int64].
func (g *Generator) Int64() int64 {
	return int64(g.r.Uint64())
}

// Byte is the seeded-generator equivalent of [Byte].
func (g *Generator) Byte() byte {
	return g.UInt8()
}

// UInt8Between is the seeded-generator equivalent of [UInt8Between].
func (g *Generator) UInt8Between(min, max uint8) uint8 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint32(max) - uint32(min) + 1
	return uint8(uint32(min) + g.r.Uint32N(rangeSize))
}

// UInt8 is the seeded-generator equivalent of [UInt8].
func (g *Generator) UInt8() uint8 {
	return uint8(g.r.Uint32())
}

// UInt16Between is the seeded-generator equivalent of [UInt16Between].
func (g *Generator) UInt16Between(min, max uint16) uint16 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint32(max) - uint32(min) + 1
	return uint16(uint32(min) + g.r.Uint32N(rangeSize))
}

// UInt16 is the seeded-generator equivalent of [UInt16].
func (g *Generator) UInt16() uint16 {
	return uint16(g.r.Uint32())
}

// UInt32Between is the seeded-generator equivalent of [UInt32Between].
func (g *Generator) UInt32Between(min, max uint32) uint32 {
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1
	if rangeSize > math.MaxUint32 {
		return g.r.Uint32()
	}
	return min + g.r.Uint32N(uint32(rangeSize))
}

// UInt32 is the seeded-generator equivalent of [UInt32].
func (g *Generator) UInt32() uint32 {
	return g.r.Uint32()
}

// UInt64Between is the seeded-generator equivalent of [UInt64Between].
func (g *Generator) UInt64Between(min, max uint64) uint64 {
	if min > max {
		min, max = max, min
	}
	rangeSize := max - min + 1
	if rangeSize == 0 {
		return g.r.Uint64()
	}
	return min + g.r.Uint64N(rangeSize)
}

// UInt64 is the seeded-generator equivalent of [UInt64].
func (g *Generator) UInt64() uint64 {
	return g.r.Uint64()
}

// ---------- Floats ----------

// Float32 is the seeded-generator equivalent of [Float32].
func (g *Generator) Float32() float32 {
	return g.r.Float32()*math.MaxFloat32 + math.SmallestNonzeroFloat32
}

// Float32Between is the seeded-generator equivalent of [Float32Between].
func (g *Generator) Float32Between(min, max float32) float32 {
	if min > max {
		min, max = max, min
	}
	span := max - min
	if span > math.MaxFloat32 {
		t := g.r.Float32()
		return min*(1-t) + max*t
	}
	return g.r.Float32()*span + min
}

// Float64 is the seeded-generator equivalent of [Float64].
func (g *Generator) Float64() float64 {
	return g.r.Float64()*math.MaxFloat64 + math.SmallestNonzeroFloat64
}

// Float64Between is the seeded-generator equivalent of [Float64Between].
func (g *Generator) Float64Between(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	span := max - min
	if span > math.MaxFloat64 {
		t := g.r.Float64()
		return min*(1-t) + max*t
	}
	return g.r.Float64()*span + min
}

// ---------- Complex ----------

// Complex64 is the seeded-generator equivalent of [Complex64].
func (g *Generator) Complex64() complex64 {
	return complex(g.r.Float32(), g.r.Float32())
}

// Complex64Between is the seeded-generator equivalent of [Complex64Between].
func (g *Generator) Complex64Between(minReal, maxReal, minImag, maxImag float32) complex64 {
	return complex(g.Float32Between(minReal, maxReal), g.Float32Between(minImag, maxImag))
}

// Complex128 is the seeded-generator equivalent of [Complex128].
func (g *Generator) Complex128() complex128 {
	return complex(g.r.Float64(), g.r.Float64())
}

// Complex128Between is the seeded-generator equivalent of [Complex128Between].
func (g *Generator) Complex128Between(minReal, maxReal, minImag, maxImag float64) complex128 {
	return complex(g.Float64Between(minReal, maxReal), g.Float64Between(minImag, maxImag))
}

// ---------- Strings ----------

// String is the seeded-generator equivalent of [String].
func (g *Generator) String(length uint32, sourceChars string) string {
	if length == 0 || sourceChars == "" {
		return ""
	}

	if length > MaxStringLength {
		length = MaxStringLength
	}

	srcLen := uint64(len(sourceChars))
	result := make([]byte, length)

	if srcLen == 1 {
		for i := range result {
			result[i] = sourceChars[0]
		}
		return string(result)
	}

	bitsPerChar := uint(bits.Len64(srcLen - 1))
	mask := uint64((1 << bitsPerChar) - 1)
	var rnd uint64
	var bitsLeft uint

	for i := uint32(0); i < length; {
		if bitsLeft < bitsPerChar {
			rnd = g.r.Uint64()
			bitsLeft = 64
		}
		idx := rnd & mask
		rnd >>= bitsPerChar
		bitsLeft -= bitsPerChar
		if idx < srcLen {
			result[i] = sourceChars[idx]
			i++
		}
	}

	return string(result)
}

// StringBetween is the seeded-generator equivalent of [StringBetween].
func (g *Generator) StringBetween(min, max uint32, sourceChars string) string {
	if min > max {
		min, max = max, min
	}
	return g.String(g.UInt32Between(min, max), sourceChars)
}

// StringAllChars is the seeded-generator equivalent of [StringAllChars].
func (g *Generator) StringAllChars(length uint32) string {
	return g.String(length, AllChars)
}

// StringAlphanumeric is the seeded-generator equivalent of [StringAlphanumeric].
func (g *Generator) StringAlphanumeric(length uint32) string {
	return g.String(length, Alphanumeric)
}

// StringAlphabetic is the seeded-generator equivalent of [StringAlphabetic].
func (g *Generator) StringAlphabetic(length uint32) string {
	return g.String(length, Alphabetic)
}

// StringAlphabeticUppercase is the seeded-generator equivalent of [StringAlphabeticUppercase].
func (g *Generator) StringAlphabeticUppercase(length uint32) string {
	return g.String(length, AlphabeticUppercase)
}

// StringAlphabeticLowercase is the seeded-generator equivalent of [StringAlphabeticLowercase].
func (g *Generator) StringAlphabeticLowercase(length uint32) string {
	return g.String(length, AlphabeticLowercase)
}

// StringNumeric is the seeded-generator equivalent of [StringNumeric].
func (g *Generator) StringNumeric(length uint32) string {
	return g.String(length, Numeric)
}

// StringHexadecimal is the seeded-generator equivalent of [StringHexadecimal].
func (g *Generator) StringHexadecimal(length uint32) string {
	return g.String(length, Hexadecimal)
}

// StringSymbols is the seeded-generator equivalent of [StringSymbols].
func (g *Generator) StringSymbols(length uint32) string {
	return g.String(length, Symbols)
}

// ---------- Dates ----------

// Date is the seeded-generator equivalent of [Date].
func (g *Generator) Date() time.Time {
	return time.Unix(g.r.Int64N(dtRange)+dtMin, 0).UTC()
}

// UnixDate is the seeded-generator equivalent of [UnixDate].
func (g *Generator) UnixDate() time.Time {
	return time.Unix(g.r.Int64N(dtUnixRange)+dtUnixMin, 0).UTC()
}

// DateBetween is the seeded-generator equivalent of [DateBetween].
func (g *Generator) DateBetween(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}

	startUnix := start.Unix()
	endUnix := end.Unix()

	if startUnix == endUnix {
		return start
	}

	randomUnix := g.r.Int64N(endUnix-startUnix+1) + startUnix

	return time.Unix(randomUnix, 0).UTC()
}

// ---------- Bool ----------

// Bool is the seeded-generator equivalent of [Bool].
func (g *Generator) Bool() bool {
	return g.r.Uint32()&1 == 1
}

// ---------- Words ----------

// Word is the seeded-generator equivalent of [Word].
func (g *Generator) Word() string {
	return g.StringAlphabeticLowercase(g.UInt32Between(2, 30))
}

// WordByLengthType is the seeded-generator equivalent of [WordByLengthType].
func (g *Generator) WordByLengthType(l LengthTypeWords) (word string) {
	switch l {
	case SmallLengthWord:
		word = g.StringAlphabeticLowercase(g.UInt32Between(1, 4))
	case MediumLengthWords:
		word = g.StringAlphabeticLowercase(g.UInt32Between(5, 8))
	case BigLengthWords:
		word = g.StringAlphabeticLowercase(g.UInt32Between(9, 30))
	default:
		word = g.StringAlphabeticLowercase(g.UInt32Between(1, 30))
	}

	return word
}

// Words is the seeded-generator equivalent of [Words].
func (g *Generator) Words(length int) []string {
	if length <= 0 {
		return []string{}
	}
	result := make([]string, length)
	for i := 0; i < length; i++ {
		switch WordLengthRatio[g.r.IntN(len(WordLengthRatio))] {
		case '1':
			result[i] = g.WordByLengthType(SmallLengthWord)
		case '2':
			result[i] = g.WordByLengthType(MediumLengthWords)
		case '3':
			result[i] = g.WordByLengthType(BigLengthWords)
		}
	}
	return result
}
