#!/usr/bin/env bash

set -euo pipefail

APP_NAME="gateway"
APP_PATH="cmd/test_ui"
BUILD_DIR="bin"
FRONTEND_DIR="frontend"
FRONTEND_DEV_HOST="${FRONTEND_DEV_HOST:-0.0.0.0}"
PORT="${PORT:-8080}"
LOG_DIR="bin/logs"
TMP_DIR="bin/tmp"
AUTO_KILL_PORT=false
SKIP_BUILD=false
SKIP_LINT=false
SKIP_TEST=false
SKIP_QUALITY=false
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

# 日誌噪音計數器（關聯陣列）
declare -A LOG_NOISE_COUNTERS

info() { printf "\033[36m[INFO]\033[0m %s\n" "$*"; }
success() { printf "\033[32m[OK]\033[0m %s\n" "$*"; }
warn() { printf "\033[33m[WARN]\033[0m %s\n" "$*"; }
err() { printf "\033[31m[ERR]\033[0m %s\n" "$*"; }

# ============================================
# 日誌格式化與噪音過濾（對齊 PS1 功能）
# ============================================

# 判斷日誌行是否為噪音類別
# 回傳噪音類別名稱，若非噪音則回傳空字串
get_log_noise_category() {
  local line="$1"
  # air watcher 常見輸出
  if [[ "$line" =~ ^watching|^building\.\.\.|^\!exclude|^[[:space:]]*/|^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
    printf "air-watcher"
    return
  fi
  # modbus status polling
  if [[ "$line" == *"/api/v1/datalink/modbus-share/status"* ]]; then
    printf "status-polling"
    return
  fi
  # dashboard 輪詢請求
  if [[ "$line" =~ \/api\/v1\/datalink\/(devices|polling-groups|points|mappings|tags) ]]; then
    printf "dashboard-refresh"
    return
  fi
  # 啟動細節日誌
  if [[ "$line" == *"資料庫路徑"* || "$line" == *"Executing SQLite migration"* || "$line" == *"ConnectionManager 已初始化"* || "$line" == *"已註冊的協議"* ]]; then
    printf "startup-detail"
    return
  fi
  printf ""
}

# 格式化單行執行時日誌輸出（對齊 PS1 Write-RuntimeLogLine）
runtime_log_line() {
  local line="$1"
  [[ -z "$line" ]] && return

  # HTTP 請求行格式化：[timestamp] METHOD /path STATUS latency
  if [[ "$line" =~ \[([^\]]+)\][[:space:]]+[^[:space:]]+[[:space:]]+(GET|POST|PUT|DELETE|PATCH)[[:space:]]+([^[:space:]]+)[[:space:]]+([0-9]{3})[[:space:]]+(.+) ]]; then
    local ts="${BASH_REMATCH[1]}"
    local method="${BASH_REMATCH[2]}"
    local path="${BASH_REMATCH[3]}"
    local status="${BASH_REMATCH[4]}"
    local latency="${BASH_REMATCH[5]}"
    # 擷取時間部分 HH:MM:SS
    local time_part
    if [[ "$ts" =~ ([0-9]{2}:[0-9]{2}:[0-9]{2})$ ]]; then
      time_part="${BASH_REMATCH[1]}"
    else
      time_part="--:--:--"
    fi

    # dashboard 輪詢在非 verbose 模式下隱藏
    if [[ "$method" == "GET" && "$status" == "200" && "$path" =~ ^/api/v1/datalink/(devices|polling-groups|points|mappings|tags|modbus-share/status)$ ]]; then
      if [[ "$VERBOSE" != true ]]; then
        LOG_NOISE_COUNTERS["dashboard-refresh"]=$(( ${LOG_NOISE_COUNTERS["dashboard-refresh"]:-0} + 1 ))
        return
      fi
    fi

    local formatted
    formatted=$(printf "[HTTP] %s  %-6s %-46.46s %3s %9s" "$time_part" "$method" "$path" "$status" "$latency")
    if (( status >= 500 )); then
      printf "\033[31m%s\033[0m\n" "$formatted"
    elif (( status >= 400 )); then
      printf "\033[33m%s\033[0m\n" "$formatted"
    else
      printf "\033[90m%s\033[0m\n" "$formatted"
    fi
    return
  fi

  # Go 標準日誌行：yyyy/mm/dd HH:MM:SS file:line: msg
  if [[ "$line" =~ ^[0-9]{4}/[0-9]{2}/[0-9]{2}[[:space:]]+[0-9]{2}:[0-9]{2}:[0-9]{2}[[:space:]]+[^:]+:[0-9]+:[[:space:]]+(.+)$ ]]; then
    local msg="${BASH_REMATCH[1]}"
    local category
    category="$(get_log_noise_category "$line")"
    if [[ -n "$category" && "$VERBOSE" != true ]]; then
      LOG_NOISE_COUNTERS["$category"]=$(( ${LOG_NOISE_COUNTERS["$category"]:-0} + 1 ))
      return
    fi
    if [[ "$msg" == *"127.0.0.1:5020"* ]]; then
      printf "\033[32m[BOOT] Local Modbus share started on 127.0.0.1:5020\033[0m\n"
    elif [[ "$msg" == *"localhost:8080"* ]]; then
      printf "\033[32m[BOOT] Server started at http://localhost:8080\033[0m\n"
    elif [[ "$msg" == *"資料庫路徑"* ]]; then
      printf "\033[90m[BOOT] Database initialized\033[0m\n"
    elif [[ "$VERBOSE" == true ]]; then
      printf "\033[90m[BOOT] %s\033[0m\n" "$msg"
    fi
    return
  fi

  # 一般噪音過濾
  local category
  category="$(get_log_noise_category "$line")"
  if [[ -n "$category" && "$VERBOSE" != true ]]; then
    LOG_NOISE_COUNTERS["$category"]=$(( ${LOG_NOISE_COUNTERS["$category"]:-0} + 1 ))
    return
  fi

  # 著色輸出
  if [[ "$line" =~ ERROR|Error|panic|FATAL ]]; then
    printf "\033[31m%s\033[0m\n" "$line"
  elif [[ "$line" =~ WARN|Warning ]]; then
    printf "\033[33m%s\033[0m\n" "$line"
  elif [[ "$line" =~ 啟動於|localhost:8080|本機\ Modbus\ 分享服務已啟動 ]]; then
    printf "\033[32m%s\033[0m\n" "$line"
  else
    printf "\033[90m%s\033[0m\n" "$line"
  fi
}

# 顯示被隱藏的噪音日誌統計（對齊 PS1 Show-LogNoiseSummary）
show_log_noise_summary() {
  [[ ${#LOG_NOISE_COUNTERS[@]} -eq 0 ]] && return
  [[ "$VERBOSE" == true ]] && return
  printf "\n\033[36m============================================\033[0m\n"
  printf "\033[36m   已隱藏雜訊日誌\033[0m\n"
  printf "\033[36m============================================\033[0m\n"
  for key in "${!LOG_NOISE_COUNTERS[@]}"; do
    printf "\033[33m- %s: %s 行\033[0m\n" "$key" "${LOG_NOISE_COUNTERS[$key]}"
  done
  info "可加上 --verbose 顯示全部原始日誌。"
}

# 清理前端開發伺服器及其子進程
cleanup_frontend() {
  if [[ -n "$FRONTEND_PID" ]]; then
    info "正在停止前端開發伺服器（PID: ${FRONTEND_PID}）..."
    # 終止整個進程組
    kill -- -"$FRONTEND_PID" 2>/dev/null || kill "$FRONTEND_PID" 2>/dev/null || true
    wait "$FRONTEND_PID" 2>/dev/null || true
    FRONTEND_PID=""
    success "前端開發伺服器已停止"
  fi
}

# 清理所有子進程
cleanup_all() {
  cleanup_frontend
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
  --quick-start        一鍵啟動 (build + start)
  --start              建置後啟動服務
  --skip-build         跳過建置
  --skip-lint          跳過 golangci-lint
  --skip-test          跳過測試
  --skip-quality       跳過代碼質量檢查
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

is_managed_pid() {
  local pid="$1"
  local cmd
  cmd="$(ps -p "$pid" -o command= 2>/dev/null || true)"
  [[ "$cmd" == *"/bin/gateway"* || "$cmd" == *"cmd/test_ui"* || "$cmd" == *"start.sh"* ]]
}

show_main_menu() {
  cat <<MENU

============================================
   Go Gateway 啟動選單 (mac)
============================================
 1) 開發模式 (go run)
 2) 熱重載模式 (air)
 3) 一鍵啟動 (lint + build + start)
 4) 完整流程 (lint + build + quality + test)
 5) 僅建置
 6) 僅啟動 (skip build)
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

