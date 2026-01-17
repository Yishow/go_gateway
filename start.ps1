# ============================================
# Go Gateway 一鍵啟動腳本
# ============================================
# 功能：
#   1. 一鍵啟動專案（構建 + 啟動後端 API 服務）
#   2. 執行 golangci-lint 靜態分析
#   3. 構建後端 API 服務
#   4. 啟動服務（可選）
#   5. 代碼質量檢查
#   6. 執行單元測試（可選）
# ============================================
# 架構說明：
#   - 統一後端 API 服務：所有功能通過 REST API 提供
#   - 前端通過 HTTP 調用 API，無需命令列工具
#   - 單一可執行文件：gateway.exe（後端服務）
# ============================================
# bin 目錄說明：
#   - 存放編譯後的可執行文件（.exe）
#   - 統一管理構建產物，便於部署和執行
#   - 避免可執行文件散落在源碼目錄中，保持項目結構整潔
#   - 構建的檔案：gateway.exe（後端 API 服務）
# ============================================
# 使用範例：
#   .\start.ps1                    # 顯示主選單
#   .\start.ps1 -DevMode           # 開發模式（使用 go run，無需編譯 exe）
#   .\start.ps1 -AirMode           # 熱重載模式（使用 Air，自動檢測變更並重啟）
#   .\start.ps1 -QuickStart        # 一鍵啟動專案（構建 + 啟動）
#   .\start.ps1 -Start             # 構建後啟動服務
#   .\start.ps1 -SkipTest          # 跳過測試
#   .\start.ps1 -SkipBuild -Start  # 僅啟動（不構建）
#   .\start.ps1 -QuickStart -AutoKillPort  # 自動清理端口並啟動
#   .\start.ps1 -Start -Port 8080 -AutoKillPort  # 指定端口並自動清理
#
# 新增功能：
#   .\start.ps1 -CheckEnv          # 檢查開發環境（Go、Node.js 等）
#   .\start.ps1 -HealthCheck       # 健康檢查（檢查服務是否運行）
#   .\start.ps1 -ListProcesses    # 列出運行中的服務進程
#   .\start.ps1 -StopAll           # 停止所有運行中的服務
#   .\start.ps1 -Diagnose          # 快速診斷常見問題
# ============================================

param(
    [switch]$DevMode,             # 開發模式（使用 go run，無需編譯 exe）
    [switch]$AirMode,             # 熱重載模式（使用 Air，自動檢測變更並重啟）
    [switch]$QuickStart,          # 一鍵啟動專案（構建 + 啟動）
    [switch]$SkipLint,            # 跳過 lint 檢查
    [switch]$SkipTest,            # 跳過測試
    [switch]$SkipBuild,           # 跳過構建
    [switch]$SkipQuality,         # 跳過代碼質量檢查
    [switch]$Start,               # 構建後啟動服務
    [string]$Target = "",         # 已廢棄：統一為 gateway 服務
    [switch]$Coverage,            # 顯示詳細覆蓋率
    [switch]$Verbose,             # 詳細輸出
    [int]$Port = 8080,            # 服務端口（預設 8080）
    [switch]$AutoKillPort,        # 自動清理佔用端口的進程（不詢問）
    [switch]$CheckEnv,            # 檢查開發環境
    [switch]$HealthCheck,         # 健康檢查（檢查服務是否運行）
    [switch]$ListProcesses,       # 列出運行中的服務進程
    [switch]$StopAll,             # 停止所有運行中的服務
    [switch]$Diagnose             # 快速診斷常見問題
)

$ErrorActionPreference = "Stop"
$script:ExitCode = 0

# ============================================
# 常數定義
# ============================================
# 保存腳本啟動時的根目錄（專案根目錄）
$script:ROOT_DIR = Get-Location

$script:APP_NAME = "gateway"
$script:APP_PATH = "cmd/test_ui"
$script:BUILD_DIR = "bin"
$script:FRONTEND_DIR = "web/test-ui"
$script:STATIC_DIR = "cmd/test_ui/static"
$script:DIST_DIR = "web/test-ui/dist"
$script:NODE_MODULES_DIR = "web/test-ui/node_modules"
$script:LOG_DIR = "logs"
$script:LOG_FILE = Join-Path $script:LOG_DIR "start-$(Get-Date -Format 'yyyyMMdd').log"

# 從環境變數讀取端口配置（如果未指定）
if ($Port -eq 8080) {
    $envPort = [System.Environment]::GetEnvironmentVariable("PORT")
    if ($envPort -and $envPort -match '^\d+$') {
        $Port = [int]$envPort
        Write-Info "從環境變數讀取端口配置: $Port"
    }
}

# 顏色輸出函數
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "✅ $Message" "Green"
}

function Write-Error {
    param([string]$Message)
    Write-ColorOutput "❌ $Message" "Red"
    $script:ExitCode = 1
}

function Write-Info {
    param([string]$Message)
    Write-ColorOutput "ℹ️  $Message" "Cyan"
}

function Write-Warning {
    param([string]$Message)
    Write-ColorOutput "⚠️  $Message" "Yellow"
}

# ============================================
# 工具函數
# ============================================

# 初始化日誌目錄
function Initialize-LogDirectory {
    if (-not (Test-Path $script:LOG_DIR)) {
        New-Item -ItemType Directory -Path $script:LOG_DIR -Force | Out-Null
    }
}

# 記錄日誌
function Write-Log {
    param(
        [string]$Message,
        [string]$Level = "INFO"
    )
    
    Initialize-LogDirectory
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logMessage = "[$timestamp] [$Level] $Message"
    Add-Content -Path $script:LOG_FILE -Value $logMessage -ErrorAction SilentlyContinue
    
    if ($Verbose) {
        Write-Host $logMessage
    }
}

