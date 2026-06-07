#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

mkdir -p \
  "$TEST_TMP_DIR/scripts" \
  "$TEST_TMP_DIR/frontend/node_modules/.bin" \
  "$TEST_TMP_DIR/bin"
cp "$REPO_ROOT/start.sh" "$TEST_TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TEST_TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_TMP_DIR/scripts/start-process-utils.sh"

cat >"$TEST_TMP_DIR/frontend/node_modules/.bin/vite" <<'EOF'
#!/usr/bin/env bash
exit 23
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

cd "$TEST_TMP_DIR"
REPO_ROOT="$TEST_TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

export PATH="$TEST_TMP_DIR/bin:$PATH"
FRONTEND_DIR="frontend"
FRONTEND_DEV_PORT=4173

clear_frontend_port() {
  return 0
}

pid_on_port() {
  return 0
}

sleep() {
  :
}

set +e
output="$(start_frontend_dev_server 2>&1)"
exit_code=$?
set -e

if [[ "$exit_code" -eq 0 ]]; then
  printf 'expected start_frontend_dev_server to fail when vite exits before binding the port\n' >&2
  exit 1
fi

if [[ "$output" != *"前端開發伺服器未成功啟動"* ]]; then
  printf 'expected start_frontend_dev_server to explain frontend startup failure, got:\n%s\n' "$output" >&2
  exit 1
fi
