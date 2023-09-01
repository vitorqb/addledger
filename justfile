# Variables
GOLANGCILINT := env_var_or_default("GOLANGCILINT", "golangci-lint")
GO := env_var_or_default("GO", "go")

# Bin folder
BIN := join(justfile_directory(), "bin")

# Delve executable
DELVE_VERSION := env_var_or_default("DELVE_VERSION", "1.21.0")
DELVE := env_var_or_default("DELVE", join(BIN, "dlv-" + DELVE_VERSION))

# Mockgen
MOCKGEN_VERSION := env_var_or_default("MOCKGEN_VERSION", "1.6.0")
MOCKGEN := env_var_or_default("MOCKGEN", join(BIN, "mockgen-" + MOCKGEN_VERSION))

PATH := env_var("PATH")

#
# Developing
#

# Set's up developer workspace
setup:
    mkdir -p out
    mkdir -p bin
    ./scripts/setup-asdf.sh
    ./scripts/setup-envfile.sh
    just install-delve
    just install-mockgen
    touch out/destfile

# Runs the app
run:
    {{GO}} run cmd/addledger/main.go

# Runs all tests
test target="./...": mocks
    {{GO}} test {{target}}

# Creates all mocks
mocks:
    MOCKGEN={{MOCKGEN}} {{GO}} generate --run=MOCKGEN -x ./...

# Lints the code
lint:
    {{GOLANGCILINT}} run ./...

# Formats the code
format:
    {{GO}} fmt ./...

# Calls go tidy
tidy:
    {{GO}} mod tidy

# Starts delve
debug-init: install-delve
    echo "Run 'just debug-attach' in another terminal..."
    {{DELVE}} debug --headless --listen localhost:4040 ./cmd/addledger/

# Attaches to started delve
debug-connect:
    {{DELVE}} connect :4040

#
# Installers
#
install-delve:
    GO={{GO}} BIN={{BIN}} DELVE={{DELVE}} DELVE_VERSION={{DELVE_VERSION}} ./scripts/install-delve.sh

install-mockgen:
    GO={{GO}} BIN={{BIN}} MOCKGEN={{MOCKGEN}} MOCKGEN_VERSION={{MOCKGEN_VERSION}} ./scripts/install-mockgen.sh
