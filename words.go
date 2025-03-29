package gengo

type LengthTypeWords uint8

// SmallLengthWord represents small length words, 1 or 4 chars
const SmallLengthWord LengthTypeWords = 1

// MediumLengthWords represents medium length words, 5 to 8 chars
const MediumLengthWords LengthTypeWords = 2

// BigLengthWords represents Big length words, 9+ chars
const BigLengthWords LengthTypeWords = 3

// WordLengthRatio represents a ratio between LengthTypes
const WordLengthRatio = `11111111111112222222222222222222222222233333333333`

func Word() string {
	return StringAlphabeticLowercase(IntBetween(2, 30))
}

func WordByLengthType(l LengthTypeWords) (word string) {
	switch l {
	case SmallLengthWord:
		word = StringAlphabeticLowercase(IntBetween(1, 4))
	case MediumLengthWords:
		word = StringAlphabeticLowercase(IntBetween(5, 8))
	case BigLengthWords:
		word = StringAlphabeticLowercase(IntBetween(9, 30))
	default:
		word = StringAlphabeticLowercase(IntBetween(1, 30))
	}

	return
}

func Words(length int) (result []string) {
	result = make([]string, length, length)
	for i := 0; i < length; i++ {
		t := String(1, WordLengthRatio)
		switch t {
		case "1":
			result[i] = WordByLengthType(SmallLengthWord)
		case "2":
			result[i] = WordByLengthType(MediumLengthWords)
		case "3":
			result[i] = WordByLengthType(BigLengthWords)
		}
	}
	return
}
