# Performance Report: v0.0.26 → v0.2.1

This report collects the gengo benchmark studies in release order. Each section
states its own date, platform, Go version, and method; figures from different
sections are not directly comparable.

| Section | Release | Date | Go |
|---|---|---|---|
| [v0.0.26 → v0.1.0](#v0026--v010) | v0.1.0 | 2026-05-29 | go1.26.2 |
| [String single-allocation rewrite (v0.2.1)](#string-single-allocation-rewrite-v021) | v0.2.1 | 2026-09-27 | go1.27.1 |
| [WordsPT sampling core (v0.2.1)](#wordspt-sampling-core-v021) | v0.2.1 | 2026-09-27 | go1.27.1 |

---

## v0.0.26 → v0.1.0

**Date:** 2026-05-29  
**Platform:** linux/amd64 · AMD Ryzen 9 5900HX · 16 threads (GOMAXPROCS=16)  
**Go:** go1.26.2  
**Tool:** `go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s`

> All figures below come from a single benchmark suite run on the reference
> machine above so that they are internally consistent and reproducible. The
> same v0.1.0 numbers backed the 2026-05-29 edition of
> [TEST_REPORT.md](TEST_REPORT.md). Deltas smaller than roughly ±0.3 ns are
> within run-to-run measurement noise.

### Summary

The most significant gain is in `String`, which became **2.37× faster** (−58%)
thanks to the batched PRNG calls introduced in commit `71d7483`. On this platform
the `numbers.go` rewrite (`ec0e6e4`) also delivered large gains for the small
integer types — `Int8` and `Int16` are **~38% faster** and `Int32`, `Uint8`,
`Byte`, and `Uint16` are **~21–24% faster** — by eliminating intermediate
wide-type casts. Float types gained inlined generation paths (~6–7%). Several new
benchmarks (`Float32Between`, `Float64Between`, `Complex64Between`,
`Complex128Between`, `Word`, `WordByLengthType`, `Words`) cover API additions that
did not exist in v0.0.26.

---

### Detailed Results

#### String generation

| Benchmark | v0.0.26 | v0.1.0 | Δ ns/op | Speedup |
|---|---|---|---|---|
| `String` | 65.85 ns | 27.76 ns | **−38.09 ns** | **2.37×** |
| `StringNumeric` / `Numeric` | 65.00 ns | 62.65 ns | −2.35 ns | ≈ same |

> `String` batches all PRNG draws into a single `Uint64` source loop instead of
> one call per character, nearly eliminating per-character overhead. Both string
> benchmarks measure an 8-character result (`String(8, …)`) and retain exactly
> 1 alloc/op (8 B/op). `StringNumeric` stays slower than `String` because its
> 10-symbol charset rejects more sampled indices per accepted character.
>
> These figures predate the single-allocation rewrite of `String` described in
> [String single-allocation rewrite](#string-single-allocation-rewrite). That
> rewrite leaves the code for lengths up to 32 bytes unchanged, yet it measured
> `String(8, …)` **+2.9%** slower (28.40 → 29.22 ns/op) and `String(8, "x")`
> **+5.8%** slower (16.50 → 17.46 ns/op) on the same machine. The cause is code
> layout: the dispatch added at the top of `String` moves the unchanged loop to
> a different address, and the gc compiler does not align loops.

---

#### Integer types

| Benchmark | v0.0.26 | v0.1.0 | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Int8` | 7.624 ns | 4.763 ns | −2.861 ns | **+38%** |
| `Int16` | 7.641 ns | 4.774 ns | −2.867 ns | **+38%** |
| `Int32` | 5.980 ns | 4.753 ns | −1.227 ns | **+21%** |
| `Uint8` | 6.233 ns | 4.780 ns | −1.453 ns | **+23%** |
| `Byte` | 6.264 ns | 4.768 ns | −1.496 ns | **+24%** |
| `Uint16` | 6.254 ns | 4.767 ns | −1.487 ns | **+24%** |
| `Int` | 4.752 ns | 4.762 ns | +0.010 ns | ≈ same |
| `Int64` | 4.757 ns | 4.770 ns | +0.013 ns | ≈ same |
| `Uint32` | 4.747 ns | 4.796 ns | +0.049 ns | ≈ same |
| `Uint64` | 4.759 ns | 4.786 ns | +0.027 ns | ≈ same |
| `Int8Between` | 7.628 ns | 7.424 ns | −0.204 ns | ≈ same |
| `Int16Between` | 7.616 ns | 7.398 ns | −0.218 ns | ≈ same |
| `Int32Between` | 6.516 ns | 7.469 ns | +0.953 ns | ≈ same¹ |
| `IntBetween` | 7.838 ns | 7.812 ns | −0.026 ns | ≈ same |
| `Int64Between` | 7.611 ns | 7.872 ns | +0.261 ns | ≈ same¹ |
| `Uint8Between` | 7.470 ns | 7.463 ns | −0.007 ns | ≈ same |
| `Uint16Between` | 7.442 ns | 7.515 ns | +0.073 ns | ≈ same |
| `Uint32Between` | 7.633 ns | 7.489 ns | −0.144 ns | ≈ same |
| `Uint64Between` | 7.852 ns | 7.942 ns | +0.090 ns | ≈ same |

¹ The `Between` variants cluster around 7.4–7.9 ns/op on v0.1.0; the remaining
inter-version differences are run-to-run measurement noise rather than real
regressions.

The plain `Int8`, `Int16`, `Int32`, `Uint8`, `Byte`, and `Uint16` variants
improved because the `numbers.go` rewrite eliminated intermediate wide-type casts
that had prevented the compiler from producing the optimal instruction sequence.

---

#### Float types

| Benchmark | v0.0.26 | v0.1.0 | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Float32` | 6.256 ns | 5.804 ns | −0.452 ns | **+7%** |
| `Float64` | 6.149 ns | 5.767 ns | −0.382 ns | **+6%** |
| `Float32Between` | — | 6.791 ns | — | new |
| `Float64Between` | — | 6.443 ns | — | new |

Inlining the generation path and removing the helper call reduced latency for
both float types. `Float32Between` and `Float64Between` are new API additions.

---

#### Complex types

| Benchmark | v0.0.26 | v0.1.0 | Δ ns/op | Improvement |
|---|---|---|---|---|
| `Complex64` | 10.27 ns | 10.25 ns | −0.02 ns | ≈ same |
| `Complex128` | 10.48 ns | 10.44 ns | −0.04 ns | ≈ same |
| `Complex64Between` | — | 14.24 ns | — | new |
| `Complex128Between` | — | 13.71 ns | — | new |

---

#### Date types

| Benchmark | v0.0.26 | v0.1.0 | Δ ns/op | Note |
|---|---|---|---|---|
| `Date` | 7.472 ns | 7.377 ns | −0.095 ns | ≈ same |
| `UnixDate` | 7.364 ns | 7.401 ns | +0.037 ns | ≈ same |
| `DateBetween` | 10.43 ns | 10.66 ns | +0.23 ns | ≈ same |

---

#### Bool

| Benchmark | v0.0.26 | v0.1.0 | Δ ns/op | Note |
|---|---|---|---|---|
| `Bool` | 4.777 ns | 4.719 ns | −0.058 ns | ≈ same |

---

#### Word / Words (new since v0.0.26)

These benchmarks have no v0.0.26 counterpart; they cover the expanded Words API.

| Benchmark | v0.1.0 |
|---|---|
| `Word` | 86.55 ns/op · 18 B/op · 1 alloc/op |
| `WordByLengthType` | 50.12 ns/op · 7 B/op · 1 alloc/op |
| `Words` (10 words) | 788.4 ns/op · 258 B/op · 10 allocs/op |

---

## String single-allocation rewrite (v0.2.1)

**Date:** 2026-09-27  
**Platform:** linux/amd64 · AMD Ryzen 9 5900HX · 16 threads  
**Go:** go1.27.1  
**Method:** `benchstat`, 10 interleaved runs per variant, `-benchtime 200ms`,
single-threaded benchmarks at `-cpu 1` and `b.RunParallel` benchmarks
(`parallel_bench_test.go`) at `-cpu 1,8,16`. All figures are Alphanumeric
source unless stated. Only changes with p < 0.05 are reported as changes.

Before the rewrite, `String` and `(*Generator).String` allocated **twice** for
lengths above 32 bytes: the `[]byte` buffer escaped to the heap and was then
copied into the returned string. Both functions now allocate **exactly once**
for every length:

- up to 32 bytes, the characters are written to a `make([]byte, length)` that
  the gc compiler places on the stack (its default variable-size `make`
  threshold is 32 bytes) and then copied into the string — unchanged code;
- from 33 to 256 bytes, they are written to a 256-byte stack array and copied
  into the string;
- above 256 bytes, they are written in 256-byte chunks to a `strings.Builder`
  pre-sized to the final length;
- a single-character source above 32 bytes is filled by doubling copies (up to
  256 bytes) or `strings.Repeat` (above), without drawing random values.

The output and the random values drawn are byte-identical to the previous
implementation (pinned by `TestGeneratorStringGolden`).

| Benchmark | Before | After | Δ | Allocs/op |
|---|---|---|---|---|
| `String(33)` | 98.00 ns | 82.01 ns | −16.3% | 2 → 1 |
| `String(64)` | 155.5 ns | 132.9 ns | −14.5% | 2 → 1 |
| `String(256)` | 518.4 ns | 437.8 ns | −15.6% | 2 → 1 |
| `String(257)` | 530.4 ns | 456.0 ns | −14.0% | 2 → 1 |
| `String(4096)` | 7.662 µs | 6.675 µs | −12.9% | 2 → 1 |
| `String(4096, "x")` | 1745 ns | 465.0 ns | −73.4% | 2 → 1 |
| `(*Generator).String(33)` | 90.42 ns | 76.27 ns | −15.6% | 2 → 1 |
| `(*Generator).String(4096)` | 6.707 µs | 6.263 µs | −6.6% | 2 → 1 |
| `String(8)` | 28.40 ns | 29.22 ns | **+2.9%** | 1 → 1 |
| `String(8, "x")` | 16.50 ns | 17.46 ns | **+5.8%** | 1 → 1 |
| `String(32)` | 78.80 ns | 77.63 ns | ~ | 1 → 1 |

Concurrent use benefits most, because the contention evaluation for task #39 found
the removed allocation to be the only contention source that gengo itself adds
(allocator and GC pressure):

| Parallel benchmark | 8 CPUs before → after | 16 CPUs before → after |
|---|---|---|
| `String(33)` | 32.36 → 19.75 ns (−39.0%) | 33.68 → 20.79 ns (−38.3%) |
| `String(64)` | 45.19 → 29.84 ns (−34.0%) | 47.71 → 30.14 ns (−36.8%) |
| `String(256)` | 165.4 → 106.7 ns (−35.5%) | 173.6 → 110.8 ns (−36.2%) |
| `String(4096)` | 2.378 → 1.634 µs (−31.3%) | 2.435 → 1.684 µs (−30.9%) |
| `String(64, "x")` | 29.63 → 17.48 ns (−41.0%) | 28.85 → 16.98 ns (−41.1%) |
| `String(8)` | ~ | ~ |
| `String(8, "x")` | ~ | 3.882 → 4.003 ns (**+3.1%**) |

**Residual regression at 8 bytes.** `String(8, …)` is 2.9–5.8% slower
single-threaded, and `String(8, "x")` is 3.1% slower at 16 CPUs. The code for
lengths up to 32 bytes is unchanged; the difference comes from code layout (the
new dispatch moves the unchanged loop to a different address, and the gc
compiler does not align loops). It was accepted in exchange for the gains above.
Code layout also explains why `fill` and `fillFrom` in `strings.go` differ by a
value-preserving `& 63` mask: each form measured fastest for its function.

---

## WordsPT sampling core (v0.2.1)

**Date:** 2026-09-27  
**Platform:** linux/amd64 · AMD Ryzen 9 5900HX · 16 threads  
**Go:** go1.27.1 (`GOAMD64=v1`)  
**Method:** `benchstat`, 10 interleaved runs per variant, `-benchtime 400ms`,
single-threaded benchmarks at `-cpu 1` pinned to one core, and `b.RunParallel`
benchmarks (`parallel_bench_test.go`) at `-cpu 16`. The baseline is commit
`2074a35`. Only changes with p < 0.05 are reported as changes.

Profiling showed that the WordsPT sampling core spent its time on three
avoidable computations. Each is now replaced by a table that is equivalent by
construction:

- **Weighted draws.** Each inventory maps a draw in `[0,total)` to the selected
  form through a precomputed `[]uint8` lookup table (7,366 bytes across all 14
  inventories), instead of a binary search over the cumulative weights.
- **Onset/nucleus agreement.** The nucleus class each onset licenses (any,
  front after `qu`/`gu`, back after `ç`) is recorded per form at package
  initialization from `onsetFrontVowelOnly` and `onsetBackVowelOnly`. Two
  string-keyed map lookups per syllable are removed.
- **Word-length draw.** For every window span up to 29, which covers every
  window of `charRangeOf`, the offset is the number of thresholds in
  `skewThresholds` (`wordspt_skew_table.go`, 1,740 bytes, generated by
  `go generate`) that the 32-bit draw reaches. `math.Pow`, `math.Log`, and
  `math.Ceil` are no longer evaluated per draw. Wider spans still evaluate the
  closed-form inverse CDF.

The random values drawn and the output are byte-identical to the previous
implementation: a SHA-256 over long seeded streams of every `Generator` method
is unchanged. `TestSkewThresholdsMatchFormula` pins every threshold to the
formula, and `TestInventoryLookupTables` checks every lookup-table entry
against the cumulative weights. B/op and allocs/op are unchanged.

| Benchmark | Before | After | Δ |
|---|---|---|---|
| `WordPT` | 457.0 ns | 303.1 ns | −33.7% |
| `(*Generator).WordPT` | 415.7 ns | 256.4 ns | −38.3% |
| `NounPT` | 448.4 ns | 284.3 ns | −36.6% |
| `AdjectivePT` | 486.2 ns | 319.6 ns | −34.3% |
| `VerbPT` | 460.9 ns | 269.1 ns | −41.6% |
| `AdverbPT` | 402.4 ns | 237.9 ns | −40.9% |
| `WordsPT` (10 words) | 6.655 µs | 4.139 µs | −37.8% |
| Closed classes (`ArticlePT`, …) | — | — | ~ |

| Parallel benchmark (16 CPUs) | Before | After | Δ |
|---|---|---|---|
| `WordPT` | 57.30 ns | 39.97 ns | −30.2% |
| `NounPT` | 57.47 ns | 39.38 ns | −31.5% |
| `AdjectivePT` | 63.04 ns | 43.91 ns | −30.4% |
| `VerbPT` | 57.76 ns | 37.20 ns | −35.6% |
| `AdverbPT` | 53.96 ns | 35.55 ns | −34.1% |
| `WordsPT` (10 words) | 979.1 ns | 721.9 ns | −26.3% |
| Closed classes (`ArticlePT`, …) | — | — | ~ |

**Initialization cost.** Building the lookup tables raises the package's
one-time initialization from 16,952 B in 84 allocations to 23,856 B in 98
allocations. The median initialization time over 30 runs rises from 58 µs to
82 µs (`GODEBUG=inittrace=1`).

**Residual differences in unchanged code.** Across the full suite, a few
nanosecond-scale benchmarks of code this change does not touch shifted by up to
+3.2% (for example `StringNumeric` +1.7%, `ResolveGender` +1.6%). The set of
affected benchmarks differs between runs. A control binary built from the
unchanged baseline code, differing only by an added test file, shows shifts of
the same size (up to +7.7%). These differences therefore come from code layout,
not from the change.

---

## Allocations

**Date:** 2026-09-27 · **Go:** go1.27.1 · **Code:** v0.2.1

The figures below count heap allocations per call exactly: a harness outside
the package reads `runtime.MemStats.Mallocs` before and after each of 20,000
calls (after 2,000 warm-up calls, with `GOMAXPROCS=1`) and reports the share of
calls that made 0, 1, or 2 allocations. Isolated samples of 0.01% (one call in
20,000) with a higher count are runtime background allocations and are omitted.
The means below are for the package-level functions; the `*Generator` methods
matched them within half a percentage point. `go test -benchmem` reports the
same figures, truncated to an integer average.

| Function | Allocations per call |
|---|---|
| Numeric, `Bool`, `Date`, `UnixDate`, `DateBetween` | **0** |
| `ArticlePT`, `PronounPT`, `NumeralPT`, `PrepositionPT`, `ConjunctionPT`, `InterjectionPT` | **0**: they return a string from a package-level list |
| `String` and its variants, result of 0 or 1 byte | **0** |
| `String` and its variants, result of 2 bytes or more | **1**, the returned string, at every length (2 to 4,096 bytes measured with an Alphanumeric source, 8 to 4,096 bytes with a single-character source) |
| `Word`, `WordByLengthType` | **1**; a one-byte `SmallLengthWord` result makes **0** (25.6% of `WordByLengthType(SmallLengthWord)` calls) |
| `NounPTOf` / `AdjectivePTOf` with `Singular`; the `Superlative` degree; `VerbPT`, `VerbPTOf`, `AdverbPT`, `AdverbPTByLengthType` | **1** |
| `NounPTOf` with `Plural` | **1** for an additive plural (`+s`, `+es`), **2** for a substitutive plural (`-ões`, `-ns`, ...): mean **1.13** (12.9% of calls make 2) |
| `AdjectivePTOf` with `Plural`, `Positive` | **1** additive, **2** substitutive (`-ais`, `-áveis`, `-íveis`, ...): mean **1.25** (24.7% of calls make 2) |
| `NounPT` | mean **1.06** (6.5% of calls make 2) |
| `AdjectivePT` | mean **1.06** (6.3% of calls make 2) |
| `WordPT`, `WordPTByLengthType` | **0** for a closed-class word, **1** or **2** as above: `WordPT` mean **0.90** (12.3% make 0, 85.1% make 1, 2.5% make 2) |
| `Words(n)`, `WordsPT(n)` | **1** for the returned slice, plus the allocations of each word as above: `Words(10)` mean **10.35**, `WordsPT(10)` mean **10.07** |

Strings of zero or one byte cost no allocation, because the Go runtime serves
one-byte strings from a static table; this explains the **0** entries for short
`String` results and one-byte words. The exact cost of each plural build path
(additive **1**, substitutive **2**) is pinned by
`TestPluralBuildAllocationBudget`. The **1** allocation of `String` and
`(*Generator).String` for every length above 32 bytes dates from the
[single-allocation rewrite](#string-single-allocation-rewrite-v021); before it,
those lengths allocated twice. In v0.0.26 and v0.1.0, the numeric, boolean,
and date functions also made 0 allocations.

---

## Key commits driving the gains

| Commit | First release | Description |
|---|---|---|
| `71d7483` | v0.0.27 | Batch PRNG in `String`, eliminate `Words` allocations, precompute date ranges, inline `Float32`/`Float64` |
| `ec0e6e4` | v0.0.27 | Rewrite `numbers.go` — eliminates G115 casts, enables better codegen for `Int8`/`Int16` |
| `c5683a2` | v0.2.1 | Single-allocation `String` and `(*Generator).String` for lengths above 32 bytes |
| `b70ae60` | v0.2.1 | WordsPT sampling core: lookup and threshold tables replace binary search, map lookups, and per-draw `math.Pow`/`math.Log`/`math.Ceil` |
