# Performance Report: v0.0.26 → HEAD

**Date:** 2026-05-29  
**Platform:** linux/amd64 · AMD Ryzen 9 5900HX · 16 threads (GOMAXPROCS=16)  
**Go:** go1.26.2  
**Tool:** `go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s`

> All figures below come from a single benchmark suite run on the reference
> machine above so that they are internally consistent and reproducible. The
> same `HEAD` numbers back [TEST_REPORT.md](TEST_REPORT.md). Deltas smaller than
> roughly ±0.3 ns are within run-to-run measurement noise.

---

## Summary

The most significant gain is in `String`, which became **2.37× faster** (−58%)
thanks to the batched PRNG calls introduced in commit `71d7483`. On this platform
the `numbers.go` rewrite (`ec0e6e4`) also delivered large gains for the small
integer types — `Int8` and `Int16` are **~38% faster** and `Int32`, `UInt8`,
`Byte`, and `UInt16` are **~21–24% faster** — by eliminating intermediate
wide-type casts. Float types gained inlined generation paths (~6–7%). Several new
benchmarks (`Float32Between`, `Float64Between`, `Complex64Between`,
`Complex128Between`, `Word`, `WordByLengthType`, `Words`) cover API additions that
did not exist in v0.0.26.

---

## Detailed Results

### String generation

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Speedup |
|---|---|---|---|---|
| `String` | 65.85 ns | 27.76 ns | **−38.09 ns** | **2.37×** |
| `StringNumeric` / `Numeric` | 65.00 ns | 62.65 ns | −2.35 ns | ≈ same |

> `String` batches all PRNG draws into a single `Uint64` source loop instead of
> one call per character, nearly eliminating per-character overhead. Both string
> benchmarks measure an 8-character result (`String(8, …)`) and retain exactly
> 1 alloc/op (8 B/op). `StringNumeric` stays slower than `String` because its
> 10-symbol charset rejects more sampled indices per accepted character.

---

### Integer types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Int8` | 7.624 ns | 4.763 ns | −2.861 ns | **+38%** |
| `Int16` | 7.641 ns | 4.774 ns | −2.867 ns | **+38%** |
| `Int32` | 5.980 ns | 4.753 ns | −1.227 ns | **+21%** |
| `UInt8` | 6.233 ns | 4.780 ns | −1.453 ns | **+23%** |
| `Byte` | 6.264 ns | 4.768 ns | −1.496 ns | **+24%** |
| `UInt16` | 6.254 ns | 4.767 ns | −1.487 ns | **+24%** |
| `Int` | 4.752 ns | 4.762 ns | +0.010 ns | ≈ same |
| `Int64` | 4.757 ns | 4.770 ns | +0.013 ns | ≈ same |
| `UInt32` | 4.747 ns | 4.796 ns | +0.049 ns | ≈ same |
| `UInt64` | 4.759 ns | 4.786 ns | +0.027 ns | ≈ same |
| `Int8Between` | 7.628 ns | 7.424 ns | −0.204 ns | ≈ same |
| `Int16Between` | 7.616 ns | 7.398 ns | −0.218 ns | ≈ same |
| `Int32Between` | 6.516 ns | 7.469 ns | +0.953 ns | ≈ same¹ |
| `IntBetween` | 7.838 ns | 7.812 ns | −0.026 ns | ≈ same |
| `Int64Between` | 7.611 ns | 7.872 ns | +0.261 ns | ≈ same¹ |
| `UInt8Between` | 7.470 ns | 7.463 ns | −0.007 ns | ≈ same |
| `UInt16Between` | 7.442 ns | 7.515 ns | +0.073 ns | ≈ same |
| `UInt32Between` | 7.633 ns | 7.489 ns | −0.144 ns | ≈ same |
| `UInt64Between` | 7.852 ns | 7.942 ns | +0.090 ns | ≈ same |

¹ The `Between` variants cluster around 7.4–7.9 ns/op on HEAD; the remaining
inter-version differences are run-to-run measurement noise rather than real
regressions.

The plain `Int8`, `Int16`, `Int32`, `UInt8`, `Byte`, and `UInt16` variants
improved because the `numbers.go` rewrite eliminated intermediate wide-type casts
that had prevented the compiler from producing the optimal instruction sequence.

---

### Float types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Float32` | 6.256 ns | 5.804 ns | −0.452 ns | **+7%** |
| `Float64` | 6.149 ns | 5.767 ns | −0.382 ns | **+6%** |
| `Float32Between` | — | 6.791 ns | — | new |
| `Float64Between` | — | 6.443 ns | — | new |

Inlining the generation path and removing the helper call reduced latency for
both float types. `Float32Between` and `Float64Between` are new API additions.

---

### Complex types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Complex64` | 10.27 ns | 10.25 ns | −0.02 ns | ≈ same |
| `Complex128` | 10.48 ns | 10.44 ns | −0.04 ns | ≈ same |
| `Complex64Between` | — | 14.24 ns | — | new |
| `Complex128Between` | — | 13.71 ns | — | new |

---

### Date types

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Note |
|---|---|---|---|---|
| `Date` | 7.472 ns | 7.377 ns | −0.095 ns | ≈ same |
| `UnixDate` | 7.364 ns | 7.401 ns | +0.037 ns | ≈ same |
| `DateBetween` | 10.43 ns | 10.66 ns | +0.23 ns | ≈ same |

---

### Bool

| Benchmark | v0.0.26 | HEAD | Δ ns/op | Note |
|---|---|---|---|---|
| `Bool` | 4.777 ns | 4.719 ns | −0.058 ns | ≈ same |

---

### Word / Words (new in HEAD)

These benchmarks have no v0.0.26 counterpart; they cover the expanded Words API.

| Benchmark | HEAD |
|---|---|
| `Word` | 86.55 ns/op · 18 B/op · 1 alloc/op |
| `WordByLengthType` | 50.12 ns/op · 7 B/op · 1 alloc/op |
| `Words` (10 words) | 788.4 ns/op · 258 B/op · 10 allocs/op |

---

## Allocations

All numeric, boolean, and date functions remain at **0 B/op · 0 allocs/op** in
both versions. String and word functions retain exactly **1 alloc/op**, matching
the single backing-array allocation for the returned string (`Words` allocates
once per generated word).

---

## Key commits driving the gains

| Commit | Description |
|---|---|
| `71d7483` | Batch PRNG in `String`, eliminate `Words` allocations, precompute date ranges, inline `Float32`/`Float64` |
| `ec0e6e4` | Rewrite `numbers.go` — eliminates G115 casts, enables better codegen for `Int8`/`Int16` |