# 檢查命令是否存在
function Test-Command {
    param([string]$Command)
    $null = Get-Command $Command -ErrorAction SilentlyContinue
    return $?
}

# 獲取命令版本
function Get-CommandVersion {
    param([string]$Command)
    
    try {
        $version = & $Command --version 2>&1 | Select-Object -First 1
        return $version
    } catch {
        return "未知"
    }
}

# 檢查開發環境
function Test-DevelopmentEnvironment {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   環境檢查" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    $allOk = $true
    
    # 檢查 Go
    Write-Info "檢查 Go 環境..."
    if (Test-Command "go") {
        $goVersion = Get-CommandVersion "go"
        Write-Success "Go: $goVersion"
        Write-Log "Go 環境檢查通過: $goVersion"
    } else {
        Write-Error "Go 未安裝或不在 PATH 中"
        Write-Info "請訪問: https://golang.org/dl/"
        $allOk = $false
        Write-Log "Go 環境檢查失敗" "ERROR"
    }
    
    # 檢查 Node.js
    Write-Info "檢查 Node.js 環境..."
    if (Test-Command "node") {
        $nodeVersion = Get-CommandVersion "node"
        Write-Success "Node.js: $nodeVersion"
        Write-Log "Node.js 環境檢查通過: $nodeVersion"
    } else {
        Write-Warning "Node.js 未安裝或不在 PATH 中"
        Write-Info "前端功能可能無法使用"
        Write-Log "Node.js 環境檢查失敗" "WARN"
    }
    
    # 檢查 pnpm
    Write-Info "檢查 pnpm 環境..."
    if (Test-Command "pnpm") {
        $pnpmVersion = Get-CommandVersion "pnpm"
        Write-Success "pnpm: $pnpmVersion"
        Write-Log "pnpm 環境檢查通過: $pnpmVersion"
    } else {
        Write-Warning "pnpm 未安裝或不在 PATH 中"
        Write-Info "前端功能可能無法使用"
        Write-Log "pnpm 環境檢查失敗" "WARN"
    }
    
    # 檢查 golangci-lint
    Write-Info "檢查 golangci-lint..."
    if (Test-Command "golangci-lint") {
        $lintVersion = Get-CommandVersion "golangci-lint"
        Write-Success "golangci-lint: $lintVersion"
        Write-Log "golangci-lint 檢查通過: $lintVersion"
    } else {
        Write-Warning "golangci-lint 未安裝"
        Write-Info "可以使用: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        Write-Log "golangci-lint 未安裝" "WARN"
    }
    
    # 檢查 Air
    Write-Info "檢查 Air..."
    if (Test-Command "air") {
        $airVersion = Get-CommandVersion "air"
        Write-Success "Air: $airVersion"
        Write-Log "Air 檢查通過: $airVersion"
    } else {
        Write-Warning "Air 未安裝"
        Write-Info "可以使用: go install github.com/air-verse/air@latest"
        Write-Log "Air 未安裝" "WARN"
    }
    
    # 檢查專案結構
    Write-Info "檢查專案結構..."
    $requiredDirs = @($script:APP_PATH, $script:FRONTEND_DIR, "internal", "go.mod")
    foreach ($dir in $requiredDirs) {
        if (Test-Path $dir) {
            Write-Success "✓ $dir"
        } else {
            Write-Error "✗ $dir 不存在"
            $allOk = $false
        }
    }
    
    Write-ColorOutput ""
    if ($allOk) {
        Write-Success "環境檢查完成，所有必需項目正常"
    } else {
        Write-Warning "環境檢查完成，發現一些問題（見上方）"
    }
    
    return $allOk
}

# 健康檢查（檢查服務是否運行）
function Test-ServiceHealth {
    param([int]$Port = 8080)
    
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   健康檢查" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    # 檢查端口是否被佔用
    if (Test-PortInUse -Port $Port) {
        Write-Info "檢查端口 $Port 狀態..."
        $process = Get-ProcessByPort -Port $Port
        if ($process) {
            Write-Success "端口 $Port 被進程佔用: $($process.ProcessName) (PID: $($process.Id))"
            
            # 嘗試連接 HTTP 服務
            try {
                $response = Invoke-WebRequest -Uri "http://localhost:$Port/api/v1/test/status" -TimeoutSec 2 -ErrorAction Stop
                if ($response.StatusCode -eq 200) {
                    Write-Success "✅ 服務健康檢查通過！"
                    Write-Info "服務正常運行在 http://localhost:$Port"
                    return $true
                }
            } catch {
                Write-Warning "服務端口被佔用，但無法連接 HTTP 服務"
                Write-Info "可能不是 Go Gateway 服務，或服務未正常啟動"
            }
        }
    } else {
        Write-Warning "端口 $Port 未被佔用，服務未運行"
    }
    
    return $false
}

