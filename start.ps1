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
#   .\start.ps1 -Start -Target fatek_test  # 啟動指定服務
#   .\start.ps1 -SkipTest          # 跳過測試
#   .\start.ps1 -SkipBuild -Start  # 僅啟動（不構建）
#   .\start.ps1 -QuickStart -AutoKillPort  # 自動清理端口並啟動
#   .\start.ps1 -Start -Port 8080 -AutoKillPort  # 指定端口並自動清理
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
    [switch]$AutoKillPort          # 自動清理佔用端口的進程（不詢問）
)

$ErrorActionPreference = "Stop"
$script:ExitCode = 0

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

# 檢查命令是否存在
function Test-Command {
    param([string]$Command)
    $null = Get-Command $Command -ErrorAction SilentlyContinue
    return $?
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
    Write-ColorOutput "  [0] 退出" "Cyan"
    Write-ColorOutput ""
}

# 熱重載模式（使用 Air，自動檢測變更並重啟）
function Start-AirMode {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   🔥 熱重載模式（Air）" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    Write-Info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後會自動重新編譯運行"
    Write-Info "💡 此模式使用 Air 工具，自動檢測程式碼變更並重啟"
    Write-ColorOutput ""
    
    # 檢查 Air 是否安裝
    if (-not (Test-Command "air")) {
        Write-Warning "Air 工具未安裝，正在嘗試安裝..."
        try {
            go install github.com/cosmtrek/air@latest
            if (-not (Test-Command "air")) {
                Write-Error "無法安裝 Air，請手動安裝："
                Write-Info "  go install github.com/cosmtrek/air@latest"
                Write-Info "  或訪問: https://github.com/cosmtrek/air"
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
    
    # 處理 test-ui 的前端建置
    Write-Info "📦 檢查前端檔案..."
    $staticPath = "cmd\test_ui\static"
    $distPath = "web\test-ui\dist"
    
    if (-not (Test-Path $staticPath) -or -not (Test-Path "$staticPath\index.html")) {
        Write-Warning "前端檔案不存在，正在建置前端..."
        
        # 檢查 node_modules
        $nodeModulesPath = "web\test-ui\node_modules"
        if (-not (Test-Path $nodeModulesPath)) {
            Write-Info "📥 安裝前端依賴..."
            Set-Location web/test-ui
            npm install
            if ($LASTEXITCODE -ne 0) {
                Write-Error "前端依賴安裝失敗"
                Set-Location ../..
                return
            }
            Set-Location ../..
        }
        
        # 建置前端
        Write-Info "🔨 建置前端..."
        Set-Location web/test-ui
        npm run build
        if ($LASTEXITCODE -ne 0) {
            Write-Error "前端建置失敗"
            Set-Location ../..
            return
        }
        Set-Location ../..
        
        # 複製前端檔案
        if (Test-Path $staticPath) {
            Remove-Item -Recurse -Force $staticPath
        }
        Copy-Item -Recurse $distPath $staticPath
        Write-Success "前端建置完成"
    } else {
        Write-Success "前端檔案已存在，跳過建置"
    }
    
    # 檢查 .air.toml 是否存在
    if (-not (Test-Path ".air.toml")) {
        Write-Warning ".air.toml 配置檔案不存在，Air 將使用預設配置"
    }
    
    # 運行 Air
    Write-ColorOutput ""
    Write-Info "▶️  啟動 Air 熱重載..."
    Write-Info "💡 修改程式碼後，Air 會自動檢測並重新編譯運行"
    Write-ColorOutput ""
    
    try {
        air
    } catch {
        Write-Error "啟動 Air 失敗: $_"
    }
}

# 開發模式（使用 go run，無需編譯 exe）
function Start-DevMode {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   🚀 開發模式（無需編譯 exe）" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    Write-Info "💡 提示: 使用 Ctrl+C 停止，修改程式碼後需要手動重新運行"
    Write-Info "💡 此模式使用 go run，無需編譯 exe，適合快速開發迭代"
    Write-ColorOutput ""
    
    # 檢查應用程式路徑
    $appPath = "cmd/test_ui"
    if (-not (Test-Path $appPath)) {
        Write-Error "找不到應用程式: $appPath"
        return
    }
    
    # 處理 test-ui 的前端建置
    Write-Info "📦 檢查前端檔案..."
    $staticPath = "cmd\test_ui\static"
    $distPath = "web\test-ui\dist"
    
    if (-not (Test-Path $staticPath) -or -not (Test-Path "$staticPath\index.html")) {
        Write-Warning "前端檔案不存在，正在建置前端..."
        
        # 檢查 node_modules
        $nodeModulesPath = "web\test-ui\node_modules"
        if (-not (Test-Path $nodeModulesPath)) {
            Write-Info "📥 安裝前端依賴..."
            Set-Location web/test-ui
            npm install
            if ($LASTEXITCODE -ne 0) {
                Write-Error "前端依賴安裝失敗"
                Set-Location ../..
                return
            }
            Set-Location ../..
        }
        
        # 建置前端
        Write-Info "🔨 建置前端..."
        Set-Location web/test-ui
        npm run build
        if ($LASTEXITCODE -ne 0) {
            Write-Error "前端建置失敗"
            Set-Location ../..
            return
        }
        Set-Location ../..
        
        # 複製前端檔案
        if (Test-Path $staticPath) {
            Remove-Item -Recurse -Force $staticPath
        }
        Copy-Item -Recurse $distPath $staticPath
        Write-Success "前端建置完成"
    } else {
        Write-Success "前端檔案已存在，跳過建置"
    }
    
    # 運行應用程式
    Write-ColorOutput ""
    Write-Info "▶️  啟動應用程式（使用 go run）..."
    Write-Info "💡 修改程式碼後，請按 Ctrl+C 停止並重新運行此選項"
    Write-ColorOutput ""
    
    Set-Location $appPath
    try {
        go run .
    } catch {
        Write-Error "啟動應用程式失敗: $_"
    } finally {
        Set-Location ../..
    }
}

# 一鍵啟動專案（構建 + 啟動）
function Start-QuickStart {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   一鍵啟動專案" "Cyan"
    Write-ColorOutput "============================================`n" "Cyan"
    
    # 1. 構建可執行文件
    Write-ColorOutput "[1/2] 構建可執行文件..." "Yellow"
    
    # 統一構建單一後端 API 服務
    # 所有功能都通過 REST API 提供，前端通過 HTTP 調用
    $buildTargets = @(
        @{Name="gateway"; Path="./cmd/test_ui"}
    )
    
    # bin 目錄用途：
    # - 存放編譯後的可執行文件（.exe）
    # - 統一管理構建產物，便於部署和執行
    # - 避免可執行文件散落在源碼目錄中，保持項目結構整潔
    $buildDir = "bin"
    if (-not (Test-Path $buildDir)) {
        New-Item -ItemType Directory -Path $buildDir | Out-Null
    }
    
    $buildSuccess = $true
    foreach ($target in $buildTargets) {
        Write-Info "構建 $($target.Name)..."
        try {
            $outputPath = Join-Path $buildDir "$($target.Name).exe"
            # gateway (test-ui) 不使用 -H=windowsgui 以顯示控制台窗口，讓用戶可以點擊 X 按鈕
            # -s: 移除符號表，-w: 移除 DWARF 除錯資訊，-trimpath: 移除檔案路徑資訊
            go build -ldflags "-s -w" -trimpath -o $outputPath $target.Path
            if ($LASTEXITCODE -eq 0) {
                Write-Success "$($target.Name) 構建成功: $outputPath"
            } else {
                Write-Error "$($target.Name) 構建失敗"
                $buildSuccess = $false
                $script:ExitCode = 1
            }
        } catch {
            Write-Error "構建 $($target.Name) 時發生錯誤: $_"
            $buildSuccess = $false
            $script:ExitCode = 1
        }
    }
    
    if (-not $buildSuccess) {
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
    
    # 統一為單一後端服務
    $targetToStart = "gateway"
    
    # 啟動服務
    $exePath = Join-Path "bin" "$targetToStart.exe"
    
    if (-not (Test-Path $exePath)) {
        Write-Error "找不到可執行文件: $exePath"
        Write-Info "請先執行構建步驟"
        $script:ExitCode = 1
        return
    }
    
    Write-Info "正在啟動 $targetToStart..."
    Write-Info "服務將監聽端口: $Port"
    
    # 檢查是否為 GUI 應用（使用 windowsgui 標誌編譯的應用）
    $isGUIApp = $exePath -like "*gateway.exe" -or $exePath -like "*test-ui.exe"
    
    if ($isGUIApp) {
        Write-Info "這是一個 GUI 應用程式，將在背景運行並顯示在系統托盤"
        Write-Info "💡 請查看系統通知區（右下角）的圖示"
        Write-Info "💡 右鍵點擊圖示可以打開瀏覽器或退出應用程式"
        
        try {
            # GUI 應用在背景啟動，不等待輸出
            $process = Start-Process -FilePath $exePath -PassThru -WindowStyle Hidden
            Write-Success "應用程式已啟動（PID: $($process.Id)）"
            Write-Info "應用程式正在背景運行，請查看系統托盤圖示"
            
            # 等待一小段時間確認應用程式啟動
            Start-Sleep -Milliseconds 500
            
            # 檢查進程是否仍在運行
            if (-not (Get-Process -Id $process.Id -ErrorAction SilentlyContinue)) {
                Write-Warning "應用程式可能啟動失敗，請檢查日誌或錯誤訊息"
                $script:ExitCode = 1
            } else {
                Write-Success "應用程式運行中，可以關閉此視窗"
            }
        } catch {
            Write-Error "啟動應用程式失敗: $_"
            $script:ExitCode = 1
        }
    } else {
        Write-Info "按 Ctrl+C 可停止服務"
        Write-ColorOutput "`n--- 服務輸出開始 ---" "Cyan"
        
        try {
            # 啟動服務（前台運行）
            & $exePath
            $serviceExitCode = $LASTEXITCODE
            
            Write-ColorOutput "--- 服務輸出結束 ---`n" "Cyan"
            
            if ($serviceExitCode -eq 0) {
                Write-Success "服務正常退出"
            } else {
                Write-Warning "服務退出，退出碼: $serviceExitCode"
            }
        } catch {
            Write-Error "啟動服務失敗: $_"
            $script:ExitCode = 1
        }
    }
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
    $menuSelection = Read-Host "請輸入選項 (0-7)"
    
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
    
    # 統一構建單一後端 API 服務
    # 所有功能都通過 REST API 提供，前端通過 HTTP 調用
    $buildTargets = @(
        @{Name="gateway"; Path="./cmd/test_ui"}
    )
    
    $buildDir = "bin"
    if (-not (Test-Path $buildDir)) {
        New-Item -ItemType Directory -Path $buildDir | Out-Null
    }
    
    foreach ($target in $buildTargets) {
        Write-Info "構建 $($target.Name)..."
        try {
            $outputPath = Join-Path $buildDir "$($target.Name).exe"
            # gateway (test-ui) 不使用 -H=windowsgui 以顯示控制台窗口，讓用戶可以點擊 X 按鈕
            # -s: 移除符號表，-w: 移除 DWARF 除錯資訊，-trimpath: 移除檔案路徑資訊
            go build -ldflags "-s -w" -trimpath -o $outputPath $target.Path
            if ($LASTEXITCODE -eq 0) {
                Write-Success "$($target.Name) 構建成功: $outputPath"
            } else {
                Write-Error "$($target.Name) 構建失敗"
                $script:ExitCode = 1
            }
        } catch {
            Write-Error "構建 $($target.Name) 時發生錯誤: $_"
            $script:ExitCode = 1
        }
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
        # 統一為單一後端 API 服務
        $targetToStart = "gateway"
        $exePath = Join-Path "bin" "$targetToStart.exe"
        
        if (-not [string]::IsNullOrWhiteSpace($targetToStart)) {
            
            if (-not (Test-Path $exePath)) {
                Write-Error "找不到可執行文件: $exePath"
                Write-Info "請先執行構建步驟"
                $script:ExitCode = 1
            } else {
                Write-Info "正在啟動 $targetToStart..."
                Write-Info "服務將監聽端口: $Port"
                
                # 檢查是否為 GUI 應用（使用 windowsgui 標誌編譯的應用）
                $isGUIApp = $exePath -like "*gateway.exe" -or $exePath -like "*test-ui.exe"
                
                if ($isGUIApp) {
                    Write-Info "這是一個 GUI 應用程式，將在背景運行並顯示在系統托盤"
                    Write-Info "💡 請查看系統通知區（右下角）的圖示"
                    Write-Info "💡 右鍵點擊圖示可以打開瀏覽器或退出應用程式"
                    
                    try {
                        # GUI 應用在背景啟動，不等待輸出
                        $process = Start-Process -FilePath $exePath -PassThru -WindowStyle Hidden
                        Write-Success "應用程式已啟動（PID: $($process.Id)）"
                        Write-Info "應用程式正在背景運行，請查看系統托盤圖示"
                        
                        # 等待一小段時間確認應用程式啟動
                        Start-Sleep -Milliseconds 500
                        
                        # 檢查進程是否仍在運行
                        if (-not (Get-Process -Id $process.Id -ErrorAction SilentlyContinue)) {
                            Write-Warning "應用程式可能啟動失敗，請檢查日誌或錯誤訊息"
                            $script:ExitCode = 1
                        } else {
                            Write-Success "應用程式運行中，可以關閉此視窗"
                        }
                    } catch {
                        Write-Error "啟動應用程式失敗: $_"
                        $script:ExitCode = 1
                    }
                } else {
                    Write-Info "按 Ctrl+C 可停止服務"
                    Write-ColorOutput "`n--- 服務輸出開始 ---" "Cyan"
                    
                    try {
                        # 啟動服務（前台運行）
                        & $exePath
                        $serviceExitCode = $LASTEXITCODE
                        
                        Write-ColorOutput "--- 服務輸出結束 ---`n" "Cyan"
                        
                        if ($serviceExitCode -eq 0) {
                            Write-Success "服務正常退出"
                        } else {
                            Write-Warning "服務退出，退出碼: $serviceExitCode"
                        }
                    } catch {
                        Write-Error "啟動服務失敗: $_"
                        $script:ExitCode = 1
                    }
                }
            }
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
