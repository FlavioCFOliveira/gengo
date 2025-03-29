bench:
	go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s .

check:
	gofmt -l .
	golangci-lint run

dev-install:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

format:
	gofmt -w .

test:
	go test .

vulncheck:
	govulncheck ./...


.PHONY: bench check dev-install format test vulncheck