clear_port() {
  local pid
  pid="$(pid_on_port "$PORT")"
  if [[ -z "$pid" ]]; then
    info "端口 $PORT 可用"
    return 0
  fi

  warn "端口 $PORT 被 PID $pid 佔用"
  if [[ "$AUTO_KILL_PORT" == true ]]; then
    kill -9 "$pid" 2>/dev/null || true
    sleep 0.3
  else
    read -r -p "是否終止該進程? (y/N): " ans
    if [[ "$ans" =~ ^[Yy]$ ]]; then
      kill -9 "$pid" 2>/dev/null || true
      sleep 0.3
    else
      return 1
    fi
  fi

  if [[ -n "$(pid_on_port "$PORT")" ]]; then
    err "端口 $PORT 清理失敗"
    return 1
  fi
  success "端口 $PORT 已清理"
}

build_frontend() {
  if [[ ! -d "$FRONTEND_DIR" ]]; then
    warn "找不到 frontend 目錄，略過前端建置"
    return 0
  fi

  pushd "$FRONTEND_DIR" >/dev/null
  if [[ ! -d node_modules ]]; then
    info "安裝前端依賴..."
    pnpm install
  fi
  info "建置前端..."
  pnpm run build
  popd >/dev/null

  rm -rf "$APP_PATH/static"
  mkdir -p "$APP_PATH/static"
  cp -R "$FRONTEND_DIR/dist/." "$APP_PATH/static/"
  success "前端建置完成"
}

