# TANGLE Protocol Reference — Alpha

## Transaction

A transaction is represented by a sender, recipient, asset, amount, nonce, expiry/commitment domain and routing target. The alpha stores public addresses for demonstration. The production protocol must replace public amount/address handling with the confidential commitment + encrypted payload model specified by the whitepaper.

## Fragmentation

The protocol samples a final fragment target uniformly in [100,1000]. It then recursively processes a queue of fragments. Each local operation is PASS, SPLIT or POOL. SPLIT creates 2–10 children in this alpha. Production must preserve exact value conservation across child commitments.

## Vaults

A Vault is a routing object with a TGL stake. It is not liquidity. Vault work is counted from accepted routing transitions. The target reward rule is Work × sqrt(Stake), normalized across eligible Vaults.

## Pools

Pool is a protocol state, not an external mixer. Production Pool residence must be represented in block heights and selected by protocol randomness. Alpha records the POOL transition but does not implement real confidential mixing.

## Proof of Routing

The alpha produces a deterministic proof hash bound to task, input commitment, output commitment, operation and next hop. Production PoR must be a formal ZK relation and must not be replaced by this hash construction.

## Nullifiers

The state contains a NullifierSet. Production nullifiers must be derived from the transaction secret/domain as defined in the whitepaper and verified atomically in the state transition.

## Settlement

The alpha credits the recipient after routing. Production settlement must verify all final fragments, conservation, fee accounting, nullifiers and expiry before crediting the recipient.
