# Test Report - gengo

**Date:** 2026-05-29
**Package:** github.com/FlavioCFOliveira/gengo
**Platform:** linux/amd64 (AMD Ryzen 9 5900HX, 16 threads) · Go go1.26.2

---

## Executive Summary

| Metric | Value |
|--------|-------|
| **Unit Tests** | 45/45 passed (100%) |
| **Benchmarks** | 36 executed |
| **Total Test Time** | ~0.01s |
| **Total Benchmark Time** | ~227s |

---

## Unit Tests

All 45 unit tests passed successfully.

### Detailed Results

| Test | Status | Duration |
|------|--------|----------|
| `TestBool` | PASS | 0.00s |
| `TestDate` | PASS | 0.00s |
| `TestUnixDate` | PASS | 0.00s |
| `TestDateBetween` | PASS | 0.00s |
| `TestInt8` | PASS | 0.00s |
| `TestInt8Between` | PASS | 0.00s |
| `TestInt16` | PASS | 0.00s |
| `TestInt16Between` | PASS | 0.00s |
| `TestInt32` | PASS | 0.00s |
| `TestInt32Between` | PASS | 0.00s |
| `TestInt` | PASS | 0.00s |
| `TestIntBetween` | PASS | 0.00s |
| `TestInt64` | PASS | 0.00s |
| `TestInt64Between` | PASS | 0.00s |
| `TestUint8` | PASS | 0.00s |
| `TestUint8Between` | PASS | 0.00s |
| `TestByte` | PASS | 0.00s |
| `TestUint16` | PASS | 0.00s |
| `TestUint16Between` | PASS | 0.00s |
| `TestUint32` | PASS | 0.00s |
| `TestUint32Between` | PASS | 0.00s |
| `TestUint64` | PASS | 0.00s |
| `TestUint64Between` | PASS | 0.00s |
| `TestFloat32` | PASS | 0.00s |
| `TestFloat32Between` | PASS | 0.00s |
| `TestFloat64` | PASS | 0.00s |
| `TestFloat64Between` | PASS | 0.00s |
| `TestComplex64` | PASS | 0.00s |
| `TestComplex64Between` | PASS | 0.00s |
| `TestComplex128` | PASS | 0.00s |
| `TestComplex128Between` | PASS | 0.00s |
| `TestString` | PASS | 0.00s |
| `TestStringByteOriented` | PASS | 0.00s |
| `TestStringAllChars` | PASS | 0.00s |
| `TestStringAlphanumeric` | PASS | 0.00s |
| `TestStringAlphabetic` | PASS | 0.00s |
| `TestStringAlphabeticUppercase` | PASS | 0.00s |
| `TestStringAlphabeticLowercase` | PASS | 0.00s |
| `TestStringNumeric` | PASS | 0.00s |
| `TestStringHexadecimal` | PASS | 0.00s |
| `TestStringSymbols` | PASS | 0.00s |
| `TestStringBetween` | PASS | 0.00s |
| `TestWord` | PASS | 0.00s |
| `TestWordByLengthType` | PASS | 0.00s |
| `TestWords` | PASS | 0.00s |

### Subtests
- `TestWordByLengthType/SmallLengthWord` - PASS
- `TestWordByLengthType/MediumLengthWords` - PASS
- `TestWordByLengthType/BigLengthWords` - PASS
- `TestWordByLengthType/Default` - PASS
- `TestWords/PositiveLength` - PASS
- `TestWords/ZeroLength` - PASS
- `TestWords/NegativeLength` - PASS

---

## Performance Tests (Benchmarks)

### Summary by Category

#### Booleans and Dates
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkBool` | 1,000,000,000 | 4.719 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkDate` | 821,951,491 | 7.377 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUnixDate` | 813,031,130 | 7.401 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkDateBetween` | 564,340,304 | 10.66 ns/op | 0 B/op, 0 allocs/op |

