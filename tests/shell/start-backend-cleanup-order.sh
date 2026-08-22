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
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TEST_TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_TMP_DIR/scripts/start-process-utils.sh"

cd "$TEST_TMP_DIR"
REPO_ROOT="$TEST_TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

BACKEND_PID=222
BACKEND_LOG_PID=111
BACKEND_LOG_FILE="$TEST_TMP_DIR/backend.log"
touch "$BACKEND_LOG_FILE"

kill_calls=()
kill() {
  kill_calls+=("$1 $2")
}

wait() {
  return 0
}

info() { :; }
success() { :; }

cleanup_backend

# 現行 kill_process_tree 以 `kill -TERM <pid>` 先子後父；
# cleanup_backend 必須先停後端（222）再停 log tail（111）
backend_term_index=-1
log_term_index=-1
for i in "${!kill_calls[@]}"; do
  if [[ "${kill_calls[$i]}" == "-TERM 222" && "$backend_term_index" -lt 0 ]]; then
    backend_term_index="$i"
  fi
  if [[ "${kill_calls[$i]}" == "-TERM 111" ]]; then
    log_term_index="$i"
  fi
done

if [[ "$backend_term_index" -lt 0 || "$log_term_index" -lt 0 ]]; then
  printf 'expected cleanup_backend to send TERM to backend(222) and log tail(111), got %s\n' "$(printf '%s|' "${kill_calls[@]}")" >&2
  exit 1
fi

if (( backend_term_index > log_term_index )); then
  printf 'expected cleanup_backend to stop backend before log tail, got %s\n' "$(printf '%s|' "${kill_calls[@]}")" >&2
  exit 1
fi
