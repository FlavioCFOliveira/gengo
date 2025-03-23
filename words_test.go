package GoGen

import (
	"strings"
	"testing"
)

func TestGenerateWord(t *testing.T) {
	word := GenerateWord()
	println(word)
}

func TestGenerateGenerateWordByLengthType(t *testing.T) {

	t.Run("Small-Length-Words", func(t *testing.T) {
		word := GenerateWordByLengthType(SmallLengthWord)
		println(word)
	})

	t.Run("Medium-Length-Words", func(t *testing.T) {
		word := GenerateWordByLengthType(MediumLengthWords)
		println(word)
	})

	t.Run("big-Length-Words", func(t *testing.T) {
		word := GenerateWordByLengthType(BigLengthWords)
		println(word)
	})

}

func TestGenerateWords(t *testing.T) {
	words := GenerateWords(100)
	println(strings.Join(words, " "))
}
