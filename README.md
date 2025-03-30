# gengo

[![build](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml/badge.svg)](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml)

gengo is a go Module containing utility functions to generate random data. 

```sh
go get -u github.com/FlavioCFOliveira/gengo
```

Examples

```go
package main

import (
	"fmt"
	"github.com/FlavioCFOliveira/gengo"
)

func main() {

	// Generate a random alfanumeric string with 16 chars
	randString := gengo.StringAlphanumeric(16)
	fmt.Println(randString)

}
```