# 列出運行中的服務進程
function Get-RunningServices {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   運行中的服務" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    $found = $false
    
    # 檢查 gateway.exe
    $gatewayProcesses = Get-Process -Name "gateway" -ErrorAction SilentlyContinue
    if ($gatewayProcesses) {
        $found = $true
        Write-Info "Gateway 服務:"
        foreach ($proc in $gatewayProcesses) {
            Write-ColorOutput "  PID: $($proc.Id) | 記憶體: $([math]::Round($proc.WS / 1MB, 2)) MB | 啟動時間: $($proc.StartTime)" "Cyan"
        }
    }
    
    # 檢查 test-ui.exe
    $testUiProcesses = Get-Process -Name "test-ui" -ErrorAction SilentlyContinue
    if ($testUiProcesses) {
        $found = $true
        Write-Info "Test-UI 服務:"
        foreach ($proc in $testUiProcesses) {
            Write-ColorOutput "  PID: $($proc.Id) | 記憶體: $([math]::Round($proc.WS / 1MB, 2)) MB | 啟動時間: $($proc.StartTime)" "Cyan"
        }
    }
    
    # 檢查端口佔用
    Write-Info "端口佔用情況:"
    for ($p = 8080; $p -le 8090; $p++) {
        if (Test-PortInUse -Port $p) {
            $proc = Get-ProcessByPort -Port $p
            if ($proc) {
                $found = $true
                Write-ColorOutput "  端口 $p : $($proc.ProcessName) (PID: $($proc.Id))" "Cyan"
            }
        }
    }
    
    if (-not $found) {
        Write-Info "未發現運行中的服務"
    }
}

# 停止所有運行中的服務
function Stop-AllServices {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   停止所有服務" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    $stopped = 0
    
    # 停止 gateway.exe
    $gatewayProcesses = Get-Process -Name "gateway" -ErrorAction SilentlyContinue
    if ($gatewayProcesses) {
        foreach ($proc in $gatewayProcesses) {
            try {
                Stop-Process -Id $proc.Id -Force
                Write-Success "已停止 Gateway 服務 (PID: $($proc.Id))"
                $stopped++
            } catch {
                Write-Warning "無法停止進程 $($proc.Id): $_"
            }
        }
    }
    
    # 停止 test-ui.exe
    $testUiProcesses = Get-Process -Name "test-ui" -ErrorAction SilentlyContinue
    if ($testUiProcesses) {
        foreach ($proc in $testUiProcesses) {
            try {
                Stop-Process -Id $proc.Id -Force
                Write-Success "已停止 Test-UI 服務 (PID: $($proc.Id))"
                $stopped++
            } catch {
                Write-Warning "無法停止進程 $($proc.Id): $_"
            }
        }
    }
    
    # 清理端口
    for ($p = 8080; $p -le 8090; $p++) {
        if (Test-PortInUse -Port $p) {
            $proc = Get-ProcessByPort -Port $p
            if ($proc -and ($proc.ProcessName -eq "gateway" -or $proc.ProcessName -eq "test-ui")) {
                try {
                    Stop-ProcessByPort -Port $p -Force
                    $stopped++
                } catch {
                    # 忽略錯誤
                }
            }
        }
    }
    
    if ($stopped -eq 0) {
        Write-Info "未發現需要停止的服務"
    } else {
        Write-Success "已停止 $stopped 個服務/進程"
    }
}

# 快速診斷
function Start-Diagnose {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   快速診斷" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    $issues = @()
    
    # 檢查環境
    Write-Info "1. 檢查開發環境..."
    if (-not (Test-Command "go")) {
        $issues += "Go 未安裝"
    }
    if (-not (Test-Command "node")) {
        $issues += "Node.js 未安裝（前端功能可能無法使用）"
    }
    
    # 檢查專案結構
    Write-Info "2. 檢查專案結構..."
    if (-not (Test-Path "go.mod")) {
        $issues += "go.mod 不存在，可能不是 Go 專案根目錄"
    }
    if (-not (Test-Path $script:APP_PATH)) {
        $issues += "應用程式目錄不存在: $script:APP_PATH"
    }
    
    # 檢查端口
    Write-Info "3. 檢查端口狀態..."
    if (Test-PortInUse -Port $Port) {
        $proc = Get-ProcessByPort -Port $Port
        if ($proc) {
            if ($proc.ProcessName -ne "gateway" -and $proc.ProcessName -ne "test-ui") {
                $issues += "端口 $Port 被其他程序佔用: $($proc.ProcessName)"
            }
        }
    }
    
    # 檢查前端
    Write-Info "4. 檢查前端..."
    if (-not (Test-Path $script:NODE_MODULES_DIR)) {
        $issues += "前端依賴未安裝，需要執行: cd web/test-ui && pnpm install"
    }
    
    # 檢查構建產物
    Write-Info "5. 檢查構建產物..."
    $exePath = Join-Path $script:BUILD_DIR "$($script:APP_NAME).exe"
    if (-not (Test-Path $exePath)) {
        Write-Info "  構建產物不存在（這是正常的，如果尚未構建）"
    }
    
    # 輸出診斷結果
    Write-ColorOutput ""
    if ($issues.Count -eq 0) {
        Write-Success "✅ 診斷完成，未發現問題"
    } else {
        Write-Warning "發現 $($issues.Count) 個潛在問題："
        foreach ($issue in $issues) {
            Write-ColorOutput "  ⚠️  $issue" "Yellow"
        }
    }
}

