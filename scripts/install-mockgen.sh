#!/bin/bash
if ! [ -f "$MOCKGEN" ]
then
    echo "Installing mockgen..."
    mkdir -p "$BIN"
    GOBIN="$BIN" "$GO" install "github.com/golang/mock/mockgen@v$MOCKGEN_VERSION"
    mv "$BIN"/mockgen "$MOCKGEN"
else
    echo "Mockgen already installed."
fi
