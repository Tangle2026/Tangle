#!/usr/bin/env bash
set -euo pipefail
mkdir -p data
exec ./bin/tangled --data ./data/tangle.json --http "${TANGLE_HTTP:-:8080}"