build_app() {
  mkdir -p "$BUILD_DIR"
  info "建置後端..."
  go build -ldflags "-s -w" -trimpath -o "$BUILD_DIR/$APP_NAME" "./$APP_PATH"
  success "建置完成: $BUILD_DIR/$APP_NAME"
}

build_single() {
  build_frontend
  mkdir -p "$BUILD_DIR"
  info "建置單一執行檔 (embed)..."
  go build -tags embed -ldflags "-s -w" -trimpath -o "$BUILD_DIR/$APP_NAME" "./$APP_PATH"
  success "單一執行檔完成: $BUILD_DIR/$APP_NAME"
}

start_app() {
  local exe="$BUILD_DIR/$APP_NAME"
  [[ -x "$exe" ]] || { err "找不到可執行檔: $exe"; return 1; }
  clear_port
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

  pids="$(pgrep -f "/bin/gateway|cmd/test_ui|start.sh" 2>/dev/null || true)"
  if [[ -n "$pids" ]]; then
    found=true
    info "相關進程:"
    while IFS= read -r pid; do
      [[ -z "$pid" ]] && continue
      ps -p "$pid" -o pid,etime,command 2>/dev/null || true
    done <<<"$pids"
  fi

  info "端口佔用情況 (8080-8090):"
  local p pid
  for p in {8080..8090}; do
    pid="$(pid_on_port "$p")"
    if [[ -n "$pid" ]]; then
      found=true
      printf "  %s -> PID %s\n" "$p" "$pid"
    fi
  done

  [[ "$found" == true ]] || info "未發現運行中的服務"
}

