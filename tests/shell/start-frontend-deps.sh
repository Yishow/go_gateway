#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"

cleanup() {
  if [[ -n "${FRONTEND_PID:-}" ]]; then
    kill "$FRONTEND_PID" 2>/dev/null || true
    wait "$FRONTEND_PID" 2>/dev/null || true
  fi
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_TMP_DIR/scripts" "$TEST_TMP_DIR/frontend" "$TEST_TMP_DIR/bin"
cp "$REPO_ROOT/start.sh" "$TEST_TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TEST_TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_TMP_DIR/scripts/start-process-utils.sh"

cat >"$TEST_TMP_DIR/frontend/package.json" <<'EOF'
{
  "name": "frontend",
  "private": true
}
EOF

cat >"$TEST_TMP_DIR/bin/pnpm" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

case "${1:-}" in
  install)
    touch .install-ran
    mkdir -p node_modules/.bin
    cat >node_modules/.bin/vite <<'INNER'
#!/usr/bin/env bash
sleep 30
INNER
    chmod +x node_modules/.bin/vite
    ;;
  exec)
    touch .exec-ran
    shift
    if [[ ! -x "./node_modules/.bin/$1" ]]; then
      touch .exec-missing-vite
      exit 1
    fi
    touch .exec-found-vite
    exec "./node_modules/.bin/$1" "${@:2}"
    ;;
  *)
    printf 'unexpected pnpm args: %s\n' "$*" >&2
    exit 1
    ;;
esac
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
  if [[ "$1" == "$FRONTEND_DEV_PORT" && -n "${FRONTEND_PID:-}" ]]; then
    printf '%s\n' "$FRONTEND_PID"
  fi
}

start_frontend_dev_server >/dev/null 2>&1
sleep 0.1

if [[ ! -f "$TEST_TMP_DIR/frontend/.install-ran" ]]; then
  printf 'expected start_frontend_dev_server to install frontend dependencies when vite is missing\n' >&2
  exit 1
fi
if [[ ! -f "$TEST_TMP_DIR/frontend/.exec-found-vite" ]]; then
  printf 'expected start_frontend_dev_server to run vite only after dependencies are ready\n' >&2
  exit 1
fi
