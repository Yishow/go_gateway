#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
mkdir -p "$TMP_DIR/scripts" "$TMP_DIR/frontend" "$TMP_DIR/cmd/test_ui"
cp "$REPO_ROOT/start.sh" "$TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TMP_DIR/scripts/start-process-utils.sh"
touch "$TMP_DIR/go.mod" "$TMP_DIR/.air.toml" "$TMP_DIR/frontend/package.json"

cd "$TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

order=()
record() { order+=("$1"); }
validate_dev_environment() { :; }
assert_ports_free() { :; }
ensure_air() { :; }
ensure_frontend_dependencies() { :; }
ensure_embedded_frontend_snapshot() { :; }
install_cleanup_traps() { :; }
acquire_run_lock() { :; }
start_backend_process() { record start-backend; BACKEND_PID=101; }
wait_for_port_ready() { record "wait-$2"; :; }
start_frontend_dev_server() { record start-frontend; FRONTEND_PID=202; }
wait_for_backend_exit() { record wait-exit; }
info() { :; }
success() { :; }

start_dev_mode >/dev/null
actual_order="$(IFS='>'; printf '%s' "${order[*]}")"
expected_order="start-backend>wait-後端>start-frontend>wait-前端>wait-exit"
if [[ "$actual_order" != "$expected_order" ]]; then
  printf 'expected backend readiness before frontend startup, got %s\n' "$actual_order" >&2
  exit 1
fi
