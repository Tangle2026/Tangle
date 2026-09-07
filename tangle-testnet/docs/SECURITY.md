# Security status

This repository is **testnet alpha**. Do not use real funds.

Known non-production components:
- deterministic local block production instead of BFT consensus;
- public/plain transaction amounts in the HTTP alpha API;
- deterministic PoR hash instead of a ZK proof system;
- no production P2P routing transport;
- no audited bridge;
- no production validator key management/slashing;
- simplified pool semantics;
- simplified accounting precision and asset model;
- no independent audit.

The alpha is useful for architecture demos, frontend integration, routing experiments and protocol-state testing. It is not a privacy guarantee.
