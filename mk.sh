#!/bin/sh

name="$1"

if [ "$name" = "run" ]; then
    echo "Running test debug model..."
    go run src/*.go
else
    echo "Build for bin/bhyveapid"
    gofmt -w src/*.go
    go build -o bin/bhyveapid src/*.go
fi

