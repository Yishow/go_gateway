# 執行期日誌降噪合約測試（對應 scripts/start-log-utils.sh 行為）
# 執行：pwsh -NoProfile -File tests/powershell/start-log-noise.ps1
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
. (Join-Path $repoRoot "scripts/start-log-utils.ps1")

$output = @()
function Write-ColorOutput { param([string]$Message, [string]$Color = "White") $script:output += @{ Text = $Message; Color = $Color } }
function Write-Section { param([string]$Title) $script:output += @{ Text = "SECTION: $Title"; Color = "Cyan" } }
function Write-Info { param([string]$Message) $script:output += @{ Text = "INFO: $Message"; Color = "Cyan" } }

$Verbose = $false
$LogNoiseCounters = @{}
$exitCode = 0

try {
    # 情境 1：dashboard 輪詢請求被計數且不輸出
    $line = "[2026-08-22 21:00:00] ::1 GET /api/v1/datalink/devices 200 1ms"
    Write-RuntimeLogLine -Line $line
    if ($LogNoiseCounters["dashboard-refresh"] -ne 1) { throw "dashboard-refresh 應計數 1" }
    if ($output.Count -ne 0) { throw "降噪行不應輸出，got: $($output | ConvertTo-Json -Compress)" }

    # 情境 2：5xx 回應輸出為紅色 [HTTP] 格式
    $output = @()
    Write-RuntimeLogLine -Line "[2026-08-22 21:00:01] ::1 GET /api/v1/broken 500 3ms"
    if ($output.Count -ne 1 -or $output[0].Text -notmatch "^\[HTTP\].*GET\s+/api/v1/broken\s+500") {
        throw "5xx 應輸出 [HTTP] 格式，got: $($output | ConvertTo-Json -Compress)"
    }
    if ($output[0].Color -ne "Red") { throw "5xx 應為紅色" }

    # 情境 3：非 Go-log 格式的錯誤行輸出紅色（與 start.sh 對等契約：
    # Go log 格式行走 BOOT 分支，非 Verbose 時靜默、僅保留於 backend log 檔）
    $output = @()
    Write-RuntimeLogLine -Line "ERROR something failed"
    if ($output.Count -ne 1 -or $output[0].Color -ne "Red") { throw "ERROR 行應輸出紅色" }

    $output = @()
    Write-RuntimeLogLine -Line "2026/08/22 21:00:02 main.go:100: ERROR swallowed in BOOT branch"
    if ($output.Count -ne 0) { throw "Go-log 格式行非 Verbose 應靜默" }
    $Verbose = $true
    $output = @()
    Write-RuntimeLogLine -Line "2026/08/22 21:00:02 main.go:100: ERROR shown when verbose"
    if ($output.Count -ne 1) { throw "Go-log 格式行 Verbose 應輸出" }
    $Verbose = $false

    # 情境 4：air watcher 啟動橫幅被歸類為 air-watcher
    $output = @()
    Write-RuntimeLogLine -Line "v1.52.3"
    if ($LogNoiseCounters["air-watcher"] -ne 1) { throw "air-watcher 應計數 1" }

    # 情境 5：Verbose 模式不降噪
    $Verbose = $true
    $output = @()
    Write-RuntimeLogLine -Line $line
    if ($output.Count -ne 1) { throw "Verbose 模式應輸出原始行" }
    $Verbose = $false

    # 情境 6：Show-LogNoiseSummary 列出計數
    $output = @()
    Show-LogNoiseSummary
    $summary = ($output | ForEach-Object { $_.Text }) -join "`n"
    if ($summary -notmatch "dashboard-refresh: 1 行") { throw "摘要應列出 dashboard-refresh 計數，got: $summary" }

    Write-Host "PASS: start-log-noise.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
exit $exitCode
