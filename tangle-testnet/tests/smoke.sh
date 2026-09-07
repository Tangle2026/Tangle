#!/usr/bin/env bash
set -euo pipefail
BASE=${BASE:-http://localhost:8080}
curl -fsS "$BASE/api/health" >/dev/null
curl -fsS -X POST "$BASE/api/faucet" -H 'content-type: application/json' -d '{"address":"alice"}' >/dev/null
curl -fsS -X POST "$BASE/api/vaults" -H 'content-type: application/json' -d '{"owner":"alice","stake":1000}' >/dev/null
curl -fsS -X POST "$BASE/api/transactions" -H 'content-type: application/json' -d '{"from":"alice","to":"bob","asset":"TGL","amount":100000}' >/dev/null
curl -fsS -X POST "$BASE/api/mine" >/dev/null
curl -fsS "$BASE/api/status"
