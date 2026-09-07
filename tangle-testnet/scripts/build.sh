#!/usr/bin/env bash
set -euo pipefail
mkdir -p bin
go test ./...
go build -trimpath -ldflags='-s -w' -o bin/tangled ./cmd/tangled
echo "Built bin/tangled"