stop_all_services() {
  info "停止所有運行中的服務..."
  local stopped=0
  local pids pid
  pids="$(pgrep -f "/bin/gateway|cmd/test_ui|start.sh --skip-lint --start|start.sh --start" 2>/dev/null || true)"

  if [[ -n "$pids" ]]; then
    while IFS= read -r pid; do
      [[ -z "$pid" ]] && continue
      kill -9 "$pid" 2>/dev/null || true
      stopped=$((stopped + 1))
    done <<<"$pids"
  fi

  local p
  for p in {8080..8090}; do
    pid="$(pid_on_port "$p")"
    if [[ -n "$pid" ]] && is_managed_pid "$pid"; then
      kill -9 "$pid" 2>/dev/null || true
      stopped=$((stopped + 1))
    fi
  done

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
  local pid
  pid="$(pid_on_port "$PORT")"
  if [[ -n "$pid" ]]; then
    warn "端口 $PORT 被 PID $pid 佔用"
  else
    success "端口 $PORT 可用"
  fi

  info "4) 檢查前端"
  if [[ -d "$FRONTEND_DIR" && ! -d "$FRONTEND_DIR/node_modules" ]]; then
    warn "前端依賴未安裝，建議執行: cd frontend && pnpm install"
    issues=$((issues + 1))
  fi

  info "5) 檢查構建產物"
  if [[ -x "$BUILD_DIR/$APP_NAME" ]]; then
    success "已存在可執行檔: $BUILD_DIR/$APP_NAME"
  else
    warn "尚未發現可執行檔: $BUILD_DIR/$APP_NAME（若尚未建置屬正常）"
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

  # 設定清理 trap
  trap cleanup_all EXIT INT TERM

  if [[ -d "$FRONTEND_DIR" ]]; then
    info "🎨 啟動前端開發伺服器... (host=$FRONTEND_DEV_HOST)"
    # 使用 setsid 建立新進程組，方便整體清理
    if has_cmd setsid; then
      setsid bash -c "cd '$FRONTEND_DIR' && pnpm run dev --host '$FRONTEND_DEV_HOST'" &
    else
      (cd "$FRONTEND_DIR" && pnpm run dev --host "$FRONTEND_DEV_HOST") &
    fi
    FRONTEND_PID=$!
    success "前端開發伺服器已啟動（PID: ${FRONTEND_PID}）"
    info "💡 前端開發伺服器通常運行在 http://localhost:5173"
  fi

  clear_port
  info "▶️  啟動後端應用程式（使用 go run）..."

  # 重置噪音計數器
  LOG_NOISE_COUNTERS=()

  # 使用日誌格式化管道
  (cd "$APP_PATH" && PORT="$PORT" go run .) 2>&1 | while IFS= read -r line; do
    runtime_log_line "$line"
  done

  show_log_noise_summary
}

start_air_mode() {
  info "🔥 熱重載模式（Air）"
  info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後會自動重新運行"
  info "💡 將同時啟動前端開發伺服器和後端服務"

  # 設定清理 trap
  trap cleanup_all EXIT INT TERM

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

  # 啟動前端開發伺服器
  if [[ -d "$FRONTEND_DIR" ]]; then
    info "🎨 啟動前端開發伺服器... (host=$FRONTEND_DEV_HOST)"
    if has_cmd setsid; then
      setsid bash -c "cd '$FRONTEND_DIR' && pnpm run dev --host '$FRONTEND_DEV_HOST'" &
    else
      (cd "$FRONTEND_DIR" && pnpm run dev --host "$FRONTEND_DEV_HOST") &
    fi
    FRONTEND_PID=$!
    success "前端開發伺服器已啟動（PID: ${FRONTEND_PID}）"
    info "💡 前端開發伺服器通常運行在 http://localhost:5173"
    info "💡 前端修改會自動熱重載"
  fi

  clear_port

  # 檢查 .air.toml 是否存在
  if [[ ! -f ".air.toml" ]]; then
    warn ".air.toml 配置檔案不存在，Air 將使用預設配置"
  fi

  info "▶️  啟動 Air 熱重載... ($air_cmd)"
  info "💡 修改程式碼後，Air 會自動檢測並使用 go run 重新運行"

  # 重置噪音計數器
  LOG_NOISE_COUNTERS=()

  # 使用日誌格式化管道（對齊 PS1 行為）
  PORT="$PORT" "$air_cmd" 2>&1 | while IFS= read -r line; do
    runtime_log_line "$line"
  done

  show_log_noise_summary
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
      SKIP_BUILD=true
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
  start_app
fi

if [[ "$SKIP_QUALITY" != true ]]; then
  run_quality
fi

if [[ "$SKIP_TEST" != true ]]; then
  run_test
fi

success "完成"
