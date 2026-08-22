# ============================================
# start.ps1 端口管理工具（對應 start.sh 內的 port 管理區塊）
# ============================================
# 說明：
#   - 由 start.ps1 以 dot-source 載入；函數經 dynamic scoping 讀取
#     主腳本的 $Port、$script:FrontendDevPort、輸出函數，
#     以及 start-process-utils.ps1 的 Test-PortInUse / Get-ProcessByPort
#   - Get-ManagedPorts：列出本腳本管理的端口集合
#   - Clear-PortForService：啟動前清理端口（-AutoKill 免詢問）
#   - Test-ManagedProcess：判斷佔用端口的進程是否屬於本專案
# ============================================

function Get-ManagedPorts {
    $ports = [System.Collections.Generic.List[int]]::new()
    # 注意：@($a, $b, 8080..8090) 會被解析為 ($a, $b, 8080)..8090（逗號優先於 range），
    # 必須以 + 串接才能正確展開 range；Modbus share 端口與 start.sh managed_ports 對齊
    foreach ($candidate in @($Port, $script:FrontendDevPort, $script:ModbusSharePort) + (8080..8090)) {
        if ($null -eq $candidate) {
            continue
        }

        $text = "$candidate"
        if ($text -notmatch '^\d+$') {
            continue
        }

        $value = [int]$text
        if (-not $ports.Contains($value)) {
            $ports.Add($value)
        }
    }

    return $ports
}

function Get-ProcessExecutablePath {
    param([System.Diagnostics.Process]$Process)

    if (-not $Process) {
        return $null
    }

    if ($Process.Path) {
        return $Process.Path
    }

    try {
        return $Process.MainModule.FileName
    }
    catch {
        try {
            $processInfo = Get-CimInstance Win32_Process -Filter "ProcessId = $($Process.Id)" -ErrorAction Stop
            return $processInfo.ExecutablePath
        }
        catch {
            return $null
        }
    }
}

function Test-ManagedProcess {
    param([System.Diagnostics.Process]$Process)

    if (-not $Process) {
        return $false
    }

    $name = $Process.ProcessName.ToLowerInvariant()
    if ($name -in @("gateway", "test-ui")) {
        return $true
    }

    $executablePath = Get-ProcessExecutablePath -Process $Process
    if ([string]::IsNullOrWhiteSpace($executablePath)) {
        return $false
    }

    return (
        $executablePath -match '[\\/]gateway(\.exe)?$' -or
        $executablePath -match '[\\/]test_ui(\.exe)?$' -or
        $executablePath -like '*frontend*node_modules*vite*'
    )
}

# 終止佔用指定端口的進程
function Stop-ProcessByPort {
    param(
        [int]$Port,
        [switch]$Force
    )

    $process = Get-ProcessByPort -Port $Port
    if ($process) {
        Write-Warning "發現端口 $Port 被進程佔用：$($process.ProcessName) (PID: $($process.Id))"

        try {
            if ($Force) {
                Stop-Process -Id $process.Id -Force -ErrorAction Stop
                Write-Success "已強制終止進程 $($process.ProcessName) (PID: $($process.Id))"
            }
            else {
                Stop-Process -Id $process.Id -ErrorAction Stop
                Write-Success "已終止進程 $($process.ProcessName) (PID: $($process.Id))"
            }

            # 等待進程完全終止
            Start-Sleep -Milliseconds 500

            # 驗證端口是否已釋放
            if (Test-PortInUse -Port $Port) {
                Write-Warning "端口 $Port 仍被佔用，嘗試強制終止..."
                $process = Get-ProcessByPort -Port $Port
                if ($process) {
                    Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
                    Start-Sleep -Milliseconds 500
                }
            }

            return $true
        }
        catch {
            Write-Error "無法終止進程 $($process.ProcessName) (PID: $($process.Id)): $_"
            return $false
        }
    }
    else {
        Write-Info "端口 $Port 未被佔用"
        return $true
    }
}

# 清理端口並準備啟動服務
function Clear-PortForService {
    param(
        [int]$Port = 8080,
        [switch]$AutoKill
    )

    if (Test-PortInUse -Port $Port) {
        Write-Info "檢測到端口 $Port 被佔用，正在清理..."

        if ($AutoKill) {
            $result = Stop-ProcessByPort -Port $Port -Force
            if (-not $result) {
                Write-Error "無法清理端口 $Port，請手動處理"
                return $false
            }
        }
        else {
            $process = Get-ProcessByPort -Port $Port
            if ($process) {
                Write-Warning "端口 $Port 被進程佔用：$($process.ProcessName) (PID: $($process.Id))"
                Write-Info "是否要終止該進程？(Y/N)"
                $response = Read-Host

                if ($response -eq "Y" -or $response -eq "y") {
                    $result = Stop-ProcessByPort -Port $Port -Force
                    if (-not $result) {
                        Write-Error "無法清理端口 $Port"
                        return $false
                    }
                }
                else {
                    Write-Warning "跳過端口清理，服務可能無法啟動"
                    return $false
                }
            }
        }

        # 再次檢查端口是否已釋放
        if (Test-PortInUse -Port $Port) {
            Write-Error "端口 $Port 清理失敗，仍被佔用"
            return $false
        }

        Write-Success "端口 $Port 已清理完成"
    }
    else {
        Write-Info "端口 $Port 可用"
    }

    return $true
}
