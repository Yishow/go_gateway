# Readiness polling must stream redirected logs without rereading the full file.
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$env:GATEWAY_START_PS1_LIBRARY_ONLY = "1"
. (Join-Path $repoRoot "start.ps1")
Remove-Item Env:\GATEWAY_START_PS1_LIBRARY_ONLY -ErrorAction SilentlyContinue

$logPath = Join-Path ([System.IO.Path]::GetTempPath()) ("start-log-stream-" + [System.IO.Path]::GetRandomFileName())
$script:BackendLogFile = $logPath
$script:BackendErrorLogFile = $null
$script:FrontendLogFile = $null
$script:FrontendErrorLogFile = $null
$script:LogPositions = @{}
$script:output = @()
$script:BackendProcess = $null
New-Item -ItemType File -Path $logPath -Force | Out-Null

function Write-Host { param([object]$Object) $script:output += [string]$Object }
function Test-PortInUse { param([int]$Port) return $true }

$exitCode = 0
try {
    [System.IO.File]::WriteAllText($logPath, "first`npartial", [System.Text.Encoding]::UTF8)
    Write-ProcessLogs
    if ($script:output.Count -ne 1 -or $script:output[0] -ne "first") { throw "應只輸出完整的第一行" }

    [System.IO.File]::AppendAllText($logPath, " tail`n", [System.Text.Encoding]::UTF8)
    Write-ProcessLogs
    if ($script:output.Count -ne 2 -or $script:output[1] -ne "partial tail") { throw "不完整行應在補上換行後只輸出一次" }

    [System.IO.File]::AppendAllText($logPath, "next", [System.Text.Encoding]::UTF8)
    Write-ProcessLogs
    if ($script:output.Count -ne 2) { throw "未完成行不應提前輸出" }
    [System.IO.File]::AppendAllText($logPath, "`n", [System.Text.Encoding]::UTF8)
    Write-ProcessLogs
    if ($script:output.Count -ne 3 -or $script:output[2] -ne "next") { throw "補上換行後應輸出未完成行一次" }

    # Wait-PortReady must invoke the same stream hook before checking readiness.
    Remove-Item Function:\Write-ProcessLogs -Force
    function Write-ProcessLogs { $script:logHookCalls++ }
    $script:logHookCalls = 0
    if (-not (Wait-PortReady -TargetPort 4545 -Attempts 3 -DelayMilliseconds 0)) { throw "listening port should be ready" }
    if ($script:logHookCalls -ne 1) { throw "readiness should stream logs before checking the port" }
    Remove-Item Function:\Write-Host -Force
    Write-Host "PASS: start-process-log-stream.ps1"
}
catch {
    Remove-Item Function:\Write-Host -Force -ErrorAction SilentlyContinue
    Write-Host "FAIL: $($_.Exception.Message)"
    $exitCode = 1
}
finally { Remove-Item -LiteralPath $logPath -Force -ErrorAction SilentlyContinue }
exit $exitCode
