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

TMP_DIR="bin/tmp"
mkdir -p "$TMP_DIR"

tail() {
  while [[ $# -gt 0 ]]; do
    shift
  done
}

bash() {
  return 0
}

start_backend_process "true"

if [[ -z "${BACKEND_LOG_FILE:-}" ]]; then
  printf 'expected start_backend_process to set BACKEND_LOG_FILE\n' >&2
  exit 1
fi
if [[ ! -f "$BACKEND_LOG_FILE" ]]; then
  printf 'expected start_backend_process to create backend log file, got %s\n' "$BACKEND_LOG_FILE" >&2
  exit 1
fi
if [[ "$(basename "$BACKEND_LOG_FILE")" == "backend.XXXXXX.log" ]]; then
  printf 'expected start_backend_process to create a unique backend log file name, got %s\n' "$BACKEND_LOG_FILE" >&2
  exit 1
fi
