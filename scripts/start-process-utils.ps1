# Shared PowerShell process, port-readiness, and run-lock helpers.
# The helpers intentionally never terminate a process merely because a port is
# occupied. The launcher owns only processes it started in this invocation.

function Get-ListeningPids {
    param([Parameter(Mandatory = $true)][int]$Port)

    $tcpCommand = Get-Command Get-NetTCPConnection -ErrorAction SilentlyContinue
    if ($tcpCommand) {
        try {
            return @(
                Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction Stop |
                    Where-Object { $_.State -eq "Listen" } |
                    Select-Object -ExpandProperty OwningProcess
            )
        }
        catch {
            # Fall through to netstat on older Windows or restricted hosts.
        }
    }

    $pids = @()
    foreach ($line in @(netstat -ano 2>$null)) {
        $parts = ([string]$line).Trim() -split '\s+'
        if ($parts.Count -lt 5) { continue }
        if ($parts[3] -ne "LISTENING") { continue }
        $localEndpoint = $parts[1]
        if ($localEndpoint -notmatch (":{0}$" -f [regex]::Escape([string]$Port))) { continue }
        if ($parts[4] -match '^\d+$') { $pids += [int]$parts[4] }
    }
    return @($pids | Select-Object -Unique)
}

function Test-PortInUse {
    param([Parameter(Mandatory = $true)][int]$Port)
    return @(Get-ListeningPids -Port $Port).Count -gt 0
}

function Get-ProcessByPort {
    param([Parameter(Mandatory = $true)][int]$Port)
    $processId = @(Get-ListeningPids -Port $Port) | Select-Object -First 1
    if ($null -eq $processId) { return $null }
    return Get-Process -Id ([int]$processId) -ErrorAction SilentlyContinue
}

function Get-RunLockPath {
    return (Join-Path $script:TMP_DIR ("start-{0}.lock" -f $script:Port))
}

function Acquire-RunLock {
    New-Item -ItemType Directory -Path $script:TMP_DIR -Force | Out-Null
    $script:RunLockDir = Get-RunLockPath
    try {
        New-Item -ItemType Directory -Path $script:RunLockDir -ErrorAction Stop | Out-Null
        Set-Content -LiteralPath (Join-Path $script:RunLockDir "pid") -Value $PID -Encoding ascii
        return $true
    }
    catch {
        $pidPath = Join-Path $script:RunLockDir "pid"
        $existingPid = if (Test-Path -LiteralPath $pidPath) {
            Get-Content -LiteralPath $pidPath -ErrorAction SilentlyContinue | Select-Object -First 1
        }
        else { $null }
        if ($existingPid -match '^\d+$' -and (Get-Process -Id ([int]$existingPid) -ErrorAction SilentlyContinue)) {
            Write-Error "已有 start.ps1 管理此 repo/port（PID: $existingPid, lock: $script:RunLockDir）"
            return $false
        }
        Write-Warning "發現過期啟動鎖，將清理: $script:RunLockDir"
        Remove-Item -LiteralPath $script:RunLockDir -Recurse -Force -ErrorAction SilentlyContinue
        try {
            New-Item -ItemType Directory -Path $script:RunLockDir -ErrorAction Stop | Out-Null
            Set-Content -LiteralPath (Join-Path $script:RunLockDir "pid") -Value $PID -Encoding ascii
            return $true
        }
        catch {
            Write-Error "無法取得啟動鎖：$($_.Exception.Message)"
            return $false
        }
    }
}

