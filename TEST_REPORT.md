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
| `TestUInt8` | PASS | 0.00s |
| `TestUInt8Between` | PASS | 0.00s |
| `TestByte` | PASS | 0.00s |
| `TestUInt16` | PASS | 0.00s |
| `TestUInt16Between` | PASS | 0.00s |
| `TestUInt32` | PASS | 0.00s |
| `TestUInt32Between` | PASS | 0.00s |
| `TestUInt64` | PASS | 0.00s |
| `TestUInt64Between` | PASS | 0.00s |
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

#### Unsigned Integers (UInt8/16/32/64) and Byte
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkUInt8` | 1,000,000,000 | 4.780 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt8Between` | 801,432,346 | 7.463 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkByte` | 1,000,000,000 | 4.768 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt16` | 1,000,000,000 | 4.767 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt16Between` | 811,150,370 | 7.515 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt32` | 1,000,000,000 | 4.796 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt32Between` | 803,114,101 | 7.489 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt64` | 1,000,000,000 | 4.786 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUInt64Between` | 757,757,079 | 7.942 ns/op | 0 B/op, 0 allocs/op |

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
   - `BenchmarkUInt16`: 4.767 ns/op

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
=== RUN   TestUInt8
--- PASS: TestUInt8 (0.00s)
=== RUN   TestUInt8Between
--- PASS: TestUInt8Between (0.00s)
=== RUN   TestByte
--- PASS: TestByte (0.00s)
=== RUN   TestUInt16
--- PASS: TestUInt16 (0.00s)
=== RUN   TestUInt16Between
--- PASS: TestUInt16Between (0.00s)
=== RUN   TestUInt32
--- PASS: TestUInt32 (0.00s)
=== RUN   TestUInt32Between
--- PASS: TestUInt32Between (0.00s)
=== RUN   TestUInt64
--- PASS: TestUInt64 (0.00s)
=== RUN   TestUInt64Between
--- PASS: TestUInt64Between (0.00s)
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
BenchmarkUInt8-16                	1000000000	         4.780 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt8Between-16         	801432346	         7.463 ns/op	       0 B/op	       0 allocs/op
BenchmarkByte-16                 	1000000000	         4.768 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt16-16               	1000000000	         4.767 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt16Between-16        	811150370	         7.515 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt32-16               	1000000000	         4.796 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt32Between-16        	803114101	         7.489 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt64-16               	1000000000	         4.786 ns/op	       0 B/op	       0 allocs/op
BenchmarkUInt64Between-16        	757757079	         7.942 ns/op	       0 B/op	       0 allocs/op
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
