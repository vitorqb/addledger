# Variables
GOLANGCILINT := env_var_or_default("GOLANGCILINT", "golangci-lint")
GO := env_var_or_default("GO", "go")

# Set's up developer workspace
setup:
    mkdir -p out
    touch out/destfile
    ./scripts/setup-asdf.sh
    ./scripts/setup-envfile.sh

run:
    {{GO}} run cmd/addledger/main.go

test: mocks
    {{GO}} test ./...

mocks:
    {{GO}} generate --run=mockgen -x ./...

lint:
    {{GOLANGCILINT}} run ./...

format:
    {{GO}} fmt ./...

tidy:
    {{GO}} mod tidy
