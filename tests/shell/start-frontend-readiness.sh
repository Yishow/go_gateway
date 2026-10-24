#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
mkdir -p "$TMP_DIR/scripts"
cp "$REPO_ROOT/start.sh" "$TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TMP_DIR/scripts/start-process-utils.sh"
cd "$TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

pid_on_port() { :; }
kill() { return 1; }
set +e
output="$(wait_for_port_ready 4173 前端 9876 120 2>&1)"
status=$?
set -e

if [[ "$status" -eq 0 || "$output" != *"前端 未成功啟動"* ]]; then
  printf 'expected readiness failure when frontend exits, status=%s output=%s\n' "$status" "$output" >&2
  exit 1
fi
