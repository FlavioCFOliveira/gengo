# gengo

[![build](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml/badge.svg)](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/FlavioCFOliveira/gengo.svg)](https://pkg.go.dev/github.com/FlavioCFOliveira/gengo)

A high-performance Go module for generating random data. Built with production environments in mind, offering excellent performance and minimal resource usage.

## Installation

```bash
go get -u github.com/FlavioCFOliveira/gengo
```

Requirements: Go 1.22+ (uses `math/rand/v2` with automatic seed initialization)

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/FlavioCFOliveira/gengo"
)

func main() {
    // Generate a random alphanumeric string
    fmt.Println(gengo.StringAlphanumeric(20))
    // Output: k9mP2vLxQr5tWnB8aJc

    // Generate a random number
    fmt.Println(gengo.IntBetween(1, 100))
    // Output: 42

    // Generate a random date
    fmt.Println(gengo.Date())
    // Output: 1847-03-15 08:23:17 UTC
}
```

---

## Table of Contents

- [String Generation](#string-generation)
- [Number Generation](#number-generation)
  - [Integers](#integers)
  - [Floating Point](#floating-point)
  - [Complex Numbers](#complex-numbers)
- [Date & Time](#date--time)
- [Booleans](#booleans)
- [Words & Text](#words--text)
- [Character Sets](#character-sets)

---

## String Generation

All string functions accept a `length` parameter (capped at 1MB for safety).

### Basic String Generation

```go
// Generate from custom character set
custom := gengo.String(10, "ABC123")
// Output: C1BA3C2A1B

// Generate with random length between min and max
variable := gengo.StringBetween(5, 15, gengo.Alphanumeric)
```

### Predefined String Types

| Function | Description | Example Output |
|----------|-------------|----------------|
| `StringAllChars(n)` | Letters + numbers + symbols | `aB3$kP9@mQ` |
| `StringAlphanumeric(n)` | Letters + numbers | `xK8pM2vLqR` |
| `StringAlphabetic(n)` | Upper + lower case | `bGtRvNmKlP` |
| `StringAlphabeticUppercase(n)` | A-Z only | `XKMTPQRLVN` |
| `StringAlphabeticLowercase(n)` | a-z only | `mkptqrvbgs` |
| `StringNumeric(n)` | 0-9 only | `8472936105` |
| `StringHexadecimal(n)` | 0-9, A-F | `A7F3B9E2D1` |
| `StringSymbols(n)` | Special chars only | `@#$%&*()-_` |

**Examples:**

```go
// Password-like string
password := gengo.StringAllChars(16)
// Output: k9@mP2#vLx$Qr5tW

// Token generation
token := gengo.StringHexadecimal(32)
// Output: A7F3B9E2D108C4E5A67F9B2C1D0E8F5A

// Readable identifier
id := gengo.StringAlphanumeric(8)
// Output: xK8pM2vL
```

---

## Number Generation

### Integers

Each integer type has two variants:
- `Type()` - full range of the type
- `TypeBetween(min, max)` - custom range (automatically swaps if min > max)

| Type | Function | Full Range | Example |
|------|----------|------------|---------|
| `int8` | `Int8()` | -128 to 127 | `Int8Between(-50, 50)` |
| `int16` | `Int16()` | -32768 to 32767 | `Int16Between(100, 1000)` |
| `int32` | `Int32()` | ~±2 billion | `Int32Between(0, 999999)` |
| `int64` | `Int64()` | ~±9 quintillion | `Int64Between(0, time.Now().Unix())` |
| `int` | `Int()` | platform dependent | `IntBetween(1, 100)` |
| `uint8` / `byte` | `UInt8()` / `Byte()` | 0 to 255 | `UInt8Between(0, 255)` |
| `uint16` | `UInt16()` | 0 to 65535 | `UInt16Between(1000, 5000)` |
| `uint32` | `UInt32()` | 0 to ~4 billion | `UInt32Between(0, 100000)` |
| `uint64` | `UInt64()` | 0 to ~18 quintillion | `UInt64Between(0, math.MaxUint32)` |

**Examples:**

```go
// Random age between 18 and 100
age := gengo.IntBetween(18, 100)

// Random percentage (0-100)
percent := gengo.UInt8Between(0, 100)

// Random port number (1024-65535)
port := gengo.UInt16Between(1024, 65535)

// Random byte for binary data
b := gengo.Byte()
```

### Floating Point

| Function | Range | Description |
|----------|-------|-------------|
| `Float32()` | Smallest to MaxFloat32 | Full range float32 |
| `Float32Between(min, max)` | custom | Bounded range |
| `Float64()` | Smallest to MaxFloat64 | Full range float64 |
| `Float64Between(min, max)` | custom | Bounded range |

**Examples:**

