#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_TMP_DIR/scripts"
cp "$REPO_ROOT/start.sh" "$TEST_TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_TMP_DIR/scripts/load-env.sh"

cd "$TEST_TMP_DIR"
REPO_ROOT="$TEST_TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

BACKEND_PID=222
BACKEND_LOG_PID=111
BACKEND_LOG_FILE="$TEST_TMP_DIR/backend.log"
touch "$BACKEND_LOG_FILE"

kill() {
  return 0
}

wait_calls=()
wait() {
  wait_calls+=("$1")
  return 0
}

info() { :; }
success() { :; }

cleanup_backend

actual_waits="$(printf '%s|' "${wait_calls[@]}")"
if [[ "$actual_waits" == *"111"* ]]; then
  printf 'expected cleanup_backend not to wait on BACKEND_LOG_PID during shutdown, got waits %s\n' "$actual_waits" >&2
  exit 1
fi
