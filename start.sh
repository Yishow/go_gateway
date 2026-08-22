#!/usr/bin/env bash

set -euo pipefail

load_env_script="$(dirname "${BASH_SOURCE[0]}")/scripts/load-env.sh"
if [[ -f "$load_env_script" ]]; then
  source "$load_env_script"
elif [[ -f ".env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source ".env"
  set +a
fi

APP_NAME="gateway"
APP_PATH="cmd/test_ui"
BUILD_DIR="bin"
FRONTEND_DIR="frontend"
FRONTEND_DEV_HOST="${FRONTEND_DEV_HOST:-0.0.0.0}"
FRONTEND_DEV_PORT="${FRONTEND_DEV_PORT:-${VITE_DEV_PORT:-5173}}"
VITE_DEV_PORT="${VITE_DEV_PORT:-$FRONTEND_DEV_PORT}"
PORT="${PORT:-8080}"
LOG_DIR="bin/logs"
TMP_DIR="bin/tmp"
AIR_BIN="$TMP_DIR/gateway-air.exe"
MODBUS_SHARE_PORT="${MODBUS_SHARE_PORT:-5020}"
AUTO_KILL_PORT=false
SKIP_BUILD=false
SKIP_LINT=false
SKIP_TEST=false
SKIP_QUALITY=false
SYNC_EMBED=false
START=false
DEV_MODE=false
AIR_MODE=false
QUICK_START=false
CHECK_ENV=false
HEALTH_CHECK=false
BUILD_SINGLE=false
LIST_PROCESSES=false
STOP_ALL=false
DIAGNOSE=false
COVERAGE=false
VERBOSE=false
HAS_ANY_PARAM=false
FRONTEND_PID=""
BACKEND_PID=""
BACKEND_LOG_PID=""
BACKEND_LOG_FILE=""
RUN_LOCK_DIR=""

# 日誌噪音計數器（關聯陣列）
# 在 bash 5.3 + set -u 下，僅 declare 而未初始化的空關聯陣列
# 於 `${#array[@]}` 存取時會觸發 unbound variable。
declare -A LOG_NOISE_COUNTERS=()

info() { printf "\033[36m[INFO]\033[0m %s\n" "$*"; }
success() { printf "\033[32m[OK]\033[0m %s\n" "$*"; }
warn() { printf "\033[33m[WARN]\033[0m %s\n" "$*"; }
err() { printf "\033[31m[ERR]\033[0m %s\n" "$*"; }

# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/scripts/start-log-utils.sh"
# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/scripts/start-process-utils.sh"

show_frontend_host_hint() {
  if [[ "$FRONTEND_DEV_HOST" == "0.0.0.0" || "$FRONTEND_DEV_HOST" == "::" ]]; then
    info "💡 前端開發伺服器已監聽所有介面，區網請使用本機 LAN IP 存取（port ${FRONTEND_DEV_PORT}）"
  else
    info "💡 前端開發伺服器通常運行在 http://$FRONTEND_DEV_HOST:${FRONTEND_DEV_PORT}"
  fi
}

show_embedded_frontend_hint() {
  if [[ "$SYNC_EMBED" == true ]]; then
    info "💡 :$PORT 會使用本次啟動前同步好的 embedded 前端快照"
    info "💡 前端在啟動後若還有新修改，:${FRONTEND_DEV_PORT} 會即時更新；:$PORT 需重新執行 start.sh --sync-embed 才會刷新"
    return
  fi

  info "💡 開發模式預設不重建 embedded 前端；即時開發請使用 :${FRONTEND_DEV_PORT}"
  info "💡 若需要刷新 :$PORT 的 embedded 前端，重新執行時加上 --sync-embed"
}

frontend_proxy_target() {
  printf "http://127.0.0.1:%s" "$PORT"
}

ensure_frontend_dependencies() {
  if [[ ! -d "$FRONTEND_DIR" ]]; then
    return 0
  fi

  if [[ -x "$FRONTEND_DIR/node_modules/.bin/vite" ]]; then
    return 0
  fi

  info "安裝前端依賴..."
  pushd "$FRONTEND_DIR" >/dev/null
  if ! pnpm install; then
    popd >/dev/null
    err "前端依賴安裝失敗，請檢查網路連線或 frontend/package.json"
    return 1
  fi
  popd >/dev/null

  if [[ ! -x "$FRONTEND_DIR/node_modules/.bin/vite" ]]; then
    err "前端依賴安裝後仍找不到 vite，請檢查 frontend/package.json"
    return 1
  fi
}

start_frontend_dev_server() {
  local proxy_target
  local node_options
  proxy_target="$(frontend_proxy_target)"
  node_options="${NODE_OPTIONS:---max-old-space-size=4096}"

  ensure_frontend_dependencies || return 1
  clear_frontend_port true

  (
    cd "$FRONTEND_DIR"
    export NODE_OPTIONS="$node_options"
    export PORT="$PORT"
    export VITE_API_PROXY_TARGET="$proxy_target"
    export VITE_DEV_PORT="$FRONTEND_DEV_PORT"
    exec pnpm exec vite --host "$FRONTEND_DEV_HOST" --port "$FRONTEND_DEV_PORT" --strictPort
  ) &
  FRONTEND_PID=$!

  local i
  for ((i = 0; i < 100; i++)); do
    if [[ -n "$(pid_on_port "$FRONTEND_DEV_PORT")" ]]; then
      return 0
    fi
    if ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
      FRONTEND_PID=""
      err "前端開發伺服器未成功啟動（port ${FRONTEND_DEV_PORT}）"
      return 1
    fi
    sleep 0.1
  done

  stop_process_tree "$FRONTEND_PID" "前端開發伺服器"
  wait "$FRONTEND_PID" 2>/dev/null || true
  FRONTEND_PID=""
  err "前端開發伺服器未成功啟動（port ${FRONTEND_DEV_PORT}）"
  return 1
}

start_backend_process() {
  local command="$1"

  mkdir -p "$TMP_DIR"
  BACKEND_LOG_FILE="$(mktemp "$TMP_DIR/backend.XXXXXX")"

  tail -n +1 -F "$BACKEND_LOG_FILE" 2>/dev/null | while IFS= read -r line; do
    runtime_log_line "$line"
  done &
  BACKEND_LOG_PID=$!

  bash -c "$command" >>"$BACKEND_LOG_FILE" 2>&1 &
  BACKEND_PID=$!
}

wait_for_port_ready() {
  local port="$1"
  local label="${2:-端口}"
  local attempts="${3:-100}"
  local i

  for ((i = 0; i < attempts; i++)); do
    if [[ -n "$(pid_on_port "$port")" ]]; then
      return 0
    fi
    if [[ -n "$BACKEND_PID" ]] && ! kill -0 "$BACKEND_PID" 2>/dev/null; then
      err "$label $port 未成功啟動"
      return 1
    fi
    sleep 0.1
  done

  err "$label $port 在預期時間內未成功啟動"
  return 1
}

wait_for_backend_exit() {
  local exit_code=0
  local pid="$BACKEND_PID"

  if [[ -z "$pid" ]]; then
    return 0
  fi

  wait "$pid" || exit_code=$?
  BACKEND_PID=""
  cleanup_all
  clear_cleanup_traps
  return "$exit_code"
}

# 清理前端開發伺服器及其子進程
cleanup_frontend() {
  if [[ -n "$FRONTEND_PID" ]]; then
    stop_process_tree "$FRONTEND_PID" "前端開發伺服器"
    wait "$FRONTEND_PID" 2>/dev/null || true
    FRONTEND_PID=""
    success "前端開發伺服器已停止"
  fi
}

cleanup_backend() {
  if [[ -n "$BACKEND_PID" ]]; then
    stop_process_tree "$BACKEND_PID" "後端服務"
    wait "$BACKEND_PID" 2>/dev/null || true
    BACKEND_PID=""
    success "後端服務已停止"
  fi

  if [[ -n "$BACKEND_LOG_PID" ]]; then
    kill_process_tree "$BACKEND_LOG_PID" TERM
    BACKEND_LOG_PID=""
  fi

  if [[ -n "$BACKEND_LOG_FILE" && -f "$BACKEND_LOG_FILE" ]]; then
    rm -f "$BACKEND_LOG_FILE"
    BACKEND_LOG_FILE=""
  fi
}

# 清理所有子進程
cleanup_all() {
  cleanup_frontend
  cleanup_backend
  release_run_lock
  show_log_noise_summary
}

usage() {
  cat <<USAGE
Go Gateway mac 啟動腳本

用法:
  ./start.sh [選項]

選項:
  --dev-mode           開發模式 (go run)
  --air-mode           熱重載模式 (air)
  --quick-start        一鍵啟動 (lint + sync frontend + build + start)
  --start              建置後啟動服務
  --skip-build         跳過一般建置（若需刷新 embedded 前端會自動補建）
  --skip-lint          跳過 golangci-lint
  --skip-test          跳過測試
  --skip-quality       跳過代碼質量檢查
  --sync-embed         開發/熱重載啟動前同步 embedded 前端快照
  --coverage           輸出詳細覆蓋率與 coverage 報告
  --verbose            顯示詳細輸出
  --port <port>        指定端口 (預設: 8080)
  --auto-kill-port     自動清理佔用端口
  --check-env          檢查環境
  --health-check       健康檢查
  --list-processes     列出運行中的服務進程
  --stop-all           停止所有運行中的服務
  --diagnose           快速診斷常見問題
  --build-single       建置單一執行檔 (embed)
  -h, --help           顯示說明
USAGE
}

has_cmd() { command -v "$1" >/dev/null 2>&1; }

find_air_cmd() {
  if has_cmd air; then
    command -v air
    return 0
  fi

  local gobin gopath
  gobin="$(go env GOBIN 2>/dev/null || true)"
  if [[ -n "$gobin" && -x "$gobin/air" ]]; then
    printf "%s\n" "$gobin/air"
    return 0
  fi

  gopath="$(go env GOPATH 2>/dev/null || true)"
  if [[ -n "$gopath" && -x "$gopath/bin/air" ]]; then
    printf "%s\n" "$gopath/bin/air"
    return 0
  fi

  if [[ -x "$HOME/go/bin/air" ]]; then
    printf "%s\n" "$HOME/go/bin/air"
    return 0
  fi

  return 1
}

show_main_menu() {
  cat <<MENU

============================================
   Go Gateway 啟動選單 (mac)
============================================
  1) 開發模式 (go run)
  2) 熱重載模式 (air)
  3) 一鍵啟動 (lint + sync frontend + build + start)
  4) 完整流程 (lint + build + quality + test)
  5) 僅建置 (sync frontend + build)
  6) 快速啟動 (sync frontend + build + start，跳過 lint/test/quality)
 7) 僅測試
 8) 檢查環境
 9) 健康檢查
 A) 列出運行中的服務
 B) 停止所有服務
 C) 快速診斷
 D) 建置單一執行檔 (embed)
 0) 退出
