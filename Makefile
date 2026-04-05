BENCH_FLAGS := -benchmem -run=^$$ -bench ^Benchmark -benchtime 5s

## help: print this help message
.PHONY: help
help:
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## bench: run benchmarks
.PHONY: bench
bench:
	go test $(BENCH_FLAGS) .

## check: run formatter check and linter
.PHONY: check
check:
	gofmt -l .
	golangci-lint run

## coverage: run tests and open HTML coverage report
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## dev-install: install development tools
.PHONY: dev-install
dev-install:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## format: format source code with gofmt
.PHONY: format
format:
	gofmt -w .

## test: run tests
.PHONY: test
test:
	go test ./...

## vulncheck: run vulnerability scanner
.PHONY: vulncheck
vulncheck:
	govulncheck ./...
