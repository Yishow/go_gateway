#!/usr/bin/env bash

set -euo pipefail

APP_NAME="gateway"
APP_PATH="cmd/test_ui"
BUILD_DIR="bin"
FRONTEND_DIR="frontend"
FRONTEND_DEV_HOST="${FRONTEND_DEV_HOST:-0.0.0.0}"
PORT="${PORT:-8080}"
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

info() { printf "\033[36m[INFO]\033[0m %s\n" "$*"; }
success() { printf "\033[32m[OK]\033[0m %s\n" "$*"; }
warn() { printf "\033[33m[WARN]\033[0m %s\n" "$*"; }
err() { printf "\033[31m[ERR]\033[0m %s\n" "$*"; }

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
  if [[ -d "$FRONTEND_DIR" ]]; then
    info "啟動前端開發伺服器... (host=$FRONTEND_DEV_HOST)"
    (cd "$FRONTEND_DIR" && pnpm run dev -- --host "$FRONTEND_DEV_HOST") &
    FRONTEND_PID=$!
    trap 'kill "$FRONTEND_PID" 2>/dev/null || true' EXIT
  fi
  clear_port
  info "啟動後端 go run..."
  (cd "$APP_PATH" && PORT="$PORT" go run .)
}

start_air_mode() {
  local air_cmd
  air_cmd="$(find_air_cmd)" || {
    err "未安裝 air，請先執行: go install github.com/air-verse/air@latest"
    warn "若已安裝，請確認 PATH 包含 \$(go env GOPATH)/bin"
    return 1
  }
  if [[ -d "$FRONTEND_DIR" ]]; then
    info "啟動前端開發伺服器... (host=$FRONTEND_DEV_HOST)"
    (cd "$FRONTEND_DIR" && pnpm run dev -- --host "$FRONTEND_DEV_HOST") &
    FRONTEND_PID=$!
    trap 'kill "$FRONTEND_PID" 2>/dev/null || true' EXIT
  fi
  clear_port
  info "啟動 air 熱重載... ($air_cmd)"
  PORT="$PORT" "$air_cmd"
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
