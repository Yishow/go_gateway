#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

PORT=3333
TMP_DIR="$TEST_TMP_DIR/bin/tmp"
AIR_BIN="$TEST_TMP_DIR/bin/tmp/gateway-air"
APP_PATH="cmd/test_ui"
FRONTEND_DIR="frontend"

mkdir -p "$TMP_DIR/start-${PORT}.lock"
printf '%s\n' "$$" >"$TMP_DIR/start-${PORT}.lock/pid"

err() { printf '%s\n' "$*"; }
info() { printf '%s\n' "$*"; }
warn() { printf '%s\n' "$*"; }

# shellcheck disable=SC1091
source "$REPO_ROOT/scripts/start-process-utils.sh"

set +e
output="$(acquire_run_lock 2>&1)"
exit_code=$?
set -e

if [[ "$exit_code" -eq 0 ]]; then
  printf 'expected acquire_run_lock to fail when lock owner is alive, got success\n' >&2
  exit 1
fi

if [[ "$output" != *"已有 start.sh 管理此 repo/port"* ]]; then
  printf 'expected acquire_run_lock to explain active lock, got:\n%s\n' "$output" >&2
  exit 1
fi

if [[ "$output" == *"unbound variable"* ]]; then
  printf 'expected acquire_run_lock to avoid unbound variable, got:\n%s\n' "$output" >&2
  exit 1
fi