#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT
mkdir -p "$TEST_ROOT/repo/scripts"
cp "$REPO_ROOT/start.sh" "$TEST_ROOT/repo/start.sh"
for helper in load-env.sh start-process-utils.sh; do
  cp "$REPO_ROOT/scripts/$helper" "$TEST_ROOT/repo/scripts/$helper"
done
source "$TEST_ROOT/repo/start.sh"
backend_definition="$(declare -f start_backend_process)"
FRONTEND_DEV_PORT=5175
PORT=8085
validate_dev_environment() { :; }
assert_ports_free() { :; }
ensure_embedded_frontend_snapshot() { :; }
ensure_air() { :; }
ensure_frontend_dependencies() { :; }
install_cleanup_traps() { :; }
acquire_run_lock() { :; }
start_backend_process() { BACKEND_PID=101; }
start_frontend_dev_server() { FRONTEND_PID=202; }
wait_for_port_ready() { [[ "$1" != "$FRONTEND_DEV_PORT" ]] || frontend_ready=true; return 0; }
wait_for_backend_exit() { printf 'monitor\n' >>"$TEST_ROOT/events"; }
uname() { printf '%s\n' "$test_os"; }
open() {
  [[ "$frontend_ready" == true ]] || printf 'opened-too-soon\n' >>"$TEST_ROOT/events"
  printf '%s\n' "$1" >>"$TEST_ROOT/browser"
}
xdg-open() { open "$@"; }

for answer in '' n N y Y yes; do
  MODE=""
  test_os=Darwin
  frontend_ready=false
  OPEN_BROWSER=true
  rm -f "$TEST_ROOT/browser" "$TEST_ROOT/events"
  main <<<"$(printf '1\n%s\n' "$answer")" >/dev/null
  case "$answer" in
    y|Y|yes)
      [[ "$(cat "$TEST_ROOT/browser")" == 'http://127.0.0.1:5175/studio/v2' ]]
      ;;
    *) [[ ! -f "$TEST_ROOT/browser" ]] ;;
  esac
  [[ "$(cat "$TEST_ROOT/events")" == monitor ]]
done

# EOF also defaults to no, and Linux dispatches through xdg-open.
request_open_browser </dev/null >/dev/null
[[ "$OPEN_BROWSER" == false ]]
rm -f "$TEST_ROOT/browser"
test_os=Linux
request_open_browser <<<y >/dev/null
open_dev_browser
[[ "$(cat "$TEST_ROOT/browser")" == 'http://127.0.0.1:5175/studio/v2' ]]

# An opener failure must not stop the dev session.
xdg-open() { return 27; }
open_dev_browser >/dev/null 2>&1

# Air and every backend reload must ignore inherited auto-open configuration.
eval "$backend_definition"
cat >"$TEST_ROOT/air" <<'EOF'
#!/usr/bin/env bash
printf '%s' "$AUTO_OPEN_BROWSER" >"$BROWSER_WITNESS"
EOF
chmod +x "$TEST_ROOT/air"
AIR_BIN="$TEST_ROOT/air"
AIR_CONFIG="$TEST_ROOT/.air.toml"
HOST_GOOS=darwin
HOST_GOARCH=arm64
export AUTO_OPEN_BROWSER=true BROWSER_WITNESS="$TEST_ROOT/child-env"
start_backend_process
wait "$BACKEND_PID"
[[ "$(cat "$BROWSER_WITNESS")" == false ]]
[[ "$AUTO_OPEN_BROWSER" == true ]]
printf 'PASS: browser prompt defaults, menu, readiness, openers, backend env\n'
