# ============================================
# start.ps1 執行期日誌工具（與 scripts/start-log-utils.sh 對齊）
# ============================================
# 說明：
#   - 由 start.ps1 以 dot-source 載入，函數經 dynamic scoping
#     讀取主腳本的 $Verbose、$script:LogNoiseCounters 與輸出函數
#   - Write-RuntimeLogLine 負責過濾/格式化後端與 Air 的執行期日誌
#   - Show-LogNoiseSummary 於結束時顯示被隱藏的雜訊統計
# ============================================

# dashboard 輪詢的 datalink 實體清單（單一來源，兩處比對共用；與 start.sh 對齊）
$script:DatalinkDashboardEntities = 'devices|polling-groups|points|mappings|tags'

# Modbus share 端口由主腳本依環境變數覆寫；此處提供 StrictMode 安全的預設值
if (-not (Test-Path Variable:Script:ModbusSharePort)) {
    $script:ModbusSharePort = '5020'
}

function Get-LogNoiseCategory {
    param([string]$Line)
    if ($Line -match "CMD will not recognize non \.exe file for execution") { return "air-warning" }
    if ($Line -match "^watching\b|^building\.\.\.|^!exclude\b|^\s*/|^v\d+\.\d+\.\d+") { return "air-watcher" }
    if ($Line -match "/api/v1/datalink/modbus-share/status") { return "status-polling" }
    if ($Line -match ("\] ::1 GET /api/v1/datalink/(" + $script:DatalinkDashboardEntities + ")\b")) { return "dashboard-refresh" }
    if ($Line -match "資料庫路徑|Executing SQLite migration|ConnectionManager 已初始化|已註冊的協議") { return "startup-detail" }
    return $null
}

function Write-RuntimeLogLine {
    param([string]$Line)
    if ([string]::IsNullOrWhiteSpace($Line)) { return }

    if ($Line -match "^\[(?<ts>[^\]]+)\]\s+\S+\s+(?<method>GET|POST|PUT|DELETE|PATCH)\s+(?<path>\S+)\s+(?<status>\d{3})\s+(?<latency>\S+)") {
        $method = $Matches.method
        $path = $Matches.path
        $status = [int]$Matches.status
        $latency = $Matches.latency
        $tsRaw = $Matches.ts
        $timePart = if ($tsRaw -match "(?<hh>\d{2}:\d{2}:\d{2})$") { $Matches.hh } else { "--:--:--" }
        $category = if ($method -eq "GET" -and $status -eq 200 -and $path -match ("^/api/v1/datalink/(" + $script:DatalinkDashboardEntities + "|modbus-share/status)$")) { "dashboard-refresh" } else { $null }
        if ($category -and -not $Verbose) {
            if (-not $script:LogNoiseCounters.ContainsKey($category)) { $script:LogNoiseCounters[$category] = 0 }
            $script:LogNoiseCounters[$category]++
            return
        }
        $methodCol = Format-FixedColumn -Text $method -Width 6
        $statusCol = ("{0,3}" -f $status)
        $latencyCol = ("{0,9}" -f $latency)
        $formatted = "[HTTP] $timePart  $methodCol $path $statusCol $latencyCol"
        $color = if ($status -ge 500) { "Red" } elseif ($status -ge 400) { "Yellow" } else { "DarkGray" }
        Write-ColorOutput $formatted $color
        return
    }

    if ($Line -match "^(?<ts>\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2})\s+[^:]+:\d+:\s+(?<msg>.+)$") {
        $msg = $Matches.msg
        $category = Get-LogNoiseCategory -Line $Line
        if ($category -and -not $Verbose) {
            if (-not $script:LogNoiseCounters.ContainsKey($category)) { $script:LogNoiseCounters[$category] = 0 }
            $script:LogNoiseCounters[$category]++
            return
        }
        if ($msg -match "127\.0\.0\.1:$script:ModbusSharePort") {
            Write-ColorOutput "[BOOT] Local Modbus share started on 127.0.0.1:$script:ModbusSharePort" "Green"
        } elseif ($msg -match "測試工具伺服器啟動於\s*(?<url>https?://\S+)") {
            Write-ColorOutput ("[BOOT] Server started at {0}" -f $Matches.url) "Green"
        } elseif ($msg -match "測試工具伺服器啟動於") {
            Write-ColorOutput "[BOOT] Server started" "Green"
        } elseif ($msg -match "資料庫路徑") {
            Write-ColorOutput "[BOOT] Database initialized" "DarkGray"
        } elseif ($Verbose) {
            Write-ColorOutput ("[BOOT] {0}" -f $msg) "DarkGray"
        }
        return
    }

    $category = Get-LogNoiseCategory -Line $Line
    if ($category -and -not $Verbose) {
        if (-not $script:LogNoiseCounters.ContainsKey($category)) { $script:LogNoiseCounters[$category] = 0 }
        $script:LogNoiseCounters[$category]++
        return
    }

    if ($Line -match "ERROR|Error|panic|FATAL|❌") {
        Write-ColorOutput $Line "Red"
    } elseif ($Line -match "WARN|Warning|⚠") {
        Write-ColorOutput $Line "Yellow"
    } elseif ($Line -match "啟動於|本機 Modbus 分享服務已啟動") {
        Write-ColorOutput $Line "Green"
    } else {
        Write-ColorOutput $Line "DarkGray"
    }
}

function Format-FixedColumn {
    param(
        [string]$Text,
        [int]$Width
    )
    if ($null -eq $Text) { $Text = "" }
    if ($Text.Length -gt $Width) {
        return ($Text.Substring(0, [Math]::Max(0, $Width - 1)) + "…")
    }
    return $Text.PadRight($Width)
}

function Show-LogNoiseSummary {
    if ($script:LogNoiseCounters.Count -eq 0 -or $Verbose) { return }
    Write-Section "已隱藏雜訊日誌"
    $script:LogNoiseCounters.GetEnumerator() | ForEach-Object {
        Write-ColorOutput "- $($_.Key): $($_.Value) 行" "DarkYellow"
    }
    Write-Info "可加上 -Verbose 顯示全部原始日誌。"
}