============================================
MENU
}

check_env() {
  info "檢查環境..."
  local ok=true
  for c in go node pnpm; do
    if has_cmd "$c"; then
      success "$c: $($c --version 2>/dev/null | head -n1 || true)"
    else
      warn "$c 未安裝"
      [[ "$c" == "go" ]] && ok=false
    fi
  done
  if has_cmd golangci-lint; then
    success "golangci-lint: $(golangci-lint version 2>/dev/null | head -n1 || true)"
  else
    warn "golangci-lint 未安裝"
  fi
  local air_cmd
  if air_cmd="$(find_air_cmd)"; then
    success "air: $("$air_cmd" -v 2>/dev/null | head -n1 || true)"
  else
    warn "air 未安裝"
  fi
  [[ "$ok" == true ]]
}

pid_on_port() {
  lsof -tiTCP:"$1" -sTCP:LISTEN 2>/dev/null | head -n1 || true
}

managed_ports() {
  local port
  declare -A seen=()

  for port in "$PORT" "$FRONTEND_DEV_PORT" "$VITE_DEV_PORT" "$MODBUS_SHARE_PORT"; do
    if [[ "$port" =~ ^[0-9]+$ && -z "${seen[$port]+x}" ]]; then
      seen["$port"]=1
      printf "%s\n" "$port"
    fi
  done

  for port in {8080..8090}; do
    if [[ -z "${seen[$port]+x}" ]]; then
      seen["$port"]=1
      printf "%s\n" "$port"
    fi
  done
}

