package GoGen

type LengthTypeWords uint8

// SmallLengthWord represents small length words, 1 or 4 chars
const SmallLengthWord LengthTypeWords = 1

// MediumLengthWords represents medium length words, 5 to 8 chars
const MediumLengthWords LengthTypeWords = 2

// BigLengthWords represents Big length words, 9+ chars
const BigLengthWords LengthTypeWords = 3

// WordLengthRatio represents a ratio between LengthTypes
const WordLengthRatio = `11111111111112222222222222222222222222233333333333`

func GenerateWord() string {
	return GenerateAlphabetic(GenerateIntBetween(2, 30))
}

func GenerateWordByLengthType(l LengthTypeWords) (word string) {

	switch l {
	case SmallLengthWord:
		word = GenerateAlphabetic(GenerateIntBetween(1, 4))
	case MediumLengthWords:
		word = GenerateAlphabetic(GenerateIntBetween(5, 8))
	case BigLengthWords:
		word = GenerateAlphabetic(GenerateIntBetween(9, 30))
	default:
		word = GenerateAlphabetic(GenerateIntBetween(1, 30))
	}

	return
}

func GenerateWords(length int) (result []string) {
	result = make([]string, length, length)
	for i := 0; i < length; i++ {
		t := GenerateString(1, WordLengthRatio)
		switch t {
		case "1":
			result[i] = GenerateWordByLengthType(SmallLengthWord)
		case "2":
			result[i] = GenerateWordByLengthType(MediumLengthWords)
		case "3":
			result[i] = GenerateWordByLengthType(BigLengthWords)
		}
	}
	return
}
