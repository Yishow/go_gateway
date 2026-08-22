# 進程樹終止合約測試（驗證 taskkill /T 不留孤兒進程）
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

try {
    # 情境 1：Stop-ProcessTree 終止 cmd → ping 整棵樹
    $parent = Start-Process -FilePath "cmd.exe" -ArgumentList "/c", "ping -n 60 127.0.0.1 > NUL 2>&1" -PassThru -WindowStyle Hidden
    $childPid = $null
    for ($i = 0; $i -lt 30; $i++) {
        $child = Get-CimInstance Win32_Process -Filter "ParentProcessId = $($parent.Id)" -ErrorAction SilentlyContinue |
            Where-Object { $_.Name -match "ping|cmd" } | Select-Object -First 1
        if ($child) { $childPid = [int]$child.ProcessId; break }
        Start-Sleep -Milliseconds 100
    }
    if (-not $childPid) { throw "測試前置失敗：找不到 cmd 子進程" }

    Stop-ProcessTree -ProcId $parent.Id -Label "測試進程樹"

    if (Get-Process -Id $parent.Id -ErrorAction SilentlyContinue) { throw "父進程應被終止" }
    if (Get-Process -Id $childPid -ErrorAction SilentlyContinue) { throw "子進程應被一併終止（不可留孤兒）" }

    # 情境 2：Wait-PortReady 在端口未開時逾時失敗、端口開啟時成功
    $freePort = 0
    $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, 0)
    $listener.Start()
    $freePort = ($listener.LocalEndpoint).Port

    $script:BackendProcess = $null
    if (-not (Wait-PortReady -TargetPort $freePort -Label "測試端口" -Attempts 3)) {
        throw "端口已監聽時 Wait-PortReady 應回傳成功"
    }
    $listener.Stop()

    if (Wait-PortReady -TargetPort $freePort -Label "測試端口" -Attempts 3) {
        throw "端口未監聽時 Wait-PortReady 應逾時失敗"
    }
    $errText = $messages -join "`n"
    if ($errText -notmatch "在預期時間內未成功啟動") { throw "預期逾時錯誤訊息，got: $errText" }

    # 情境 3：後端進程已退出時 Wait-PortReady 快速失敗
    $messages = @()
    $deadProc = Start-Process -FilePath "cmd.exe" -ArgumentList "/c", "exit 0" -PassThru -WindowStyle Hidden
    $deadProc.WaitForExit()
    $script:BackendProcess = $deadProc
    if (Wait-PortReady -TargetPort $freePort -Label "測試端口" -Attempts 100) {
        throw "後端進程已退出時 Wait-PortReady 應快速失敗"
    }
    $errText = $messages -join "`n"
    if ($errText -notmatch "後端進程已退出") { throw "預期後端退出錯誤訊息，got: $errText" }

    Write-Host "PASS: start-process-tree.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
finally {
    if ($listener) { $listener.Stop() }
}
exit $exitCode
