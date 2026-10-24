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

FRONTEND_PID=111
BACKEND_PID=222
stopped=()
stop_process_tree() { stopped+=("$1:$2"); }
wait() { :; }
release_run_lock() { :; }

cleanup_all

if [[ "$(IFS='|'; printf '%s' "${stopped[*]}")" != "111:前端開發伺服器|222:後端 Air" ]]; then
  printf 'expected cleanup to stop only owned frontend/backend in order, got %s\n' "$(IFS='|'; printf '%s' "${stopped[*]}")" >&2
  exit 1
fi
if [[ -n "$FRONTEND_PID" || -n "$BACKEND_PID" ]]; then
  printf 'expected cleanup to clear owned process IDs\n' >&2
  exit 1
fi
