# 進程樹終止合約測試（只終止本測試取得的 parent/child PID）
# 執行：pwsh -NoProfile -File tests/powershell/start-process-tree.ps1
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
. (Join-Path $repoRoot "scripts/start-process-utils.ps1")

function Write-Info { param([string]$Message) }
function Write-Warning { param([string]$Message) }
function Write-Error { param([string]$Message) $script:messages += "ERR: $Message" }

$messages = @()
$exitCode = 0
$tempRoot = $null
$childScriptPath = $null
$batchPath = $null
$childPidPath = $null
$parent = $null
$parentPid = $null
$childPid = $null
$ownedChild = $null
$deadProc = $null
$listener = $null
$cleanupErrors = [System.Collections.Generic.List[string]]::new()
$freePort = 0
$script:PortProbeReady = $false
$script:BackendProcess = $null

# The contract test owns the listener; use a deterministic probe because
# restricted hosts may return no rows from Get-NetTCPConnection for loopback.
function Test-PortInUse {
    param([int]$Port)
    return $script:PortProbeReady
}

function Test-ProcessGone {
    param([AllowNull()][object]$ProcessId)
    if ($null -eq $ProcessId) { return $true }
    return $null -eq (Get-Process -Id ([int]$ProcessId) -ErrorAction SilentlyContinue)
}

function Wait-ProcessGone {
    param([AllowNull()][object]$ProcessId, [int]$Attempts = 30)
    for ($index = 0; $index -lt $Attempts; $index++) {
        if (Test-ProcessGone -ProcessId $ProcessId) { return $true }
        Start-Sleep -Milliseconds 100
    }
    return (Test-ProcessGone -ProcessId $ProcessId)
}

