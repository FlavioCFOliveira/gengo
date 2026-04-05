# Changelog

All notable changes to **gengo** are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) conventions,
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.0.27] — 2026-04-05

### Added
- `Float32Between(min, max)` and `Float64Between(min, max)` — bounded float generation.
- `Complex64Between` and `Complex128Between` — bounded complex number generation.
- `Word()`, `WordByLengthType()`, and `Words()` — word and word-list generation with length-category control.
- Open-source tooling: `govulncheck`, `golangci-lint` integration via `make dev-install`.
- CI/CD: vulnerability scanning (`govulncheck`) and security linting in GitHub Actions.
- `BENCHMARKS.md` — performance report comparing v0.0.26 to HEAD.

### Changed
- Go minimum version raised to **1.26**.
- GitHub Actions updated to `actions/checkout@v4` and `actions/setup-go@v5` with module cache enabled.
- `String` rewritten to batch all PRNG draws per character into a single `Uint64` loop — **2.49× faster** (46 ns → 18 ns/op).
- `numbers.go` rewritten to eliminate redundant wide-type casts (G115), improving `Int8` by **+16%** and `Int16` by **+19%**.
- `Float32` and `Float64` generation paths inlined — **+10%** and **+12%** respectively.
- `Complex64` and `Complex128` generation paths inlined — **+4%** each.
- Date range constants (`dtUnixMin`, `dtUnixMax`) precomputed at init time to avoid repeated calculation.
- Test coverage raised to **98.4%**; edge-case branches in `IntBetween` and `String` now covered.
- Makefile `golangci-lint` and `govulncheck` commands simplified.

### Fixed
- `Words` panicked on negative length input — now returns an empty slice.
- Critical correctness bugs in range-bounded integer generation.
- Naked return statements, inconsistent comment style, and doc-comment formatting flagged by `golangci-lint`.
- `G115` (integer overflow conversion) lint findings scoped and resolved throughout `numbers.go`.

---

## [v0.0.26] — 2026-02-18

### Changed
- Go version in `go.mod` stabilised at **1.25.0** after a round-trip through 1.26.0 → 1.22.0 → 1.25.0.

---

## [v0.0.25] — 2026-02-17

> No code changes. Re-tag of [v0.0.19](#v0019--2026-02-17).

---

## [v0.0.24] — 2026-02-17

> No code changes. Re-tag of [v0.0.19](#v0019--2026-02-17).

---

## [v0.0.23] — 2026-02-17

> No code changes. Re-tag of [v0.0.19](#v0019--2026-02-17).

---

## [v0.0.22] — 2026-02-17

> No code changes. Re-tag of [v0.0.19](#v0019--2026-02-17).

---

## [v0.0.21] — 2026-02-17

> No code changes. Re-tag of [v0.0.19](#v0019--2026-02-17).

---

## [v0.0.20] — 2026-02-17

> No code changes. Re-tag of [v0.0.19](#v0019--2026-02-17).

---

## [v0.0.19] — 2026-02-17

### Added
- `CLAUDE.md` — project guidance and architecture documentation for Claude Code.

### Changed
- `Bool` and `Complex` generation paths optimised for better throughput.
- Random number generation functions updated; maximum string length increased.
- README expanded with module description and usage examples.

---

## [v0.0.18] — 2025-04-15

### Added
- `Bool()` — random boolean generation.

---

## [v0.0.17] — 2025-04-14

### Changed
- Migrated from `math/rand` to **`math/rand/v2`** (Go 1.22+), which provides automatic seed initialisation — manual seeding removed.
- Random number generation functions further simplified following the new RNG API.

### Added
- Benchmarks for date and number generation functions.

### Removed
- `go.work.sum` from `.gitignore`.

---

## [v0.0.16] — 2025-04-10

### Changed
- `Between` function parameters changed from `int` to **unsigned integer types** to better reflect valid range semantics.
- Internal random number generation logic simplified and randomness quality improved.

---

## [v0.0.15] — 2025-04-05

### Changed
- Random number generation functions refactored to enhance randomness and simplify value extraction.
- Test helpers refactored to use a shared `loop` constant (1 000 iterations) for consistent randomness validation.

---

## [v0.0.14] — 2025-04-05

### Changed
- All random number generation functions migrated to a new random source for improved randomness distribution.

---

## [v0.0.13] — 2025-04-05

### Changed
- Range-handling logic in number generation functions simplified and made consistent.
- Date generation functions refactored to simplify random date range handling.

---

## [v0.0.12] — 2025-04-04

### Changed
- `init` function simplified: date range constants now computed once at startup.
- `Complex64` and `Complex128` refactored to accept `min`/`max` parameters.
- Overflow handling in integer generation functions improved.

---

## [v0.0.11] — 2025-04-04

### Changed
- Date functions refactored to share a common `randomUnixBetween` helper, improving readability and reducing duplication.

---

## [v0.0.10] — 2025-04-04

### Added
- `StringBetween(minLen, maxLen, sourceChars)` — generates a string of variable length drawn from a given character set.

---

## [v0.0.9] — 2025-04-04

### Added
- `Byte()` — generates a random `byte` value.

---

## [v0.0.8] — 2025-04-04

### Added
- `Makefile` with targets for formatting, linting, benchmarking, and running tests.
- GitHub Actions CI workflow covering lint, format check, and test stages.
- README badges (build status, coverage).
- Usage examples for string generation in the README.

### Changed
- `String` parameter renamed from `size` to `length` for clarity.
- Random number generation functions simplified internally.
- Test assertions simplified to operate directly on `rune` values.

---

## [v0.0.7] — 2025-03-25

### Changed
- Project and package renamed to **gengo**.
- README updated with usage example.

---

## [v0.0.6] — 2025-03-23

### Added
- `Date()`, `UnixDate()`, and `DateBetween(start, end)` — random `time.Time` generation.

### Changed
- Major refactoring of method signatures and package names for a more convenient public API.

---

## [v0.0.5] — 2025-03-23

### Fixed
- Incorrect values generated for `UInt8`, `UInt16`, `UInt32`, and `UInt64` functions.
- Added tests to validate unsigned integer generation across the full range.

---

## [v0.0.4] — 2025-03-23

### Added
- `Word()` and word-length category types — generates random alphabetic sequences categorised by length (Small / Medium / Big).

---

## [v0.0.3] — 2025-03-23

### Changed
- `String` now evaluates `len(sourceChars)` once before the generation loop instead of on every iteration.

---

## [v0.0.2] — 2025-03-21

### Added
- Number generation functions for all integer and float types: `Int8`, `Int16`, `Int32`, `Int64`, `Int`, `UInt8`, `UInt16`, `UInt32`, `UInt64`, `Float32`, `Float64`, and their `Between` variants.

---

## [v0.0.1] — 2025-03-20

### Added
- Initial release.
- `String(length, sourceChars)` — core random string generation function.
- Predefined character set constants: `AllChars`, `Alphanumeric`, `Alphabetic`, `Numeric`, `Hexadecimal`, `Symbols`.
- Basic benchmark tests.
