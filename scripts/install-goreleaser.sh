#!/bin/bash
if ! [ -f "$GORELEASER" ]
then
    echo "Installing goreleaser..."
    mkdir -p "$BIN"
    GOBIN="$BIN" "$GO" install "github.com/goreleaser/goreleaser@v$GORELEASER_VERSION"
    mv "$BIN"/goreleaser "$GORELEASER"
else
    echo "goreleaser already installed."
fi
