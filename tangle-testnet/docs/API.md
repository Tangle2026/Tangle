# TANGLE Testnet Alpha API

Base URL: `http://HOST:8080`

## GET /api/status
Returns chain height, supply and object counts.

## GET /api/blocks
Returns blocks produced by the alpha node.

## GET /api/vaults
Returns registered routing Vaults.

## GET /api/account/{address}
Returns an account balance and nonce.

## GET /api/tx/{id}
Returns transaction state.

## POST /api/faucet
Body:
```json
{"address":"alice","amount":1000000000}
```

## POST /api/vaults
Body:
```json
{"owner":"alice","stake":1000}
```

## POST /api/transactions
Body:
```json
{"from":"alice","to":"bob","asset":"TGL","amount":100000}
```

The alpha charges 0.01% once per transaction and allocates 30/30/30/10 across Vaults/Validators/Treasury/Insurance.

## POST /api/mine
Forces a block attempt. The node also attempts a block automatically every 3 seconds.

## GET /api/health
Returns node health.
