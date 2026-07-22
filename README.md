<p align="center">
  <img src="Assets/gengo-logo.png" alt="gengo logo" width="300" height="300" />
</p>

# gengo — Random Data Generator for Go

[![build](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml/badge.svg)](https://github.com/FlavioCFOliveira/gengo/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/FlavioCFOliveira/gengo.svg)](https://pkg.go.dev/github.com/FlavioCFOliveira/gengo)
[![Go Report Card](https://goreportcard.com/badge/github.com/FlavioCFOliveira/gengo)](https://goreportcard.com/report/github.com/FlavioCFOliveira/gengo)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.26-00ADD8?logo=go)](https://go.dev/doc/devel/release)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/FlavioCFOliveira/gengo)](https://github.com/FlavioCFOliveira/gengo/releases)

**gengo** is a zero-dependency Go library for generating random test data — strings, numbers, dates, booleans, and words — with a single function call. No configuration, no boilerplate. Just fast, reliable random data whenever you need it.

## Why gengo?

Generating realistic test data in Go usually means writing repetitive helper functions or pulling in heavy dependencies. gengo solves this with a clean, idiomatic API that covers the most common use cases out of the box:

- **One function call** to generate any random value
- **Zero external dependencies** — only the Go standard library
- **Blazing fast** — single-digit-nanosecond generation for primitive types on modern hardware
- **Safe by default** — built-in overflow protection, size limits, and uniform distribution

## Installation

```bash
go get -u github.com/FlavioCFOliveira/gengo
```

Requires Go 1.26 or later.

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/FlavioCFOliveira/gengo"
)

func main() {
    // Random alphanumeric string — handy for test fixtures and IDs
    fmt.Println(gengo.StringAlphanumeric(20))
    // Output: k9mP2vLxQr5tWnB8aJcT

    // Random integer in a custom range
    fmt.Println(gengo.IntBetween(1, 100))
    // Output: 42

    // Random date anywhere in history or the future
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
- [European Portuguese Words (pt-PT)](#european-portuguese-words-pt-pt)
- [Reproducible Output](#reproducible-output)
- [Character Sets](#character-sets)
- [Safety & Limits](#safety--limits)
- [Performance](#performance)

---

## String Generation

All string functions accept a `length` parameter — the number of bytes, capped at 1 MB. Pick a predefined type or pass your own character set. String generation is byte-oriented, so custom character sets should contain only single-byte (ASCII) characters; multibyte input (emoji, accents, CJK) produces invalid UTF-8.

### Custom and Variable-Length Strings

```go
// Generate from a custom character set
custom := gengo.String(10, "ABC123")
// Output: C1BA3C2A1B

// Generate a string with a random length between min and max
variable := gengo.StringBetween(5, 15, gengo.Alphanumeric)
```

### Predefined String Types

| Function | Output characters | Example |
|----------|-------------------|---------|
| `StringAllChars(n)` | Letters + numbers + symbols | `aB3$kP9@mQ` |
| `StringAlphanumeric(n)` | Letters + numbers | `xK8pM2vLqR` |
| `StringAlphabetic(n)` | Upper and lower case letters | `bGtRvNmKlP` |
| `StringAlphabeticUppercase(n)` | A–Z only | `XKMTPQRLVN` |
| `StringAlphabeticLowercase(n)` | a–z only | `mkptqrvbgs` |
| `StringNumeric(n)` | Digits 0–9 only | `8472936105` |
| `StringHexadecimal(n)` | Hex digits 0–9, A–F | `A7F3B9E2D1` |
| `StringSymbols(n)` | Special characters only | `@#$%&*()-_` |

### Security: Not for Passwords or Secrets

> [!WARNING]
> **gengo is not cryptographically secure.** It is built on
> [`math/rand/v2`](https://pkg.go.dev/math/rand/v2), whose output is fast but
> predictable and therefore unsuitable for any security-sensitive purpose.
> **Never** use gengo to generate passwords, API tokens, session IDs,
> encryption keys, salts, or any other secret — the results are guessable. For
> secrets, use Go's [`crypto/rand`](https://pkg.go.dev/crypto/rand) package
> instead.

See [SECURITY.md](SECURITY.md) for the full security policy and how to report a vulnerability.

### Common Use Cases

```go
// Short, human-readable identifier for test fixtures
id := gengo.StringAlphanumeric(8)
// Output: xK8pM2vL

// Random label for seeding a database column
label := gengo.StringAlphabetic(12)
// Output: bGtRvNmKlPqX

// Random hex value for non-secret test data (e.g. a fake colour code)
hex := gengo.StringHexadecimal(6)
// Output: A7F3B9
```

---

## Number Generation

### Integers

Every integer type comes in two flavours:
- `Type()` — full range of that type
- `TypeBetween(min, max)` — your custom range (min and max are swapped automatically if reversed)

| Type | Function | Full Range | Range example |
|------|----------|------------|---------------|
| `int8` | `Int8()` | −128 to 127 | `Int8Between(-50, 50)` |
| `int16` | `Int16()` | −32 768 to 32 767 | `Int16Between(100, 1000)` |
| `int32` | `Int32()` | ±2 billion | `Int32Between(0, 999999)` |
| `int64` | `Int64()` | ±9 quintillion | `Int64Between(0, time.Now().Unix())` |
| `int` | `Int()` | Platform dependent | `IntBetween(1, 100)` |
| `uint8` / `byte` | `Uint8()` / `Byte()` | 0 to 255 | `Uint8Between(0, 255)` |
| `uint16` | `Uint16()` | 0 to 65 535 | `Uint16Between(1000, 5000)` |
| `uint32` | `Uint32()` | 0 to 4 billion | `Uint32Between(0, 100000)` |
| `uint64` | `Uint64()` | 0 to 18 quintillion | `Uint64Between(0, math.MaxUint32)` |

```go
// Random age between 18 and 100
age := gengo.IntBetween(18, 100)

// Random percentage (0–100)
percent := gengo.Uint8Between(0, 100)

// Random unprivileged port number
port := gengo.Uint16Between(1024, 65535)

// Random byte for binary data
b := gengo.Byte()
```

### Floating Point

| Function | Range | Description |
|----------|-------|-------------|
| `Float32()` | (0, MaxFloat32) | Positive values only — use `Float32Between` for negatives or a custom range |
| `Float32Between(min, max)` | Custom | Bounded float32 |
| `Float64()` | (0, MaxFloat64) | Positive values only — use `Float64Between` for negatives or a custom range |
| `Float64Between(min, max)` | Custom | Bounded float64 |

```go
// Random price for e-commerce testing
price := gengo.Float64Between(0.99, 999.99)

// Random percentage with decimal precision
percentage := gengo.Float32Between(0.0, 100.0)

// Random geographic coordinates
lat := gengo.Float64Between(-90.0, 90.0)
lng := gengo.Float64Between(-180.0, 180.0)
```

### Complex Numbers

| Function | Description |
|----------|-------------|
| `Complex64()` | Random complex64 with random real and imaginary parts |
| `Complex64Between(rMin, rMax, iMin, iMax)` | Custom ranges for both parts |
| `Complex128()` | Random complex128 with random real and imaginary parts |
| `Complex128Between(rMin, rMax, iMin, iMax)` | Custom ranges for both parts |

```go
// Random complex number
c := gengo.Complex64()
// Output: (12.34+56.78i)

// Bounded unit circle range
c = gengo.Complex64Between(-1.0, 1.0, -1.0, 1.0)
// Output: (0.45-0.92i)
```

---

## Date & Time

Generate random timestamps for testing date-sensitive logic, filling databases with realistic data, or simulating time-based events.

| Function | Description | Range |
|----------|-------------|-------|
| `Date()` | Random date and time | Year 1 to 9999 |
| `UnixDate()` | Random Unix timestamp | 1970 to 2038 (32-bit range) |
| `DateBetween(start, end)` | Random date within a specific range | Custom |

> All date functions work at whole-second granularity and return UTC times with a zero sub-second component. For `DateBetween`, the bounds' sub-second parts are ignored; as a special case, when `start` and `end` fall within the same second, `start` is returned unchanged.

```go
// Random historical or future date
d := gengo.Date()
// Output: 1847-03-15 08:23:17 UTC

// Random date in the Unix epoch range
d = gengo.UnixDate()
// Output: 1995-11-30 14:45:02 UTC

// Random date within a custom range
start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
d = gengo.DateBetween(start, end)
// Output: 2022-07-14 09:33:45 UTC

// Random birthdate for someone aged 18 to 65
now := time.Now()
minAge := now.AddDate(-65, 0, 0)
maxAge := now.AddDate(-18, 0, 0)
birthdate := gengo.DateBetween(minAge, maxAge)
```

---

## Booleans

```go
// 50/50 coin flip
if gengo.Bool() {
    fmt.Println("Heads")
} else {
    fmt.Println("Tails")
}

// Randomly toggle a feature flag in tests
if gengo.Bool() {
    enableNewFeature()
}
```

---

## Words & Text

Generate random words for populating text fields, testing search functionality, or seeding databases with human-readable content.

### Single Word

```go
// Random word between 2 and 30 lowercase letters
word := gengo.Word()
// Output: galhjhyqxfmvwqxqywuigohxhvclk
```

### Words by Length Category

| Category | Length | Example output |
|----------|--------|----------------|
| `SmallLengthWord` | 1–4 chars | `cat`, `dog`, `sun` |
| `MediumLengthWords` | 5–8 chars | `house`, `garden` |
| `BigLengthWords` | 9–30 chars | `wonderful`, `beautiful` |

```go
small  := gengo.WordByLengthType(gengo.SmallLengthWord)
medium := gengo.WordByLengthType(gengo.MediumLengthWords)
big    := gengo.WordByLengthType(gengo.BigLengthWords)
```

### Multiple Words

```go
// Generate a slice of random words
words := gengo.Words(5)
// Output: [ok then computer building wonderful]

// Natural-language-like distribution (from WordLengthRatio, 13/26/11 of 50):
// ~26% small words (1–4 chars)
// ~52% medium words (5–8 chars)
// ~22% big words (9–30 chars)
```

---

## European Portuguese Words (pt-PT)

Beyond the ASCII words above, gengo generates European-Portuguese (pt-PT) words with full grammatical inflection — handy for seeding realistic Portuguese text, testing search and collation, or filling forms.

Words fall into two families:

- **Open classes** (noun, adjective, verb, adverb) are **generated pseudo-words**: morphologically well-formed, phonotactically valid, and correctly accented, but not guaranteed to be real dictionary entries.
- **Closed classes** (article, pronoun, numeral, preposition, conjunction, interjection) are **real, curated words** drawn uniformly from hand-checked lists.

Every result is always **lowercase** and valid UTF-8. No function ever returns an error or panics.

### Inflection Options

Inflection is controlled by six enums. Each has an `Any…` zero value that means "choose a valid value at random"; the length category adds its own `AnyLengthWord` zero value that means "any length".

| Enum | Zero value (random) | Concrete values |
|------|---------------------|-----------------|
| `Gender` | `AnyGender` | `Masculine`, `Feminine` |
| `Number` | `AnyNumber` | `Singular`, `Plural` |
| `Degree` | `AnyDegree` | `Positive`, `Superlative` (synthetic absolute superlative, -íssimo) |
| `Mood` | `AnyMood` | `Indicative`, `Subjunctive`, `ImperativeAffirmative`, `ImperativeNegative`, `Infinitive`, `Gerund`, `Participle` |
| `Tense` | `AnyTense` | `Present`, `Imperfect`, `Preterite`, `Future`, `Conditional` |
| `Person` | `AnyPerson` | `First`, `Second`, `Third` |
| `LengthTypeWords` | `AnyLengthWord` | `SmallLengthWord`, `MediumLengthWords`, `BigLengthWords` |

A grammatically impossible request never panics: it normalizes to the nearest valid form (for example, a second-person-plural verb becomes third-person plural, matching current pt-PT usage), and a length category too small for the chosen word is widened **upward**, never downward.

### Open Classes (generated pseudo-words)

```go
// Nouns — random gender, number, and length, or fully specified
n := gengo.NounPT()
n = gengo.NounPTOf(gengo.Feminine, gengo.Plural, gengo.MediumLengthWords)

// Adjectives — including the synthetic superlative (-íssimo)
a := gengo.AdjectivePT()
a = gengo.AdjectivePTOf(gengo.Masculine, gengo.Singular, gengo.Superlative, gengo.AnyLengthWord)

// Verbs — any mood, tense, person, and number
v := gengo.VerbPT()
v = gengo.VerbPTOf(gengo.Indicative, gengo.Present, gengo.First, gengo.Singular, gengo.AnyLengthWord)

// Adverbs — productive -mente forms (invariable)
adv := gengo.AdverbPT()
adv = gengo.AdverbPTByLengthType(gengo.MediumLengthWords)
```

### Closed Classes (real curated words)

Closed-class selectors return genuine pt-PT words in European spelling. They take **no length parameter**, because the length concept does not apply.

```go
art  := gengo.ArticlePT()      // e.g. "a", "do", "pelas"
pron := gengo.PronounPT()      // e.g. "eu", "aquilo", "cujo"
num  := gengo.NumeralPT()      // e.g. "três", "primeiro", "meio"
prep := gengo.PrepositionPT()  // e.g. "de", "para", "perante"
conj := gengo.ConjunctionPT()  // e.g. "e", "mas", "porque"
intj := gengo.InterjectionPT() // e.g. "olá", "caramba", "psiu"
```

### Any Word

`WordPT` returns a random word of a random class, following a natural distribution across word classes. `WordsPT` returns a slice of them.

```go
w := gengo.WordPT()                                // any class, any length
w = gengo.WordPTByLengthType(gengo.BigLengthWords) // length applies to open classes only
list := gengo.WordsPT(5)                           // slice of 5 random pt-PT words
```

### Reproducible pt-PT Output

Every function above is also a `*Generator` method, so a seeded generator produces a reproducible sequence — ideal for golden tests:

```go
g := gengo.New(7)
fmt.Println(g.NounPTOf(gengo.Feminine, gengo.Plural, gengo.MediumLengthWords))
fmt.Println(g.VerbPTOf(gengo.Indicative, gengo.Present, gengo.First, gengo.Singular, gengo.AnyLengthWord))
// Same seed → same words, every run.
```

---

## Reproducible Output

By default, gengo's package-level functions draw from Go's global `math/rand/v2` source, which is automatically seeded — great for variety, but the output cannot be reproduced. When you need a deterministic sequence (for example, to reproduce a failing test), create a seeded `Generator`:

```go
// Same seed → same sequence, every run
g := gengo.New(42)
fmt.Println(g.IntBetween(1, 100))            // deterministic for seed 42
fmt.Println(g.String(8, gengo.Alphanumeric)) // deterministic for seed 42

// A second generator with the same seed reproduces the sequence exactly
g2 := gengo.New(42)
// g2 yields the identical sequence as g
```

A `Generator` exposes a method for every package-level function (`g.Int8()`, `g.Float64Between(...)`, `g.Date()`, `g.Word()`, and so on), so it is a drop-in replacement when you need reproducibility. For full control over the underlying algorithm, build one from any `math/rand/v2` source:

```go
import "math/rand/v2"

g := gengo.NewSource(rand.NewPCG(1, 2))
```

> A `Generator` is **not** safe for concurrent use by multiple goroutines (it wraps a `math/rand/v2.Rand`). Use one generator per goroutine, or the package-level functions — which are safe for concurrent use — when you don't need reproducibility.

---

## Character Sets

Use these built-in constants to build your own character sets or combine them for custom string generation.

| Constant | Characters | Size |
|----------|------------|------|
| `AllChars` | `a-zA-Z0-9!@#$%...` | 87 |
| `Alphanumeric` | `a-zA-Z0-9` | 62 |
| `Alphabetic` | `a-zA-Z` | 52 |
| `AlphabeticUppercase` | `A-Z` | 26 |
| `AlphabeticLowercase` | `a-z` | 26 |
| `Numeric` | `0-9` | 10 |
| `Hexadecimal` | `0-9A-F` | 16 |
| `Symbols` | `!@#$%...` | 25 |

```go
// Letters and underscores — safe for usernames
usernameSafe := gengo.AlphabeticLowercase + "_"
name := gengo.String(12, usernameSafe)

// URL-safe slugs
urlSafe := gengo.Alphanumeric + "-_"
slug := gengo.String(20, urlSafe)

// Fun custom charset (ASCII) — e.g. DNA bases
dna := "ACGT"
sequence := gengo.String(20, dna)
// Output: GATTACAGATTACAGGCTAG
```

---

## Safety & Limits

gengo is designed to be safe to use in production code and CI pipelines without any extra configuration:

- **String length cap**: Maximum 1 MB (`MaxStringLength`) prevents accidental memory exhaustion
- **Automatic range correction**: All `Between` functions swap min and max if they are reversed
- **Overflow protection**: Integer calculations use wider intermediate types to prevent wrap-around
- **Uniform distribution**: Backed by `math/rand/v2` for statistically unbiased results

---

## Performance

Benchmarks run on an AMD Ryzen 9 5900HX (Go 1.26.2, `-benchtime 5s`). See `make bench` to reproduce on your machine.

| Operation | ns/op |
|-----------|-------|
| `IntBetween` | ~7.8 |
| `Float64` | ~5.8 |
| `String` (8 chars) | ~27.8 |
| `Date` | ~7.4 |
| `Uint64` | ~4.8 |

For the full benchmark report and detailed test results, see [BENCHMARKS.md](BENCHMARKS.md) and [TEST_REPORT.md](TEST_REPORT.md).

---

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for the build, test, and lint workflow, the coding conventions, and the pull-request process. For security issues, follow [SECURITY.md](SECURITY.md).

---

## License

Released under the [MIT License](LICENSE).