```go
// Random price between 0.99 and 999.99
price := gengo.Float64Between(0.99, 999.99)

// Random percentage with decimals
percentage := gengo.Float32Between(0.0, 100.0)

// Random coordinate
lat := gengo.Float64Between(-90.0, 90.0)
lng := gengo.Float64Between(-180.0, 180.0)
```

### Complex Numbers

| Function | Description |
|----------|-------------|
| `Complex64()` | Random complex64 with random real/imag parts |
| `Complex64Between(rMin, rMax, iMin, iMax)` | Custom ranges |
| `Complex128()` | Random complex128 with random real/imag parts |
| `Complex128Between(rMin, rMax, iMin, iMax)` | Custom ranges |

**Examples:**

```go
// Random complex number
c := gengo.Complex64()
// Output: (12.34+56.78i)

// Custom range
c = gengo.Complex64Between(-1.0, 1.0, -1.0, 1.0)
// Output: (0.45-0.92i)
```

---

## Date & Time

| Function | Description | Range |
|----------|-------------|-------|
| `Date()` | Random date/time | Year 1 to 9999 |
| `UnixDate()` | Random Unix timestamp | 1970 to 2038 (32-bit range) |
| `DateBetween(start, end)` | Date within specific range | Custom |

**Examples:**

```go
// Random date (year 1 to 9999)
d := gengo.Date()
// Output: 1847-03-15 08:23:17 UTC

// Random Unix date (1970-2038)
d = gengo.UnixDate()
// Output: 1995-11-30 14:45:02 UTC

// Random date in specific range
start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
d = gengo.DateBetween(start, end)
// Output: 2022-07-14 09:33:45 UTC

// Random birthdate (18-65 years ago)
now := time.Now()
minAge := now.AddDate(-65, 0, 0)
maxAge := now.AddDate(-18, 0, 0)
birthdate := gengo.DateBetween(minAge, maxAge)
```

---

## Booleans

```go
// Simple coin flip
if gengo.Bool() {
    fmt.Println("Heads")
} else {
    fmt.Println("Tails")
}

// Randomly decide feature flag
if gengo.Bool() {
    enableNewFeature()
}
```

---

## Words & Text

### Single Word

```go
// Random word (2-30 lowercase letters)
word := gengo.Word()
// Output: galhjhyqxfmvwqxqywuigohxhvclk
```

### Words by Length Category

| Category | Length | Example |
|----------|--------|---------|
| `SmallLengthWord` | 1-4 chars | `cat`, `dog`, `sun` |
| `MediumLengthWords` | 5-8 chars | `house`, `garden` |
| `BigLengthWords` | 9+ chars | `wonderful`, `beautiful` |

```go
// Generate word by category
small := gengo.WordByLengthType(gengo.SmallLengthWord)
medium := gengo.WordByLengthType(gengo.MediumLengthWords)
big := gengo.WordByLengthType(gengo.BigLengthWords)
```

### Multiple Words

```go
// Generate slice of random words
words := gengo.Words(5)
// Output: [ok then computer building wonderful]

// Distribution is weighted:
// - ~22% small words (1-4 chars)
// - ~56% medium words (5-8 chars)
// - ~22% big words (9+ chars)
```

---

## Character Sets

Use these constants for custom string generation:

| Constant | Characters | Length |
|----------|------------|--------|
| `AllChars` | `a-zA-Z0-9!@#$%...` | 94 |
| `Alphanumeric` | `a-zA-Z0-9` | 62 |
| `Alphabetic` | `a-zA-Z` | 52 |
| `AlphabeticUppercase` | `A-Z` | 26 |
| `AlphabeticLowercase` | `a-z` | 26 |
| `Numeric` | `0-9` | 10 |
| `Hexadecimal` | `0-9A-F` | 16 |
| `Symbols` | `!@#$%...` | 32 |

**Custom combinations:**

```go
// Letters + underscores only
usernameSafe := gengo.AlphabeticLowercase + "_"
name := gengo.String(12, usernameSafe)

// URL-safe characters
urlSafe := gengo.Alphanumeric + "-_"
slug := gengo.String(20, urlSafe)

// Custom charset
emojiLike := "😀😎🎉🚀💡"
fun := gengo.String(5, emojiLike)
```

---

## Safety & Limits

- **String length**: Maximum 1MB (`MaxStringLength`) to prevent memory exhaustion
- **Range handling**: All `Between` functions automatically swap min/max if reversed
- **Overflow protection**: Integer calculations use wider types to prevent overflow
- **Uniform distribution**: Uses `math/rand/v2` for statistically unbiased randomness

---

## Performance

Benchmarks on Apple M4:

| Operation | ns/op |
|-----------|-------|
| IntBetween | ~4.5 |
| Float64 | ~4.4 |
| String (100 chars) | ~42 |
| Date | ~4.6 |
| UInt64 | ~3.9 |

See `make bench` to run benchmarks locally.

---

## License

This project is licensed under the terms specified in the repository.