#### Signed Integers (Int8/16/32/64)
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkInt8` | 1,000,000,000 | 4.763 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt8Between` | 801,699,829 | 7.424 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt16` | 1,000,000,000 | 4.774 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt16Between` | 828,890,362 | 7.398 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt32` | 1,000,000,000 | 4.753 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt32Between` | 810,285,100 | 7.469 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt` | 1,000,000,000 | 4.762 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkIntBetween` | 767,547,925 | 7.812 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt64` | 1,000,000,000 | 4.770 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt64Between` | 761,782,174 | 7.872 ns/op | 0 B/op, 0 allocs/op |

#### Unsigned Integers (Uint8/16/32/64) and Byte
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkUint8` | 1,000,000,000 | 4.780 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint8Between` | 801,432,346 | 7.463 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkByte` | 1,000,000,000 | 4.768 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint16` | 1,000,000,000 | 4.767 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint16Between` | 811,150,370 | 7.515 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint32` | 1,000,000,000 | 4.796 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint32Between` | 803,114,101 | 7.489 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint64` | 1,000,000,000 | 4.786 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint64Between` | 757,757,079 | 7.942 ns/op | 0 B/op, 0 allocs/op |

#### Floating Point and Complex Numbers
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkFloat32` | 1,000,000,000 | 5.804 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkFloat32Between` | 885,918,817 | 6.791 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkFloat64` | 1,000,000,000 | 5.767 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkFloat64Between` | 935,214,699 | 6.443 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex64` | 586,800,524 | 10.25 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex64Between` | 421,643,383 | 14.24 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex128` | 575,015,457 | 10.44 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex128Between` | 438,102,889 | 13.71 ns/op | 0 B/op, 0 allocs/op |

#### Strings and Words
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkString` | 217,566,918 | 27.76 ns/op | 8 B/op, 1 allocs/op |
| `BenchmarkStringNumeric` | 92,571,711 | 62.65 ns/op | 8 B/op, 1 allocs/op |
| `BenchmarkWord` | 68,938,351 | 86.55 ns/op | 18 B/op, 1 allocs/op |
| `BenchmarkWordByLengthType` | 100,000,000 | 50.12 ns/op | 7 B/op, 1 allocs/op |
| `BenchmarkWords` | 7,636,587 | 788.4 ns/op | 258 B/op, 10 allocs/op |

---

## Performance Analysis

### Highlights

1. **Zero Allocations in Primitive Types**: All number, boolean, and date generation functions perform no memory allocations (0 allocs/op), adhering to the library's performance philosophy.

2. **Fast Execution**: Numeric and boolean functions execute in ~4.7–8 nanoseconds per operation, achieving hundreds of millions to over a billion operations in 5 seconds of benchmarking.

3. **Fastest Functions** (top 5):
   - `BenchmarkBool`: 4.719 ns/op
   - `BenchmarkInt32`: 4.753 ns/op
   - `BenchmarkInt`: 4.762 ns/op
   - `BenchmarkInt8`: 4.763 ns/op
   - `BenchmarkUint16`: 4.767 ns/op

4. **Complex Numbers Slower**: `BenchmarkComplex64` and `BenchmarkComplex128` are naturally slower (~10.3 ns/op), and their `Between` variants slower still (~14 ns/op), as they generate two parts (real and imaginary).

5. **Strings with Controlled Allocation**: String and word functions (`BenchmarkString`, `BenchmarkStringNumeric`, `BenchmarkWord`, `BenchmarkWordByLengthType`, `BenchmarkWords`) are the only ones that allocate (1 alloc/op for single results; `BenchmarkWords` allocates once per generated word), as they must allocate memory for the returned string(s).

### Aggregated Metrics

| Category | Best Time | Worst Time | Approximate Average |
|----------|-----------|------------|---------------------|
| Booleans & Dates | 4.719 ns/op | 10.66 ns/op | ~7.5 ns/op |
| Integers | 4.753 ns/op | 7.942 ns/op | ~6.0 ns/op |
| Floats & Complex | 5.767 ns/op | 14.24 ns/op | ~8.7 ns/op |
| Strings | 27.76 ns/op | 62.65 ns/op | ~45 ns/op |

---

## Evidence (Raw Output)

