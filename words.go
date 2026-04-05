package gengo

import "math/rand/v2"

// LengthTypeWords classifies a word by its character length category.
type LengthTypeWords uint8

// SmallLengthWord represents small length words, 1 or 4 chars.
const SmallLengthWord LengthTypeWords = 1

// MediumLengthWords represents medium length words, 5 to 8 chars.
const MediumLengthWords LengthTypeWords = 2

// BigLengthWords represents Big length words, 9+ chars.
const BigLengthWords LengthTypeWords = 3

// WordLengthRatio represents a ratio between LengthTypes.
const WordLengthRatio = `11111111111112222222222222222222222222233333333333`

// Word returns a random lowercase alphabetic word between 2 and 30 characters.
func Word() string {
	return StringAlphabeticLowercase(UInt32Between(2, 30))
}

// WordByLengthType returns a random lowercase alphabetic word of the given length category.
func WordByLengthType(l LengthTypeWords) (word string) {
	switch l {
	case SmallLengthWord:
		word = StringAlphabeticLowercase(UInt32Between(1, 4))
	case MediumLengthWords:
		word = StringAlphabeticLowercase(UInt32Between(5, 8))
	case BigLengthWords:
		word = StringAlphabeticLowercase(UInt32Between(9, 30))
	default:
		word = StringAlphabeticLowercase(UInt32Between(1, 30))
	}

	return word
}

// Words returns a slice of length random words with a natural distribution across length categories.
// Returns an empty slice when length <= 0.
func Words(length int) (result []string) {
	if length <= 0 {
		return []string{}
	}
	result = make([]string, length)
	for i := 0; i < length; i++ {
		switch WordLengthRatio[rand.IntN(len(WordLengthRatio))] {
		case '1':
			result[i] = WordByLengthType(SmallLengthWord)
		case '2':
			result[i] = WordByLengthType(MediumLengthWords)
		case '3':
			result[i] = WordByLengthType(BigLengthWords)
		}
	}
	return result
}