clear_listen_port() {
  local port="$1"
  local label="${2:-端口}"
  local force_cleanup="${3:-false}"
  local pid
  pid="$(pid_on_port "$port")"
  if [[ -z "$pid" ]]; then
    info "$label $port 可用"
    return 0
  fi

  if [[ "$AUTO_KILL_PORT" == true || "$force_cleanup" == true ]]; then
    warn "$label $port 被 PID $pid 佔用，將先清理再啟動"
    stop_process_tree "$pid" "$label 佔用進程"
  else
    warn "$label $port 被 PID $pid 佔用"
    read -r -p "是否終止 $label ${port} 的進程? (y/N): " ans
    if [[ "$ans" =~ ^[Yy]$ ]]; then
      stop_process_tree "$pid" "$label 佔用進程"
    else
      return 1
    fi
  fi

  if [[ -n "$(pid_on_port "$port")" ]]; then
    err "$label $port 清理失敗"
    return 1
  fi
  success "$label $port 已清理"
}

clear_port() {
  clear_listen_port "$PORT" "後端端口" "${1:-false}"
}

clear_frontend_port() {
  clear_listen_port "$FRONTEND_DEV_PORT" "前端端口" "${1:-false}"
}

build_frontend() {
  if [[ ! -d "$FRONTEND_DIR" || ! -f "$FRONTEND_DIR/package.json" ]]; then
    err "找不到 frontend source 或 package.json"
    return 1
  fi

  pushd "$FRONTEND_DIR" >/dev/null
  if [[ ! -d node_modules ]]; then
    info "安裝前端依賴..."
    pnpm install
  fi
  info "建置前端..."
  pnpm run build
  popd >/dev/null

  mkdir -p "$APP_PATH/static"
  find "$APP_PATH/static" -mindepth 1 -maxdepth 1 ! -name 'embed-placeholder.txt' -exec rm -rf {} +
  cp -R "$FRONTEND_DIR/dist/." "$APP_PATH/static/"
  success "前端建置完成"
}

