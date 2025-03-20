package GoGen

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateString(t *testing.T) {

	// arrange
	size := 100
	srcChars := "abc"

	// act
	result := GenerateString(size, srcChars)

	// assert
	for _, char := range result {
		if strings.Index(srcChars, string(char)) == -1 {
			t.Errorf("char:%s not exist", string(char))
		}
	}

	t.Run("zero size", func(tt *testing.T) {
		// arrange
		size := 0
		srcChars := "abc"

		// act
		result := GenerateString(size, srcChars)

		// assert
		if len(result) != 0 {
			t.Errorf("GenerateString result length %d != 0", len(result))
		}
	})

	t.Run("negative size", func(tt *testing.T) {
		// arrange
		size := -3
		srcChars := "abc"

		// act
		result := GenerateString(size, srcChars)

		// assert
		if len(result) != 0 {
			t.Errorf("GenerateString result length %d != %d", len(result), size)
		}
	})

}

func ExampleGenerateString() {
	HexColor := GenerateString(6, "0123456789ABCDEF")
	fmt.Println(HexColor)
}

func BenchmarkGenerateString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateString(24, AllChars)
	}
}
