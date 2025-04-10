package gengo

import (
	"fmt"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	// arrange
	size := uint32(100)
	srcChars := "abc"

	// act
	result := String(size, srcChars)

	// assert
	for _, char := range result {
		if !strings.ContainsRune(srcChars, char) {
			t.Errorf("char: %q doesn't exist", char)
		}
	}

	t.Run("zero size", func(tt *testing.T) {
		// arrange
		size := uint32(0)
		srcChars := "abc"

		// act
		result := String(size, srcChars)

		// assert
		if len(result) != 0 {
			t.Errorf("String result length %d != 0", len(result))
		}
	})

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

func BenchmarkNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringNumeric(8)
	}
}
