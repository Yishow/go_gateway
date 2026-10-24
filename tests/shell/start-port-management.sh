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

PORT=3333
FRONTEND_DEV_PORT=4173
pid_on_port() {
  [[ "$1" == 3333 ]] && printf '7777\n'
}
killed=false
prompted=false
kill() { killed=true; }
read() { prompted=true; }

set +e
output="$(assert_ports_free 2>&1)"
status=$?
set -e

if [[ "$status" -eq 0 || "$output" != *"後端埠 3333 已被 PID 7777 佔用"* ]]; then
  printf 'expected occupied port refusal, status=%s output=%s\n' "$status" "$output" >&2
  exit 1
fi
if [[ "$killed" == true || "$prompted" == true ]]; then
  printf 'occupied port must not kill or prompt for an unrelated process\n' >&2
  exit 1
fi

if validate_port port 0 >/dev/null 2>&1 || validate_port port 65536 >/dev/null 2>&1 || ! validate_port port 65535; then
  printf 'expected port range validation\n' >&2
  exit 1
fi
FRONTEND_DEV_PORT=3333
if validate_dev_ports >/dev/null 2>&1; then
  printf 'expected duplicate ports to be rejected\n' >&2
  exit 1
fi
PORT=08080
FRONTEND_DEV_PORT=8080
if validate_dev_ports >/dev/null 2>&1; then
  printf 'expected numerically equal ports with leading zeroes to be rejected\n' >&2
  exit 1
fi