### Unit Test Output
```
=== RUN   TestBool
--- PASS: TestBool (0.00s)
=== RUN   TestDate
--- PASS: TestDate (0.00s)
=== RUN   TestUnixDate
--- PASS: TestUnixDate (0.00s)
=== RUN   TestDateBetween
--- PASS: TestDateBetween (0.00s)
=== RUN   TestInt8
--- PASS: TestInt8 (0.00s)
=== RUN   TestInt8Between
--- PASS: TestInt8Between (0.00s)
=== RUN   TestInt16
--- PASS: TestInt16 (0.00s)
=== RUN   TestInt16Between
--- PASS: TestInt16Between (0.00s)
=== RUN   TestInt32
--- PASS: TestInt32 (0.00s)
=== RUN   TestInt32Between
--- PASS: TestInt32Between (0.00s)
=== RUN   TestInt
--- PASS: TestInt (0.00s)
=== RUN   TestIntBetween
--- PASS: TestIntBetween (0.00s)
=== RUN   TestInt64
--- PASS: TestInt64 (0.00s)
=== RUN   TestInt64Between
--- PASS: TestInt64Between (0.00s)
=== RUN   TestUint8
--- PASS: TestUint8 (0.00s)
=== RUN   TestUint8Between
--- PASS: TestUint8Between (0.00s)
=== RUN   TestByte
--- PASS: TestByte (0.00s)
=== RUN   TestUint16
--- PASS: TestUint16 (0.00s)
=== RUN   TestUint16Between
--- PASS: TestUint16Between (0.00s)
=== RUN   TestUint32
--- PASS: TestUint32 (0.00s)
=== RUN   TestUint32Between
--- PASS: TestUint32Between (0.00s)
=== RUN   TestUint64
--- PASS: TestUint64 (0.00s)
=== RUN   TestUint64Between
--- PASS: TestUint64Between (0.00s)
=== RUN   TestFloat32
--- PASS: TestFloat32 (0.00s)
=== RUN   TestFloat32Between
--- PASS: TestFloat32Between (0.00s)
=== RUN   TestFloat64
--- PASS: TestFloat64 (0.00s)
=== RUN   TestFloat64Between
--- PASS: TestFloat64Between (0.00s)
=== RUN   TestComplex64
--- PASS: TestComplex64 (0.00s)
=== RUN   TestComplex64Between
--- PASS: TestComplex64Between (0.00s)
=== RUN   TestComplex128
--- PASS: TestComplex128 (0.00s)
=== RUN   TestComplex128Between
--- PASS: TestComplex128Between (0.00s)
=== RUN   TestString
--- PASS: TestString (0.00s)
=== RUN   TestStringByteOriented
--- PASS: TestStringByteOriented (0.00s)
=== RUN   TestStringAllChars
--- PASS: TestStringAllChars (0.00s)
=== RUN   TestStringAlphanumeric
--- PASS: TestStringAlphanumeric (0.00s)
=== RUN   TestStringAlphabetic
--- PASS: TestStringAlphabetic (0.00s)
=== RUN   TestStringAlphabeticUppercase
--- PASS: TestStringAlphabeticUppercase (0.00s)
=== RUN   TestStringAlphabeticLowercase
--- PASS: TestStringAlphabeticLowercase (0.00s)
=== RUN   TestStringNumeric
--- PASS: TestStringNumeric (0.00s)
=== RUN   TestStringHexadecimal
--- PASS: TestStringHexadecimal (0.00s)
=== RUN   TestStringSymbols
--- PASS: TestStringSymbols (0.00s)
=== RUN   TestStringBetween
--- PASS: TestStringBetween (0.00s)
=== RUN   TestWord
--- PASS: TestWord (0.00s)
=== RUN   TestWordByLengthType
=== RUN   TestWordByLengthType/SmallLengthWord
=== RUN   TestWordByLengthType/MediumLengthWords
=== RUN   TestWordByLengthType/BigLengthWords
=== RUN   TestWordByLengthType/Default
--- PASS: TestWordByLengthType (0.00s)
    --- PASS: TestWordByLengthType/SmallLengthWord (0.00s)
    --- PASS: TestWordByLengthType/MediumLengthWords (0.00s)
    --- PASS: TestWordByLengthType/BigLengthWords (0.00s)
    --- PASS: TestWordByLengthType/Default (0.00s)
=== RUN   TestWords
=== RUN   TestWords/PositiveLength
=== RUN   TestWords/ZeroLength
=== RUN   TestWords/NegativeLength
--- PASS: TestWords (0.00s)
    --- PASS: TestWords/PositiveLength (0.00s)
    --- PASS: TestWords/ZeroLength (0.00s)
    --- PASS: TestWords/NegativeLength (0.00s)
PASS
ok  	github.com/FlavioCFOliveira/gengo	0.014s
```

