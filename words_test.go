package gengo

import (
	"testing"
)

func TestWord(t *testing.T) {
	for i := 0; i < loop; i++ {
		word := Word()
		if len(word) < 2 || len(word) > 30 {
			t.Fatalf("Word() length %d out of range [2, 30]", len(word))
		}
	}
}

func TestWordByLengthType(t *testing.T) {
	t.Run("SmallLengthWord", func(t *testing.T) {
		for i := 0; i < loop; i++ {
			word := WordByLengthType(SmallLengthWord)
			if len(word) < 1 || len(word) > 4 {
				t.Fatalf("SmallLengthWord length %d out of range [1, 4]", len(word))
			}
		}
	})
	t.Run("MediumLengthWords", func(t *testing.T) {
		for i := 0; i < loop; i++ {
			word := WordByLengthType(MediumLengthWords)
			if len(word) < 5 || len(word) > 8 {
				t.Fatalf("MediumLengthWords length %d out of range [5, 8]", len(word))
			}
		}
	})
	t.Run("BigLengthWords", func(t *testing.T) {
		for i := 0; i < loop; i++ {
			word := WordByLengthType(BigLengthWords)
			if len(word) < 9 || len(word) > 30 {
				t.Fatalf("BigLengthWords length %d out of range [9, 30]", len(word))
			}
		}
	})
	t.Run("Default", func(t *testing.T) {
		for i := 0; i < loop; i++ {
			word := WordByLengthType(LengthTypeWords(99))
			if len(word) < 1 || len(word) > 30 {
				t.Fatalf("Default length %d out of range [1, 30]", len(word))
			}
		}
	})
}

func TestWords(t *testing.T) {
	const n = 100
	words := Words(n)
	if len(words) != n {
		t.Fatalf("Words(%d) returned %d words", n, len(words))
	}
	for i, word := range words {
		if len(word) < 1 || len(word) > 30 {
			t.Fatalf("Words(%d)[%d] length %d out of range [1, 30]", n, i, len(word))
		}
	}
}

func BenchmarkWord(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Word()
	}
}

func BenchmarkWordByLengthType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WordByLengthType(MediumLengthWords)
	}
}

func BenchmarkWords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Words(10)
	}
}