sync_embedded_frontend_if_requested() {
  if [[ "$SYNC_EMBED" != true ]]; then
    info "略過 embedded 前端同步（開發模式預設使用 Vite 即時伺服器）"
    return 0
  fi

  info "📦 同步 embedded 前端資產..."
  build_frontend
}

build_app() {
  build_frontend || return 1
  mkdir -p "$BUILD_DIR"
  info "建置後端..."
  go build -ldflags "-s -w" -trimpath -o "$BUILD_DIR/$APP_NAME" "./$APP_PATH"
  success "建置完成: $BUILD_DIR/$APP_NAME"
}

build_single() {
  build_frontend || return 1
  mkdir -p "$BUILD_DIR"
  info "建置單一執行檔 (embed)..."
  go build -tags embed -ldflags "-s -w" -trimpath -o "$BUILD_DIR/$APP_NAME" "./$APP_PATH"
  success "單一執行檔完成: $BUILD_DIR/$APP_NAME"
}

start_app() {
  local exe="$BUILD_DIR/$APP_NAME"
  [[ -x "$exe" ]] || { err "找不到可執行檔: $exe"; return 1; }
  clear_port true
  info "啟動服務於 :$PORT"
  PORT="$PORT" "$exe"
}

health_check() {
  info "健康檢查 http://localhost:$PORT/api/v1/test/status"
  if curl -fsS "http://localhost:$PORT/api/v1/test/status" >/dev/null 2>&1; then
    success "服務正常"
  else
    err "服務不可用"
    return 1
  fi
}

list_processes() {
  info "列出運行中的服務..."
  local found=false
  local pids

  pids="$(related_process_pids | sort -u || true)"
  if [[ -n "$pids" ]]; then
    found=true
    info "相關進程:"
    while IFS= read -r pid; do
      [[ -z "$pid" ]] && continue
      ps -p "$pid" -o pid,ppid,etime,rss,command 2>/dev/null || true
    done <<<"$pids"
  fi

  info "端口佔用情況:"
  local p pid
  while IFS= read -r p; do
    pid="$(pid_on_port "$p")"
    if [[ -n "$pid" ]]; then
      found=true
      printf "  %s -> PID %s\n" "$p" "$pid"
    fi
  done < <(managed_ports)

  [[ "$found" == true ]] || info "未發現運行中的服務"
}

