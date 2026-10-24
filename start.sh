#!/usr/bin/env bash

set -euo pipefail

# Resolve the project from this file so the launcher also works when invoked
# from another directory. load-env.sh intentionally runs after this cd.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
ROOT_DIR="$SCRIPT_DIR"
export REPO_ROOT="$ROOT_DIR"
cd "$ROOT_DIR"

ENV_PORT_SET="${PORT+x}"
ENV_PORT_VALUE="${PORT-}"
ENV_FRONTEND_HOST_SET="${FRONTEND_DEV_HOST+x}"
ENV_FRONTEND_HOST_VALUE="${FRONTEND_DEV_HOST-}"
ENV_FRONTEND_PORT_SET="${FRONTEND_DEV_PORT+x}"
ENV_FRONTEND_PORT_VALUE="${FRONTEND_DEV_PORT-}"
ENV_VITE_PORT_SET="${VITE_DEV_PORT+x}"
ENV_VITE_PORT_VALUE="${VITE_DEV_PORT-}"
if [[ -f "$ROOT_DIR/scripts/load-env.sh" ]]; then
  # shellcheck disable=SC1091
  source "$ROOT_DIR/scripts/load-env.sh"
elif [[ -f "$ROOT_DIR/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ROOT_DIR/.env"
  set +a
fi
# A caller supplied environment value takes precedence over the optional file.
[[ "$ENV_PORT_SET" == x ]] && PORT="$ENV_PORT_VALUE"
[[ "$ENV_FRONTEND_HOST_SET" == x ]] && FRONTEND_DEV_HOST="$ENV_FRONTEND_HOST_VALUE"
[[ "$ENV_FRONTEND_PORT_SET" == x ]] && FRONTEND_DEV_PORT="$ENV_FRONTEND_PORT_VALUE"
if [[ "$ENV_FRONTEND_PORT_SET" != x && "$ENV_VITE_PORT_SET" == x ]]; then
  FRONTEND_DEV_PORT="$ENV_VITE_PORT_VALUE"
fi
[[ "$ENV_VITE_PORT_SET" == x ]] && VITE_DEV_PORT="$ENV_VITE_PORT_VALUE"

FRONTEND_DIR="frontend"
BUILD_HELPER="scripts/build-single.mjs"
TOOLS_DIR="$ROOT_DIR/bin/tools"
TMP_DIR="$ROOT_DIR/bin/tmp"
AIR_CONFIG="$ROOT_DIR/.air.toml"
PORT="${PORT:-8080}"
FRONTEND_DEV_HOST="${FRONTEND_DEV_HOST:-0.0.0.0}"
FRONTEND_DEV_PORT="${FRONTEND_DEV_PORT:-${VITE_DEV_PORT:-5173}}"
VITE_DEV_PORT="$FRONTEND_DEV_PORT"
READY_TIMEOUT="${START_READY_TIMEOUT:-120}"

HOST_GOOS=""
HOST_GOARCH=""
AIR_BIN="$TOOLS_DIR/air"
FRONTEND_PID=""
BACKEND_PID=""
RUN_LOCK_DIR=""
CLEANING_UP=false
MODE=""
BUILD_OS=""
BUILD_ARCH=""

has_cmd() { command -v "$1" >/dev/null 2>&1; }
info() { printf '[INFO] %s\n' "$*"; }
success() { printf '[OK] %s\n' "$*"; }
warn() { printf '[WARN] %s\n' "$*" >&2; }
err() { printf '[ERR] %s\n' "$*" >&2; }

# shellcheck disable=SC1091
source "$ROOT_DIR/scripts/start-process-utils.sh"

usage() {
  cat <<'USAGE'
Go Gateway 啟動器

用法:
  ./start.sh                         顯示選單
  ./start.sh --dev [--port PORT]     Air + Vite 熱重載開發
  ./start.sh --build-single [選項]   建置前後端單一執行檔

選項:
  --dev, --dev-mode, --air-mode      開發模式別名
  --build-single                     單一執行檔建置
  --os OS                            windows、mac/darwin 或 linux
  --arch ARCH                        amd64 或 arm64
  --port PORT                        後端埠（預設 8080）
  -h, --help                         顯示說明

開發模式需要 Go、Node.js、npm、lsof、pgrep；前端第一次啟動若找不到
frontend/node_modules/vite/bin/vite.js，才會執行 npm ci。
USAGE
}

validate_port() {
  local name="$1"
  local value="$2"
  if [[ ! "$value" =~ ^[0-9]+$ ]] || (( 10#$value < 1 || 10#$value > 65535 )); then
    err "${name} 必須是 1 到 65535 的整數：${value}"
    return 1
  fi
}

canonical_port() { printf '%s\n' "$((10#$1))"; }

validate_dev_ports() {
  validate_port "後端 PORT" "$PORT" || return 1
  validate_port "前端開發埠" "$FRONTEND_DEV_PORT" || return 1
  if [[ "$(canonical_port "$PORT")" == "$(canonical_port "$FRONTEND_DEV_PORT")" ]]; then
    err "後端 PORT 與前端開發埠不可相同：${PORT}"
    return 1
  fi
}

pid_on_port() {
  lsof -tiTCP:"$1" -sTCP:LISTEN 2>/dev/null | head -n 1 || true
}

assert_ports_free() {
  local backend_pid frontend_pid
  backend_pid="$(pid_on_port "$PORT")"
  frontend_pid="$(pid_on_port "$FRONTEND_DEV_PORT")"
  if [[ -n "$backend_pid" ]]; then
    err "後端埠 ${PORT} 已被 PID ${backend_pid} 佔用；請先停止該服務"
    return 1
  fi
  if [[ -n "$frontend_pid" ]]; then
    err "前端開發埠 ${FRONTEND_DEV_PORT} 已被 PID ${frontend_pid} 佔用；請先停止該服務"
    return 1
  fi
}

validate_dev_environment() {
  local command_name missing=0
  for command_name in go node npm lsof pgrep; do
    if ! has_cmd "$command_name"; then
      err "找不到必要工具：${command_name}"
      missing=1
    fi
  done
  [[ -f "$ROOT_DIR/go.mod" ]] || { err "找不到 go.mod：${ROOT_DIR}"; missing=1; }
  [[ -f "$ROOT_DIR/.air.toml" ]] || { err "找不到 Air 設定：${AIR_CONFIG}"; missing=1; }
  [[ -f "$ROOT_DIR/frontend/package.json" ]] || { err "找不到前端 package.json"; missing=1; }
  validate_dev_ports || missing=1
  if ! [[ "$READY_TIMEOUT" =~ ^[0-9]+$ ]] || (( 10#$READY_TIMEOUT < 120 )); then
    err "START_READY_TIMEOUT 必須是至少 120 秒的整數：${READY_TIMEOUT}"
    missing=1
  fi
  (( missing == 0 ))
}

ensure_air() {
  if ! HOST_GOOS="$(go env GOHOSTOS 2>/dev/null)" || ! HOST_GOARCH="$(go env GOHOSTARCH 2>/dev/null)"; then
    err "無法讀取 Go 主機 GOOS/GOARCH"
    return 1
  fi
  [[ -n "$HOST_GOOS" && -n "$HOST_GOARCH" ]] || {
    err "Go 未回傳有效的主機 GOOS/GOARCH"
    return 1
  }

  if [[ "$HOST_GOOS" == "windows" ]]; then
    AIR_BIN="$TOOLS_DIR/air.exe"
  else
    AIR_BIN="$TOOLS_DIR/air"
  fi
  if [[ ! -x "$AIR_BIN" ]]; then
    info "安裝 Air v1.63.0 至 ${TOOLS_DIR}"
    mkdir -p "$TOOLS_DIR"
    if ! GOBIN="$TOOLS_DIR" GOOS="$HOST_GOOS" GOARCH="$HOST_GOARCH" CGO_ENABLED=0 GOTOOLCHAIN=local \
      go install github.com/air-verse/air@v1.63.0; then
      err "Air 安裝失敗"
      return 1
    fi
  fi
  if [[ ! -x "$AIR_BIN" ]]; then
    err "找不到可執行的 Air：${AIR_BIN}"
    return 1
  fi
}

ensure_frontend_dependencies() {
  local vite_entry="$ROOT_DIR/$FRONTEND_DIR/node_modules/vite/bin/vite.js"
  [[ -f "$vite_entry" ]] && return 0

  info "找不到本機 Vite，執行 npm ci"
  if ! (cd "$ROOT_DIR/$FRONTEND_DIR" && npm ci); then
    err "npm ci 失敗，無法啟動前端"
    return 1
  fi
  [[ -f "$vite_entry" ]] || {
    err "npm ci 完成後仍找不到 Vite：${vite_entry}"
    return 1
  }
}

ensure_embedded_frontend_snapshot() {
  local index="$ROOT_DIR/cmd/test_ui/static/index.html"
  [[ -s "$index" ]] && return 0
  info "找不到 embedded 前端快照，先準備開發啟動資產"
  if ! node "$ROOT_DIR/$BUILD_HELPER" --prepare-dev; then
    err "無法準備 embedded 前端快照"
    return 1
  fi
  [[ -s "$index" ]] || {
    err "開發資產準備後仍找不到：${index}"
    return 1
  }
}

start_backend_process() {
  (
    export PORT="$PORT"
    export AUTO_OPEN_BROWSER=false
    export GOOS="$HOST_GOOS" GOARCH="$HOST_GOARCH" CGO_ENABLED=0 GOTOOLCHAIN=local
    exec "$AIR_BIN" -c "$AIR_CONFIG"
  ) &
  BACKEND_PID=$!
}

start_frontend_dev_server() {
  local vite_entry="$ROOT_DIR/$FRONTEND_DIR/node_modules/vite/bin/vite.js"
  (
    cd "$ROOT_DIR/$FRONTEND_DIR"
    export PORT="$PORT"
    export VITE_DEV_PORT="$FRONTEND_DEV_PORT"
    export VITE_API_PROXY_TARGET="http://127.0.0.1:$PORT"
    export BROWSERSLIST_IGNORE_OLD_DATA=true
    exec node "$vite_entry" --host "$FRONTEND_DEV_HOST" --port "$FRONTEND_DEV_PORT" --strictPort
  ) &
  FRONTEND_PID=$!
}

cleanup_frontend() {
  if [[ -n "$FRONTEND_PID" ]]; then
    stop_process_tree "$FRONTEND_PID" "前端開發伺服器"
    wait "$FRONTEND_PID" 2>/dev/null || true
    FRONTEND_PID=""
  fi
}

cleanup_backend() {
  if [[ -n "$BACKEND_PID" ]]; then
    stop_process_tree "$BACKEND_PID" "後端 Air"
    wait "$BACKEND_PID" 2>/dev/null || true
    BACKEND_PID=""
  fi
}

cleanup_all() {
  [[ "$CLEANING_UP" == true ]] && return 0
  CLEANING_UP=true
  cleanup_frontend
  cleanup_backend
  release_run_lock
}

install_cleanup_traps() {
  trap cleanup_all EXIT
  trap 'exit 130' INT
  trap 'exit 143' TERM
}

monitor_processes() {
  local backend_code frontend_code
  while :; do
    if [[ -z "$BACKEND_PID" ]] || ! kill -0 "$BACKEND_PID" 2>/dev/null; then
      backend_code=1
      if [[ -n "$BACKEND_PID" ]]; then
        if wait "$BACKEND_PID"; then backend_code=0; else backend_code=$?; fi
        BACKEND_PID=""
      fi
      err "後端 Air 已停止（exit ${backend_code}）"
      cleanup_frontend
      (( backend_code == 0 )) && return 1
      return "$backend_code"
    fi
    if [[ -z "$FRONTEND_PID" ]] || ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
      frontend_code=1
      if [[ -n "$FRONTEND_PID" ]]; then
        if wait "$FRONTEND_PID"; then frontend_code=0; else frontend_code=$?; fi
        FRONTEND_PID=""
      fi
      err "前端開發伺服器已停止（exit ${frontend_code}），將停止後端"
      cleanup_backend
      return 1
    fi
    sleep 0.25
  done
}

# Kept as descriptive compatibility names for focused shell tests.
wait_for_backend_exit() { monitor_processes; }
clear_port() { assert_ports_free; }
clear_frontend_port() { assert_ports_free; }

start_dev_mode() {
  validate_dev_environment || return 1
  PORT="$(canonical_port "$PORT")"
  FRONTEND_DEV_PORT="$(canonical_port "$FRONTEND_DEV_PORT")"
  VITE_DEV_PORT="$FRONTEND_DEV_PORT"
  READY_TIMEOUT="$((10#$READY_TIMEOUT))"
  # Refuse an occupied listener before installing anything or starting a
  # child. A launcher must never guess whether that process is safe to kill.
  assert_ports_free || return 1
  request_open_browser
  ensure_embedded_frontend_snapshot || return 1
  ensure_air || return 1
  ensure_frontend_dependencies || return 1
  assert_ports_free || return 1

  install_cleanup_traps
  acquire_run_lock || return 1
  info "啟動 Air 後端（GOOS=${HOST_GOOS} GOARCH=${HOST_GOARCH} CGO_ENABLED=0）"
  start_backend_process
  if ! wait_for_port_ready "$PORT" "後端" "$BACKEND_PID" "$READY_TIMEOUT"; then
    return 1
  fi
  success "後端已就緒：http://127.0.0.1:${PORT}"

  info "啟動 Vite 前端（host=${FRONTEND_DEV_HOST} port=${FRONTEND_DEV_PORT}）"
  start_frontend_dev_server
  if ! wait_for_port_ready "$FRONTEND_DEV_PORT" "前端" "$FRONTEND_PID" "$READY_TIMEOUT"; then
    return 1
  fi
  success "前端已就緒：${FRONTEND_DEV_PORT}"
  info "開發入口：http://127.0.0.1:${FRONTEND_DEV_PORT}/studio/v2"
  info "後端 API：http://127.0.0.1:${PORT}"
  info "修改 Go 或前端程式碼後會自動重建／更新；按 Ctrl+C 停止"
  open_dev_browser
  wait_for_backend_exit
}

normalise_build_os() {
  case "$1" in
    windows|darwin|linux) printf '%s\n' "$1" ;;
    mac) printf 'darwin\n' ;;
    *) err "--os 必須是 windows、mac/darwin 或 linux：${1}"; return 1 ;;
  esac
}

validate_build_target() {
  if [[ -n "$BUILD_OS" ]]; then BUILD_OS="$(normalise_build_os "$BUILD_OS")" || return 1; fi
  if [[ -n "$BUILD_ARCH" && "$BUILD_ARCH" != "amd64" && "$BUILD_ARCH" != "arm64" ]]; then
    err "--arch 必須是 amd64 或 arm64：${BUILD_ARCH}"
    return 1
  fi
}

build_single() {
  has_cmd node || { err "找不到必要工具：node"; return 1; }
  [[ -f "$ROOT_DIR/$BUILD_HELPER" ]] || {
    err "找不到單檔建置 helper：${ROOT_DIR}/${BUILD_HELPER}"
    return 1
  }
  validate_build_target || return 1
  if [[ -n "$BUILD_OS" && -n "$BUILD_ARCH" ]]; then
    info "建置單一執行檔（OS=${BUILD_OS} ARCH=${BUILD_ARCH}）"
  elif [[ -n "$BUILD_OS" ]]; then
    info "建置單一執行檔（OS=${BUILD_OS}）"
  elif [[ -n "$BUILD_ARCH" ]]; then
    info "建置單一執行檔（ARCH=${BUILD_ARCH}）"
  else
    info "建置單一執行檔"
  fi
  if [[ -n "$BUILD_OS" && -n "$BUILD_ARCH" ]]; then
    node "$ROOT_DIR/$BUILD_HELPER" --os "$BUILD_OS" --arch "$BUILD_ARCH"
  elif [[ -n "$BUILD_OS" ]]; then
    node "$ROOT_DIR/$BUILD_HELPER" --os "$BUILD_OS"
  elif [[ -n "$BUILD_ARCH" ]]; then
    node "$ROOT_DIR/$BUILD_HELPER" --arch "$BUILD_ARCH"
  else
    node "$ROOT_DIR/$BUILD_HELPER"
  fi
}

show_main_menu() {
  cat <<'MENU'

Go Gateway
  1) 開發模式（Air + Vite 熱重載）
  2) 建置單一執行檔
  0) 離開
MENU
}

