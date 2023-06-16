#!/bin/bash
if ! [ -f "$DELVE" ]
then
    echo "Installing delve..."
    mkdir -p "$BIN"
    GOBIN="$BIN" "$GO" install "github.com/go-delve/delve/cmd/dlv@v$DELVE_VERSION"
    mv "$BIN"/dlv "$DELVE"
else
    echo "Delve already installed."
fi