stop_all_services() {
  info "停止所有運行中的服務..."
  local stopped=0
  local pids pid
  pids="$(related_process_pids | sort -u || true)"

  if [[ -n "$pids" ]]; then
    while IFS= read -r pid; do
      [[ -z "$pid" ]] && continue
      [[ "$pid" == "$$" || "$pid" == "$BASHPID" ]] && continue
      stop_process_tree "$pid" "服務進程"
      stopped=$((stopped + 1))
    done <<<"$pids"
  fi

  local p
  while IFS= read -r p; do
    pid="$(pid_on_port "$p")"
    if [[ -n "$pid" ]] && is_managed_pid "$pid"; then
      stop_process_tree "$pid" "端口 $p 進程"
      stopped=$((stopped + 1))
    fi
  done < <(managed_ports)

  if [[ -d "$TMP_DIR" ]]; then
    find "$TMP_DIR" -maxdepth 1 -type d -name 'start-*.lock' -exec rm -rf {} + 2>/dev/null || true
  fi

  if [[ "$stopped" -eq 0 ]]; then
    info "未發現需要停止的服務"
  else
    success "已停止 $stopped 個服務/進程"
  fi
}

diagnose() {
  info "快速診斷..."
  local issues=0

  info "1) 檢查開發環境"
  if ! has_cmd go; then err "Go 未安裝"; issues=$((issues + 1)); fi
  if ! has_cmd node; then warn "Node.js 未安裝（前端功能可能無法使用）"; issues=$((issues + 1)); fi
  if ! has_cmd pnpm; then warn "pnpm 未安裝（前端功能可能無法使用）"; issues=$((issues + 1)); fi

  info "2) 檢查專案結構"
  [[ -f go.mod ]] || { err "go.mod 不存在，可能不在專案根目錄"; issues=$((issues + 1)); }
  [[ -d "$APP_PATH" ]] || { err "應用程式目錄不存在: $APP_PATH"; issues=$((issues + 1)); }

  info "3) 檢查端口狀態"
  local pid p
  while IFS= read -r p; do
    pid="$(pid_on_port "$p")"
    if [[ -n "$pid" ]]; then
      warn "端口 $p 被 PID $pid 佔用"
      ps -p "$pid" -o pid,ppid,etime,rss,command 2>/dev/null || true
    else
      success "端口 $p 可用"
    fi
  done < <(managed_ports)

  info "4) 檢查前端"
  if [[ -d "$FRONTEND_DIR" && ! -d "$FRONTEND_DIR/node_modules" ]]; then
    warn "前端依賴未安裝，建議執行: cd frontend && pnpm install"
    issues=$((issues + 1))
  fi

  info "5) 檢查構建產物"
  if [[ -x "$BUILD_DIR/$APP_NAME" ]]; then
    success "已存在可執行檔: $BUILD_DIR/$APP_NAME"
  else
    warn "尚未發現可執行檔: ${BUILD_DIR}/${APP_NAME}（若尚未建置屬正常）"
  fi

  info "6) 檢查相關進程與記憶體"
  local pids
  pids="$(related_process_pids | sort -u || true)"
  if [[ -n "$pids" ]]; then
    while IFS= read -r pid; do
      [[ -z "$pid" ]] && continue
      ps -p "$pid" -o pid,ppid,etime,rss,command 2>/dev/null || true
    done <<<"$pids"
  else
    success "未發現相關長駐進程"
  fi

  if [[ "$issues" -eq 0 ]]; then
    success "診斷完成，未發現阻塞問題"
  else
    warn "診斷完成，發現 $issues 個潛在問題"
  fi
}

start_dev_mode() {
  info "🚀 開發模式（無需編譯）"
  info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後需要手動重新運行"
  info "💡 將同時啟動前端開發伺服器和後端服務"

  install_cleanup_traps
  acquire_run_lock || {
    clear_cleanup_traps
    return 1
  }
  sync_embedded_frontend_if_requested

  clear_port true
  info "▶️  啟動後端應用程式（使用 go run）..."

  # 重置噪音計數器
  LOG_NOISE_COUNTERS=()
  start_backend_process "cd '$APP_PATH' && PORT='$PORT' exec go run ."
  wait_for_port_ready "$PORT" "後端端口" 300 || {
    cleanup_all
    clear_cleanup_traps
    return 1
  }

  if [[ -d "$FRONTEND_DIR" ]]; then
    info "🎨 啟動前端開發伺服器... (host=$FRONTEND_DEV_HOST)"
    # 顯式傳遞 backend PORT / proxy target，避免 dev server 代理到錯的埠。
    start_frontend_dev_server || return 1
    success "前端開發伺服器已啟動（PID: ${FRONTEND_PID}）"
    show_frontend_host_hint
    show_embedded_frontend_hint
  fi

  wait_for_backend_exit
}