# 啟動前端開發伺服器
function Start-FrontendDevServer {
    [CmdletBinding()]
    param()
    
    Write-Info "🎨 檢查前端開發環境..."
    
    # 檢查前端目錄
    if (-not (Test-Path $script:FRONTEND_DIR)) {
        Write-Error "找不到前端目錄: $script:FRONTEND_DIR"
        return $null
    }
    
    # 檢查 package.json
    $packageJsonPath = Join-Path $script:FRONTEND_DIR "package.json"
    if (-not (Test-Path $packageJsonPath)) {
        Write-Error "找不到 package.json: $packageJsonPath"
        return $null
    }
    
    # 檢查 node_modules
    if (-not (Test-Path $script:NODE_MODULES_DIR)) {
        Write-Warning "前端依賴未安裝，正在安裝..."
        $originalLocation = Get-Location
        try {
            Push-Location $script:FRONTEND_DIR
            pnpm install
            if ($LASTEXITCODE -ne 0) {
                Write-Error "前端依賴安裝失敗"
                return $null
            }
            Write-Success "前端依賴安裝完成"
        } catch {
            Write-Error "安裝前端依賴時發生錯誤: $_"
            return $null
        } finally {
            Pop-Location
            Set-Location $originalLocation
        }
    }
    
    # 啟動前端開發伺服器
    Write-Info "🚀 啟動前端開發伺服器..."
    $originalLocation = Get-Location
    
    try {
        Push-Location $script:FRONTEND_DIR
        
        # 使用 Start-Process 在背景啟動前端伺服器
        # 使用 cmd.exe 來正確處理 pnpm 命令，避免 PowerShell 的問題
        $frontendProcess = Start-Process -FilePath "cmd.exe" -ArgumentList "/c", "pnpm run dev" -PassThru -WindowStyle Hidden -WorkingDirectory (Get-Location).Path
        
        if ($frontendProcess) {
            Write-Success "前端開發伺服器已啟動（PID: $($frontendProcess.Id)）"
            Write-Info "💡 前端開發伺服器通常運行在 http://localhost:5173"
            Write-Info "💡 前端修改會自動熱重載"
            
            # 等待一小段時間確認伺服器啟動
            Start-Sleep -Milliseconds 1000
            
            # 檢查進程是否仍在運行
            if (-not (Get-Process -Id $frontendProcess.Id -ErrorAction SilentlyContinue)) {
                Write-Warning "前端開發伺服器可能啟動失敗"
                return $null
            }
            
            return $frontendProcess
        } else {
            Write-Error "無法啟動前端開發伺服器"
            return $null
        }
    } catch {
        Write-Error "啟動前端開發伺服器時發生錯誤: $_"
        return $null
    } finally {
        Pop-Location
        Set-Location $originalLocation
    }
}

# 停止前端開發伺服器
function Stop-FrontendDevServer {
    [CmdletBinding()]
    param(
        [System.Diagnostics.Process]$Process
    )
    
    if ($Process -and -not $Process.HasExited) {
        Write-Info "正在停止前端開發伺服器（PID: $($Process.Id)）..."
        try {
            $Process.Kill()
            $Process.WaitForExit(3000)
            Write-Success "前端開發伺服器已停止"
        } catch {
            Write-Warning "停止前端開發伺服器時發生錯誤: $_"
        }
    }
}

# 建置前端（如果需要）
function Build-Frontend {
    [CmdletBinding()]
    param(
        [switch]$Force  # 強制重新建置
    )
    
    Write-Info "📦 檢查前端檔案..."
    
    # 檢查前端檔案是否已存在
    $staticExists = (Test-Path $script:STATIC_DIR) -and (Test-Path "$script:STATIC_DIR\index.html")
    
    if ($staticExists -and -not $Force) {
        # 檢查前端原始碼是否有修改
        $needsRebuild = $false
        
        # 獲取前端原始碼目錄的最新修改時間
        $frontendSrcPath = Join-Path $script:FRONTEND_DIR "src"
        $packageJsonPath = Join-Path $script:FRONTEND_DIR "package.json"
        
        if (Test-Path $frontendSrcPath) {
            $srcFiles = Get-ChildItem -Path $frontendSrcPath -Recurse -File -ErrorAction SilentlyContinue
            if ($srcFiles) {
                $srcLastWrite = ($srcFiles | Measure-Object -Property LastWriteTime -Maximum).Maximum
            } else {
                $srcLastWrite = $null
            }
        } else {
            $srcLastWrite = $null
        }
        
        # 獲取 package.json 的修改時間（依賴變更）
        if (Test-Path $packageJsonPath) {
            $packageJsonTime = (Get-Item $packageJsonPath).LastWriteTime
            if ($srcLastWrite -and $packageJsonTime -gt $srcLastWrite) {
                $srcLastWrite = $packageJsonTime
            } elseif (-not $srcLastWrite) {
                $srcLastWrite = $packageJsonTime
            }
        }
        
        # 獲取建置產物的最新修改時間
        $staticFiles = Get-ChildItem -Path $script:STATIC_DIR -Recurse -File -ErrorAction SilentlyContinue
        if ($staticFiles) {
            $staticLastWrite = ($staticFiles | Measure-Object -Property LastWriteTime -Maximum).Maximum
        } else {
            $staticLastWrite = $null
        }
        
        # 如果原始碼比建置產物新，需要重新建置
        if ($srcLastWrite -and $staticLastWrite -and $srcLastWrite -gt $staticLastWrite) {
            $needsRebuild = $true
            Write-Info "檢測到前端原始碼有修改（原始碼: $($srcLastWrite.ToString('yyyy-MM-dd HH:mm:ss')), 建置產物: $($staticLastWrite.ToString('yyyy-MM-dd HH:mm:ss'))）"
        } elseif (-not $staticLastWrite) {
            # 如果建置產物目錄存在但沒有檔案，也需要重新建置
            $needsRebuild = $true
            Write-Info "建置產物目錄存在但沒有檔案，需要重新建置"
        }
        
        if (-not $needsRebuild) {
            Write-Success "前端檔案已存在且為最新，跳過建置"
            return $true
        } else {
            Write-Warning "前端原始碼有修改，需要重新建置..."
        }
    } else {
        if ($Force) {
            Write-Warning "強制重新建置前端..."
        } else {
            Write-Warning "前端檔案不存在，正在建置前端..."
        }
    }
    
    # 確保在專案根目錄
    $originalLocation = Get-Location
    
    try {
        # 檢查並安裝依賴
        if (-not (Test-Path $script:NODE_MODULES_DIR)) {
            Write-Info "📥 安裝前端依賴..."
            Push-Location $script:FRONTEND_DIR
            try {
                pnpm install
                if ($LASTEXITCODE -ne 0) {
                    Write-Error "前端依賴安裝失敗"
                    return $false
                }
            } finally {
                Pop-Location
            }
        }
        
        # 建置前端
        Write-Info "🔨 建置前端..."
        Push-Location $script:FRONTEND_DIR
        try {
            pnpm run build
            if ($LASTEXITCODE -ne 0) {
                Write-Error "前端建置失敗"
                return $false
            }
        } finally {
            Pop-Location
        }
        
        # 複製前端檔案
        if (Test-Path $script:STATIC_DIR) {
            Remove-Item -Recurse -Force $script:STATIC_DIR
        }
        Copy-Item -Recurse $script:DIST_DIR $script:STATIC_DIR
        Write-Success "前端建置完成"
        return $true
    } catch {
        Write-Error "建置前端時發生錯誤: $_"
        return $false
    } finally {
        Set-Location $originalLocation
    }
}

