package gen

const AllChars string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()-_=+[]{}<>?/|\^~`
const Alphanumeric string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
const Alphabetic string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
const AlphabeticUppercase string = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
const AlphabeticLowercase string = `abcdefghijklmnopqrstuvwxyz`
const Numeric string = `0123456789`
const Hexadecimal string = `0123456789ABCDEF`
const Symbols string = `!@#$%&*()-_=+[]{}<>?/|\^~`

func String(size int, sourceChars string) string {

	if size <= 0 || len(sourceChars) == 0 {
		return ""
	}

	srcLen := len(sourceChars)
	result := make([]byte, size, size)
	for i := 0; i < size; i++ {
		result[i] = sourceChars[rnd.Intn(srcLen)]
	}

	return string(result)
	//return *(*string)(unsafe.Pointer(&result))
}

// StringAllChars returns a string with a given size containing all characters
func StringAllChars(size int) string {
	return String(size, AllChars)
}

// StringAlphanumeric returns a string with a given size containing only alphanumeric characters
func StringAlphanumeric(size int) string {
	return String(size, Alphanumeric)
}

// StringAlphabetic returns a string with a given size containing only alphabetic characters
func StringAlphabetic(size int) string {
	return String(size, Alphabetic)
}

// StringAlphabeticUppercase returns a string with a given size containing only alphabetic characters in upper case only
func StringAlphabeticUppercase(size int) string {
	return String(size, AlphabeticUppercase)
}

// StringAlphabeticLowercase returns a string with a given size containing only alphabetic characters in lower case only
func StringAlphabeticLowercase(size int) string {
	return String(size, AlphabeticLowercase)
}

// StringNumeric returns a string with a given size containing only numeric characters
func StringNumeric(size int) string {
	return String(size, Numeric)
}

// StringHexadecimal returns a string with a given size containing only hexadecimal characters
func StringHexadecimal(size int) string {
	return String(size, Hexadecimal)
}

// StringSymbols returns a string with a given size containing only symbols characters
func StringSymbols(size int) string {
	return String(size, Symbols)
}