start_air_mode() {
  info "🔥 熱重載模式（Air）"
  info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後會自動重新運行"
  info "💡 將同時啟動前端開發伺服器和後端服務"

  install_cleanup_traps
  acquire_run_lock || {
    clear_cleanup_traps
    return 1
  }
  sync_embedded_frontend_if_requested

  # 檢查 Air 是否安裝，若無則自動安裝
  local air_cmd
  if ! air_cmd="$(find_air_cmd)"; then
    warn "Air 工具未安裝，正在嘗試安裝..."
    if has_cmd go; then
      go install github.com/air-verse/air@latest
      if air_cmd="$(find_air_cmd)"; then
        success "Air 安裝成功"
      else
        err "無法安裝 Air，請手動安裝："
        info "  go install github.com/air-verse/air@latest"
        info "💡 建議: 使用選項 1 開發模式（go run）作為替代方案"
        return 1
      fi
    else
      err "未安裝 go，無法自動安裝 air"
      return 1
    fi
  fi

  if [[ ! -f ".air.toml" ]]; then
    warn ".air.toml 配置檔案不存在，Air 將使用預設配置"
  fi

  clear_port true

  info "▶️  啟動 Air 熱重載... ($air_cmd)"
  info "💡 修改程式碼後，Air 會自動編譯 $AIR_BIN 並重新啟動"

  # 重置噪音計數器
  LOG_NOISE_COUNTERS=()
  start_backend_process "PORT='$PORT' exec '$air_cmd'"
  wait_for_port_ready "$PORT" "後端端口" 300 || {
    cleanup_all
    clear_cleanup_traps
    return 1
  }

  # 啟動前端開發伺服器
  if [[ -d "$FRONTEND_DIR" ]]; then
    info "🎨 啟動前端開發伺服器... (host=$FRONTEND_DEV_HOST)"
    start_frontend_dev_server || return 1
    success "前端開發伺服器已啟動（PID: ${FRONTEND_PID}）"
    show_frontend_host_hint
    show_embedded_frontend_hint
    info "💡 前端修改會自動熱重載"
  fi

  wait_for_backend_exit
}

run_lint() {
  has_cmd golangci-lint || {
    warn "golangci-lint 未安裝，略過"
    return 0
  }
  info "執行 golangci-lint..."
  if [[ "$VERBOSE" == true ]]; then
    golangci-lint run ./...
  else
    local out
    out="$(golangci-lint run ./... 2>&1)" || {
      printf "%s\n" "$out"
      return 1
    }
  fi
  success "lint 完成"
}

run_quality() {
  info "執行代碼質量檢查..."

  if has_cmd deadcode; then
    local deadcode_out
    deadcode_out="$(deadcode ./... 2>&1)" || true
    if [[ -n "$deadcode_out" ]]; then
      warn "發現可能未使用代碼（非阻塞）"
      [[ "$VERBOSE" == true ]] && printf "%s\n" "$deadcode_out"
    else
      success "未發現未使用代碼"
    fi
  else
    warn "deadcode 未安裝，略過未使用代碼檢查"
  fi

  local fmt_out
  fmt_out="$(go fmt ./... 2>&1)" || {
    err "go fmt 檢查失敗"
    [[ -n "$fmt_out" ]] && printf "%s\n" "$fmt_out"
    return 1
  }
  if [[ -n "$fmt_out" ]]; then
    warn "go fmt 已自動修正部分檔案"
    [[ "$VERBOSE" == true ]] && printf "%s\n" "$fmt_out"
  else
    success "代碼格式正確"
  fi
}

run_test() {
  info "執行測試..."
  mkdir -p "$BUILD_DIR"
  if [[ "$COVERAGE" == true ]]; then
    go test -v -coverprofile="$BUILD_DIR/coverage.out" -covermode=atomic ./internal/protocol/...
    if [[ -f "$BUILD_DIR/coverage.out" ]]; then
      go tool cover -html="$BUILD_DIR/coverage.out" -o "$BUILD_DIR/coverage.html"
      success "覆蓋率報告已生成: $BUILD_DIR/coverage.html"
    fi
  else
    go test -cover ./internal/protocol/...
  fi
  success "測試完成"
}

