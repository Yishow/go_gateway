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
set -euo pipefail
[[ "${1:-}" == ci ]] || exit 41
touch .npm-ci-ran
mkdir -p node_modules/vite/bin
touch node_modules/vite/bin/vite.js
EOF
chmod +x "$TMP_DIR/bin/npm"

cd "$TMP_DIR"
export PATH="$TMP_DIR/bin:$PATH"
source ./start.sh >/dev/null 2>&1 || true
ensure_frontend_dependencies

[[ -f frontend/.npm-ci-ran ]] || { printf 'expected npm ci when Vite is missing\n' >&2; exit 1; }
rm -f frontend/.npm-ci-ran
ensure_frontend_dependencies
[[ ! -f frontend/.npm-ci-ran ]] || { printf 'expected no second npm ci when local Vite exists\n' >&2; exit 1; }
