#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TEST_TMP_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_TMP_DIR/scripts" "$TEST_TMP_DIR/frontend" "$TEST_TMP_DIR/bin"
cp "$REPO_ROOT/start.sh" "$TEST_TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-log-utils.sh" "$TEST_TMP_DIR/scripts/start-log-utils.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_TMP_DIR/scripts/start-process-utils.sh"

cat >"$TEST_TMP_DIR/bin/pnpm" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

if [[ "${1:-}" != "install" ]]; then
  printf 'unexpected pnpm args: %s\n' "$*" >&2
  exit 1
fi

printf 'simulated install failure\n' >&2
exit 41
EOF
chmod +x "$TEST_TMP_DIR/bin/pnpm"

cd "$TEST_TMP_DIR"
REPO_ROOT="$TEST_TMP_DIR"
source ./start.sh >/dev/null 2>&1 || true

export PATH="$TEST_TMP_DIR/bin:$PATH"
FRONTEND_DIR="frontend"

set +e
output="$(start_frontend_dev_server 2>&1)"
exit_code=$?
set -e

if [[ "$exit_code" -eq 0 ]]; then
  printf 'expected start_frontend_dev_server to fail when pnpm install fails\n' >&2
  exit 1
fi

if [[ "$output" != *"前端依賴安裝失敗"* ]]; then
  printf 'expected install failure to be reported clearly, got:\n%s\n' "$output" >&2
  exit 1
fi
