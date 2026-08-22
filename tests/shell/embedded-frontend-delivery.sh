#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TEST_TMP_DIR"' EXIT
case_failed=0; failed_cases=0; run_code=0; run_output=""
export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=0
native_path() { command -v cygpath >/dev/null 2>&1 && cygpath -w "$1" || printf '%s\n' "$1"; }
record_failure() { case_failed=1; printf '    FAIL: %s\n' "$1" >&2; }
reset_marker() { [[ -z "${FAKE_COPY_MARKER:-}" ]] || rm -f -- "$FAKE_COPY_MARKER"; }
append_marker() { [[ -z "${FAKE_COPY_MARKER:-}" || ! -f "$FAKE_COPY_MARKER" ]] || run_output+=$'\n'"$(<"$FAKE_COPY_MARKER")"; }
write_fake_tools() {
  local tools_dir="$1"
  mkdir -p "$tools_dir"
  cat >"$tools_dir/npm.cmd" <<'EOF'
@echo off
if /I "%~1"=="install" goto :install
if /I not "%~1"=="run" exit /b 90
if /I not "%~2"=="build" exit /b 91
if "%FAKE_FRONTEND_BUILD_FAIL%"=="1" exit /b 23
if "%FAKE_FRONTEND_NO_DIST%"=="1" exit /b 0
if not exist dist mkdir dist
if not exist dist\assets mkdir dist\assets
if not exist dist\routes mkdir dist\routes
>dist\index.html echo ^<script type="module" src="./assets/app.js"^>^</script^>
>dist\assets\app.js echo console.log("fixture");
>dist\assets\style.css echo body { color: black; }
>dist\routes\studio-v2.js echo export default "studio-v2";
>dist\routes\runtime.json echo {"route":"runtime"}
exit /b 0
:install
if not exist node_modules mkdir node_modules
exit /b 0
EOF
  cp "$tools_dir/npm.cmd" "$tools_dir/pnpm.cmd"
  cat >"$tools_dir/pnpm" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  install)
    mkdir -p node_modules
    ;;
  run)
    [[ "${2:-}" == "build" ]] || exit 91
    [[ "${FAKE_FRONTEND_BUILD_FAIL:-0}" == 1 ]] && exit 23
    [[ "${FAKE_FRONTEND_NO_DIST:-0}" == 1 ]] && exit 0
    mkdir -p dist/assets dist/routes
    printf '%s\n' '<script type="module" src="./assets/app.js"></script>' >dist/index.html
    printf '%s\n' 'console.log("fixture");' >dist/assets/app.js
    printf '%s\n' 'body { color: black; }' >dist/assets/style.css
    printf '%s\n' 'export default "studio-v2";' >dist/routes/studio-v2.js
    printf '%s\n' '{"route":"runtime"}' >dist/routes/runtime.json
    ;;
  *)
    exit 90
    ;;
esac
EOF
  chmod +x "$tools_dir/pnpm"
  cat >"$tools_dir/go.cmd" <<'EOF'
