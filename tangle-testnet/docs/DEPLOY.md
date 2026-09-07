# Testnet deployment

## Docker

```bash
docker compose up --build -d
```

Then open port 8080.

## Native

```bash
./scripts/build.sh
./scripts/start-testnet.sh
```

For a public VPS, put the node behind a reverse proxy with TLS, restrict the faucet endpoint, add rate limiting, backups, metrics and a firewall. Do not expose the development faucet with unrestricted minting on a value-bearing network.

## GitHub

Push the repository to a GitHub repository. GitHub Actions will run `go test ./...` and build the node on every push/PR.