select_build_target() {
  printf '%s\n' '選擇目標作業系統：' '  1) windows' '  2) mac (darwin)' '  3) linux' '  0) 取消'
  local os_choice arch_choice
  IFS= read -r os_choice || return 1
  case "$os_choice" in
    1) BUILD_OS=windows ;;
    2) BUILD_OS=darwin ;;
    3) BUILD_OS=linux ;;
    0) return 2 ;;
    *) err "無效的作業系統選項"; return 1 ;;
  esac
  printf '%s\n' '選擇目標架構：' '  1) amd64' '  2) arm64' '  0) 取消'
  IFS= read -r arch_choice || return 1
  case "$arch_choice" in
    1) BUILD_ARCH=amd64 ;;
    2) BUILD_ARCH=arm64 ;;
    0) return 2 ;;
    *) err "無效的架構選項"; return 1 ;;
  esac
}

main() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --dev|--dev-mode|--air-mode)
        [[ -z "$MODE" || "$MODE" == dev ]] || { err "--dev 與 --build-single 不可同時使用"; return 1; }
        MODE=dev
        ;;
      --build-single)
        [[ -z "$MODE" || "$MODE" == build ]] || { err "--dev 與 --build-single 不可同時使用"; return 1; }
        MODE=build
        ;;
      --os|--arch|--port)
        [[ $# -ge 2 && -n "$2" ]] || { err "${1} 需要非空值"; return 1; }
        case "$1" in
          --os) BUILD_OS="$2" ;;
          --arch) BUILD_ARCH="$2" ;;
          --port) PORT="$2" ;;
        esac
        shift
        ;;
      --os=*) BUILD_OS="${1#*=}"; [[ -n "$BUILD_OS" ]] || { err "--os 需要非空值"; return 1; } ;;
      --arch=*) BUILD_ARCH="${1#*=}"; [[ -n "$BUILD_ARCH" ]] || { err "--arch 需要非空值"; return 1; } ;;
      --port=*) PORT="${1#*=}"; [[ -n "$PORT" ]] || { err "--port 需要非空值"; return 1; } ;;
      -h|--help) usage; return 0 ;;
      *) err "未知參數：${1}"; usage >&2; return 1 ;;
    esac
    shift
  done

  if [[ -z "$MODE" ]]; then
    show_main_menu
    local menu_choice
    IFS= read -r menu_choice || return 0
    case "$menu_choice" in
      1) MODE=dev ;;
      2)
        MODE=build
        select_build_target || {
          local selection_status=$?
          (( selection_status == 2 )) && return 0
          return "$selection_status"
        }
        ;;
      0) return 0 ;;
      *) err "無效的選項"; return 1 ;;
    esac
  fi

  case "$MODE" in
    dev)
      [[ -z "$BUILD_OS" && -z "$BUILD_ARCH" ]] || { err "--os/--arch 只能搭配 --build-single"; return 1; }
      start_dev_mode
      ;;
    build)
      build_single
      ;;
  esac
}

if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then
  return 0
fi

main "$@"
