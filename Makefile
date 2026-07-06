.PHONY: build web dev test lint

# Build the binary with the current web/dist embedded.
build:
	go build -o tama ./cmd/tama

# Rebuild the frontend into web/dist (committed, so plain go build works).
web:
	cd web && npm install && npm run build

# Dev loop: run `make dev` in one terminal and `cd web && npm run dev` in
# another; Vite proxies /api to the Go server.
dev:
	go run ./cmd/tama serve

test:
	go test -race ./...
	cd web && npm run typecheck

lint:
	gofmt -s -l . && go vet ./...
