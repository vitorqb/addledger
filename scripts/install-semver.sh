#!/bin/bash
URL="https://raw.githubusercontent.com/fsaintjacques/semver-tool/${SEMVER_VERSION}/src/semver"
if ! $SEMVER --version 2>/dev/null 1>/dev/null;
then
    echo "Installing semver"
    curl -sL -o "$SEMVER" "$URL"
    chmod +x "$SEMVER"
    $SEMVER --version
else
    echo "semver already installed"
    $SEMVER --version
fi
    
