#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_TMP_DIR/scripts" "$TEST_TMP_DIR/frontend"
cp "$REPO_ROOT/start.sh" "$TEST_TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TEST_TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_TMP_DIR/scripts/start-process-utils.sh"

cat > "$TEST_TMP_DIR/.env" <<'EOF'
PORT=3333
FRONTEND_DEV_HOST=0.0.0.0
FRONTEND_DEV_PORT=4173
VITE_DEV_PORT=4173
EOF

cd "$TEST_TMP_DIR"
REPO_ROOT="$TEST_TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true
PORT=3333
FRONTEND_DEV_HOST=0.0.0.0
FRONTEND_DEV_PORT=4173
VITE_DEV_PORT=4173

pid_on_port() {
  case "$1" in
    3333) printf '1234\n' ;;
    4173) printf '5678\n' ;;
  esac
}

pgrep() {
  return 1
}

list_output="$(list_processes 2>&1 || true)"
if [[ "$list_output" != *"3333 -> PID 1234"* ]]; then
  printf 'expected list_processes to report configured backend port, got:\n%s\n' "$list_output" >&2
  exit 1
fi
if [[ "$list_output" != *"4173 -> PID 5678"* ]]; then
  printf 'expected list_processes to report configured frontend port, got:\n%s\n' "$list_output" >&2
  exit 1
fi

killed_pids=""
kill() {
  killed_pids+="$*;"
}

ps() {
  case "$*" in
    *"1234"*)
      printf '/Users/test/Library/Caches/go-build/anonymous-binary\n'
      ;;
    *"5678"*)
      printf '/Users/test/go_gateway/frontend/node_modules/.bin/vite --host 0.0.0.0 --port 4173\n'
      ;;
  esac
}

lsof() {
  case "$*" in
    *"-p 1234 -d txt -Fn"*)
      printf 'p1234\nftxt\nn/Users/test/Library/Caches/go-build/68/test_ui\n'
      ;;
    *"-p 1234 -d cwd -Fn"*)
      printf 'p1234\nfcwd\nn/Users/test/go_gateway\n'
      ;;
  esac
}

stop_all_services >/dev/null 2>&1 || true

if [[ "$killed_pids" != *"1234"* ]]; then
  printf 'expected stop_all_services to kill configured backend port listener, got %s\n' "$killed_pids" >&2
  exit 1
fi
if [[ "$killed_pids" != *"5678"* ]]; then
  printf 'expected stop_all_services to kill configured frontend port listener, got %s\n' "$killed_pids" >&2
  exit 1
fi

# 現行 start_frontend_dev_server 以子 shell 執行 `pnpm exec vite`；
# 以 fake pnpm/vite 捕獲實際指令參數
captured_command_file="$TEST_TMP_DIR/vite-args.txt"
mkdir -p "$TEST_TMP_DIR/frontend/node_modules/.bin" "$TEST_TMP_DIR/bin"
cat >"$TEST_TMP_DIR/frontend/node_modules/.bin/vite" <<EOF
#!/usr/bin/env bash
printf '%s\n' "\$*" >"$captured_command_file"
exit 0
EOF
chmod +x "$TEST_TMP_DIR/frontend/node_modules/.bin/vite"

cat >"$TEST_TMP_DIR/bin/pnpm" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

if [[ "${1:-}" != "exec" ]]; then
  printf 'unexpected pnpm args: %s\n' "$*" >&2
  exit 1
fi

shift
exec "./node_modules/.bin/$1" "${@:2}"
EOF
chmod +x "$TEST_TMP_DIR/bin/pnpm"

export PATH="$TEST_TMP_DIR/bin:$PATH"
clear_frontend_port() {
  return 0
}

start_frontend_dev_server
for _ in $(seq 1 20); do
  [[ -f "$captured_command_file" ]] && break
  sleep 0.1
done
frontend_command="$(cat "$captured_command_file" 2>/dev/null || true)"
if [[ -z "$frontend_command" ]]; then
  printf 'expected frontend dev server to invoke vite, got no command captured\n' >&2
  exit 1
fi
if [[ "$frontend_command" != *"--strictPort"* ]]; then
  printf 'expected frontend dev server command to enforce a fixed port, got:\n%s\n' "$frontend_command" >&2
  exit 1
fi
if [[ "$frontend_command" != *"--port 4173"* ]]; then
  printf 'expected frontend dev server command to use configured port, got:\n%s\n' "$frontend_command" >&2
  exit 1
fi

occupied_pid=7777
prompted=false
killed_pids=""
pid_on_port() {
  case "$1" in
    4173)
      if [[ -n "$occupied_pid" ]]; then
        printf '%s\n' "$occupied_pid"
      fi
      ;;
  esac
}
read() {
  # 僅將帶 -p 的互動提示視為 prompting；
  # kill_process_tree 內部也會用 `read -r child` 迭代子進程，不可計入
  case "$*" in
    *-p*)
      prompted=true
      return 1
      ;;
    *)
      return 1
      ;;
  esac
}
kill() {
  killed_pids+="$*;"
  occupied_pid=""
}

auto_clear_output_file="$TEST_TMP_DIR/auto-clear-output.txt"
clear_listen_port 4173 "前端端口" true >"$auto_clear_output_file" 2>&1 || true
auto_clear_output="$(cat "$auto_clear_output_file")"
if [[ "$prompted" == true ]]; then
  printf 'expected forced startup cleanup to avoid prompting, got:\n%s\n' "$auto_clear_output" >&2
  exit 1
fi
if [[ "$killed_pids" != *"7777"* ]]; then
  printf 'expected forced startup cleanup to kill occupied frontend port listener, got %s\n' "$killed_pids" >&2
  exit 1
fi
if [[ "$auto_clear_output" != *"已清理"* ]]; then
  printf 'expected forced startup cleanup to report deletion, got:\n%s\n' "$auto_clear_output" >&2
  exit 1
fi
