# 啟動鎖管理合約測試（對應 tests/shell/start-lock-management.sh）
# 執行：pwsh -NoProfile -File tests/powershell/start-lock-management.ps1
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
. (Join-Path $repoRoot "scripts/start-process-utils.ps1")

$testTmp = Join-Path ([System.IO.Path]::GetTempPath()) ("start-lock-test-" + [System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $testTmp -Force | Out-Null

$messages = @()
function Write-Error { param([string]$Message) $script:messages += "ERR: $Message" }
function Write-Info { param([string]$Message) $script:messages += "INFO: $Message" }
function Write-Warning { param([string]$Message) $script:messages += "WARN: $Message" }

$exitCode = 0
try {
    $Port = 3333
    $TMP_DIR = Join-Path $testTmp "bin/tmp"

    # 情境 1：鎖持有者存活（本測試進程）→ 第二次取得必須失敗並說明原因
    $script:RunLockDir = $null
    if (-not (Acquire-RunLock)) { throw "預期第一次 Acquire-RunLock 成功" }
    $lockDir = $script:RunLockDir

    $script:RunLockDir = $null
    if (Acquire-RunLock) { throw "預期鎖被存活進程持有時 Acquire-RunLock 失敗，卻成功了" }
    $activeMsg = $messages -join "`n"
    if ($activeMsg -notmatch "已有 start\.ps1 管理此 repo/port") {
        throw "預期失敗訊息說明鎖持有者，got: $activeMsg"
    }
    if (-not (Test-Path (Join-Path $lockDir "pid"))) { throw "鎖目錄應保留原持有者的 pid 檔" }

    # 情境 2：過期鎖（持有者已死）→ 清理後可重新取得
    $messages = @()
    Set-Content -Path (Join-Path $lockDir "pid") -Value "99999999"
    $script:RunLockDir = $null
    if (-not (Acquire-RunLock)) { throw "預期過期鎖會被清理並重新取得" }
    $staleMsg = $messages -join "`n"
    if ($staleMsg -notmatch "發現過期啟動鎖") { throw "預期顯示過期鎖清理訊息，got: $staleMsg" }

    # 情境 3：Release-RunLock 移除自己的鎖
    Release-RunLock
    if (Test-Path $lockDir) { throw "預期 Release-RunLock 移除鎖目錄" }

    # 情境 4：非持有者的鎖不可被 Release 清掉
    New-Item -ItemType Directory -Path $lockDir -Force | Out-Null
    Set-Content -Path (Join-Path $lockDir "pid") -Value "99999998"
    $script:RunLockDir = $lockDir
    Release-RunLock
    if (-not (Test-Path $lockDir)) { throw "非持有者的鎖不應被 Release-RunLock 移除" }

    Write-Host "PASS: start-lock-management.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
finally {
    Remove-Item -Recurse -Force $testTmp -ErrorAction SilentlyContinue
}
exit $exitCode
