package gengo

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestString(t *testing.T) {
	const size = 100
	const srcChars = "abc"
	result := String(size, srcChars)
	if len(result) != size {
		t.Fatalf("String(%d, ...) length = %d", size, len(result))
	}
	for _, ch := range result {
		if !strings.ContainsRune(srcChars, ch) {
			t.Fatalf("String: char %q not in source %q", ch, srcChars)
		}
	}
	// zero size
	if got := String(0, srcChars); got != "" {
		t.Fatalf("String(0, ...) length = %d, want 0", len(got))
	}
	// empty source
	if got := String(10, ""); got != "" {
		t.Fatalf("String(10, \"\") length = %d, want 0", len(got))
	}
	// length exceeding MaxStringLength must be capped
	over := String(MaxStringLength+1, "ab")
	if len(over) != 1024*1024 {
		t.Fatalf("String(MaxStringLength+1) length = %d, want %d", len(over), MaxStringLength)
	}
	// exact MaxStringLength boundary must not be capped
	exact := String(MaxStringLength, "ab")
	if len(exact) != 1024*1024 {
		t.Fatalf("String(MaxStringLength) length = %d, want %d", len(exact), MaxStringLength)
	}
	// single-character source must produce a string of only that character
	const singleN = 20
	single := String(singleN, "a")
	if len(single) != singleN {
		t.Fatalf("String(%d, \"a\") length = %d", singleN, len(single))
	}
	for _, ch := range single {
		if ch != 'a' {
			t.Fatalf("String(%d, \"a\") contains non-'a' char %q", singleN, ch)
		}
	}
}

// TestStringByteOriented documents the byte-oriented (ASCII) contract: length
// counts bytes, and an ASCII source always yields valid UTF-8 composed solely
// of bytes drawn from the charset.
func TestStringByteOriented(t *testing.T) {
	const charset = "ACGT0123456789xyz!@#"
	const n = 64
	for i := 0; i < loop; i++ {
		result := String(n, charset)
		if len(result) != n {
			t.Fatalf("String(%d, ...) = %d bytes, want %d", n, len(result), n)
		}
		if !utf8.ValidString(result) {
			t.Fatalf("String(%d, %q) = %q: invalid UTF-8", n, charset, result)
		}
		for j := 0; j < len(result); j++ {
			if strings.IndexByte(charset, result[j]) < 0 {
				t.Fatalf("String produced byte %q not in charset %q", result[j], charset)
			}
		}
	}
}

// testStringCharset verifies that fn(size) always returns a string of exactly
// size bytes whose every character belongs to charset.
func testStringCharset(t *testing.T, fn func(uint32) string, charset, name string) {
	t.Helper()
	const size = 50
	for i := 0; i < loop; i++ {
		result := fn(size)
		if len(result) != size {
			t.Fatalf("%s(%d) length = %d", name, size, len(result))
		}
		for _, ch := range result {
			if !strings.ContainsRune(charset, ch) {
				t.Fatalf("%s: char %q not in charset", name, ch)
			}
		}
	}
}

func TestStringAllChars(t *testing.T) {
	testStringCharset(t, StringAllChars, AllChars, "StringAllChars")
}

func TestStringAlphanumeric(t *testing.T) {
	testStringCharset(t, StringAlphanumeric, Alphanumeric, "StringAlphanumeric")
}

func TestStringAlphabetic(t *testing.T) {
	testStringCharset(t, StringAlphabetic, Alphabetic, "StringAlphabetic")
}

func TestStringAlphabeticUppercase(t *testing.T) {
	testStringCharset(t, StringAlphabeticUppercase, AlphabeticUppercase, "StringAlphabeticUppercase")
}

func TestStringAlphabeticLowercase(t *testing.T) {
	testStringCharset(t, StringAlphabeticLowercase, AlphabeticLowercase, "StringAlphabeticLowercase")
}

func TestStringNumeric(t *testing.T) {
	testStringCharset(t, StringNumeric, Numeric, "TestStringNumeric")
}

func TestStringHexadecimal(t *testing.T) {
	testStringCharset(t, StringHexadecimal, Hexadecimal, "StringHexadecimal")
}

func TestStringSymbols(t *testing.T) {
	testStringCharset(t, StringSymbols, Symbols, "StringSymbols")
}

func TestStringBetween(t *testing.T) {
	const min, max = 5, 20
	for i := 0; i < loop; i++ {
		result := StringBetween(min, max, Alphanumeric)
		l := len(result)
		if l < min || l > max {
			t.Fatalf("StringBetween(%d, %d) length = %d: out of range", min, max, l)
		}
	}
	// swapped args must produce same valid range
	for i := 0; i < loop; i++ {
		result := StringBetween(max, min, Alphanumeric)
		l := len(result)
		if l < min || l > max {
			t.Fatalf("StringBetween(%d, %d) length = %d: out of range", max, min, l)
		}
	}
	// min == max must return exactly min length
	if got := StringBetween(10, 10, Alphanumeric); len(got) != 10 {
		t.Fatalf("StringBetween(10, 10) length = %d, want 10", len(got))
	}
	// diff == 1 must return either min or max length
	for i := 0; i < loop; i++ {
		result := StringBetween(5, 6, Alphanumeric)
		l := len(result)
		if l < 5 || l > 6 {
			t.Fatalf("StringBetween(5, 6) length = %d: out of range [5, 6]", l)
		}
	}
	// both bounds zero must return ""
	if got := StringBetween(0, 0, Alphanumeric); got != "" {
		t.Fatalf("StringBetween(0, 0) length = %d, want 0", len(got))
	}
}

func ExampleString() {
	HexColor := String(8, "0123456789ABCDEF")
	fmt.Println(HexColor)
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(8, Alphanumeric)
	}
}

func BenchmarkStringNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringNumeric(8)
	}
}
