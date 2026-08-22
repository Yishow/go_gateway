# ============================================
# start.ps1 進程/端口/啟動鎖工具（與 scripts/start-process-utils.sh 對齊）
# ============================================
# 說明：
#   - 由 start.ps1 以 dot-source 載入；函數經 dynamic scoping 讀取
#     主腳本的 $Port、$script:TMP_DIR、$script:ROOT_DIR 與輸出函數
#   - Acquire-RunLock / Release-RunLock：以 bin/tmp/start-<port>.lock
#     目錄鎖確保同一 repo/port 同時間只有一個 start.ps1 管理
#   - Stop-ProcessTree：Windows 以 taskkill /T 終止整棵進程樹，
#     避免 cmd → node / air → go → exe 殘留孤兒進程
# ============================================

function Test-PortInUse {
    param([int]$Port)

    try {
        $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        return $null -ne $connection
    }
    catch {
        # 如果 Get-NetTCPConnection 不可用，使用 netstat
        $netstatOutput = netstat -ano | Select-String ":$Port\s"
        return $null -ne $netstatOutput
    }
}

function Get-ProcessByPort {
    param([int]$Port)

    try {
        $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        if ($connection) {
            $processId = $connection.OwningProcess | Select-Object -First 1
            if ($processId) {
                return Get-Process -Id $processId -ErrorAction SilentlyContinue
            }
        }
    }
    catch {
        # 回退到使用 netstat
        $netstatLine = netstat -ano | Select-String ":$Port\s" | Select-Object -First 1
        if ($netstatLine) {
            $parts = $netstatLine -split '\s+'
            $processId = $parts[-1]
            if ($processId -match '^\d+$') {
                return Get-Process -Id ([int]$processId) -ErrorAction SilentlyContinue
            }
        }
    }
    return $null
}

function Get-RunLockPath {
    return (Join-Path $script:TMP_DIR "start-$Port.lock")
}

function Acquire-RunLock {
    New-Item -ItemType Directory -Path $script:TMP_DIR -Force | Out-Null
    $script:RunLockDir = Get-RunLockPath

    $created = $false
    try {
        # 目錄已存在時 New-Item 拋錯，等同 mkdir 的原子語意
        New-Item -ItemType Directory -Path $script:RunLockDir -ErrorAction Stop | Out-Null
        $created = $true
    }
    catch {
        $created = $false
    }

    if (-not $created) {
        $pidFile = Join-Path $script:RunLockDir "pid"
        $existingPid = $null
        if (Test-Path -LiteralPath $pidFile) {
            $existingPid = (Get-Content -LiteralPath $pidFile -ErrorAction SilentlyContinue | Select-Object -First 1)
        }
        if ($existingPid -match '^\d+$' -and (Get-Process -Id ([int]$existingPid) -ErrorAction SilentlyContinue)) {
            Write-Error "已有 start.ps1 管理此 repo/port（PID: $existingPid, lock: $script:RunLockDir）"
            Write-Info "若確認是殘留狀態，請先執行: .\start.ps1 -StopAll"
            return $false
        }
        Write-Warning "發現過期啟動鎖，將清理: $script:RunLockDir"
        Remove-Item -LiteralPath $script:RunLockDir -Recurse -Force -ErrorAction SilentlyContinue
        New-Item -ItemType Directory -Path $script:RunLockDir -ErrorAction Stop | Out-Null
    }

    Set-Content -LiteralPath (Join-Path $script:RunLockDir "pid") -Value $PID
    return $true
}

function Release-RunLock {
    if ([string]::IsNullOrEmpty($script:RunLockDir) -or -not (Test-Path -LiteralPath $script:RunLockDir)) {
        return
    }

    $pidFile = Join-Path $script:RunLockDir "pid"
    $owner = $null
    if (Test-Path -LiteralPath $pidFile) {
        $owner = (Get-Content -LiteralPath $pidFile -ErrorAction SilentlyContinue | Select-Object -First 1)
    }
    if ([string]::IsNullOrWhiteSpace($owner) -or "$owner" -eq "$PID") {
        Remove-Item -LiteralPath $script:RunLockDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    $script:RunLockDir = $null
}

function Wait-PidExit {
    param(
        [int]$ProcId,
        [int]$Tenths = 30
    )

    for ($i = 0; $i -lt $Tenths; $i++) {
        if (-not (Get-Process -Id $ProcId -ErrorAction SilentlyContinue)) {
            return $true
        }
        Start-Sleep -Milliseconds 100
    }
    return $false
}

function Stop-ProcessTree {
    param(
        [int]$ProcId,
        [string]$Label = "進程"
    )

    if ($ProcId -le 0) { return }
    if ($ProcId -eq $PID) { return }
    if (-not (Get-Process -Id $ProcId -ErrorAction SilentlyContinue)) { return }

    Write-Info "正在停止 $Label（PID: $ProcId）..."
    # 與 start.sh 的 TERM → 等待 → KILL 升級對齊：
    # 先嘗試優雅停止（taskkill /T 對視窗進程送 WM_CLOSE），讓後端有機會執行
    # graceful shutdown（cmd/test_ui 的 server.Shutdown）；Console 進程通常
    # 不回應 WM_CLOSE，逾時後仍以 /F 強制終止整棵樹
    $null = & taskkill.exe /PID $ProcId /T 2>$null
    if (-not (Wait-PidExit -ProcId $ProcId -Tenths 30)) {
        Write-Warning "$Label 未在預期時間內停止，改用強制終止"
        $null = & taskkill.exe /PID $ProcId /T /F 2>$null
        $null = (Wait-PidExit -ProcId $ProcId -Tenths 20)
    }
}

function Wait-PortReady {
    param(
        [int]$TargetPort,
        [string]$Label = "端口",
        [int]$Attempts = 100
    )

    for ($i = 0; $i -lt $Attempts; $i++) {
        if (Test-PortInUse -Port $TargetPort) { return $true }
        if ($script:BackendProcess -and $script:BackendProcess.HasExited) {
            Write-Error "$Label $TargetPort 未成功啟動（後端進程已退出）"
            return $false
        }
        Start-Sleep -Milliseconds 100
    }

    Write-Error "$Label $TargetPort 在預期時間內未成功啟動"
    return $false
}

function Get-RelatedProcessPids {
    # 以命令列特徵找出本 repo 相關進程（後端 exe、air、vite 等），
    # 供 -StopAll 清理孤兒進程使用；排除腳本自身。
    # 所有 pattern 一律錨定到本 repo 的絕對路徑，
    # 避免誤殺機器上其他專案的同名進程（vite、其他 checkout 的 test_ui）
    $root = "$script:ROOT_DIR"
    $frontendNodeModules = [regex]::Escape((Join-Path $root $script:FRONTEND_DIR)) + '[\\/]node_modules'
    $patterns = @(
        [regex]::Escape((Join-Path $root 'start.ps1')),
        [regex]::Escape((Join-Path $root $script:APP_PATH)),
        [regex]::Escape((Join-Path $root $script:TMP_DIR 'gateway-air.exe')),
        $frontendNodeModules
    )

    $found = @()
    foreach ($proc in (Get-CimInstance Win32_Process -ErrorAction SilentlyContinue)) {
        if ($proc.ProcessId -eq $PID) { continue }
        $cmdLine = $proc.CommandLine
        if ([string]::IsNullOrWhiteSpace($cmdLine)) { continue }
        foreach ($pattern in $patterns) {
            if ($cmdLine -match $pattern) {
                $found += [int]$proc.ProcessId
                break
            }
        }
    }
    return $found
}