function Release-RunLock {
    if ([string]::IsNullOrWhiteSpace($script:RunLockDir)) { return }
    $pidPath = Join-Path $script:RunLockDir "pid"
    $owner = if (Test-Path -LiteralPath $pidPath) {
        Get-Content -LiteralPath $pidPath -ErrorAction SilentlyContinue | Select-Object -First 1
    }
    else { $null }
    if ("$owner" -eq "$PID") {
        Remove-Item -LiteralPath $script:RunLockDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    $script:RunLockDir = $null
}

function Wait-PidExit {
    param([Parameter(Mandatory = $true)][int]$ProcId, [int]$Tenths = 30)
    for ($index = 0; $index -lt $Tenths; $index++) {
        if (-not (Get-Process -Id $ProcId -ErrorAction SilentlyContinue)) { return $true }
        Start-Sleep -Milliseconds 100
    }
    return $false
}

function Stop-ProcessTree {
    param(
        [Parameter(Mandatory = $true)][int]$ProcId,
        [string]$Label = "進程"
    )
    if ($ProcId -le 0 -or $ProcId -eq $PID) { return }
    if (-not (Get-Process -Id $ProcId -ErrorAction SilentlyContinue)) { return }
    Write-Info "正在停止 $Label（PID: $ProcId）..."
    $taskkill = Get-Command taskkill.exe -ErrorAction SilentlyContinue
    if ($taskkill) {
        try { & $taskkill.Source /PID $ProcId /T 2>$null | Out-Null }
        catch { Write-Warning "$Label 優雅停止失敗，將嘗試強制終止" }
        if (-not (Wait-PidExit -ProcId $ProcId -Tenths 30)) {
            Write-Warning "$Label 未在預期時間內停止，改用強制終止"
            try { & $taskkill.Source /PID $ProcId /T /F 2>$null | Out-Null }
            catch { Write-Warning "$Label 強制終止失敗：$($_.Exception.Message)" }
            [void](Wait-PidExit -ProcId $ProcId -Tenths 20)
        }
    }
    else {
        try { Stop-Process -Id $ProcId -Force -ErrorAction SilentlyContinue }
        catch { Write-Warning "$Label 停止失敗：$($_.Exception.Message)" }
        [void](Wait-PidExit -ProcId $ProcId -Tenths 20)
    }
}

function Test-ProcessTreeOwnership {
    param(
        [Parameter(Mandatory = $true)][int]$CandidateProcessId,
        [Parameter(Mandatory = $true)][int]$OwnerProcessId
    )
    if ($CandidateProcessId -le 0 -or $OwnerProcessId -le 0) { return $false }
    $visited = [System.Collections.Generic.HashSet[int]]::new()
    $currentProcessId = $CandidateProcessId
    for ($depth = 0; $depth -lt 64; $depth++) {
        if ($currentProcessId -eq $OwnerProcessId) { return $true }
        if (-not $visited.Add($currentProcessId)) { return $false }
        $processInfo = @(Get-CimInstance Win32_Process -Filter "ProcessId = $currentProcessId" -ErrorAction SilentlyContinue) |
            Select-Object -First 1
        if ($null -eq $processInfo) { return $false }
        try { $parentProcessId = [int]$processInfo.ParentProcessId }
        catch { return $false }
        if ($parentProcessId -le 0 -or $parentProcessId -eq $currentProcessId) { return $false }
        $currentProcessId = $parentProcessId
    }
    return $false
}

function Test-PortOwnedByProcess {
    param(
        [Parameter(Mandatory = $true)][int]$TargetPort,
        [Parameter(Mandatory = $true)][object]$Process
    )
    try { $ownerProcessId = [int]$Process.Id }
    catch { return $false }
    $listener = @(Get-ProcessByPort -Port $TargetPort) | Select-Object -First 1
    if ($null -eq $listener) { return $false }
    try { $listenerProcessId = [int]$listener.Id }
    catch { return $false }
    return Test-ProcessTreeOwnership -CandidateProcessId $listenerProcessId -OwnerProcessId $ownerProcessId
}

function Wait-PortReady {
    param(
        [Parameter(Mandatory = $true)][int]$TargetPort,
        [string]$Label = "端口",
        [int]$Attempts = 1200,
        [int]$DelayMilliseconds = 100,
        [AllowNull()][object]$Process
    )
    if ($null -eq $Process) { $Process = $script:BackendProcess }
    for ($index = 0; $index -lt $Attempts; $index++) {
        if (Get-Command Write-ProcessLogs -ErrorAction SilentlyContinue) { Write-ProcessLogs }
        if ($Process) {
            try {
                if ($Process.HasExited) {
                    Write-Error "$Label $TargetPort 未成功啟動（進程已退出）"
                    return $false
                }
            }
            catch {
                Write-Error "$Label $TargetPort 未成功啟動（進程狀態無法取得）"
                return $false
            }
        }
        if (Test-PortInUse -Port $TargetPort) {
            if ($null -eq $Process -or (Test-PortOwnedByProcess -TargetPort $TargetPort -Process $Process)) {
                return $true
            }
        }
        Start-Sleep -Milliseconds $DelayMilliseconds
    }
    Write-Error "$Label $TargetPort 在預期時間內未成功啟動"
    return $false
}

function Get-RelatedProcessPids {
    $root = [regex]::Escape([string]$script:ROOT_DIR)
    $patterns = @(
        "$root[\\/]start\.ps1",
        "$root[\\/]$([regex]::Escape($script:APP_PATH))",
        "$root[\\/]bin[\\/]tools[\\/]air\.exe",
        "$root[\\/]frontend[\\/]node_modules"
    )
    $found = @()
    foreach ($process in @(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue)) {
        if ($process.ProcessId -eq $PID -or [string]::IsNullOrWhiteSpace($process.CommandLine)) { continue }
        foreach ($pattern in $patterns) {
            if ($process.CommandLine -match $pattern) {
                $found += [int]$process.ProcessId
                break
            }
        }
    }
    return @($found | Select-Object -Unique)
}
