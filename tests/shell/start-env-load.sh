#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cp "$REPO_ROOT/start.sh" "$TMP_DIR/start.sh"
mkdir -p "$TMP_DIR/scripts"
cp "$REPO_ROOT/scripts/load-env.sh" "$TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TMP_DIR/scripts/start-process-utils.sh"
mkdir -p "$TMP_DIR/frontend"
cat > "$TMP_DIR/.env" <<'EOF'
PORT=3333
FRONTEND_DEV_HOST=192.168.1.50
EOF

output=$(cd "$TMP_DIR" && bash -c 'source ./start.sh >/dev/null 2>&1; printf "%s|%s" "$PORT" "$FRONTEND_DEV_HOST"' || true)

if [[ "$output" != "3333|192.168.1.50" ]]; then
  printf 'expected env values to load from .env, got %s\n' "$output" >&2
  exit 1
fi
