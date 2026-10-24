#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT
mkdir -p "$TEST_ROOT/scripts" "$TEST_ROOT/frontend/node_modules/vite/bin" "$TEST_ROOT/cmd/test_ui/static" "$TEST_ROOT/fake-bin"
cp "$REPO_ROOT/start.sh" "$TEST_ROOT/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_ROOT/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_ROOT/scripts/start-process-utils.sh"
printf 'module go\n' >"$TEST_ROOT/go.mod"
printf '{}' >"$TEST_ROOT/frontend/package.json"
touch "$TEST_ROOT/.air.toml" "$TEST_ROOT/frontend/node_modules/vite/bin/vite.js"
printf '<html>ready</html>\n' >"$TEST_ROOT/cmd/test_ui/static/index.html"

cat >"$TEST_ROOT/fake-bin/go" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  env)
    if [[ "${2:-}" == GOHOSTOS ]]; then printf 'darwin\n'; else printf 'arm64\n'; fi
    ;;
  install)
    mkdir -p "$GOBIN"
    cp "$FAKE_AIR_SOURCE" "$GOBIN/air"
    chmod +x "$GOBIN/air"
    printf '%s|%s|%s|%s\n' "$GOBIN" "$GOOS" "$GOARCH" "$GOTOOLCHAIN" >"$FAKE_ROOT/air-install"
    ;;
  *) exit 2 ;;
esac
EOF
cat >"$TEST_ROOT/fake-bin/lsof" <<'EOF'
#!/usr/bin/env bash
case "$*" in
  *TCP:4381*) [[ -f "$FAKE_ROOT/backend-ready" ]] && cat "$FAKE_ROOT/backend-pid" ;;
  *TCP:5381*) [[ -f "$FAKE_ROOT/frontend-ready" ]] && cat "$FAKE_ROOT/frontend-pid" ;;
esac
EOF
cat >"$TEST_ROOT/fake-bin/pgrep" <<'EOF'
#!/usr/bin/env bash
exit 1
EOF
cat >"$TEST_ROOT/fake-bin/npm" <<'EOF'
#!/usr/bin/env bash
exit 41
EOF
cat >"$TEST_ROOT/fake-bin/air-source" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$$" >"$FAKE_ROOT/backend-pid"
touch "$FAKE_ROOT/backend-ready"
trap 'touch "$FAKE_ROOT/backend-stopped"; exit 0' TERM INT
while :; do sleep 1; done
EOF
cat >"$TEST_ROOT/fake-bin/node" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$$" >"$FAKE_ROOT/frontend-pid"
touch "$FAKE_ROOT/frontend-ready"
sleep 1
EOF
chmod +x "$TEST_ROOT/fake-bin"/*

export FAKE_ROOT="$TEST_ROOT"
export FAKE_AIR_SOURCE="$TEST_ROOT/fake-bin/air-source"
export PATH="$TEST_ROOT/fake-bin:$PATH"
export PORT=4381 FRONTEND_DEV_PORT=5381
set +e
output="$(cd "$TEST_ROOT" && ./start.sh --dev 2>&1)"
result=$?
set -e

if [[ "$result" -eq 0 || "$output" != *"前端開發伺服器已停止"* ]]; then
  printf 'expected frontend failure to fail the dev session, result=%s output=%s\n' "$result" "$output" >&2
  exit 1
fi
if [[ ! -f "$TEST_ROOT/backend-stopped" || -d "$TEST_ROOT/bin/tmp/start-4381.lock" ]]; then
  printf 'expected cleanup to stop backend and release lock\n' >&2
  exit 1
fi
if [[ "$(cat "$TEST_ROOT/air-install")" != *"|darwin|arm64|local"* ]]; then
  printf 'expected pinned Air install to use host Go and local toolchain\n' >&2
  exit 1
fi
