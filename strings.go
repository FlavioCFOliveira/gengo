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
func String(length int, sourceChars string) string {
	if length <= 0 || len(sourceChars) == 0 {
		return ""
	}

	srcLen := len(sourceChars)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = sourceChars[rnd.Intn(srcLen)]
	}

	return string(result)
}

// StringBetween Generates a string with a variable length using only characters of a given source.
func StringBetween(min, max int, sourceChars string) string {
	if min > max {
		min, max = max, min // swap to ensure valid range
	}

	if max < 1 {
		return ""
	}

	length := IntBetween(min, max)
	for length < 0 {
		length = IntBetween(min, max)
	}

	return String(length, sourceChars)
}

// StringAllChars returns a string with a given length containing all the predefined alphabetic, alphanumeric and symbols characters
func StringAllChars(length int) string {
	return String(length, AllChars)
}

// StringAlphanumeric returns a string with a given length containing only alphanumeric characters
func StringAlphanumeric(length int) string {
	return String(length, Alphanumeric)
}

// StringAlphabetic returns a string with a given length containing only alphabetic characters
func StringAlphabetic(length int) string {
	return String(length, Alphabetic)
}

// StringAlphabeticUppercase returns a string with a given length containing only alphabetic characters in upper case only
func StringAlphabeticUppercase(length int) string {
	return String(length, AlphabeticUppercase)
}

// StringAlphabeticLowercase returns a string with a given length containing only alphabetic characters in lower case only
func StringAlphabeticLowercase(length int) string {
	return String(length, AlphabeticLowercase)
}

// StringNumeric returns a string with a given length containing only numeric characters
func StringNumeric(length int) string {
	return String(length, Numeric)
}

// StringHexadecimal returns a string with a given length containing only hexadecimal characters
func StringHexadecimal(length int) string {
	return String(length, Hexadecimal)
}

// StringSymbols returns a string with a given length containing only symbols characters
func StringSymbols(length int) string {
	return String(length, Symbols)
}
