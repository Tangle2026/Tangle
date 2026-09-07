# Interoperability Layer — Alpha Contract

This directory reserves the production interfaces for:

1. Asset Gateway: lock/mint/burn/release accounting with source-chain proofs.
2. Message Bridge: authenticated cross-chain messages with replay protection.
3. Chain-specific Proof Adapters: Ethereum/Cosmos/other light-client or ZK verification.
4. Privacy Adapter: lets an external application delegate private state/routing to TANGLE.

No external bridge is enabled in the alpha. This is deliberate: a bridge must not be presented as trustless until its verifier, finality, replay, reorg and emergency model are specified and audited.
