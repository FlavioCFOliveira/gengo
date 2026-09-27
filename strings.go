package gengo

import (
	"math/bits"
	"math/rand/v2"
	"strings"
)

// Character sets available for string generation.
const (
	AllChars            string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()-_=+[]{}<>?/|\^~`
	Alphanumeric        string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	Alphabetic          string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	AlphabeticUppercase string = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	AlphabeticLowercase string = `abcdefghijklmnopqrstuvwxyz`
	Numeric             string = `0123456789`
	Hexadecimal         string = `0123456789ABCDEF`
	Symbols             string = `!@#$%&*()-_=+[]{}<>?/|\^~`
	MaxStringLength     uint32 = 1024 * 1024 // 1MB maximum string length
)

// String build thresholds. Every String call allocates exactly once, for the
// returned string, whatever the length:
//
//   - up to stackStringMaxLength bytes, the characters are written to a
//     make([]byte, length) that the gc compiler places on the stack (its
//     default variable-size make threshold is 32 bytes) and then copied into
//     the string;
//   - up to stringChunkLength bytes, they are written to a fixed stack array
//     and copied into the string;
//   - above that, they are written chunk by chunk to a strings.Builder
//     pre-sized to the final length.
const (
	stackStringMaxLength = 32
	stringChunkLength    = 256
)

