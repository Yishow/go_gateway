#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT
mkdir -p "$TEST_ROOT/scripts" "$TEST_ROOT/bin/tools"
cp "$REPO_ROOT/start.sh" "$TEST_ROOT/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_ROOT/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_ROOT/scripts/start-process-utils.sh"
cat >"$TEST_ROOT/bin/tools/fake-air" <<'EOF'
#!/usr/bin/env bash
printf '%s|%s|%s|%s\n' "$PORT" "$GOOS" "$GOARCH" "$CGO_ENABLED" > "$AIR_WITNESS"
exit 0
EOF
chmod +x "$TEST_ROOT/bin/tools/fake-air"

cd "$TEST_ROOT"
source ./start.sh >/dev/null 2>&1 || true
AIR_BIN="$TEST_ROOT/bin/tools/fake-air"
AIR_CONFIG="$TEST_ROOT/.air.toml"
HOST_GOOS=darwin
HOST_GOARCH=arm64
PORT=4381
AIR_WITNESS="$TEST_ROOT/air-witness"
export AIR_WITNESS
start_backend_process
wait "$BACKEND_PID" || true

if [[ "$(cat "$AIR_WITNESS")" != "4381|darwin|arm64|0" ]]; then
  printf 'expected Air to inherit host build environment, got %s\n' "$(cat "$AIR_WITNESS")" >&2
  exit 1
fi
