#!/usr/bin/env bash

run_lock_path() {
  printf "%s/start-%s.lock" "$TMP_DIR" "$PORT"
}

acquire_run_lock() {
  mkdir -p "$TMP_DIR"
  RUN_LOCK_DIR="$(run_lock_path)"

  if mkdir "$RUN_LOCK_DIR" 2>/dev/null; then
    printf "%s\n" "$$" >"$RUN_LOCK_DIR/pid"
    return 0
  fi

  local existing_pid=""
  [[ -f "$RUN_LOCK_DIR/pid" ]] && existing_pid="$(cat "$RUN_LOCK_DIR/pid" 2>/dev/null || true)"
  if [[ "$existing_pid" =~ ^[0-9]+$ ]] && kill -0 "$existing_pid" 2>/dev/null; then
    err "已有 start.sh 管理此 repo/port（PID: $existing_pid, lock: $RUN_LOCK_DIR）"
    info "若確認是殘留狀態，請先執行 ./start.sh --stop-all"
    return 1
  fi

  warn "發現過期啟動鎖，將清理: $RUN_LOCK_DIR"
  rm -rf "$RUN_LOCK_DIR"
  mkdir "$RUN_LOCK_DIR"
  printf "%s\n" "$$" >"$RUN_LOCK_DIR/pid"
}

release_run_lock() {
  [[ -n "$RUN_LOCK_DIR" && -d "$RUN_LOCK_DIR" ]] || return 0

  local owner=""
  [[ -f "$RUN_LOCK_DIR/pid" ]] && owner="$(cat "$RUN_LOCK_DIR/pid" 2>/dev/null || true)"
  if [[ -z "$owner" || "$owner" == "$$" ]]; then
    rm -rf "$RUN_LOCK_DIR"
  fi
  RUN_LOCK_DIR=""
}

install_cleanup_traps() {
  trap cleanup_all EXIT
  trap 'cleanup_all; trap - EXIT INT TERM; exit 130' INT
  trap 'cleanup_all; trap - EXIT INT TERM; exit 143' TERM
}

clear_cleanup_traps() {
  trap - EXIT INT TERM
}

kill_process_tree() {
  local pid="$1"
  local signal="${2:-TERM}"
  local child

  [[ "$pid" =~ ^[0-9]+$ ]] || return 0
  [[ "$pid" == "$$" || "$pid" == "$BASHPID" ]] && return 0

  while IFS= read -r child; do
    [[ -z "$child" ]] && continue
    kill_process_tree "$child" "$signal"
  done < <(pgrep -P "$pid" 2>/dev/null || true)

  kill "-$signal" "$pid" 2>/dev/null || true
}

wait_for_pid_exit() {
  local pid="$1"
  local timeout_tenths="${2:-30}"
  local i

  for ((i = 0; i < timeout_tenths; i++)); do
    if ! kill -0 "$pid" 2>/dev/null; then
      return 0
    fi
    sleep 0.1
  done
  return 1
}

stop_process_tree() {
  local pid="$1"
  local label="$2"

  [[ "$pid" =~ ^[0-9]+$ ]] || return 0
  kill -0 "$pid" 2>/dev/null || return 0

  info "正在停止 $label（PID: $pid）..."
  kill_process_tree "$pid" TERM
  if ! wait_for_pid_exit "$pid" 30; then
    warn "$label 未在預期時間內停止，改用強制終止"
    kill_process_tree "$pid" KILL
    wait_for_pid_exit "$pid" 20 || true
  fi
}

is_managed_pid() {
  local pid="$1"
  local cmd
  local exe_path
  cmd="$(ps -p "$pid" -o command= 2>/dev/null || true)"
  if [[ "$cmd" == *"/bin/gateway"* ||
    "$cmd" == *"$AIR_BIN"* ||
    "$cmd" == *"cmd/test_ui"* ||
    "$cmd" == *"/exe/test_ui"* ||
    "$cmd" == *"pnpm exec vite"* ||
    "$cmd" == *"/frontend/node_modules/.bin/vite"* ||
    "$cmd" == *"/vite/bin/vite.js"* ]]; then
    return 0
  fi

  exe_path="$(lsof -a -p "$pid" -d txt -Fn 2>/dev/null | sed -n 's/^n//p' | head -n1)"
  [[ "$exe_path" == *"/test_ui"* || "$exe_path" == *"/gateway"* ]]
}

related_process_pids() {
  local pattern
  pattern="$PWD/(start\\.sh|$APP_PATH|$AIR_BIN|$FRONTEND_DIR)|cmd/test_ui|/exe/test_ui|pnpm exec vite|/frontend/node_modules/.bin/vite|/vite/bin/vite.js"
  pgrep -f "$pattern" 2>/dev/null | while IFS= read -r pid; do
    [[ -z "$pid" ]] && continue
    [[ "$pid" == "$$" || "$pid" == "$BASHPID" ]] && continue
    if is_managed_pid "$pid"; then
      printf "%s\n" "$pid"
    fi
  done
}
