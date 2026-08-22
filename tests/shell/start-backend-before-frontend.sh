#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_TMP_DIR/scripts" "$TEST_TMP_DIR/frontend" "$TEST_TMP_DIR/cmd/test_ui"
cp "$REPO_ROOT/start.sh" "$TEST_TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TEST_TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_TMP_DIR/scripts/start-process-utils.sh"

cd "$TEST_TMP_DIR"
REPO_ROOT="$TEST_TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

APP_PATH="cmd/test_ui"
FRONTEND_DIR="frontend"
# sync-embed 改為 opt-in 後，測試需顯式啟用才會觸發 build_frontend 步驟
SYNC_EMBED=true

order=()
record_step() {
  order+=("$1")
}

build_frontend() {
  record_step "build"
}

clear_port() {
  record_step "clear-backend"
}

start_backend_process() {
  record_step "start-backend"
}

wait_for_port_ready() {
  record_step "wait-backend"
}

start_frontend_dev_server() {
  record_step "start-frontend"
  sleep 0.01 &
}

wait_for_backend_exit() {
  record_step "wait-exit"
}

show_frontend_host_hint() { :; }
show_embedded_frontend_hint() { :; }
show_log_noise_summary() { :; }
cleanup_all() { :; }
runtime_log_line() { :; }
trap() { :; }
go() { :; }

set +e
start_dev_mode >/dev/null 2>&1
set -e

actual_order="$(IFS='>'; printf '%s' "${order[*]}")"
expected_order="build>clear-backend>start-backend>wait-backend>start-frontend>wait-exit"

if [[ "$actual_order" != "$expected_order" ]]; then
  printf 'expected start_dev_mode order %s, got %s\n' "$expected_order" "$actual_order" >&2
  exit 1
fi
