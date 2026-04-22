#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

mkdir -p "$TMP_DIR/scripts"
cat > "$TMP_DIR/.env" <<'EOF'
PORT=3333
FRONTEND_DEV_HOST=192.168.1.50
FRONTEND_DEV_PORT=4173
EOF

output=$(REPO_ROOT="$TMP_DIR" bash -c 'source "$1"; printf "%s|%s|%s" "$PORT" "$FRONTEND_DEV_HOST" "$FRONTEND_DEV_PORT"' _ "$REPO_ROOT/scripts/load-env.sh" || true)

if [[ "$output" != "3333|192.168.1.50|4173" ]]; then
  printf 'expected env values to load from .env, got %s\n' "$output" >&2
  exit 1
fi
