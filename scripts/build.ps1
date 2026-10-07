# 明確選擇 Windows desktop 或既有 console 開發產物。
param(
    [ValidateSet("console", "desktop")][string]$Mode = "console",
    [ValidatePattern('^[0-9A-Za-z._+-]+$')][string]$Version = "dev",
    [ValidatePattern('^[0-9A-Za-z._+-]*$')][string]$Commit = ""
)
$ErrorActionPreference = "Stop"

function Stop-Build {
    param(
        [Parameter(Mandatory = $true)][string]$Message,
        [int]$ExitCode = 1
    )

    Write-Error $Message -ErrorAction Continue
    exit $ExitCode
}

Set-Location -LiteralPath (Split-Path -Parent $PSScriptRoot)
$targetOS = & go env GOOS
if ($LASTEXITCODE -ne 0) { Stop-Build "Go toolchain is unavailable" }
if ($Mode -eq "desktop" -and $targetOS -ne "windows") {
    Stop-Build "desktop requires GOOS=windows"
}
if (-not $Commit) {
    $Commit = "unknown"
    if (Get-Command git -ErrorAction SilentlyContinue) {
        $gitCommit = & git rev-parse --short=12 HEAD 2>$null
        if ($LASTEXITCODE -eq 0) {
            $Commit = $gitCommit
            $gitChanges = & git status --porcelain --untracked-files=normal 2>$null
            if ($LASTEXITCODE -eq 0 -and $gitChanges) { $Commit += "+dirty" }
        }
    }
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

$indexPath = Join-Path $distPath "index.html"
if (-not (Test-Path -LiteralPath $indexPath -PathType Leaf) -or (Get-Item -LiteralPath $indexPath).Length -eq 0) {
    Stop-Build "frontend index.html is missing or empty; refusing incomplete embedded UI"
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
$output = "bin/test-ui.exe"
$ldflags = "-s -w -X go-gateway/internal/desktop.Version=$Version -X go-gateway/internal/desktop.Commit=$Commit"
$buildArgs = @()
if ($Mode -eq "desktop") {
    $output = "bin/gateway-desktop.exe"
    $buildArgs += @("-tags", "desktop")
    $ldflags += " " + "-H=windowsgui"
}
$buildArgs += @("-trimpath", "-ldflags", $ldflags, "-o", $output, "./cmd/test_ui")
& go build @buildArgs
$backendBuildExitCode = $LASTEXITCODE
if ($backendBuildExitCode -ne 0) {
    Stop-Build "backend build failed (go build exit code: $backendBuildExitCode)" $backendBuildExitCode
}

Write-Host "建置完成！" -ForegroundColor Green
