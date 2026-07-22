# Changelog

All notable changes to **gengo** are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) conventions,
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## Versioning Policy

gengo follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **MAJOR** — incompatible (breaking) API changes.
- **MINOR** — new, backward-compatible functionality.
- **PATCH** — backward-compatible bug fixes.

The project is in the **0.x** series. Per SemVer, while the major version is `0`
the public API is not yet guaranteed stable: backward-compatible features and
fixes bump the **minor** and **patch** numbers respectively, and an incompatible
change may still occur in a minor release before `1.0.0`. Starting with
**v0.1.0**, releases follow this policy; the earlier run of `0.0.x` tags — which
shipped features as patches and included the no-op re-tags `v0.0.20`–`v0.0.25`
(all pointing at the same commit) — predates it.

Each version is tagged exactly once; tags are never moved or re-pointed. A
mistaken or empty release is superseded by a new, higher version, never by
re-tagging an existing one.

---

## [v0.2.0] — 2026-07-22

A **minor** release under the Versioning Policy above. This cycle adds
**WordsPT**, a European Portuguese (pt-PT) word generator. Every change is
additive: the existing public API is unchanged, so upgrading from `v0.1.0`
requires no code changes.

### Added
- **WordsPT — European Portuguese (pt-PT) word generation.** A syllabic engine
  that produces morphologically inflected pseudo-words for the open word classes
  and draws from curated single-word lists for the closed classes. It is
  correct-by-construction, applies full pt-PT graphic accentuation (Acordo
  Ortográfico), performs roughly one allocation per generated word, and adds no
  external dependencies. The public API adds 17 functions, each mirrored by a
  `*Generator` method for reproducible, seeded output:
  - Top-level orchestration — `WordPT`, `WordPTByLengthType`, and `WordsPT`.
  - Open classes (morphologically inflected):
    - `NounPT` / `NounPTOf` — nouns with gender and number inflection.
    - `AdjectivePT` / `AdjectivePTOf` — adjectives, including the synthetic
      absolute superlative (`-íssimo`).
    - `VerbPT` / `VerbPTOf` — regular verb conjugation across all moods, tenses,
      and persons (indicative, subjunctive, imperative, and the non-finite
      forms); no *vós*, and the pt-PT `falámos` vs. `falamos` distinction is
      preserved.
    - `AdverbPT` / `AdverbPTByLengthType` — productive `-mente` adverbs.
  - Closed classes (curated single words) — `ArticlePT`, `PronounPT`,
    `NumeralPT`, `PrepositionPT`, `ConjunctionPT`, and `InterjectionPT`.
- Inflection enums, each with an `Any…` zero value meaning "choose a valid value
  at random": `Gender`, `Number`, `Degree`, `Mood`, `Tense`, and `Person`.
- `AnyLengthWord` — the zero value of `LengthTypeWords`, selecting the full
  length range; the length counterpart of the `Any…` inflection zero values.
- Functional specification for WordsPT under `specification/`, and a project
  knowledge model (`knowledge-model.md`).

### Changed
- `README.md` gained a WordsPT section with runnable examples; `TEST_REPORT.md`
  and `CLAUDE.md` were updated. No public API was changed.

---

## [v0.1.0] — 2026-05-29

The first release under the Versioning Policy above. This cycle bumps the
**minor** version. Per the Versioning Policy, while still in `0.x` it carries one
incompatible API change (the `UInt*` → `Uint*` rename below) alongside
backward-compatible additions and fixes.

### Added
- `Generator` type with `New(seed)` and `NewSource(rand.Source)` constructors for
  reproducible, seeded output; its methods mirror every package-level function.
- `SECURITY.md` — security policy with a non-cryptographic disclaimer and a
  vulnerability-reporting process.
- `CONTRIBUTING.md` — build, test, and lint workflow, coding and commit
  conventions, and the pull-request process.

### Fixed
- `Float32Between` / `Float64Between` no longer return `+Inf` for very wide
  ranges; results are always finite and within `[min, max]`.

### Changed
- **BREAKING:** renamed the unsigned-integer API from `UInt*` to Go's idiomatic
  `Uint*` spelling — `Uint8`, `Uint16`, `Uint32`, `Uint64`, their `Between`
  variants, and the matching `Generator` methods. This matches the standard
  library (`math.MaxUint64`, `math/rand/v2.Uint64`, `sync/atomic.Uint64`). The
  old `UInt*` names were removed; update call sites accordingly (for example,
  `gengo.UInt8Between` → `gengo.Uint8Between`).
- Documentation corrected to match the implementation: `Words` length
  distribution (26/52/22), `Float32`/`Float64` positive-only range, and the
  `String` performance label (8-char input).
- `String` documented as byte/ASCII-oriented (multibyte charsets yield invalid
  UTF-8); the broken README emoji example was replaced.
- `DateBetween` whole-second granularity documented.
- `BENCHMARKS.md` and `TEST_REPORT.md` regenerated on a single reference machine
  so figures are consistent and reproducible (resolves the contradictory
  `BenchmarkString` numbers).

### Security
- Removed insecure password / API-token / session-ID recommendations from the
  README and added a prominent warning to use `crypto/rand` for secrets, since
  gengo is built on the non-cryptographic `math/rand/v2`.

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