# 構建 Go 應用程式
function Build-Application {
    [CmdletBinding()]
    param(
        [string]$OutputPath,
        [string]$SourcePath
    )
    
    Write-Info "構建應用程式: $SourcePath -> $OutputPath"
    
    try {
        # 確保輸出目錄存在
        $outputDir = Split-Path -Parent $OutputPath
        if (-not (Test-Path $outputDir)) {
            New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
        }
        
        # 構建參數：-s 移除符號表，-w 移除 DWARF 除錯資訊，-trimpath 移除檔案路徑資訊
        go build -ldflags "-s -w" -trimpath -o $OutputPath $SourcePath
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "構建成功: $OutputPath"
            return $true
        } else {
            Write-Error "構建失敗"
            return $false
        }
    } catch {
        Write-Error "構建時發生錯誤: $_"
        return $false
    }
}

# 啟動應用程式（GUI 或 CLI）
function Start-Application {
    [CmdletBinding()]
    param(
        [string]$ExePath,
        [int]$Port = 8080,
        [switch]$IsGUI
    )
    
    if (-not (Test-Path $ExePath)) {
        Write-Error "找不到可執行文件: $ExePath"
        return $false
    }
    
    Write-Info "正在啟動應用程式..."
    Write-Info "服務將監聽端口: $Port"
    
    if ($IsGUI) {
        Write-Info "這是一個 GUI 應用程式，將在背景運行並顯示在系統托盤"
        Write-Info "💡 請查看系統通知區（右下角）的圖示"
        Write-Info "💡 右鍵點擊圖示可以打開瀏覽器或退出應用程式"
        
        try {
            $process = Start-Process -FilePath $ExePath -PassThru -WindowStyle Hidden
            Write-Success "應用程式已啟動（PID: $($process.Id)）"
            Write-Info "應用程式正在背景運行，請查看系統托盤圖示"
            
            # 等待一小段時間確認應用程式啟動
            Start-Sleep -Milliseconds 500
            
            # 檢查進程是否仍在運行
            if (-not (Get-Process -Id $process.Id -ErrorAction SilentlyContinue)) {
                Write-Warning "應用程式可能啟動失敗，請檢查日誌或錯誤訊息"
                return $false
            } else {
                Write-Success "應用程式運行中，可以關閉此視窗"
                return $true
            }
        } catch {
            Write-Error "啟動應用程式失敗: $_"
            return $false
        }
    } else {
        Write-Info "按 Ctrl+C 可停止服務"
        Write-ColorOutput "`n--- 服務輸出開始 ---" "Cyan"
        
        try {
            & $ExePath
            $serviceExitCode = $LASTEXITCODE
            
            Write-ColorOutput "--- 服務輸出結束 ---`n" "Cyan"
            
            if ($serviceExitCode -eq 0) {
                Write-Success "服務正常退出"
            } else {
                Write-Warning "服務退出，退出碼: $serviceExitCode"
            }
            return $true
        } catch {
            Write-Error "啟動服務失敗: $_"
            return $false
        }
    }
}

# 檢查端口是否被佔用
function Test-PortInUse {
    param([int]$Port)
    
    try {
        $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        return $null -ne $connection
    } catch {
        # 如果 Get-NetTCPConnection 不可用，使用 netstat
        $netstatOutput = netstat -ano | Select-String ":$Port\s"
        return $null -ne $netstatOutput
    }
}

# 取得佔用指定端口的進程資訊
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
    } catch {
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
            } else {
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
        } catch {
            Write-Error "無法終止進程 $($process.ProcessName) (PID: $($process.Id)): $_"
            return $false
        }
    } else {
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
        } else {
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
                } else {
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
    } else {
        Write-Info "端口 $Port 可用"
    }
    
    return $true
}

# 顯示主選單
function Show-MainMenu {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   Go Gateway 一鍵啟動腳本" "Cyan"
    Write-ColorOutput "============================================" "Cyan"
    Write-ColorOutput ""
    Write-ColorOutput "請選擇要執行的操作：" "Yellow"
    Write-ColorOutput "  [1] 🚀 開發模式（go run，無需編譯 exe）" "Green"
    Write-ColorOutput "  [2] 🔥 熱重載模式（Air，自動檢測變更並重啟）" "Magenta"
    Write-ColorOutput "  [3] 一鍵啟動專案（構建 + 啟動服務）" "Cyan"
    Write-ColorOutput "  [4] 完整流程（Lint + 構建 + 測試）" "Cyan"
    Write-ColorOutput "  [5] 僅構建可執行文件" "Cyan"
    Write-ColorOutput "  [6] 僅啟動服務" "Cyan"
    Write-ColorOutput "  [7] 執行測試" "Cyan"
    Write-ColorOutput ""
    Write-ColorOutput "  [8] 🔍 環境檢查" "Yellow"
    Write-ColorOutput "  [9] ❤️  健康檢查" "Yellow"
    Write-ColorOutput "  [A] 📋 列出運行中的服務" "Yellow"
    Write-ColorOutput "  [B] 🛑 停止所有服務" "Yellow"
    Write-ColorOutput "  [C] 🔧 快速診斷" "Yellow"
    Write-ColorOutput ""
    Write-ColorOutput "  [0] 退出" "Cyan"
    Write-ColorOutput ""
}