@echo off
if not "%FAKE_GO_MARKER%"=="" echo invoked>"%FAKE_GO_MARKER%"
exit /b 0
EOF
}
make_fixture() {
  local fixture
  fixture="$(mktemp -d "$TEST_TMP_DIR/fixture.XXXXXX")"
  mkdir -p "$fixture/frontend" "$fixture/cmd/test_ui/static"
  cp "$REPO_ROOT/cmd/test_ui/static/embed-placeholder.txt" "$fixture/cmd/test_ui/static/embed-placeholder.txt"
  cp "$fixture/cmd/test_ui/static/embed-placeholder.txt" "$fixture/placeholder.expected"
  printf '%s\n' stale >"$fixture/cmd/test_ui/static/stale-sentinel.js"
  printf '%s\n' '{"name":"fixture","private":true}' >"$fixture/frontend/package.json"
  write_fake_tools "$fixture/tools"
  printf '%s\n' "$fixture"
}
prepare_start_fixture() {
  local fixture="$1"
  mkdir -p "$fixture/scripts"; cp "$REPO_ROOT/start.sh" "$fixture/start.sh"; cp "$REPO_ROOT/scripts/load-env.sh" "$fixture/scripts/load-env.sh"; cp "$REPO_ROOT/scripts/start-log-utils.sh" "$fixture/scripts/start-log-utils.sh"; cp "$REPO_ROOT/scripts/start-process-utils.sh" "$fixture/scripts/start-process-utils.sh"
}
run_ps_build() {
  local fixture="$1"
  if ! command -v pwsh >/dev/null 2>&1; then
    run_code=125
    run_output="pwsh not found"
    return
  fi
  export EMBED_TEST_ROOT="$(native_path "$fixture")" EMBED_TEST_SCRIPT="$(native_path "$REPO_ROOT/scripts/build.ps1")" EMBED_TEST_TOOLS="$(native_path "$fixture/tools")"
  export FAKE_GO_MARKER="$(native_path "$fixture/go-build-ran.marker")"
  reset_marker
  if [[ -n "${FAKE_COPY_MARKER:-}" ]]; then export EMBED_TEST_COPY_MARKER="$(native_path "$FAKE_COPY_MARKER")"; fi
  set +e
  run_output="$(pwsh -NoProfile -NonInteractive -Command '
    Set-Location -LiteralPath $env:EMBED_TEST_ROOT
    $env:PATH = $env:EMBED_TEST_TOOLS + [IO.Path]::PathSeparator + $env:PATH
    if ($env:FAKE_COPY_FAIL -eq "1") { function Copy-Item { param([switch]$Recurse, [switch]$Force, [object]$Path, [object]$Destination); if ($env:EMBED_TEST_COPY_MARKER) { Set-Content -LiteralPath $env:EMBED_TEST_COPY_MARKER -Value "fixture copy failure" }; throw "fixture copy failure" } }
    & $env:EMBED_TEST_SCRIPT
    exit $LASTEXITCODE
  ' 2>&1)"
  run_code=$?
  append_marker
  set -e
}
run_ps_start() {
  local fixture="$1"
  if [[ -z "${FAKE_COPY_MARKER:-}" ]]; then export FAKE_COPY_MARKER="$fixture/powershell-error.marker"; fi
  if ! command -v pwsh >/dev/null 2>&1; then
    run_code=125
    run_output="pwsh not found"
    return
  fi
  export EMBED_TEST_ROOT="$(native_path "$fixture")" EMBED_TEST_SCRIPT="$(native_path "$REPO_ROOT/start.ps1")" EMBED_TEST_TOOLS="$(native_path "$fixture/tools")"
  reset_marker
  if [[ -n "${FAKE_COPY_MARKER:-}" ]]; then export EMBED_TEST_COPY_MARKER="$(native_path "$FAKE_COPY_MARKER")"; fi
  set +e
  run_output="$(pwsh -NoProfile -NonInteractive -Command '
    $ErrorActionPreference = "Stop"
    Set-Location -LiteralPath $env:EMBED_TEST_ROOT
    $env:PATH = $env:EMBED_TEST_TOOLS + [IO.Path]::PathSeparator + $env:PATH
    $script:FRONTEND_DIR = "frontend"
    $script:STATIC_DIR = "cmd/test_ui/static"
    $script:DIST_DIR = "frontend/dist"
    $script:NODE_MODULES_DIR = "frontend/node_modules"
    function Write-Info { param([string]$Message) }
    function Write-Success { param([string]$Message) }
    function Write-Warning { param([string]$Message) }
    function Write-Error { param([string]$Message); if ($env:EMBED_TEST_COPY_MARKER) { Add-Content -LiteralPath $env:EMBED_TEST_COPY_MARKER -Value $Message }; Write-Output "[ERROR] $Message" }
    if ($env:FAKE_COPY_FAIL -eq "1") { function Copy-Item { param([switch]$Recurse, [switch]$Force, [object]$Path, [object]$Destination); if ($env:EMBED_TEST_COPY_MARKER) { Set-Content -LiteralPath $env:EMBED_TEST_COPY_MARKER -Value "fixture copy failure" }; throw "fixture copy failure" } }
    $tokens = $null
    $parseErrors = $null
    $ast = [System.Management.Automation.Language.Parser]::ParseFile(
      $env:EMBED_TEST_SCRIPT, [ref]$tokens, [ref]$parseErrors)
    $functionAst = $ast.Find({ param($node)
      $node -is [System.Management.Automation.Language.FunctionDefinitionAst] -and
        $node.Name -eq "Build-Frontend"
    }, $true)
    if ($null -eq $functionAst) { throw "Build-Frontend function not found" }
    . ([scriptblock]::Create($functionAst.Extent.Text))
    $result = Build-Frontend -Force
    if ($result -ne $true) { exit 1 }
    exit 0
  ' 2>&1)"
  run_code=$?
  append_marker
  set -e
}
run_posix_start() {
  local fixture="$1"
  prepare_start_fixture "$fixture"
  reset_marker
  set +e
  run_output="$( { cd "$fixture"; export REPO_ROOT="$fixture" PATH="$fixture/tools:$PATH"; go() { [[ -z "${FAKE_GO_MARKER:-}" ]] || printf '%s\n' invoked >"$FAKE_GO_MARKER"; }; if [[ "${FAKE_COPY_FAIL:-0}" == 1 ]]; then cp() { printf '%s\n' 'fixture copy failure' >&2; [[ -z "${FAKE_COPY_MARKER:-}" ]] || printf '%s\n' 'fixture copy failure' >"$FAKE_COPY_MARKER"; return 73; }; fi; source ./start.sh; ${2:-build_frontend}; } 2>&1)"
  run_code=$?
  append_marker
  set -e
}
expect_success() {
  local label="$1"
  if [[ "$run_code" -ne 0 ]]; then
    record_failure "$label expected exit 0, got $run_code"
    printf '      output:\n%s\n' "${run_output:-<empty>}" >&2
  fi
}
expect_failure() {
  local label="$1"
  local needle="$2"
  local output_lower="${run_output,,}"
  if [[ "$run_code" -eq 0 ]]; then
    record_failure "$label expected nonzero exit, got 0"
  fi
  if [[ "$output_lower" != *"$needle"* ]]; then
    record_failure "$label output did not identify $needle"
    printf '      output:\n%s\n' "${run_output:-<empty>}" >&2
  fi
}
assert_no_backend_build() {
  local label="$1"
  local fixture="$2"
  if [[ -f "$fixture/go-build-ran.marker" ]]; then
    record_failure "$label invoked backend build after frontend delivery failure"
  fi
}
assert_synced() {
  local label="$1"
  local fixture="$2"
  local expected="$TEST_TMP_DIR/expected.$RANDOM" actual="$TEST_TMP_DIR/actual.$RANDOM"
  (cd "$fixture/frontend/dist" && find . -mindepth 1 -printf '%y %P\n' | sort) >"$expected"
  (cd "$fixture/cmd/test_ui/static" &&
    find . -mindepth 1 ! -path './embed-placeholder.txt' -printf '%y %P\n' | sort) >"$actual"
  if ! diff -u "$expected" "$actual" >&2; then
    record_failure "$label synchronized tree differs from frontend/dist"
  fi
  while IFS= read -r entry; do
    [[ "$entry" == f\ * ]] || continue
    local relative="${entry#f }"
    if ! cmp -s "$fixture/frontend/dist/$relative" "$fixture/cmd/test_ui/static/$relative"; then
      record_failure "$label content differs at $relative"
    fi
  done <"$expected"
  if [[ ! -f "$fixture/cmd/test_ui/static/embed-placeholder.txt" ]] ||
    ! cmp -s "$fixture/placeholder.expected" "$fixture/cmd/test_ui/static/embed-placeholder.txt"; then
    record_failure "$label did not preserve embed-placeholder.txt"
  fi
  if [[ -e "$fixture/cmd/test_ui/static/stale-sentinel.js" ]]; then
    record_failure "$label did not remove stale sentinel"
  fi
}
case_powershell_build_sync() {
  local fixture="$(make_fixture)"
  export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=0
  run_ps_build "$fixture"
  expect_success 'scripts/build.ps1 powershell-build-sync'
  if [[ "$run_code" -eq 0 ]]; then
    assert_synced 'scripts/build.ps1 powershell-build-sync' "$fixture"
  fi
  [[ -f "$fixture/go-build-ran.marker" ]] || record_failure 'scripts/build.ps1 did not reach stub backend build'
}
case_powershell_start_sync() {
  local fixture="$(make_fixture)"
  export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=0
  run_ps_start "$fixture"
  expect_success 'start.ps1 powershell-start-sync'
  if [[ "$run_code" -eq 0 ]]; then
    assert_synced 'start.ps1 powershell-start-sync' "$fixture"
  fi
}
case_posix_start_sync() {
  local fixture="$(make_fixture)"
  export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=0
  run_posix_start "$fixture"
  expect_success 'start.sh posix-start-sync'
  if [[ "$run_code" -eq 0 ]]; then
    assert_synced 'start.sh posix-start-sync' "$fixture"
  fi
}
case_stale_file_removal() {
  local fixture="$(make_fixture)"
  mkdir -p "$fixture/cmd/test_ui/static/old/assets"
  printf '%s\n' stale >"$fixture/cmd/test_ui/static/old/assets/stale.js"
  export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=0
  run_posix_start "$fixture"
  expect_success 'start.sh stale-file-removal'
  if [[ "$run_code" -eq 0 ]]; then
    assert_synced 'start.sh stale-file-removal' "$fixture"
  fi
}
case_copy_failure() {
  local fixture="$(make_fixture)"; export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=1
  export FAKE_COPY_MARKER="$fixture/copy-failure.marker"
  run_ps_build "$fixture"
  expect_failure 'scripts/build.ps1 copy-failure' 'copy'
  assert_no_backend_build 'scripts/build.ps1 copy-failure' "$fixture"
  run_ps_start "$fixture"
  expect_failure 'start.ps1 copy-failure' 'copy'
  run_posix_start "$fixture"
  expect_failure 'start.sh copy-failure' 'copy'
}
case_missing_dist() {
  local fixture="$(make_fixture)"
  export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=1 FAKE_COPY_FAIL=0
  run_ps_build "$fixture"
  expect_failure 'scripts/build.ps1 missing-dist' 'dist'
  assert_no_backend_build 'scripts/build.ps1 missing-dist' "$fixture"
  run_ps_start "$fixture"
  expect_failure 'start.ps1 missing-dist' 'dist'
  run_posix_start "$fixture"
  expect_failure 'start.sh missing-dist' 'dist'
}
case_missing_frontend_source() { local fixture="$(make_fixture)"; rm -rf "$fixture/frontend"
  export FAKE_FRONTEND_BUILD_FAIL=0 FAKE_FRONTEND_NO_DIST=0 FAKE_COPY_FAIL=0 FAKE_GO_MARKER="$fixture/go-build-ran.marker"
  run_posix_start "$fixture" build_app
  expect_failure 'start.sh missing-frontend-source' 'frontend source'
  assert_no_backend_build 'start.sh missing-frontend-source' "$fixture"
}
for selected_case in powershell-build-sync powershell-start-sync posix-start-sync \
  stale-file-removal copy-failure missing-dist missing-frontend-source; do
  if [[ "$#" -gt 0 && "$1" != "$selected_case" ]]; then
    continue
  fi
  case_failed=0
  printf 'CASE %s\n' "$selected_case"
  "case_${selected_case//-/_}"
  if [[ "$case_failed" -eq 0 ]]; then
    printf '  PASS\n'
  else
    printf '  FAIL\n'
    failed_cases=$((failed_cases + 1))
  fi
done
if [[ "$#" -gt 0 ]]; then
  case "$1" in
    powershell-build-sync|powershell-start-sync|posix-start-sync|stale-file-removal|copy-failure|missing-dist|missing-frontend-source) ;;
    *) printf 'unknown case: %s\n' "$1" >&2; exit 2 ;;
  esac
fi
if [[ "$failed_cases" -ne 0 ]]; then
  printf 'RED embedded-frontend-delivery contract failures: %s case(s)\n' "$failed_cases"
  exit 1
fi
printf 'GREEN embedded-frontend-delivery contract\n'
