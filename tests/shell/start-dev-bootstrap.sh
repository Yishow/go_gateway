#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT
mkdir -p "$TEST_ROOT/scripts" "$TEST_ROOT/cmd/test_ui/static" "$TEST_ROOT/frontend" "$TEST_ROOT/bin"
cp "$REPO_ROOT/start.sh" "$TEST_ROOT/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_ROOT/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_ROOT/scripts/start-process-utils.sh"
touch "$TEST_ROOT/frontend/package.json" "$TEST_ROOT/scripts/build-single.mjs"
printf 'placeholder\n' >"$TEST_ROOT/cmd/test_ui/static/embed-placeholder.txt"
cat >"$TEST_ROOT/bin/node" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
[[ "$*" == *"--prepare-dev"* ]] || exit 41
touch "$FAKE_ROOT/prepare-dev-called"
printf '<html>dev</html>\n' > "$FAKE_ROOT/cmd/test_ui/static/index.html"
EOF
chmod +x "$TEST_ROOT/bin/node"
export FAKE_ROOT="$TEST_ROOT" PATH="$TEST_ROOT/bin:$PATH"

cd "$TEST_ROOT"
source ./start.sh >/dev/null 2>&1 || true
ensure_embedded_frontend_snapshot

[[ -f prepare-dev-called && -s cmd/test_ui/static/index.html ]] || {
  printf 'expected fresh dev startup to prepare the embedded bootstrap snapshot\n' >&2
  exit 1
}
