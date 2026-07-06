.PHONY: build web dev test lint

# Build the binary with a fresh web bundle embedded.
build: web
	go build -o tama ./cmd/tama

# Rebuild the frontend into web/dist (committed, so plain go build works).
web:
	cd web && npm ci && npm run build

# Dev loop: Vite serves the hot-reloading UI, the Go server proxies to it.
dev:
	@trap 'kill 0' EXIT INT TERM; \
	(cd web && npm run dev) & \
	go run ./cmd/tama serve --dev

test:
	go test -race ./...
	cd web && npm test

lint:
	@unformatted=$$(gofmt -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "These files need gofmt -s -w:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	go vet ./...
	cd web && npx tsc --noEmit