### Benchmark Output
```
goos: linux
goarch: amd64
pkg: github.com/FlavioCFOliveira/gengo
cpu: AMD Ryzen 9 5900HX with Radeon Graphics        
BenchmarkBool-16                 	1000000000	         4.719 ns/op	       0 B/op	       0 allocs/op
BenchmarkDate-16                 	821951491	         7.377 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnixDate-16             	813031130	         7.401 ns/op	       0 B/op	       0 allocs/op
BenchmarkDateBetween-16          	564340304	        10.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt8-16                 	1000000000	         4.763 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt8Between-16          	801699829	         7.424 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt16-16                	1000000000	         4.774 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt16Between-16         	828890362	         7.398 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt32-16                	1000000000	         4.753 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt32Between-16         	810285100	         7.469 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt-16                  	1000000000	         4.762 ns/op	       0 B/op	       0 allocs/op
BenchmarkIntBetween-16           	767547925	         7.812 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt64-16                	1000000000	         4.770 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt64Between-16         	761782174	         7.872 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint8-16                	1000000000	         4.780 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint8Between-16         	801432346	         7.463 ns/op	       0 B/op	       0 allocs/op
BenchmarkByte-16                 	1000000000	         4.768 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint16-16               	1000000000	         4.767 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint16Between-16        	811150370	         7.515 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint32-16               	1000000000	         4.796 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint32Between-16        	803114101	         7.489 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint64-16               	1000000000	         4.786 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint64Between-16        	757757079	         7.942 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat32-16              	1000000000	         5.804 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat32Between-16       	885918817	         6.791 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat64-16              	1000000000	         5.767 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat64Between-16       	935214699	         6.443 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex64-16            	586800524	        10.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex64Between-16     	421643383	        14.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex128-16           	575015457	        10.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex128Between-16    	438102889	        13.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkString-16               	217566918	        27.76 ns/op	       8 B/op	       1 allocs/op
BenchmarkStringNumeric-16        	92571711	        62.65 ns/op	       8 B/op	       1 allocs/op
BenchmarkWord-16                 	68938351	        86.55 ns/op	      18 B/op	       1 allocs/op
BenchmarkWordByLengthType-16     	100000000	        50.12 ns/op	       7 B/op	       1 allocs/op
BenchmarkWords-16                	 7636587	       788.4 ns/op	     258 B/op	      10 allocs/op
PASS
ok  	github.com/FlavioCFOliveira/gengo	227.251s
```

---

## Conclusion

The **gengo** library demonstrates excellent test coverage and optimized performance:

- **100% tests passing** (45/45)
- **Zero memory allocations** in all primitive functions
- **Execution times in the nanosecond range** for numeric types
- **Consistency** between functions with and without custom ranges

The library achieves its design goals of simplicity and performance, prioritizing execution speed and memory efficiency as documented in the project philosophy.

---

# WordsPT Allocation & Benchmark Audit (Task #37)

**Date:** 2026-07-22
**Package:** github.com/FlavioCFOliveira/gengo
**Platform:** linux/amd64 (AMD Ryzen 9 5900HX, 16 threads) · Go go1.26.5
**Command:** `go test -run=^$ -bench 'PT' -benchmem -benchtime 2s -count 1 .`

This section covers the pt-PT word generators (the `WordsPT` feature) that were added after the original report above. It records the full allocation/benchmark surface of every public `WordsPT` function and documents the Task #37 optimization: building an **additive** noun/adjective plural in a single allocation.

## 1. Optimization: single-allocation additive plural

