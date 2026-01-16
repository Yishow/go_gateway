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
#   .\start.ps1 -QuickStart        # 一鍵啟動專案（第一項）
#   .\start.ps1 -Start             # 構建後啟動服務
#   .\start.ps1 -Start -Target fatek_test  # 啟動指定服務
#   .\start.ps1 -SkipTest          # 跳過測試
#   .\start.ps1 -SkipBuild -Start  # 僅啟動（不構建）
# ============================================

param(
    [switch]$QuickStart,          # 一鍵啟動專案（構建 + 啟動）
    [switch]$SkipLint,            # 跳過 lint 檢查
    [switch]$SkipTest,            # 跳過測試
    [switch]$SkipBuild,           # 跳過構建
    [switch]$SkipQuality,         # 跳過代碼質量檢查
    [switch]$Start,               # 構建後啟動服務
    [string]$Target = "",         # 已廢棄：統一為 gateway 服務
    [switch]$Coverage,            # 顯示詳細覆蓋率
    [switch]$Verbose              # 詳細輸出
)

$ErrorActionPreference = "Stop"
$script:ExitCode = 0

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

# 顯示主選單
function Show-MainMenu {
    Write-ColorOutput "`n============================================" "Cyan"
    Write-ColorOutput "   Go Gateway 一鍵啟動腳本" "Cyan"
    Write-ColorOutput "============================================" "Cyan"
    Write-ColorOutput ""
    Write-ColorOutput "請選擇要執行的操作：" "Yellow"
    Write-ColorOutput "  [1] 一鍵啟動專案（構建 + 啟動服務）" "Cyan"
    Write-ColorOutput "  [2] 完整流程（Lint + 構建 + 測試）" "Cyan"
    Write-ColorOutput "  [3] 僅構建可執行文件" "Cyan"
    Write-ColorOutput "  [4] 僅啟動服務" "Cyan"
    Write-ColorOutput "  [5] 執行測試" "Cyan"
    Write-ColorOutput "  [0] 退出" "Cyan"
    Write-ColorOutput ""
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
            go build -o $outputPath $target.Path
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

# 如果指定了 QuickStart，直接執行一鍵啟動
if ($QuickStart) {
    Start-QuickStart
    exit $script:ExitCode
}

# 如果沒有提供任何參數，顯示主選單
$hasAnyParam = $SkipLint -or $SkipTest -or $SkipBuild -or $SkipQuality -or $Start -or $Coverage -or $Verbose -or (-not [string]::IsNullOrWhiteSpace($Target))
if (-not $hasAnyParam) {
    Show-MainMenu
    $menuSelection = Read-Host "請輸入選項 (0-5)"
    
    switch ($menuSelection) {
        "1" {
            Start-QuickStart
            exit $script:ExitCode
        }
        "2" {
            # 完整流程，繼續執行後續步驟
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   執行完整流程" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "3" {
            $SkipLint = $true
            $SkipTest = $true
            $SkipQuality = $true
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   僅構建可執行文件" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "4" {
            $SkipLint = $true
            $SkipTest = $true
            $SkipBuild = $true
            $SkipQuality = $true
            $Start = $true
            Write-ColorOutput "`n============================================" "Cyan"
            Write-ColorOutput "   僅啟動服務" "Cyan"
            Write-ColorOutput "============================================`n" "Cyan"
        }
        "5" {
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
            go build -o $outputPath $target.Path
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
