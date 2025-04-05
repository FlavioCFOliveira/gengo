package gengo

import (
	"testing"
)

func TestWord(t *testing.T) {
	word := Word()
	println(word)
}

func TestWordByLengthType(t *testing.T) {
	t.Run("Small-Length-Words", func(t *testing.T) {
		WordByLengthType(SmallLengthWord)
	})

	t.Run("Medium-Length-Words", func(t *testing.T) {
		WordByLengthType(MediumLengthWords)
	})

	t.Run("big-Length-Words", func(t *testing.T) {
		WordByLengthType(BigLengthWords)
	})
}

func TestWords(t *testing.T) {
	Words(100)
}
