package GoGen

const AllChars string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()-_=+[]{}<>?/|\^~`
const Alphanumeric string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
const Alphabetic string = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
const AlphabeticUppercase string = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
const AlphabeticLowercase string = `abcdefghijklmnopqrstuvwxyz`
const Numeric string = `0123456789`
const Hexadecimal string = `0123456789ABCDEF`
const Symbols string = `!@#$%&*()-_=+[]{}<>?/|\^~`

func GenerateString(size int, sourceChars string) string {

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

// GenerateAllChars returns a string with a given size containing all characters
func GenerateAllChars(size int) string {
	return GenerateString(size, AllChars)
}

// GenerateAlphanumeric returns a string with a given size containing only alphanumeric characters
func GenerateAlphanumeric(size int) string {
	return GenerateString(size, Alphanumeric)
}

// GenerateAlphabetic returns a string with a given size containing only alphabetic characters
func GenerateAlphabetic(size int) string {
	return GenerateString(size, Alphabetic)
}

// GenerateAlphabeticUppercase returns a string with a given size containing only alphabetic characters in upper case only
func GenerateAlphabeticUppercase(size int) string {
	return GenerateString(size, AlphabeticUppercase)
}

// GenerateAlphabeticLowercase returns a string with a given size containing only alphabetic characters in lower case only
func GenerateAlphabeticLowercase(size int) string {
	return GenerateString(size, AlphabeticLowercase)
}

// GenerateNumeric returns a string with a given size containing only numeric characters
func GenerateNumeric(size int) string {
	return GenerateString(size, Numeric)
}

// GenerateHexadecimal returns a string with a given size containing only hexadecimal characters
func GenerateHexadecimal(size int) string {
	return GenerateString(size, Hexadecimal)
}

// GenerateSymbols returns a string with a given size containing only symbols characters
func GenerateSymbols(size int) string {
	return GenerateString(size, Symbols)
}