The noun and adjective plural path used to be two allocations: assemble the singular word (one allocation), then apply a tail string transform (`pluralize`) to it (a second allocation). Task #37 splits the plural into two build strategies decided up front from each ending's own metadata (`additivePluralSuffix`):

- **Additive plural** — a pure suffix append: a regular `+s` for vowel/diphthong endings, or `+es` for `-r`/`-z`/`-n`/oxytone-`-s` endings. The suffix is now written into the **same** `strings.Builder` that assembles the word (`accentedWordSuffixed`), so the plural is **one allocation**. Reachable additive endings: nouns `-o`, `-a`, `-eiro`, `-eira`, `-or` (the only `+es` case), `-ora`, `-mento`, `-dade`, `-ista`; adjectives `-o`, `-a`, `-oso`, `-osa`, `-ico`, `-ica`, `-ivo`, `-iva`, `-ente`, `-ante`.
- **Substitutive plural** — a tail rewrite that cannot be a suffix append: `-ão→-ões/-ães/-ãos` (incl. the fixed `-ção→-ções`), `-m→-ns` (`-agem→-agens`), and the vowel+`l` endings `-al/-ável/-ível→-ais/-áveis/-íveis`. These **keep the post-assembly `pluralize` transform** and stay **two allocations**, because they modify graphemes inside the word (e.g. `-agens` carries the `ns` coda cluster that no single-coda syllable can represent).

The produced plural **string is unchanged** — only how it is built changed. The additive classifier draws no randomness, and the substitutive endings it defers (fixed `-ção`, `-m`, `-l`) draw none either, so the random stream is byte-for-byte identical: same-seed reproducibility is preserved. This equivalence is locked by regression tests in `wordspt_plural_optimization_test.go` (`TestAdditivePluralClassificationMatchesPluralize`, `TestAccentedWordSuffixedMatchesPluralize`, `TestPluralBuildAllocationBudget`, `TestAdditivePluralSuffixClassification`).

### Before / After (plural benchmarks, `-benchtime 2s`)

| Benchmark (Generator surface) | Allocs before | Allocs after | B/op before | B/op after | ns/op before | ns/op after |
|-------------------------------|:-------------:|:------------:|:-----------:|:----------:|:------------:|:-----------:|
| `BenchmarkNounPTOfPlural` | 2 | **1** | 25 | **15** | 451.3 | **404.4** |
| `BenchmarkAdjectivePTOfPositivePlural` | 2 | **1** | 27 | **17** | 477.9 | **450.1** |

`testing`'s reported `allocs/op` is the integer-truncated **average** over the benchmark's random mix of endings (both benchmarks force `Plural` but pick the ending at random). The average now truncates to 1 because additive endings dominate (~87% of noun plurals, ~75% of adjective plurals). The exact per-path costs are pinned by `TestPluralBuildAllocationBudget`: an additive build is **exactly 1** allocation and a substitutive build is **exactly 2**.

## 2. Full public-function benchmark surface

`ns/op`, `B/op` and `allocs/op` for every public `WordsPT` function, on the package-level (global source) and `*Generator` surfaces. `*Of…` and `Generator…` rows use a seeded `New(1)` generator.

### Open classes (noun, adjective, verb, adverb)

| Benchmark | Surface | ns/op | B/op | allocs/op |
|-----------|---------|------:|-----:|:---------:|
| `BenchmarkNounPT` | package | 457.4 | 13 | 1 |
| `BenchmarkNounPTOfSingular` | Generator | 405.4 | 12 | 1 |
| `BenchmarkNounPTOfPlural` | Generator | 405.3 | 15 | 1 (avg; additive 1 / substitutive 2) |
| `BenchmarkAdjectivePT` | package | 491.4 | 17 | 1 |
| `BenchmarkAdjectivePTOfPositiveSingular` | Generator | 430.6 | 13 | 1 |
| `BenchmarkAdjectivePTOfPositivePlural` | Generator | 448.2 | 17 | 1 (avg; additive 1 / substitutive 2) |
| `BenchmarkAdjectivePTOfSuperlativeSingular` | Generator | 398.7 | 19 | 1 |
| `BenchmarkVerbPT` | package | 444.1 | 13 | 1 |
| `BenchmarkVerbPTOfPresentIndicative` | Generator | 381.0 | 12 | 1 |
| `BenchmarkVerbPTOfImperfectSubjunctive` | Generator | 363.5 | 18 | 1 |
| `BenchmarkVerbPTOfInfinitive` | Generator | 371.6 | 12 | 1 |
| `BenchmarkAdverbPT` | package | 405.6 | 18 | 1 |
| `BenchmarkAdverbPTByLengthType` | Generator (Big) | 385.6 | 18 | 1 |

