# ============================================
# Fatek Gateway 超強一鍵構建腳本
# ============================================

param(
    [switch]$Clean = $false,
    [switch]$Test = $true,
    [switch]$Lint = $true,
    [switch]$Cross = $false,
    [string]$Output = "bin",
    [string]$Version = "",
    [switch]$Verbose = $false
)

# 設置錯誤處理
$ErrorActionPreference = "Stop"

# 顏色輸出函數
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    $colorMap = @{
        "Red" = [ConsoleColor]::Red
        "Green" = [ConsoleColor]::Green
        "Yellow" = [ConsoleColor]::Yellow
        "Blue" = [ConsoleColor]::Blue
        "Cyan" = [ConsoleColor]::Cyan
        "Magenta" = [ConsoleColor]::Magenta
    }
    $originalColor = $Host.UI.RawUI.ForegroundColor
    if ($colorMap.ContainsKey($Color)) {
        $Host.UI.RawUI.ForegroundColor = $colorMap[$Color]
    }
    Write-Host $Message
    $Host.UI.RawUI.ForegroundColor = $originalColor
}

function Write-Success { param([string]$msg) Write-ColorOutput "✓ $msg" "Green" }
function Write-Error { param([string]$msg) Write-ColorOutput "✗ $msg" "Red" }
function Write-Info { param([string]$msg) Write-ColorOutput "ℹ $msg" "Cyan" }
function Write-Warning { param([string]$msg) Write-ColorOutput "⚠ $msg" "Yellow" }
function Write-Step { param([string]$msg) Write-ColorOutput "▶ $msg" "Blue" }

# 顯示標題
Write-ColorOutput @"
╔═══════════════════════════════════════════════════════════╗
║         Fatek Gateway 超強一鍵構建腳本 v1.0              ║
╚═══════════════════════════════════════════════════════════╝
"@ "Cyan"

# 檢查 Go 環境
Write-Step "檢查 Go 環境..."
try {
    $goVersionOutput = & go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Go 環境正常: $goVersionOutput"
    } else {
        throw "Go command failed"
    }
} catch {
    Write-Error "未找到 Go 環境，請先安裝 Go"
    exit 1
}

# 獲取項目信息
$projectRoot = $PSScriptRoot
$moduleName = (Get-Content "$projectRoot\go.mod" | Select-String "^module\s+(\S+)" | ForEach-Object { $_.Matches.Groups[1].Value })
$goVersion = (Get-Content "$projectRoot\go.mod" | Select-String "^go\s+(\S+)" | ForEach-Object { $_.Matches.Groups[1].Value })

Write-Info "項目根目錄: $projectRoot"
Write-Info "模組名稱: $moduleName"
Write-Info "Go 版本要求: $goVersion"

# 自動檢測版本號
if ([string]::IsNullOrEmpty($Version)) {
    try {
        $gitTag = git describe --tags --exact-match 2>$null
        if ($LASTEXITCODE -eq 0) {
            $Version = $gitTag
        } else {
            $gitCommit = git rev-parse --short HEAD 2>$null
            if ($LASTEXITCODE -eq 0) {
                $Version = "dev-$gitCommit"
            } else {
                $Version = "dev-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
            }
        }
    } catch {
        $Version = "dev-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
    }
}
Write-Info "構建版本: $Version"

# 設置構建變數
$ldFlags = "-X main.Version=$Version -X main.BuildTime=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') -s -w"
$buildTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

# 清理舊構建
if ($Clean) {
    Write-Step "清理舊構建文件..."
    $cleanPaths = @("$Output", "*.exe", "*.test")
    foreach ($path in $cleanPaths) {
        if (Test-Path $path) {
            Remove-Item -Path $path -Recurse -Force -ErrorAction SilentlyContinue
            Write-Info "已清理: $path"
        }
    }
    Write-Success "清理完成"
}

# 創建輸出目錄
if (-not (Test-Path $Output)) {
    New-Item -ItemType Directory -Path $Output -Force | Out-Null
    Write-Info "創建輸出目錄: $Output"
}

# 運行測試
if ($Test) {
    Write-Step "運行測試..."
    try {
        $testOutput = go test ./... -v 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success "所有測試通過"
            if ($Verbose) {
                Write-Host $testOutput
            }
        } else {
            Write-Error "測試失敗"
            Write-Host $testOutput
            exit 1
        }
    } catch {
        Write-Error "測試執行失敗: $_"
        exit 1
    }
}