# ============================================
# 開發模式函數
# ============================================

# 熱重載模式（使用 Air，自動檢測變更並重啟）
function Start-AirMode {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   🔥 熱重載模式（Air）" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    Write-Info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後會自動重新運行"
    Write-Info "💡 此模式使用 Air 工具，自動檢測程式碼變更並重啟"
    Write-Info "💡 使用 go run，無需編譯 exe，與開發模式一致"
    Write-Info "💡 將同時啟動前端開發伺服器和後端服務"
    Write-ColorOutput ""
    
    # 詢問是否開啟瀏覽器
    Request-OpenBrowser | Out-Null
    
    # 檢查 Air 是否安裝
    if (-not (Test-Command "air")) {
        Write-Warning "Air 工具未安裝，正在嘗試安裝..."
        try {
            # Air 專案已遷移到新倉庫: github.com/air-verse/air
            go install github.com/air-verse/air@latest
            if (-not (Test-Command "air")) {
                Write-Error "無法安裝 Air，請手動安裝："
                Write-Info "  go install github.com/air-verse/air@latest"
                Write-Info "  或訪問: https://github.com/air-verse/air"
                Write-Info ""
                Write-Info "💡 建議: 使用選項 [1] 開發模式（go run）作為替代方案"
                return
            } else {
                Write-Success "Air 安裝成功"
            }
        } catch {
            Write-Error "安裝 Air 失敗: $_"
            Write-Info "💡 建議: 使用選項 [1] 開發模式（go run）作為替代方案"
            return
        }
    }
    
    # 啟動前端開發伺服器
    $frontendProcess = Start-FrontendDevServer
    if (-not $frontendProcess) {
        Write-Warning "前端開發伺服器啟動失敗，繼續啟動後端服務..."
    }
    
    Write-ColorOutput ""
    
    # 檢查 .air.toml 是否存在
    if (-not (Test-Path ".air.toml")) {
        Write-Warning ".air.toml 配置檔案不存在，Air 將使用預設配置"
    }
    
    # 運行 Air（Air 會在專案根目錄運行）
    Write-Info "▶️  啟動 Air 熱重載..."
    Write-Info "💡 修改程式碼後，Air 會自動檢測並使用 go run 重新運行"
    Write-Info "💡 注意: 如果看到 'CMD will not recognize' 警告，可忽略（Air 在 Windows 上的已知提示）"
    Write-ColorOutput ""
    
    try {
        # 確保在專案根目錄運行 Air
        Set-Location $script:ROOT_DIR
        # 將環境變數傳遞給 Air（Air 會傳遞給子進程）
        air
    } catch {
        Write-Error "啟動 Air 失敗: $_"
    } finally {
        # 確保回到專案根目錄
        Set-Location $script:ROOT_DIR
        
        # 停止前端開發伺服器
        if ($frontendProcess) {
            Stop-FrontendDevServer -Process $frontendProcess
        }
        
        # 清理環境變數
        Remove-Item Env:\AUTO_OPEN_BROWSER -ErrorAction SilentlyContinue
    }
}

# 詢問是否自動開啟瀏覽器
function Request-OpenBrowser {
    Write-ColorOutput ""
    Write-Info "是否要在啟動後自動開啟瀏覽器？"
    $response = Read-Host "請輸入 (Y/N，預設為 N)"
    
    if ($response -eq "Y" -or $response -eq "y") {
        $env:AUTO_OPEN_BROWSER = "true"
        Write-Info "將自動開啟瀏覽器"
        return $true
    } else {
        $env:AUTO_OPEN_BROWSER = "false"
        Write-Info "不會自動開啟瀏覽器，請手動點擊系統托盤圖示打開"
        return $false
    }
}

# 開發模式（使用 go run，無需編譯 exe）
function Start-DevMode {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   🚀 開發模式（無需編譯 exe）" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    Write-Info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後需要手動重新運行"
    Write-Info "💡 此模式使用 go run，無需編譯 exe，適合快速開發迭代"
    Write-Info "💡 將同時啟動前端開發伺服器和後端服務"
    Write-ColorOutput ""
    
    # 詢問是否開啟瀏覽器
    Request-OpenBrowser | Out-Null
    
    # 檢查應用程式路徑
    if (-not (Test-Path $script:APP_PATH)) {
        Write-Error "找不到應用程式: $script:APP_PATH"
        return
    }
    
    # 啟動前端開發伺服器
    $frontendProcess = Start-FrontendDevServer
    if (-not $frontendProcess) {
        Write-Warning "前端開發伺服器啟動失敗，繼續啟動後端服務..."
    }
    
    Write-ColorOutput ""
    
    # 運行應用程式
    Write-Info "▶️  啟動後端應用程式（使用 go run）..."
    Write-Info "💡 修改程式碼後，請按 Ctrl+C 停止並重新運行此選項"
    Write-ColorOutput ""
    
    try {
        Push-Location $script:APP_PATH
        go run .
    } catch {
        Write-Error "啟動應用程式失敗: $_"
    } finally {
        # 確保回到專案根目錄
        Pop-Location
        Set-Location $script:ROOT_DIR
        
        # 停止前端開發伺服器
        if ($frontendProcess) {
            Stop-FrontendDevServer -Process $frontendProcess
        }
        
        # 清理環境變數
        Remove-Item Env:\AUTO_OPEN_BROWSER -ErrorAction SilentlyContinue
    }
}

