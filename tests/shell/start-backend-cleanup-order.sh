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

BACKEND_PID=301
FRONTEND_PID=302
stopped=""
stop_process_tree() { stopped="$stopped$1:$2|"; }
wait() { return 23; }
kill() {
  [[ "$1" == "-0" && "$2" == "301" ]] && return 0
  return 1
}

set +e
monitor_processes >/dev/null 2>&1
status=$?
set -e

if [[ "$status" -eq 0 ]]; then
  printf 'expected frontend exit to fail the monitored dev session\n' >&2
  exit 1
fi
if [[ "$stopped" != *"301:後端 Air|"* ]]; then
  printf 'expected frontend exit to stop backend, got %s\n' "$stopped" >&2
  exit 1
fi