### Closed classes (curated selectors)

| Benchmark | Surface | ns/op | B/op | allocs/op |
|-----------|---------|------:|-----:|:---------:|
| `BenchmarkArticlePT` | package | 8.567 | 0 | 0 |
| `BenchmarkGeneratorArticlePT` | Generator | 4.584 | 0 | 0 |
| `BenchmarkPrepositionPT` | package | 8.430 | 0 | 0 |
| `BenchmarkConjunctionPT` | package | 8.922 | 0 | 0 |
| `BenchmarkPronounPT` | package | 8.487 | 0 | 0 |
| `BenchmarkGeneratorPronounPT` | Generator | 4.579 | 0 | 0 |
| `BenchmarkInterjectionPT` | package | 8.582 | 0 | 0 |
| `BenchmarkNumeralPT` | package | 8.240 | 0 | 0 |
| `BenchmarkGeneratorNumeralPT` | Generator | 4.740 | 0 | 0 |

### Orchestration (`WordPT` / `WordsPT`)

| Benchmark | Surface | ns/op | B/op | allocs/op |
|-----------|---------|------:|-----:|:---------:|
| `BenchmarkWordPT` | package | 454.2 | 11 | 0 (avg; see note) |
| `BenchmarkGeneratorWordPT` | Generator | 414.8 | 11 | 0 (avg; see note) |
| `BenchmarkWordsPT` | Generator (16 words) | 6646 | 436 | 15 (≈0.94/word) |

## 3. Allocation classification (0 / 1 / 2 allocations, and why)

- **0 allocations — the six closed-class selectors** (`ArticlePT`, `PrepositionPT`, `ConjunctionPT`, `PronounPT`, `InterjectionPT`, `NumeralPT`). They return a string that already lives in a package-level `[]string`; the selector only indexes into the shared backing array, copying no bytes (asserted at 0 allocs by `TestClosedClassNoAllocation`). The `*Generator` variants are ~2× faster than the package-level ones (~4.6 ns vs ~8.5 ns) because the package-level path goes through the `globalSource` wrapper over the concurrency-safe global `math/rand/v2` generator, while a `*Generator` draws from its own unsynchronized PCG source.
- **1 allocation — every open-class singular, the additive plural, the superlative, verbs, and adverbs.** The syllable buffer is a stack array (reused across length attempts with no heap cost), so the only allocation is the final assembled string. The additive plural (Task #37) keeps this to one allocation by appending its suffix into the assembling builder.
- **2 allocations — the substitutive noun/adjective plural only.** The singular is assembled (one allocation) and then its tail is rewritten by `pluralize`/`pluralizeNoun` (a second allocation). This is inherent: `-ões`, `-ns`, `-ais` and friends change graphemes inside the word, so they cannot be produced by appending to the singular.
- **`WordPT` averages 0 allocations (truncated).** `WordPT` dispatches across the open and closed classes by type; the frequent zero-allocation closed-class words pull the truncated per-word average below one (`TestWordPTAllocationBudget` asserts the average stays below two). `WordsPT(16)` reports 15 allocs for 16 words (≈0.94/word) for the same reason, reusing one syllable buffer across the whole slice.

## 4. Gate results (Task #37)

| Gate | Result |
|------|--------|
| `gofmt -l .` | clean (no files) |
| `go vet ./...` | clean |
| `golangci-lint run` | 0 issues |
| `go test -race -count=1 .` | ok (full suite, no regression) |
| `go mod tidy -diff` | clean |
| `govulncheck ./...` | No vulnerabilities found |
