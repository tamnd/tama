#!/bin/sh
# Smoke test: build the binary, seed the demo data, boot the server, and hit
# healthz plus a real login. Finishes well under 30 seconds.
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d)
addr=127.0.0.1:8137
pid=

cleanup() {
    [ -n "$pid" ] && kill "$pid" 2>/dev/null && wait "$pid" 2>/dev/null
    rm -rf "$tmp"
}
trap cleanup EXIT INT TERM

echo "== build"
go build -C "$root" -o "$tmp/tama" ./cmd/tama

echo "== seed"
"$tmp/tama" seed --demo --data "$tmp/data"

echo "== boot"
"$tmp/tama" serve --addr "$addr" --data "$tmp/data" &
pid=$!

up=0
for _ in $(seq 1 50); do
    if curl -fsS "http://$addr/api/healthz" >/dev/null 2>&1; then
        up=1
        break
    fi
    sleep 0.1
done
[ "$up" = 1 ] || { echo "server never came up" >&2; exit 1; }

echo "== healthz"
body=$(curl -fsS "http://$addr/api/healthz")
echo "$body"
case $body in
*'"status":"ok"'*) ;;
*) echo "unexpected healthz body" >&2; exit 1 ;;
esac

echo "== login"
body=$(curl -fsS -X POST "http://$addr/api/auth/login" \
    -H 'Content-Type: application/json' \
    -d '{"username":"demo","password":"demo1234"}')
echo "$body"
case $body in
*'"username":"demo"'*) ;;
*) echo "unexpected login body" >&2; exit 1 ;;
esac

echo "== ok"
