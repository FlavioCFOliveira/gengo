# Contributing to gengo

Thanks for your interest in improving **gengo**! This guide covers how to set up
your environment, the checks your change must pass, and the conventions the
project follows.

## Prerequisites

- **Go 1.26 or later** (see [`go.mod`](go.mod)).
- A working `git` and `make`.

## Getting Started

1. Fork the repository and clone your fork.
2. Create a topic branch for your change:
   ```bash
   git checkout -b fix/short-description
   ```
3. Install the development tools (linter and vulnerability scanner):
   ```bash
   make dev-install
   ```
   This installs `golangci-lint` and `govulncheck` into your `GOBIN`.

## Build, Test, and Lint

Run these before opening a pull request — they mirror the CI pipeline:

```bash
make test                 # run the test suite (go test ./...)
go test -race ./...       # run the tests under the race detector
make check                # gofmt -l . and golangci-lint run
make vulncheck            # govulncheck ./...
make coverage             # generate and open the HTML coverage report
make bench                # run benchmarks (go test -benchmem -bench, 5s)
```

Run `make help` to see every available target.

A change is ready when:

- `make test` and `go test -race ./...` pass.
- `make check` reports no formatting or lint issues.
- `make vulncheck` reports no vulnerabilities.

## Coding Conventions

gengo is a small, focused library; please keep contributions aligned with its
philosophy:

- **Zero dependencies** — use only the Go standard library. New third-party
  runtime dependencies are not accepted.
- **Performance matters** — primitive generators should avoid heap allocations
  (`0 allocs/op`). Validate performance claims with `make bench`, and include
  before/after numbers from a single machine when changing a hot path.
- **Idiomatic Go** — code must be `gofmt`-clean and pass `golangci-lint`.
- **Documentation** — every exported identifier needs a doc comment that starts
  with its name and ends with a period (enforced by the linter). Keep the
  `README.md` and the godoc accurate and in sync with behaviour.
- **Tests** — add tests for new functionality and bug fixes; the suite is kept
  at 100% statement coverage, so cover new branches (run `make coverage` to
  check). gengo is **not** cryptographically secure — never present it as
  suitable for passwords, tokens, or other secrets (see [SECURITY.md](SECURITY.md)).

## Commit Messages

Write clear, self-contained commit messages in the imperative mood. The project
uses [Conventional Commits](https://www.conventionalcommits.org/)-style subjects:

```
type(scope): short summary

Optional body explaining what changed and why.
```

Common types: `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `chore`.
Example: `fix(numbers): prevent Float64Between from returning +Inf`.

## Pull Requests

1. Make sure all checks above pass locally.
2. Keep each pull request focused on a single change.
3. Update documentation (`README.md`, godoc, `CHANGELOG.md`) when behaviour
   changes. Add your entry under the `## [Unreleased]` section of
   [`CHANGELOG.md`](CHANGELOG.md) following the existing categories.
4. Describe what changed and why in the pull request, and link any related
   issue.

## Versioning

gengo follows [Semantic Versioning](https://semver.org/). See the
[Versioning Policy](CHANGELOG.md#versioning-policy) in the changelog for how
features, fixes, and breaking changes map to version bumps.

## Reporting Security Issues

Please do not open public issues for security vulnerabilities. Follow the process
in [SECURITY.md](SECURITY.md) instead.
