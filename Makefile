.PHONY: fmt lint test

fmt:
	go fmt ./...

lint:
	# Checks if golangci-lint is installed, otherwise falls back to go vet
	if command -v golangci-lint >/dev/null; then golangci-lint run; else go vet ./...; fi

test:
	go test -v ./...
