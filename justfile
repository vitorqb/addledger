# Variables
GOLANGCILINT := env_var_or_default("GOLANGCILINT", "golangci-lint")
GO := env_var_or_default("GO", "go")

# Bin folder
BIN := join(justfile_directory(), "bin")

# Delve executable
DELVE_VERSION := env_var_or_default("DELVE_VERSION", "1.22.1")
DELVE := env_var_or_default("DELVE", join(BIN, "dlv-" + DELVE_VERSION))

# Mockgen
MOCKGEN_VERSION := env_var_or_default("MOCKGEN_VERSION", "1.6.0")
MOCKGEN := env_var_or_default("MOCKGEN", join(BIN, "mockgen-" + MOCKGEN_VERSION))

# goreleaser
GORELEASER_VERSION := env_var_or_default("GORELEASER_VERSION", "1.20.0")
GORELEASER := env_var_or_default("GORELEASER", join(BIN, "goreleaser-" + GORELEASER_VERSION))

# semver
SEMVER_VERSION := env_var_or_default("SEMVER_VERSION", "3c76a6f9d113f4045f693845131185611a62162e")
SEMVER := env_var_or_default("SEMVER", join(BIN, "semver-" + SEMVER_VERSION + ".sh"))

# helpers
PATH := env_var("PATH")

# known paths
resources := join(justfile_directory(), "resources")


#
# Developing
#

# Set's up developer workspace. Must be ran once per runtime.
setup:
    mkdir -p out bin    
    touch out/destfile
    ./scripts/setup-asdf.sh
    ./scripts/setup-envfile.sh

# Runs the app with a specific envfile and arguments
run envfile="configs/default.env" args="":
    #!/bin/bash
    echo "Running app with envfile: {{envfile}} and args: {{args}}"
    set -a; source {{envfile}}; set +a
    cmd="{{GO}} run cmd/addledger/main.go {{args}}"
    echo "Running command: ${cmd}"
    eval "${cmd}"

# Runs all tests
test target="./...": mocks
    {{GO}} test {{target}}

# Creates all mocks
mocks: install-mockgen
    rm -rf ./mocks
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
debug-init envfile="configs/default.env": install-delve
    #!/bin/bash
    echo "Starting delve with envfile: {{envfile}}"
    set -a; source {{envfile}}; set +a
    echo "Run 'just debug-attach' in another terminal..."
    {{DELVE}} debug --headless --listen localhost:4040 ./cmd/addledger/

# Attaches to started delve
debug-connect: install-delve
    {{DELVE}} connect :4040

# Builds the app
build args="": install-goreleaser
    {{GORELEASER}} build {{args}}

# Releases a new version
_release bumpType="": install-semver install-goreleaser setup-github-token
    SEMVER={{SEMVER}} GORELEASER={{GORELEASER}} ./scripts/release.sh {{bumpType}}

release-patch: (_release "patch")
release-minor: (_release "minor")
release-major: (_release "major")

# Sets up the github token
setup-github-token:
    ./scripts/setup-github-token.sh

#
# Installers
#
install-delve:
    GO={{GO}} BIN={{BIN}} DELVE={{DELVE}} DELVE_VERSION={{DELVE_VERSION}} ./scripts/install-delve.sh

install-mockgen:
    GO={{GO}} BIN={{BIN}} MOCKGEN={{MOCKGEN}} MOCKGEN_VERSION={{MOCKGEN_VERSION}} ./scripts/install-mockgen.sh

install-goreleaser:
    GO={{GO}} BIN={{BIN}} GORELEASER={{GORELEASER}} GORELEASER_VERSION={{GORELEASER_VERSION}} ./scripts/install-goreleaser.sh

install-semver:
    SEMVER={{SEMVER}} SEMVER_VERSION={{SEMVER_VERSION}} ./scripts/install-semver.sh
