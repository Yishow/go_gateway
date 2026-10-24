#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
mkdir -p "$TMP_DIR/scripts" "$TMP_DIR/frontend" "$TMP_DIR/bin"
cp "$REPO_ROOT/start.sh" "$TMP_DIR/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TMP_DIR/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TMP_DIR/scripts/start-process-utils.sh"
printf '{"name":"frontend"}\n' >"$TMP_DIR/frontend/package.json"
cat >"$TMP_DIR/bin/npm" <<'EOF'
#!/usr/bin/env bash
printf 'simulated npm ci failure\n' >&2
exit 41
EOF
chmod +x "$TMP_DIR/bin/npm"

cd "$TMP_DIR"
export PATH="$TMP_DIR/bin:$PATH"
source ./start.sh >/dev/null 2>&1 || true
set +e
output="$(ensure_frontend_dependencies 2>&1)"
status=$?
set -e

if [[ "$status" -eq 0 || "$output" != *"npm ci 失敗"* ]]; then
  printf 'expected clear npm ci failure, status=%s output=%s\n' "$status" "$output" >&2
  exit 1
fi
