#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
OUTSIDE_DIR="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT" "$OUTSIDE_DIR"' EXIT
mkdir -p "$TEST_ROOT/scripts" "$TEST_ROOT/bin" "$TEST_ROOT/frontend"
cp "$REPO_ROOT/start.sh" "$TEST_ROOT/start.sh"
cp "$REPO_ROOT/scripts/load-env.sh" "$TEST_ROOT/scripts/load-env.sh"
cp "$REPO_ROOT/scripts/start-process-utils.sh" "$TEST_ROOT/scripts/start-process-utils.sh"
touch "$TEST_ROOT/scripts/build-single.mjs"
cat >"$TEST_ROOT/bin/node" <<'EOF'
#!/usr/bin/env bash
printf '%s|%s\n' "$PWD" "$*" > "$NODE_WITNESS"
EOF
chmod +x "$TEST_ROOT/bin/node"

export PATH="$TEST_ROOT/bin:$PATH"
export NODE_WITNESS="$TEST_ROOT/node-witness"
(cd "$OUTSIDE_DIR" && "$TEST_ROOT/start.sh" --build-single --os mac --arch arm64)

CANONICAL_ROOT="$(cd "$TEST_ROOT" && pwd -P)"
expected_prefix="$CANONICAL_ROOT|$CANONICAL_ROOT/scripts/build-single.mjs --os darwin --arch arm64"
actual="$(cat "$NODE_WITNESS")"
if [[ "$actual" != "$expected_prefix" ]]; then
  printf 'expected build helper to run from script root with normalized args, got %s\n' "$actual" >&2
  exit 1
fi

(cd "$OUTSIDE_DIR" && "$TEST_ROOT/start.sh" --build-single)
actual="$(cat "$NODE_WITNESS")"
if [[ "$actual" != "$CANONICAL_ROOT|$CANONICAL_ROOT/scripts/build-single.mjs" ]]; then
  printf 'expected no-target build to delegate host defaults without an empty argument, got %s\n' "$actual" >&2
  exit 1
fi

menu_output="$(printf '0\n' | (cd "$OUTSIDE_DIR" && "$TEST_ROOT/start.sh"))"
if [[ "$menu_output" != *'1) 開發模式'* || "$menu_output" != *'2) 建置單一執行檔'* || "$menu_output" != *'0) 離開'* ]]; then
  printf 'expected two-item launcher menu, got:\n%s\n' "$menu_output" >&2
  exit 1
fi