# 一鍵啟動專案（構建 + 啟動）
function Start-QuickStart {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   一鍵啟動專案" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    # 1. 構建可執行文件
    Write-ColorOutput "[1/2] 構建可執行文件..." "Yellow"
    
    $outputPath = Join-Path $script:BUILD_DIR "$($script:APP_NAME).exe"
    $sourcePath = "./$script:APP_PATH"
    
    if (-not (Build-Application -OutputPath $outputPath -SourcePath $sourcePath)) {
        Write-Error "構建失敗，無法啟動服務"
        return
    }
    
    # 2. 啟動後端 API 服務
    Write-ColorOutput "`n[2/2] 啟動後端 API 服務..." "Yellow"
    
    # 檢查並清理端口
    Write-Info "檢查端口 $Port 狀態..."
    if (-not (Clear-PortForService -Port $Port -AutoKill:$AutoKillPort)) {
        Write-Error "端口 $Port 清理失敗，無法啟動服務"
        $script:ExitCode = 1
        return
    }
    
    # 啟動服務（GUI 應用）
    $isGUIApp = $true  # gateway/test-ui 是 GUI 應用
    if (-not (Start-Application -ExePath $outputPath -Port $Port -IsGUI:$isGUIApp)) {
        $script:ExitCode = 1
    }
}

# 處理特殊功能參數
if ($CheckEnv) {
    Test-DevelopmentEnvironment
    exit $script:ExitCode
}

if ($HealthCheck) {
    Test-ServiceHealth -Port $Port
    exit $script:ExitCode
}

if ($ListProcesses) {
    Get-RunningServices
    exit $script:ExitCode
}

if ($StopAll) {
    Stop-AllServices
    exit $script:ExitCode
}

if ($Diagnose) {
    Start-Diagnose
    exit $script:ExitCode
}

# 如果指定了 DevMode，直接執行開發模式
if ($DevMode) {
    Start-DevMode
    exit $script:ExitCode
}

# 如果指定了 AirMode，直接執行熱重載模式
if ($AirMode) {
    Start-AirMode
    exit $script:ExitCode
}

# 如果指定了 QuickStart，直接執行一鍵啟動
if ($QuickStart) {
    Start-QuickStart
    exit $script:ExitCode
}

# 如果沒有提供任何參數，顯示主選單
$hasAnyParam = $SkipLint -or $SkipTest -or $SkipBuild -or $SkipQuality -or $Start -or $Coverage -or $Verbose -or (-not [string]::IsNullOrWhiteSpace($Target))
if (-not $hasAnyParam) {
    Show-MainMenu
    $menuSelection = Read-Host "請輸入選項 (0-7, 8-9, A-C)"
    
    switch ($menuSelection) {
        "1" {
            Start-DevMode
            exit $script:ExitCode
        }
        "2" {
            Start-AirMode
            exit $script:ExitCode
        }
        "3" {
            Start-QuickStart
            exit $script:ExitCode
        }
        "4" {
            # 完整流程，繼續執行後續步驟
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   執行完整流程" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "5" {
            $SkipLint = $true
            $SkipTest = $true
            $SkipQuality = $true
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   僅構建可執行文件" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "6" {
            $SkipLint = $true
            $SkipTest = $true
            $SkipBuild = $true
            $SkipQuality = $true
            $Start = $true
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   僅啟動服務" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "7" {
            $SkipLint = $true
            $SkipBuild = $true
            $SkipQuality = $true
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   執行測試" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "8" {
            Test-DevelopmentEnvironment
            exit $script:ExitCode
        }
        "9" {
            Test-ServiceHealth -Port $Port
            exit $script:ExitCode
        }
        "A" {
            Get-RunningServices
            exit $script:ExitCode
        }
        "a" {
            Get-RunningServices
            exit $script:ExitCode
        }
        "B" {
            Stop-AllServices
            exit $script:ExitCode
        }
        "b" {
            Stop-AllServices
            exit $script:ExitCode
        }
        "C" {
            Start-Diagnose
            exit $script:ExitCode
        }
        "c" {
            Start-Diagnose
            exit $script:ExitCode
        }
        "0" {
            Write-Info "退出"
            exit 0
        }
        default {
            Write-Error "無效的選項"
            exit 1
        }
    }
} else {
    # 有參數時顯示標題
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   Go Gateway 一鍵啟動腳本" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
}

# 1. 檢查 golangci-lint
Write-Info "檢查工具依賴..."
if (-not $SkipLint) {
    if (-not (Test-Command "golangci-lint")) {
        Write-Warning "golangci-lint 未安裝，正在嘗試安裝..."
        try {
            # 嘗試使用 go install 安裝
            go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
            if (-not (Test-Command "golangci-lint")) {
                Write-Error "無法安裝 golangci-lint，請手動安裝："
                Write-Info "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
                Write-Info "  或訪問: https://golangci-lint.run/usage/install/"
                $SkipLint = $true
            } else {
                Write-Success "golangci-lint 安裝成功"
            }
        } catch {
            Write-Error "安裝 golangci-lint 失敗: $_"
            $SkipLint = $true
        }
    }
}

