#!/usr/bin/env bash

request_open_browser() {
  OPEN_BROWSER=false
  local answer=""
  printf '是否開啟瀏覽器？[y/N]: '
  IFS= read -r answer || answer=""
  case "$answer" in y|Y|yes|YES|Yes) OPEN_BROWSER=true ;; esac
}

open_dev_browser() {
  [[ "${OPEN_BROWSER:-false}" == true ]] || return 0
  local url="http://127.0.0.1:${FRONTEND_DEV_PORT}/studio/v2"
  local opener
  case "$(uname -s)" in
    Darwin) opener=open ;;
    Linux) opener=xdg-open ;;
    *) warn "無法自動開啟瀏覽器，請手動開啟：${url}"; return 0 ;;
  esac
  if ! command -v "$opener" >/dev/null 2>&1 || ! "$opener" "$url"; then
    warn "瀏覽器開啟失敗，請手動開啟：${url}"
  fi
}

pid_is_descendant() {
  local candidate_pid="$1"
  local owner_pid="$2"
  local current_pid="$candidate_pid"
  local parent_pid
  local hops=0

  while (( hops < 64 )); do
    [[ "$current_pid" == "$owner_pid" ]] && return 0
    [[ "$current_pid" =~ ^[0-9]+$ ]] || return 1
    parent_pid="$(ps -p "$current_pid" -o ppid= 2>/dev/null || true)"
    parent_pid="${parent_pid//[[:space:]]/}"
    [[ "$parent_pid" =~ ^[0-9]+$ ]] || return 1
    [[ "$parent_pid" == "$current_pid" ]] && return 1
    current_pid="$parent_pid"
    hops=$((hops + 1))
  done
  return 1
}

port_owned_by_pid() {
  local port="$1"
  local owner_pid="$2"
  local listener_pid
  local found=0

  while IFS= read -r listener_pid; do
    listener_pid="${listener_pid//[[:space:]]/}"
    [[ -z "$listener_pid" ]] && continue
    found=1
    [[ "$listener_pid" =~ ^[0-9]+$ ]] || return 1
    if [[ "$listener_pid" != "$owner_pid" ]] && ! pid_is_descendant "$listener_pid" "$owner_pid"; then
      return 1
    fi
  done <<<"$(pid_on_port "$port")"

  (( found == 1 ))
}

wait_for_port_ready() {
  local port="$1"
  local label="$2"
  local pid="$3"
  local timeout_seconds="${4:-$READY_TIMEOUT}"
  local attempts=$((timeout_seconds * 10))
  local attempt=0

  while (( attempt < attempts )); do
    if ! kill -0 "$pid" 2>/dev/null; then
      err "${label} 未成功啟動（埠 ${port}）"
      return 1
    fi
    port_owned_by_pid "$port" "$pid" && return 0
    sleep 0.1
    attempt=$((attempt + 1))
  done
  err "${label} 在 ${timeout_seconds} 秒內未監聽埠 ${port}"
  return 1
}

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
    err "已有 start.sh 管理此 repo/port（PID: ${existing_pid}, lock: ${RUN_LOCK_DIR}）"
    info "請先停止原啟動器（PID: ${existing_pid}），確認該 PID 已結束後再重試"
    return 1
  fi

  warn "發現過期啟動鎖，將清理: ${RUN_LOCK_DIR}"
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
  # Bash 3.2 (the macOS system Bash) does not provide BASHPID. The launcher
  # only ever passes child PIDs it owns, and $$ is enough to protect itself.
  [[ "$pid" == "$$" ]] && return 0

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

  info "正在停止 ${label}（PID: ${pid}）..."
  kill_process_tree "$pid" TERM
  if ! wait_for_pid_exit "$pid" 30; then
    warn "${label} 未在預期時間內停止，改用強制終止"
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
    "$cmd" == *"/frontend/node_modules/vite/bin/vite.js"* ||
    "$cmd" == *"/node_modules/vite/bin/vite.js"* ]]; then
    return 0
  fi

  exe_path="$(lsof -a -p "$pid" -d txt -Fn 2>/dev/null | sed -n 's/^n//p' | head -n1)"
  [[ "$exe_path" == *"/test_ui"* || "$exe_path" == *"/gateway"* ]]
}

related_process_pids() {
  local pattern
  pattern="$PWD/(start\\.sh|$APP_PATH|$AIR_BIN|$FRONTEND_DIR)|cmd/test_ui|/exe/test_ui|/node_modules/vite/bin/vite.js"
  pgrep -f "$pattern" 2>/dev/null | while IFS= read -r pid; do
    [[ -z "$pid" ]] && continue
    [[ "$pid" == "$$" ]] && continue
    if is_managed_pid "$pid"; then
      printf "%s\n" "$pid"
    fi
  done
}