// String generates a length-byte string by sampling bytes from sourceChars.
//
// It is byte-oriented: length counts bytes (not runes), characters are picked
// by byte index, and sourceChars is expected to hold only single-byte (ASCII)
// characters — every predefined character set in this package qualifies.
// Passing a multibyte charset (emoji, accented letters, CJK, and so on) samples
// individual UTF-8 bytes and therefore yields invalid UTF-8, so use an ASCII
// charset instead. length is capped at MaxStringLength bytes.
func String(length uint32, sourceChars string) string {
	if length > stackStringMaxLength {
		return stringLongEntry(length, sourceChars, nil)
	}
	if length == 0 || sourceChars == "" {
		return ""
	}

	if length > MaxStringLength {
		length = MaxStringLength
	}

	srcLen := uint64(len(sourceChars))
	result := make([]byte, length)

	// Fast path: single-char source — fill without any PRNG call.
	if srcLen == 1 {
		for i := range result {
			result[i] = sourceChars[0]
		}
		return string(result)
	}

	// Batch bit extraction: pull bitsPerChar bits at a time from a single
	// rand.Uint64() word, refilling only when exhausted. This reduces PRNG
	// invocations by 6–16x compared to one rand.IntN call per character.
	bitsPerChar := uint(bits.Len64(srcLen - 1)) //nolint:gosec // bits.Len64 returns [0,64]; conversion to uint is always safe
	mask := uint64((1 << bitsPerChar) - 1)
	var rnd uint64
	var bitsLeft uint

	for i := uint32(0); i < length; {
		if bitsLeft < bitsPerChar {
			rnd = rand.Uint64()
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

// stringLongEntry is the String and Generator.String path for a length above
// stackStringMaxLength: it applies the same empty-source result and
// MaxStringLength cap as String, then builds the string with a single
// allocation, drawing from r, or from the global math/rand/v2 generator when r
// is nil. Both callers dispatch to it first, so their code for lengths up to
// stackStringMaxLength is unchanged.
func stringLongEntry(length uint32, sourceChars string, r *rand.Rand) string {
	if sourceChars == "" {
		return ""
	}
	if length > MaxStringLength {
		length = MaxStringLength
	}
	if len(sourceChars) == 1 {
		return stringSingleLong(length, sourceChars)
	}
	return stringLong(length, sourceChars, r)
}

// stringBits carries the batched bit-extraction state of a string being built
// in several fills, so that the fills consume the random source exactly as a
// single pass over the whole string would.
type stringBits struct {
	rnd      uint64
	bitsLeft uint
}

// fill writes len(dst) characters sampled from sourceChars into dst, drawing
// from the global math/rand/v2 generator. It uses the same bit-extraction loop
// as the short path of [String], so it draws the same values in the same order.
// sourceChars must hold at least two bytes.
//
// fill and fillFrom differ only in their random source and in the "& 63" mask
// that fillFrom applies to bitsPerChar. The mask never changes a value
// (len(sourceChars) < 1<<63, so bitsPerChar is at most 63); each form was
// chosen because it measured fastest for its function on the reference
// machine (see BENCHMARKS.md).
func (s *stringBits) fill(dst []byte, sourceChars string) {
	srcLen := uint64(len(sourceChars))
	bitsPerChar := uint(bits.Len64(srcLen - 1)) //nolint:gosec // bits.Len64 returns [0,64]; conversion to uint is always safe
	mask := uint64((1 << bitsPerChar) - 1)
	rnd, bitsLeft := s.rnd, s.bitsLeft
	for i := 0; i < len(dst); {
		if bitsLeft < bitsPerChar {
			rnd = rand.Uint64()
			bitsLeft = 64
		}
		idx := rnd & mask
		rnd >>= bitsPerChar
		bitsLeft -= bitsPerChar
		if idx < srcLen {
			dst[i] = sourceChars[idx]
			i++
		}
	}
	s.rnd, s.bitsLeft = rnd, bitsLeft
}

// fillFrom is fill drawing from r instead of the global generator.
func (s *stringBits) fillFrom(r *rand.Rand, dst []byte, sourceChars string) {
	srcLen := uint64(len(sourceChars))
	bitsPerChar := uint(bits.Len64(srcLen-1)) & 63 //nolint:gosec // bits.Len64 returns [0,64]; conversion to uint is always safe
	mask := uint64((1 << bitsPerChar) - 1)
	rnd, bitsLeft := s.rnd, s.bitsLeft
	for i := 0; i < len(dst); {
		if bitsLeft < bitsPerChar {
			rnd = r.Uint64()
			bitsLeft = 64
		}
		idx := rnd & mask
		rnd >>= bitsPerChar
		bitsLeft -= bitsPerChar
		if idx < srcLen {
			dst[i] = sourceChars[idx]
			i++
		}
	}
	s.rnd, s.bitsLeft = rnd, bitsLeft
}

// fillChunk calls fillFrom with r, or fill when r is nil. Selecting the source
// once per chunk keeps the per-character loops free of the choice.
func (s *stringBits) fillChunk(r *rand.Rand, dst []byte, sourceChars string) {
	if r == nil {
		s.fill(dst, sourceChars)
		return
	}
	s.fillFrom(r, dst, sourceChars)
}

// stringLong builds a random string longer than stackStringMaxLength with a
// single allocation, drawing from r, or from the global math/rand/v2 generator
// when r is nil. sourceChars must hold at least two bytes.
func stringLong(length uint32, sourceChars string, r *rand.Rand) string {
	var s stringBits
	var chunk [stringChunkLength]byte
	if length <= stringChunkLength {
		s.fillChunk(r, chunk[:length], sourceChars)
		return string(chunk[:length])
	}
	var sb strings.Builder
	sb.Grow(int(length))
	for rem := length; rem > 0; {
		n := min(rem, stringChunkLength)
		s.fillChunk(r, chunk[:n], sourceChars)
		sb.Write(chunk[:n])
		rem -= n
	}
	return sb.String()
}

// stringSingleLong returns the single-byte string c repeated length times, for
// a length above stackStringMaxLength, with a single allocation.
func stringSingleLong(length uint32, c string) string {
	if length <= stringChunkLength {
		var chunk [stringChunkLength]byte
		chunk[0] = c[0]
		for filled := uint32(1); filled < length; filled *= 2 {
			copy(chunk[filled:length], chunk[:filled])
		}
		return string(chunk[:length])
	}
	return strings.Repeat(c, int(length))
}

// StringBetween generates a string with a variable length using only characters of a given source.
func StringBetween(min, max uint32, sourceChars string) string {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return String(Uint32Between(min, max), sourceChars)
}

// StringAllChars returns a string with a given length containing all the predefined alphabetic, alphanumeric and symbols characters.
func StringAllChars(length uint32) string {
	return String(length, AllChars)
}

// StringAlphanumeric returns a string with a given length containing only alphanumeric characters.
func StringAlphanumeric(length uint32) string {
	return String(length, Alphanumeric)
}

// StringAlphabetic returns a string with a given length containing only alphabetic characters.
func StringAlphabetic(length uint32) string {
	return String(length, Alphabetic)
}

// StringAlphabeticUppercase returns a string with a given length containing only alphabetic characters in upper case only.
func StringAlphabeticUppercase(length uint32) string {
	return String(length, AlphabeticUppercase)
}

// StringAlphabeticLowercase returns a string with a given length containing only alphabetic characters in lower case only.
func StringAlphabeticLowercase(length uint32) string {
	return String(length, AlphabeticLowercase)
}

// StringNumeric returns a string with a given length containing only numeric characters.
func StringNumeric(length uint32) string {
	return String(length, Numeric)
}

// StringHexadecimal returns a string with a given length containing only hexadecimal characters.
func StringHexadecimal(length uint32) string {
	return String(length, Hexadecimal)
}

// StringSymbols returns a string with a given length containing only symbols characters.
func StringSymbols(length uint32) string {
	return String(length, Symbols)
}
