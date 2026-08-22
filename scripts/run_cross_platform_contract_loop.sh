#!/usr/bin/env bash
# cross-platform-test-contracts 能力的可重複驗證序列。
# 依 spec 順序執行：Vite focused x2 -> handler focused -> handler package
# -> dbtarget row-group -count=20（正常 GOMAXPROCS）-> frontend full suite
# -> go test ./...。每個 loop 通過後才會開始下一個，失敗時明確標示 scope，
# 便於將非本契約的 baseline failure 分離（openspec/specs/cross-platform-test-contracts）。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

run_step() {
  local name="$1"
  shift
  echo ""
  echo "==================================================================="
  echo "[cross-platform-loop] STEP: ${name}"
  echo "==================================================================="
  if ! "$@"; then
    echo "[cross-platform-loop] FAIL at step: ${name}" >&2
    exit 1
  fi
  echo "[cross-platform-loop] PASS: ${name}"
}

cd "$REPO_ROOT"

run_step "vite focused test (1/2)" bash -c 'cd frontend && npx vitest run tests/unit/utils/viteConfig.test.ts'
run_step "vite focused test (2/2)" bash -c 'cd frontend && npx vitest run tests/unit/utils/viteConfig.test.ts'
run_step "handler focused tests" go test ./internal/api/handlers -run 'TestStudioV2Workspace(Database|Audit)' -count=1
run_step "handler package full" go test ./internal/api/handlers -count=1
run_step "dbtarget row-group -count=20 (normal GOMAXPROCS)" go test ./internal/datalink/dbtarget -run RowGroups -count=20
run_step "frontend full suite" bash -c 'cd frontend && npm test -- --run'
run_step "go full suite" go test ./...

echo ""
echo "[cross-platform-loop] ALL STEPS PASSED"
