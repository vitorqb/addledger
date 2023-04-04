#!/bin/bash
# See https://asdf-vm.com/

# Check asdf installed
if ! command -v asdf &> /dev/null
then
    echo "asdf could not be found: see https://asdf-vm.com/"
    exit 1
fi

# Add needed repositories
asdf plugin add golangci-lint

# Install
asdf install
