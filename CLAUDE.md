# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

**gengo** is a Go library designed to provide a simple and practical way to generate random data structures. It offers an intuitive, straightforward API for creating random values of various types (strings, numbers, dates, booleans, words) without complex configuration or boilerplate code.

### Design Philosophy

- **Simplicity First**: One function call to get random data - no initialization, no configuration
- **Practical Defaults**: Sensible defaults that cover 90% of use cases out of the box
- **Intuitive Naming**: Functions follow the pattern `Type()` for full range, `TypeBetween(min, max)` for custom ranges
- **Zero Dependencies**: Only uses Go standard library - no external dependencies to manage

### Performance Philosophy

- **Execution Speed**: Implementations prioritize the fastest execution path, leveraging hardware capabilities where possible
- **Minimal Allocations**: Reduce memory allocations to the absolute minimum - prefer stack over heap
- **Memory Efficiency**: Keep memory footprint low; avoid unnecessary buffering or temporary structures
- **No Bloat**: Avoid abstractions that add overhead; prefer direct, efficient implementations

## Build and Test Commands

```bash
# Run all tests
go test .

# Run a single test by name
go test -run TestInt8 .

# Run benchmarks
make bench
# Or directly:
go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s .

# Format code
go fmt .
# Or:
make format

# Lint code
make check          # Runs gofmt -l and golangci-lint

# Security check
make vulncheck      # Runs govulncheck ./...

# Install dev dependencies
make dev-install    # Installs govulncheck and golangci-lint
```

## Architecture

**gengo** is a Go module for generating random data. It uses `math/rand/v2` (Go 1.22+) which has automatic seed initialization.

### Code Organization

The package follows a flat file structure where each file handles a specific data type category:

- **init.go** - Package initialization with date range constants (`dtMin`, `dtMax`, `dtUnixMin`, `dtUnixMax`) and error variables (`ErrorNilFunc`, `ErrorFuncNilResult`)
- **strings.go** - Random string generation with predefined character set constants (`AllChars`, `Alphanumeric`, `Alphabetic`, `Numeric`, `Hexadecimal`, `Symbols`)
- **numbers.go** - Random number generation for all numeric types (int8/16/32/64, uint8/16/32/64, float32/64, complex64/128). Each type has a `Type()` function for full range and `TypeBetween(min, max)` for custom ranges
- **dates.go** - Random `time.Time` generation with `Date()`, `UnixDate()`, and `DateBetween(start, end)`
- **bool.go** - Boolean generation via `Bool()`
- **words.go** - Word generation with length categories (Small/Medium/Big) via `LengthTypeWords` enum

### Key Patterns

1. **Range Handling**: Functions accepting min/max parameters swap values if `min > max` to ensure valid ranges
2. **Overflow Prevention**: Integer functions use wider types for intermediate calculations (e.g., `int16` for `int8` math)
3. **String Generation**: `String(length, sourceChars)` is the core function; other string functions are convenience wrappers using predefined constants
4. **Testing**: Tests use a `loop` variable (1000 iterations) to validate randomness through repeated execution rather than asserting specific values

### Dependencies

- No external runtime dependencies
- Uses only Go standard library (`math/rand/v2`, `math`, `time`, `errors`)
- Dev tools: `golangci-lint`, `govulncheck`

### CI/CD

GitHub Actions workflow (`.github/workflows/go.yml`) runs on every push:
- Dependency check (`go mod tidy -diff`)
- Vulnerability scan (`govulncheck`)
- Linting (`golangci-lint`)
- Format check (`gofmt`)
- Tests (`go test -v`)

## Git Workflow

**IMPORTANT:** All commits must exclusively be authored by `flaviocfo <flaviocfo@gmail.com>`. Do not include Co-Authored-By tags or any other author attribution in commits. The user will handle commit authorship.

## Commit Policy

**All commits must be authored exclusively by flaviocfo.** When creating commits, I will use your configured git identity. Do not add Co-Authored-By tags or other attribution that would change the commit author from flaviocfo.