# 運行 Lint 檢查
if ($Lint) {
    Write-Step "運行 Lint 檢查..."
    try {
        # 檢查是否安裝了 golangci-lint
        $lintCmd = Get-Command golangci-lint -ErrorAction SilentlyContinue
        if ($lintCmd) {
            golangci-lint run ./...
            if ($LASTEXITCODE -eq 0) {
                Write-Success "Lint 檢查通過"
            } else {
                Write-Warning "Lint 檢查發現問題（繼續構建）"
            }
        } else {
            Write-Warning "未安裝 golangci-lint，跳過 lint 檢查"
            Write-Info "安裝命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        }
    } catch {
        Write-Warning "Lint 檢查失敗（繼續構建）: $_"
    }
}

# 構建 Windows 版本
Write-Step "構建 Windows 版本..."
$windowsBinary = "$Output\fatek_test.exe"
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    
    go build -ldflags $ldFlags -o $windowsBinary ./cmd/fatek_test
    if ($LASTEXITCODE -eq 0) {
        $fileInfo = Get-Item $windowsBinary
        $fileSize = [math]::Round($fileInfo.Length / 1KB, 2)
        Write-Success "Windows 版本構建成功: $windowsBinary ($fileSize KB)"
    } else {
        Write-Error "Windows 版本構建失敗"
        exit 1
    }
} catch {
    Write-Error "構建 Windows 版本時發生錯誤: $_"
    exit 1
} finally {
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue
}

# 構建多平台版本（可選）
if ($Cross) {
    Write-Step "構建多平台版本..."
    
    $platforms = @(
        @{OS="linux"; ARCH="amd64"; Ext=""},
        @{OS="linux"; ARCH="arm64"; Ext=""},
        @{OS="darwin"; ARCH="amd64"; Ext=""},
        @{OS="darwin"; ARCH="arm64"; Ext=""},
        @{OS="windows"; ARCH="amd64"; Ext=".exe"},
        @{OS="windows"; ARCH="386"; Ext=".exe"}
    )
    
    foreach ($platform in $platforms) {
        $binaryName = "fatek_test$($platform.Ext)"
        $outputPath = "$Output\fatek_test-$($platform.OS)-$($platform.ARCH)$($platform.Ext)"
        
        Write-Info "構建 $($platform.OS)/$($platform.ARCH)..."
        try {
            $env:GOOS = $platform.OS
            $env:GOARCH = $platform.ARCH
            $env:CGO_ENABLED = "0"
            
            go build -ldflags $ldFlags -o $outputPath ./cmd/fatek_test
            if ($LASTEXITCODE -eq 0) {
                $fileInfo = Get-Item $outputPath
                $fileSize = [math]::Round($fileInfo.Length / 1KB, 2)
                Write-Success "  ✓ $($platform.OS)/$($platform.ARCH): $outputPath ($fileSize KB)"
            } else {
                Write-Warning "  ✗ $($platform.OS)/$($platform.ARCH): 構建失敗"
            }
        } catch {
            Write-Warning "  ✗ $($platform.OS)/$($platform.ARCH): 構建錯誤: $_"
        } finally {
            Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
            Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
            Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue
        }
    }
}

# 顯示構建摘要
Write-ColorOutput @"
╔═══════════════════════════════════════════════════════════╗
║                     構建完成摘要                          ║
╚═══════════════════════════════════════════════════════════╝
"@ "Green"

Write-Info "構建時間: $buildTime"
Write-Info "版本號: $Version"
Write-Info "輸出目錄: $Output"

$builtFiles = Get-ChildItem -Path $Output -File | Where-Object { $_.Name -like "fatek_test*" }
if ($builtFiles) {
    Write-Info "`n構建的文件:"
    foreach ($file in $builtFiles) {
        $size = [math]::Round($file.Length / 1KB, 2)
        Write-Host "  • $($file.Name) ($size KB)"
    }
}

Write-Success "`n構建完成！"

# 顯示使用說明
Write-ColorOutput @"
`n使用說明:
  運行測試工具: .\$Output\fatek_test.exe -mode=tcp -host=`"127.0.0.1`" -action=read -symbol=D -addr=0 -count=10
  查看幫助: .\$Output\fatek_test.exe -help
"@ "Cyan"
