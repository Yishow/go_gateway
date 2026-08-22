# 開發模式啟動順序合約測試（對應 tests/shell/start-backend-before-frontend.sh）
# 驗證順序：sync > clear-backend > start-backend > wait-backend > start-frontend > wait-exit
# 執行：pwsh -NoProfile -File tests/powershell/start-dev-mode-order.ps1
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
Set-Location $repoRoot
$env:GATEWAY_START_PS1_LIBRARY_ONLY = "1"
. (Join-Path $repoRoot "start.ps1")
Remove-Item Env:\GATEWAY_START_PS1_LIBRARY_ONLY

$order = @()
function Record-Step { param([string]$Name) $script:order += $Name }

# 覆寫 orchestration 依賴（語意同 shell 測試的 stub 模式）
function Sync-EmbeddedFrontendIfRequested { Record-Step "sync" }
function Request-OpenBrowser { }
function Clear-PortForService { Record-Step "clear-backend"; return $true }
function Start-BackendProcess { Record-Step "start-backend" }
function Wait-PortReady { Record-Step "wait-backend"; return $true }
function Start-FrontendDevServer {
    Record-Step "start-frontend"
    return [PSCustomObject]@{ Id = 43210 }
}
function Wait-BackendExit { Record-Step "wait-exit"; return 0 }
function Cleanup-All { Record-Step "cleanup" }
function Acquire-RunLock { return $true }

Start-DevMode | Out-Null

$actual = $order -join ">"
$expected = "sync>clear-backend>start-backend>wait-backend>start-frontend>wait-exit>cleanup"

if ($actual -ne $expected) {
    Write-Host "FAIL: expected Start-DevMode order $expected, got $actual" -ForegroundColor Red
    exit 1
}

Write-Host "PASS: start-dev-mode-order.ps1"
exit 0
