# TANGLE Testnet Alpha

Reference implementation of the TANGLE protocol mechanics described in the TANGLE whitepaper. This repository is a **public-testnet alpha/reference implementation**, not a production-secure mainnet.

## What is implemented

- Native TGL accounting, 10,000,000 max supply.
- 0.01% protocol fee, charged once per transaction.
- Fee split: 30% Vaults / 30% Validators / 30% Treasury / 10% Insurance.
- Routing Vault registration and TGL stake accounting.
- Progressive recursive fragmentation with a random final target in [100,1000].
- PASS / SPLIT / POOL transitions; SPLIT creates 1–10 children.
- Protocol-selected next hops and randomized pool residence in blocks.
- Confidential-style commitments and deterministic transition hashes.
- Nullifier protection and conservation checks.
- Proof-of-Routing records with cryptographic hashes and validator verification.
- Vault rewards based on accepted routing work and sqrt(stake) weighting.
- Validator reward accounting.
- Settlement to recipient.
- JSON persistence so a node restart keeps chain state.
- HTTP API and mobile-first frontend.
- Faucet endpoint for testnet TGL.
- Docker and local multi-node launch scripts.
- GitHub Actions CI.

## Important security boundary

This alpha deliberately does **not** claim production cryptographic privacy, BFT consensus, bridge security, or production Nova proving. The PoR object is a deterministic cryptographic reference proof record, not a replacement for a formally specified ZK circuit. The node uses a deterministic block-production loop so the entire protocol can be demonstrated without pulling a large external consensus dependency into the first public repository. The next hardening stage is CometBFT/Cosmos SDK integration, formal ZK circuits, real P2P networking, validator signatures/slashing, and independent security review.

The architecture is designed to map onto Cosmos SDK + CometBFT for the production L1. Cosmos SDK is a modular framework for sovereign application-specific L1s and CometBFT provides BFT state-machine replication. See the official projects linked in `docs/PRODUCTION_PATH.md`.

## Quick start

Requirements: Go 1.23+.

```bash
go run ./cmd/tangled --data ./data --http :8080
```

Open `http://localhost:8080/`.

Faucet:

```bash
curl -X POST http://localhost:8080/api/faucet \
  -H 'content-type: application/json' \
  -d '{"address":"alice"}'
```

Create a Vault:

```bash
curl -X POST http://localhost:8080/api/vaults \
  -H 'content-type: application/json' \
  -d '{"owner":"alice","stake":1000}'
```

Submit a transfer:

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H 'content-type: application/json' \
  -d '{"from":"alice","to":"bob","asset":"TGL","amount":10000}'
```

Inspect:

```bash
curl http://localhost:8080/api/status
curl http://localhost:8080/api/blocks
curl http://localhost:8080/api/vaults
```

## Production path

1. Replace the deterministic block loop with CometBFT consensus.
2. Move state transitions into Cosmos SDK modules.
3. Specify and audit the real commitment scheme and confidential asset model.
4. Implement formal PoR circuits and Nova-style folding/IVC.
5. Add authenticated routing tasks and validator signatures.
6. Add P2P routing transport and anti-DDoS controls.
7. Implement IBC/chain-specific proof adapters only after the native L1 is stable.
8. Run adversarial tests, fuzzing, economic simulations, cryptographic review, and an independent audit.
9. Only then promote a network to a public-value mainnet.
