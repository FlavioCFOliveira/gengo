package gengo

import (
	"math/bits"
	"math/rand/v2"
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

// String generates a string with a given length using only characters of a given source.
func String(length uint32, sourceChars string) string {
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

// StringBetween generates a string with a variable length using only characters of a given source.
func StringBetween(min, max uint32, sourceChars string) string {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}
	return String(UInt32Between(min, max), sourceChars)
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
