# gengo

[![build](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml/badge.svg)](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml)

gengo is a Go Module containing utility functions to generate random data.

This module is also developed having performance and resource usage 
considerations in mind, so it can be used in production environments.

Download and install the Go module
```sh
go get -u github.com/FlavioCFOliveira/gengo
```
## Quick start
To start generating random data you just need to import the package and 
start using the package API as follows on this quick example:  

```go
package main

import (
	"fmt"
	"github.com/FlavioCFOliveira/gengo"
)

func main() {
    // Generate a random alfanumeric string with 33 chars
    randString := gengo.StringAlphanumeric(33)
    fmt.Println(randString)
}
```

## Generate random strings

### String(length int, sourceChars string) string
The `gengo.String(length, sourceChars)` func generates a random string with a given length using only characters of a given source..
```go
    // Specify a variable containing the characters you want to be used on the new string
    sourceChars := "YourCustomSouce0fCh4r4ct3rs"
    
    customRandString := gengo.String(33, sourceChars)
    fmt.Println(customRandString)
```



### StringAllChars(length int) string
The `gengo.StringAllChars(length int)` func generates a random string with a given length containing all the predefined alphabetic, alphanumeric and symbols characters.

This uses the predefined constant string `AllChars` as characters source reference.

```go 
    // generates using letters, numbers and symbols
    randomString := gengo.StringAllChars(33)
    fmt.Println(randomString)
```



### StringAlphanumeric(length int) string
The `gengo.StringAlphanumeric(length int)` func generates a random string with a given length containing all the predefined alphabetic and alphanumeric characters.

This uses the predefined constant string `Alphanumeric` as characters source reference.

```go 
    // generates using letters and numbers
    randomString := gengo.StringAlphanumeric(33)
    fmt.Println(randomString)
```


### StringAlphabetic(length int) string
The `gengo.Alphabetic(length int)` func generates a random string with a given length containing all the predefined alphabetic characters.

This func uses the predefined constant string `Alphabetic` as characters source reference.

There are also two variants to generate upper case and lowe case only strings, that uses the predefined constant string `AlphabeticUppercase` and `AlphabeticLowercase` as characters source reference.

```go 
    // upper and lower case characters   
    randomString := gengo.Alphabetic(33)
    fmt.Println(randomString)

    // upper case only characters   
    upperString := gengo.StringAlphabeticUppercase(33)
    fmt.Println(randomString)

    // lower case only characters   
    lowerString := gengo.StringAlphabeticLowercase(33)
    fmt.Println(randomString)
```


### StringNumeric(length int) string
The `gengo.StringNumeric(length int)` func generates a random string with a given length containing all the predefined numeric characters.

This func uses the predefined constant string `Numeric` as characters source reference.

```go 
    randomString := gengo.StringNumeric(33)
    fmt.Println(randomString)
```


### StringHexadecimal(length int) string
The `gengo.StringHexadecimal(length int)` func generates a random string with a given length containing all the predefined Hexadecimal characters.

This func uses the predefined constant string `Hexadecimal` as characters source reference.

```go 
    randomString := gengo.StringHexadecimal(33)
    fmt.Println(randomString)
```

### StringSymbols(length int) string
The `gengo.StringSymbols(length int)` func generates a random string with a given length containing all the predefined Hexadecimal characters.

This func uses the predefined constant string `Symbols` as characters source reference.

```go 
    randomString := gengo.StringSymbols(33)
    fmt.Println(randomString)
```