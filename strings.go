package gengo

const (
	AllChars            string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()-_=+[]{}<>?/|\^~`
	Alphanumeric        string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	Alphabetic          string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	AlphabeticUppercase string = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	AlphabeticLowercase string = `abcdefghijklmnopqrstuvwxyz`
	Numeric             string = `0123456789`
	Hexadecimal         string = `0123456789ABCDEF`
	Symbols             string = `!@#$%&*()-_=+[]{}<>?/|\^~`
)

// String Generates a string with a given length using only characters of a given source.
func String(length uint32, sourceChars string) string {
	if length <= 0 || len(sourceChars) == 0 {
		return ""
	}

	srcLen := len(sourceChars)
	result := make([]byte, length)

	for i := 0; i < int(length); i++ {
		result[i] = sourceChars[rnd.Intn(srcLen)]
	}

	return string(result)
}

// StringBetween Generates a string with a variable length using only characters of a given source.
func StringBetween(min, max uint32, sourceChars string) string {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}

	diff := uint64(max) - uint64(min)

	if min == max || diff < 2 {
		return String(min, sourceChars)
	}

	return String(UInt32Between(min, max), sourceChars)
}

// StringAllChars returns a string with a given length containing all the predefined alphabetic, alphanumeric and symbols characters
func StringAllChars(length uint32) string {
	return String(length, AllChars)
}

// StringAlphanumeric returns a string with a given length containing only alphanumeric characters
func StringAlphanumeric(length uint32) string {
	return String(length, Alphanumeric)
}

// StringAlphabetic returns a string with a given length containing only alphabetic characters
func StringAlphabetic(length uint32) string {
	return String(length, Alphabetic)
}

// StringAlphabeticUppercase returns a string with a given length containing only alphabetic characters in upper case only
func StringAlphabeticUppercase(length uint32) string {
	return String(length, AlphabeticUppercase)
}

// StringAlphabeticLowercase returns a string with a given length containing only alphabetic characters in lower case only
func StringAlphabeticLowercase(length uint32) string {
	return String(length, AlphabeticLowercase)
}

// StringNumeric returns a string with a given length containing only numeric characters
func StringNumeric(length uint32) string {
	return String(length, Numeric)
}

// StringHexadecimal returns a string with a given length containing only hexadecimal characters
func StringHexadecimal(length uint32) string {
	return String(length, Hexadecimal)
}

// StringSymbols returns a string with a given length containing only symbols characters
func StringSymbols(length uint32) string {
	return String(length, Symbols)
}
