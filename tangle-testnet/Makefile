.PHONY: run test build docker
run:
	go run ./cmd/tangled --data ./data/tangle.json --http :8080
test:
	go test ./...
build:
	go build -trimpath -o bin/tangled ./cmd/tangled
docker:
	docker compose up --build