if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then
  return 0
fi

while [[ $# -gt 0 ]]; do
  HAS_ANY_PARAM=true
  case "$1" in
    --dev-mode) DEV_MODE=true ;;
    --air-mode) AIR_MODE=true ;;
    --quick-start) QUICK_START=true ;;
    --start) START=true ;;
    --skip-build) SKIP_BUILD=true ;;
    --skip-lint) SKIP_LINT=true ;;
    --skip-test) SKIP_TEST=true ;;
    --skip-quality) SKIP_QUALITY=true ;;
    --sync-embed) SYNC_EMBED=true ;;
    --coverage) COVERAGE=true ;;
    --verbose) VERBOSE=true ;;
    --auto-kill-port) AUTO_KILL_PORT=true ;;
    --check-env) CHECK_ENV=true ;;
    --health-check) HEALTH_CHECK=true ;;
    --list-processes) LIST_PROCESSES=true ;;
    --stop-all) STOP_ALL=true ;;
    --diagnose) DIAGNOSE=true ;;
    --build-single) BUILD_SINGLE=true ;;
    --port)
      shift
      PORT="${1:-}"
      [[ -n "$PORT" ]] || { err "--port 需要數值"; exit 1; }
      ;;
    -h|--help) usage; exit 0 ;;
    *) err "未知參數: $1"; usage; exit 1 ;;
  esac
  shift
done

if [[ "$HAS_ANY_PARAM" != true ]]; then
  show_main_menu
  read -r -p "請輸入選項 (0-9, A-D): " menu_selection
  case "$menu_selection" in
    1)
      DEV_MODE=true
      ;;
    2)
      AIR_MODE=true
      ;;
    3)
      QUICK_START=true
      ;;
    4)
      ;;
    5)
      SKIP_LINT=true
      SKIP_TEST=true
      SKIP_QUALITY=true
      ;;
    6)
      SKIP_LINT=true
      SKIP_TEST=true
      SKIP_QUALITY=true
      START=true
      ;;
    7)
      SKIP_LINT=true
      SKIP_BUILD=true
      SKIP_QUALITY=true
      ;;
    8)
      CHECK_ENV=true
      ;;
    9)
      HEALTH_CHECK=true
      ;;
    A|a)
      LIST_PROCESSES=true
      ;;
    B|b)
      STOP_ALL=true
      ;;
    C|c)
      DIAGNOSE=true
      ;;
    D|d)
      BUILD_SINGLE=true
      ;;
    0)
      info "退出"
      exit 0
      ;;
    *)
      err "無效的選項"
      exit 1
      ;;
  esac
fi

if [[ "$CHECK_ENV" == true ]]; then
  check_env
  exit $?
fi

if [[ "$HEALTH_CHECK" == true ]]; then
  health_check
  exit $?
fi

if [[ "$LIST_PROCESSES" == true ]]; then
  list_processes
  exit $?
fi

if [[ "$STOP_ALL" == true ]]; then
  stop_all_services
  exit $?
fi

if [[ "$DIAGNOSE" == true ]]; then
  diagnose
  exit $?
fi

if [[ "$DEV_MODE" == true ]]; then
  start_dev_mode
  exit $?
fi

if [[ "$AIR_MODE" == true ]]; then
  start_air_mode
  exit $?
fi

if [[ "$BUILD_SINGLE" == true ]]; then
  build_single
  exit $?
fi

if [[ "$QUICK_START" == true ]]; then
  run_lint
  build_app
  start_app
  exit $?
fi

if [[ "$SKIP_LINT" != true ]]; then
  run_lint
fi

if [[ "$SKIP_BUILD" != true ]]; then
  build_app
fi

if [[ "$START" == true ]]; then
  if [[ "$SKIP_BUILD" == true ]]; then
    warn "偵測到 --start --skip-build。為避免 8080 使用過期 embedded 前端，將自動重新建置。"
    build_app
  fi
  start_app
fi

if [[ "$SKIP_QUALITY" != true ]]; then
  run_quality
fi

if [[ "$SKIP_TEST" != true ]]; then
  run_test
fi

success "完成"
