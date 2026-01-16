# ============================================
# Go Gateway 一鍵啟動腳本
# ============================================
# 功能：
#   1. 執行 golangci-lint 靜態分析
#   2. 執行單元測試
#   3. 檢查測試覆蓋率
#   4. 構建所有可執行文件
# ============================================

param(
    [switch]$SkipLint,      # 跳過 lint 檢查
    [switch]$SkipTest,      # 跳過測試
    [switch]$SkipBuild,     # 跳過構建
    [switch]$Coverage,      # 顯示詳細覆蓋率
    [switch]$Verbose        # 詳細輸出
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

# 標題
Write-ColorOutput "`n============================================" "Cyan"
Write-ColorOutput "   Go Gateway 一鍵啟動腳本" "Cyan"
Write-ColorOutput "============================================`n" "Cyan"

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

# 2. 執行 golangci-lint 靜態分析
if (-not $SkipLint) {
    Write-ColorOutput "`n[1/4] 執行 golangci-lint 靜態分析..." "Yellow"
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
} else {
    Write-Info "[1/4] 跳過 golangci-lint 檢查"
}

# 3. 執行單元測試
if (-not $SkipTest) {
    Write-ColorOutput "`n[2/4] 執行單元測試..." "Yellow"
    
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
} else {
    Write-Info "[2/4] 跳過單元測試"
}

# 4. 構建可執行文件
if (-not $SkipBuild) {
    Write-ColorOutput "`n[3/4] 構建可執行文件..." "Yellow"
    
    $buildTargets = @(
        @{Name="fatek_test"; Path="./cmd/fatek_test"},
        @{Name="test_all"; Path="./cmd/test_all"}
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
} else {
    Write-Info "[3/4] 跳過構建"
}

# 5. 代碼質量檢查
Write-ColorOutput "`n[4/4] 代碼質量檢查..." "Yellow"

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
