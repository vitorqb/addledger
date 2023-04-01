# Variables
GOLANGCILINT := env_var_or_default("GOLANGCILINT", "golangci-lint")
GO := env_var_or_default("GO", "go")

# recipts
run:
    {{GO}} run cmd/addledger/main.go

lint:
    {{GOLANGCILINT}} run ./...

format:
    {{GO}} fmt ./...
