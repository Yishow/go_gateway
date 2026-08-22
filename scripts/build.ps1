# 建置腳本
$ErrorActionPreference = "Stop"

function Stop-Build {
    param(
        [Parameter(Mandatory = $true)][string]$Message,
        [int]$ExitCode = 1
    )

    Write-Error $Message -ErrorAction Continue
    exit $ExitCode
}

Write-Host "建置前端..." -ForegroundColor Green
Set-Location -LiteralPath frontend
npm run build
$frontendBuildExitCode = $LASTEXITCODE
Set-Location -LiteralPath ..
if ($frontendBuildExitCode -ne 0) {
    Stop-Build "frontend build failed (npm run build exit code: $frontendBuildExitCode)" $frontendBuildExitCode
}

$distPath = Join-Path (Get-Location) "frontend/dist"
$staticPath = Join-Path (Get-Location) "cmd/test_ui/static"
if (-not (Test-Path -LiteralPath $distPath -PathType Container)) {
    Stop-Build "frontend/dist is unavailable; embedded frontend synchronization cannot continue"
}

Write-Host "複製前端檔案到 embed 目錄..." -ForegroundColor Green
try {
    # 保留 git 追蹤的 embed 佔位檔（fresh clone 後 go:embed static 才能編譯）
    New-Item -ItemType Directory -Path $staticPath -Force | Out-Null
    if (-not (Test-Path -LiteralPath $staticPath -PathType Container)) {
        throw "target static directory was not created: $staticPath"
    }
    Get-ChildItem -LiteralPath $staticPath -Force |
        Where-Object { $_.Name -ne "embed-placeholder.txt" } |
        Remove-Item -Recurse -Force
    Copy-Item -Recurse -Force -Path (Join-Path $distPath "*") -Destination $staticPath
}
catch {
    Stop-Build "frontend asset copy/synchronization failed: $($_.Exception.Message)"
}

Write-Host "建置後端..." -ForegroundColor Green
# 移除 -H=windowsgui 標誌以顯示控制台窗口，讓用戶可以點擊 X 按鈕
# -s: 移除符號表，-w: 移除 DWARF 除錯資訊，-trimpath: 移除檔案路徑資訊
go build -ldflags "-s -w" -trimpath -o bin/test-ui.exe ./cmd/test_ui
$backendBuildExitCode = $LASTEXITCODE
if ($backendBuildExitCode -ne 0) {
    Stop-Build "backend build failed (go build exit code: $backendBuildExitCode)" $backendBuildExitCode
}

Write-Host "建置完成！" -ForegroundColor Green