# 計算總步驟數（實際會執行的步驟）
# 步驟順序：1. lint, 2. build, 3. start(可選), 4. quality(可選), 5. test(可選)
$totalSteps = 0
if (-not $SkipLint) { $totalSteps++ }
if (-not $SkipBuild) { $totalSteps++ }
if ($Start) { $totalSteps++ }
if (-not $SkipQuality) { $totalSteps++ }
if (-not $SkipTest) { $totalSteps++ }
$currentStep = 0

# 2. 執行 golangci-lint 靜態分析
if (-not $SkipLint) {
    $currentStep++
    Write-ColorOutput "`n[$currentStep/$totalSteps] 執行 golangci-lint 靜態分析..." "Yellow"
    Write-Info "這可能需要一些時間，請稍候..."
    
    try {
        $lintOutput = golangci-lint run ./... 2>&1
        $lintExitCode = $LASTEXITCODE
        
        if ($lintExitCode -eq 0) {
            Write-Success "golangci-lint 檢查通過，未發現問題"
        } else {
            Write-Error "golangci-lint 發現問題："
            Write-Host $lintOutput
            $script:ExitCode = 1
        }
    } catch {
        Write-Error "執行 golangci-lint 失敗: $_"
        $script:ExitCode = 1
    }
}

# 3. 構建可執行文件
if (-not $SkipBuild) {
    $currentStep++
    Write-ColorOutput "`n[$currentStep/$totalSteps] 構建可執行文件..." "Yellow"
    
    $outputPath = Join-Path $script:BUILD_DIR "$($script:APP_NAME).exe"
    $sourcePath = "./$script:APP_PATH"
    
    if (-not (Build-Application -OutputPath $outputPath -SourcePath $sourcePath)) {
        $script:ExitCode = 1
    }
}

# 4. 啟動服務（可選）
if ($Start) {
    $currentStep++
    Write-ColorOutput "`n[$currentStep/$totalSteps] 啟動服務..." "Yellow"
    
    # 檢查並清理端口
    Write-Info "檢查端口 $Port 狀態..."
    if (-not (Clear-PortForService -Port $Port -AutoKill:$AutoKillPort)) {
        Write-Error "端口 $Port 清理失敗，無法啟動服務"
        $script:ExitCode = 1
    } else {
        $exePath = Join-Path $script:BUILD_DIR "$($script:APP_NAME).exe"
        $isGUIApp = $true  # gateway/test-ui 是 GUI 應用
        
        if (-not (Start-Application -ExePath $exePath -Port $Port -IsGUI:$isGUIApp)) {
            $script:ExitCode = 1
        }
    }
}

# 5. 代碼質量檢查
if (-not $SkipQuality) {
    $currentStep++
    Write-ColorOutput "`n[$currentStep/$totalSteps] 代碼質量檢查..." "Yellow"

    # 檢查是否有未使用的導入
    Write-Info "檢查未使用的導入..."
    try {
        $unusedOutput = go run golang.org/x/tools/cmd/deadcode@latest ./... 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success "未發現未使用的代碼"
        } else {
            Write-Warning "發現未使用的代碼（這不是錯誤）"
            if ($Verbose) {
                Write-Host $unusedOutput
            }
        }
    } catch {
        Write-Info "跳過未使用代碼檢查（需要安裝 golang.org/x/tools/cmd/deadcode）"
    }

    # 檢查代碼格式
    Write-Info "檢查代碼格式..."
    try {
        $fmtOutput = go fmt ./... 2>&1
        if ($fmtOutput) {
            Write-Warning "代碼格式已自動修正"
            if ($Verbose) {
                Write-Host $fmtOutput
            }
        } else {
            Write-Success "代碼格式正確"
        }
    } catch {
        Write-Error "檢查代碼格式失敗: $_"
        $script:ExitCode = 1
    }
}

# 6. 執行單元測試（移到最後）
if (-not $SkipTest) {
    $currentStep++
    Write-ColorOutput "`n[$currentStep/$totalSteps] 執行單元測試..." "Yellow"
    
    try {
        if ($Coverage) {
            Write-Info "生成詳細覆蓋率報告..."
            $testOutput = go test -v -coverprofile=coverage.out -covermode=atomic ./internal/protocol/... 2>&1
            Write-Host $testOutput
            
            # 顯示覆蓋率摘要
            Write-ColorOutput "`n測試覆蓋率摘要：" "Cyan"
            $coverageSummary = go test -cover ./internal/protocol/... 2>&1 | Select-String -Pattern "coverage:"
            Write-Host $coverageSummary
            
            # 生成 HTML 報告
            if (Test-Path "coverage.out") {
                Write-Info "生成 HTML 覆蓋率報告: coverage.html"
                go tool cover -html=coverage.out -o coverage.html
                if ($?) {
                    Write-Success "覆蓋率報告已生成: coverage.html"
                }
            }
        } else {
            $testOutput = go test -cover ./internal/protocol/... 2>&1
            Write-Host $testOutput
            
            # 提取覆蓋率信息
            $coverageLines = $testOutput | Select-String -Pattern "coverage:"
            if ($coverageLines) {
                Write-ColorOutput "`n測試覆蓋率：" "Cyan"
                Write-Host $coverageLines
            }
        }
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "所有測試通過"
        } else {
            Write-Error "部分測試失敗"
            $script:ExitCode = 1
        }
    } catch {
        Write-Error "執行測試失敗: $_"
        $script:ExitCode = 1
    }
}

# 總結
Write-ColorOutput "`n============================================" "Cyan"
if ($script:ExitCode -eq 0) {
    Write-Success "所有檢查完成！"
    Write-ColorOutput "============================================`n" "Cyan"
} else {
    Write-Error "檢查完成，但發現問題，請查看上方錯誤信息"
    Write-ColorOutput "============================================`n" "Cyan"
}

exit $script:ExitCode
