#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
FIXTURE_DIR="$(mktemp -d)"
trap 'rm -rf "$FIXTURE_DIR"' EXIT
mkdir -p "$FIXTURE_DIR/scripts"
cp "$REPO_ROOT/start.sh" "$FIXTURE_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$FIXTURE_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$FIXTURE_DIR/scripts/start-process-utils.sh"
cd "$FIXTURE_DIR"
source ./start.sh >/dev/null 2>&1 || true

LISTENER_PIDS=""
DEAD_PID="404"

pid_on_port() { printf '%s\n' "$LISTENER_PIDS"; }
kill() {
  [[ "${2:-}" != "$DEAD_PID" ]]
}
ps() {
  case "${2:-}" in
    202) printf ' 101\n' ;;
    303) printf ' 999\n' ;;
    999) printf ' 1\n' ;;
    1) printf ' 1\n' ;;
    *) command ps "$@" ;;
  esac
}

run_ready() {
  local owner="$1"
  local expected_status="$2"
  local output status
  set +e
  output="$(wait_for_port_ready 4173 測試 "$owner" 1 2>&1)"
  status=$?
  set -e
  if [[ "$status" -ne "$expected_status" ]]; then
    printf 'unexpected readiness status owner=%s status=%s output=%s\n' "$owner" "$status" "$output" >&2
    exit 1
  fi
}

LISTENER_PIDS="101"
run_ready 101 0

LISTENER_PIDS="202"
run_ready 101 0

LISTENER_PIDS="303"
run_ready 101 1

LISTENER_PIDS="$DEAD_PID"
run_ready "$DEAD_PID" 1

PORT=4173
RUN_LOCK_DIR=""
mkdir -p "$TMP_DIR/start-4173.lock"
printf '%s\n' "$$" >"$TMP_DIR/start-4173.lock/pid"
set +e
lock_output="$(acquire_run_lock 2>&1)"
lock_status=$?
set -e
if [[ "$lock_status" -eq 0 || "$lock_output" == *"--stop-all"* ]]; then
  printf 'stale-lock guidance is not actionable: status=%s output=%s\n' "$lock_status" "$lock_output" >&2
  exit 1
fi

printf 'PASS: start-owned-readiness.sh\n'
