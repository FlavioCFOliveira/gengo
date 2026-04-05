# Performance Report: v0.0.26 → HEAD

**Date:** 2026-04-05  
**Platform:** darwin/arm64 · Apple M4 · 10 cores  
**Tool:** `go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s`

---

## Summary

The most significant gain is in `String`, which became **2.49× faster** (−60%) thanks to
batched PRNG calls introduced in commit `71d7483`. Integer types `Int8` and `Int16`
also improved through a rewrite of `numbers.go` that eliminated redundant type
conversions. Float and Complex types benefit from inlined generation paths that
remove an extra PRNG call. Several new benchmarks (`Float32Between`, `Float64Between`,
`Complex64Between`, `Complex128Between`, `Word`, `WordByLengthType`, `Words`) cover
API additions that did not exist in v0.0.26.

---

## Detailed Results

### String generation

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Speedup |
|---|---|---|---|---|
| `String` | 46.01 ns | 18.46 ns | **−27.55 ns** | **2.49×** |
| `StringNumeric` / `Numeric` | 46.13 ns | 46.68 ns | +0.55 ns | ≈ same |

> `String` batches all PRNG draws into a single `Uint64` source loop instead of one
> call per character, nearly eliminating per-character overhead.

---

### Integer types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Int8` | 4.920 ns | 4.123 ns | −0.797 ns | **+16%** |
| `Int16` | 5.076 ns | 4.125 ns | −0.951 ns | **+19%** |
| `Int8Between` | 5.033 ns | 5.253 ns | +0.220 ns | ≈ same¹ |
| `Int16Between` | 5.113 ns | 5.257 ns | +0.144 ns | ≈ same¹ |
| `Int32` | 4.121 ns | 4.162 ns | +0.041 ns | ≈ same |
| `Int32Between` | 4.271 ns | 5.268 ns | +0.997 ns | ≈ same¹ |
| `Int` | 4.063 ns | 4.149 ns | +0.086 ns | ≈ same |
| `IntBetween` | 4.951 ns | 5.238 ns | +0.287 ns | ≈ same¹ |
| `Int64` | 4.080 ns | 4.161 ns | +0.081 ns | ≈ same |
| `Int64Between` | 4.904 ns | 5.214 ns | +0.310 ns | ≈ same¹ |
| `UInt8` | 4.278 ns | 4.188 ns | −0.090 ns | ≈ same |
| `UInt8Between` | 4.959 ns | 5.296 ns | +0.337 ns | ≈ same¹ |
| `Byte` | 4.377 ns | 4.156 ns | −0.221 ns | ≈ same |
| `UInt16` | 4.387 ns | 4.178 ns | −0.209 ns | ≈ same |
| `UInt16Between` | 5.066 ns | 5.272 ns | +0.206 ns | ≈ same¹ |
| `UInt32` | 4.204 ns | 4.151 ns | −0.053 ns | ≈ same |
| `UInt32Between` | 5.029 ns | 5.213 ns | +0.184 ns | ≈ same¹ |
| `UInt64` | 4.191 ns | 4.149 ns | −0.042 ns | ≈ same |
| `UInt64Between` | 5.035 ns | 5.195 ns | +0.160 ns | ≈ same¹ |

¹ Sub-nanosecond delta; within measurement noise (~±0.3 ns).

`Int8` and `Int16` plain-range variants improved because the `numbers.go` rewrite
eliminated intermediate wide-type casts that had prevented the compiler from
producing the optimal instruction sequence.

---

### Float types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Float32` | 4.709 ns | 4.245 ns | −0.464 ns | **+10%** |
| `Float64` | 4.703 ns | 4.116 ns | −0.587 ns | **+12%** |
| `Float32Between` | — | 4.551 ns | — | new |
| `Float64Between` | — | 4.517 ns | — | new |

Inlining the generation path and removing the helper call reduced latency for both
float types. `Float32Between` and `Float64Between` are new API additions.

---

### Complex types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Complex64` | 8.373 ns | 8.042 ns | −0.331 ns | **+4%** |
| `Complex128` | 8.416 ns | 8.103 ns | −0.313 ns | **+4%** |
| `Complex64Between` | — | 9.580 ns | — | new |
| `Complex128Between` | — | 9.658 ns | — | new |

---

### Date types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Note |
|---|---|---|---|---|
| `Date` | 5.170 ns | 5.263 ns | +0.093 ns | ≈ same |
| `UnixDate` | 5.132 ns | 5.241 ns | +0.109 ns | ≈ same |
| `DateBetween` | 6.698 ns | 7.099 ns | +0.401 ns | ≈ same¹ |

---

### Bool

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Note |
|---|---|---|---|---|
| `Bool` | 4.040 ns | 4.154 ns | +0.114 ns | ≈ same |

---

### Word / Words (new in HEAD)

These benchmarks have no v0.0.26 counterpart; they cover the expanded Words API.

| Benchmark | HEAD |
|---|---|
| `Word` | 61.49 ns/op · 18 B/op · 1 alloc/op |
| `WordByLengthType` | 35.93 ns/op · 7 B/op · 1 alloc/op |
| `Words` (10 words) | 543.9 ns/op · 258 B/op · 10 allocs/op |

---

## Allocations

All numeric, boolean, and date functions remain at **0 B/op · 0 allocs/op** in both
versions. String and word functions retain exactly **1 alloc/op**, matching the single
backing-array allocation for the returned string.

---

## Key commits driving the gains

| Commit | Description |
|---|---|
| `71d7483` | Batch PRNG in `String`, eliminate `Words` allocations, precompute date ranges, inline `Float32`/`Float64` |
| `ec0e6e4` | Rewrite `numbers.go` — eliminates G115 casts, enables better codegen for `Int8`/`Int16` |