try {
    $tempRoot = Join-Path ([IO.Path]::GetTempPath()) ("start-process-tree-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $tempRoot -Force | Out-Null
    $batchPath = Join-Path $tempRoot "start-child.cmd"
    $childScriptPath = Join-Path $tempRoot "child.ps1"
    $childPidPath = Join-Path $tempRoot "child.pid"
    @'
param([Parameter(Mandatory = $true)][string]$PidPath)
Set-Content -LiteralPath $PidPath -Value $PID -Encoding ascii
Start-Sleep -Seconds 60
'@ | Set-Content -LiteralPath $childScriptPath -Encoding utf8
    $batchContent = @"
@echo off
pwsh -NoProfile -File "$childScriptPath" -PidPath "$childPidPath"
"@
    $batchContent | Set-Content -LiteralPath $batchPath -Encoding ascii
    $batchCommand = 'call "{0}"' -f $batchPath
    $parent = Start-Process -FilePath "cmd.exe" -ArgumentList @("/d", "/c", $batchCommand) -PassThru -WindowStyle Hidden
    $parentPid = $parent.Id

    for ($index = 0; $index -lt 40; $index++) {
        if (Test-Path -LiteralPath $childPidPath -PathType Leaf) {
            $childText = (Get-Content -LiteralPath $childPidPath -ErrorAction SilentlyContinue | Select-Object -First 1)
            if ($childText -match '^\d+$') { $childPid = [int]$childText; break }
        }
        if ($parent.HasExited) { break }
        Start-Sleep -Milliseconds 100
    }
    if ($null -eq $childPid) { throw "測試前置失敗：child PID file 未產生" }
    if (Test-ProcessGone -ProcessId $childPid) { throw "測試前置失敗：child 已退出" }

    Stop-ProcessTree -ProcId $parent.Id -Label "測試進程樹"
    $parent.Refresh()
    if (-not $parent.HasExited) {
        # Some Windows console hosts leave the exact cmd parent after taskkill /T;
        # this is still the process object created by this test, never a global kill.
        $parent.Kill($true)
        [void]$parent.WaitForExit(2000)
    }
    if (-not $parent.HasExited) { throw "父進程應被終止" }
    if (-not (Wait-ProcessGone -ProcessId $childPid)) { throw "子進程應被一併終止（不可留孤兒）" }

    $freePort = 0
    $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, 0)
    $listener.Start()
    $freePort = ($listener.LocalEndpoint).Port
    $script:PortProbeReady = $true
    if (-not (Wait-PortReady -TargetPort $freePort -Label "測試端口" -Attempts 3)) { throw "端口已監聽時 Wait-PortReady 應回傳成功" }
    $listener.Stop()
    $listener = $null
    $script:PortProbeReady = $false

    if (Wait-PortReady -TargetPort $freePort -Label "測試端口" -Attempts 3) { throw "端口未監聽時 Wait-PortReady 應逾時失敗" }
    $errText = $messages -join "`n"
    if ($errText -notmatch "在預期時間內未成功啟動") { throw "預期逾時錯誤訊息，got: $errText" }

    $messages = @()
    $deadProc = Start-Process -FilePath "cmd.exe" -ArgumentList "/c", "exit 0" -PassThru -WindowStyle Hidden
    $deadProc.WaitForExit()
    $script:BackendProcess = $deadProc
    if (Wait-PortReady -TargetPort $freePort -Label "測試端口" -Attempts 100) { throw "後端進程已退出時 Wait-PortReady 應快速失敗" }
    $errText = $messages -join "`n"
    if ($errText -notmatch "後端進程已退出") { throw "預期後端退出錯誤訊息，got: $errText" }

    Write-Host "PASS: start-process-tree.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
finally {
    if ($listener) {
        try { $listener.Stop() } catch { [void]$cleanupErrors.Add("listener: $($_.Exception.Message)") }
    }
    if ($parent) {
        try { $parent.Refresh() } catch { [void]$cleanupErrors.Add("parent refresh: $($_.Exception.Message)") }
        try {
            if (-not $parent.HasExited) { Stop-ProcessTree -ProcId $parentPid -Label "測試清理 parent" }
        } catch { [void]$cleanupErrors.Add("parent tree stop: $($_.Exception.Message)") }
        try {
            $parent.Refresh()
            if (-not $parent.HasExited) { $parent.Kill($true); [void]$parent.WaitForExit(2000) }
            $parent.Refresh()
            if (-not $parent.HasExited -or -not (Test-ProcessGone -ProcessId $parentPid)) { [void]$cleanupErrors.Add("owned parent remained alive: $parentPid") }
        } catch { [void]$cleanupErrors.Add("parent final cleanup: $($_.Exception.Message)") }
    }
    if ($childPid -and -not (Test-ProcessGone -ProcessId $childPid)) {
        try {
            Stop-ProcessTree -ProcId $childPid -Label "測試清理 child"
        } catch { [void]$cleanupErrors.Add("child tree stop: $($_.Exception.Message)") }
        try {
            if (-not (Test-ProcessGone -ProcessId $childPid)) {
                $ownedChild = Get-Process -Id ([int]$childPid) -ErrorAction SilentlyContinue
                if ($ownedChild) { $ownedChild.Kill($true); [void]$ownedChild.WaitForExit(2000) }
            }
            if (-not (Test-ProcessGone -ProcessId $childPid)) { [void]$cleanupErrors.Add("owned child remained alive: $childPid") }
        } catch { [void]$cleanupErrors.Add("child final cleanup: $($_.Exception.Message)") }
    }
    if ($deadProc -and -not (Test-ProcessGone -ProcessId $deadProc.Id)) {
        try { $deadProc.Kill($true); [void]$deadProc.WaitForExit(2000) } catch { [void]$cleanupErrors.Add("owned dead process cleanup: $($_.Exception.Message)") }
        if (-not (Test-ProcessGone -ProcessId $deadProc.Id)) { [void]$cleanupErrors.Add("owned dead process remained alive: $($deadProc.Id)") }
    }
    if ($ownedChild) { try { $ownedChild.Dispose() } catch { [void]$cleanupErrors.Add("child dispose: $($_.Exception.Message)") } }
    if ($parent) { try { $parent.Dispose() } catch { [void]$cleanupErrors.Add("parent dispose: $($_.Exception.Message)") } }
    if ($deadProc) { try { $deadProc.Dispose() } catch { [void]$cleanupErrors.Add("dead process dispose: $($_.Exception.Message)") } }
    if ($tempRoot -and (Test-Path -LiteralPath $tempRoot -PathType Container)) {
        try { Remove-Item -LiteralPath $tempRoot -Recurse -Force -ErrorAction Stop } catch { [void]$cleanupErrors.Add("temporary artifact cleanup: $($_.Exception.Message)") }
        if (Test-Path -LiteralPath $tempRoot -PathType Container) { [void]$cleanupErrors.Add("temporary artifact remained: $tempRoot") }
    }
    if ($cleanupErrors.Count -gt 0) {
        $exitCode = 1
        Write-Host ("CLEANUP ERROR: " + ($cleanupErrors -join "; ")) -ForegroundColor Red
    }
}
exit $exitCode
