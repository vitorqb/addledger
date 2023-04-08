# Variables
GOLANGCILINT := env_var_or_default("GOLANGCILINT", "golangci-lint")
GO := env_var_or_default("GO", "go")

# Set's up developer workspace
setup:
    ./scripts/setup-asdf.sh

run:
    {{GO}} run cmd/addledger/main.go

test:
    {{GO}} test ./...

lint:
    {{GOLANGCILINT}} run ./...

format:
    {{GO}} fmt ./...

tidy:
    {{GO}} mod tidy
