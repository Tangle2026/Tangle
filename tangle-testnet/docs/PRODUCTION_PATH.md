# Production path

The current repository is a self-contained alpha so developers can run the TANGLE state model immediately. It is intentionally dependency-light. Production TANGLE should use the Cosmos SDK application framework and CometBFT consensus rather than treating this reference block loop as BFT consensus.

The intended production stack is:

TANGLE L1 -> Cosmos SDK application modules -> CometBFT consensus -> native TGL -> Routing/Vault/Fragment/Pool/Nullifier/Settlement modules -> formal ZK/PoR subsystem -> Application Layer -> Interoperability Layer.

See official repositories:
- Cosmos SDK: https://github.com/cosmos/cosmos-sdk
- CometBFT: https://github.com/cometbft/cometbft

As of September 2026, the Cosmos SDK repository lists v0.55.0 as its latest release, and CometBFT lists v0.39.3 as its latest release. Pin exact releases during the production implementation rather than depending on `main`.
