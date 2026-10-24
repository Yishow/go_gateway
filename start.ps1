# Go Gateway launcher: -Dev (Air + Vite) or -BuildSingle (Node builder).
param(
    [Alias("DevMode", "AirMode")][switch]$Dev,
    [switch]$BuildSingle,
    [ValidateSet("windows", "mac", "darwin", "linux")][string]$OS = "",
    [ValidateSet("amd64", "arm64")][string]$Arch = "",
    [int]$Port = 0,
    [Alias("?", "h")][switch]$Help
)

$ErrorActionPreference = "Stop"
$script:ExitCode = 0
$script:ROOT_DIR = [System.IO.Path]::GetFullPath($PSScriptRoot)
$script:APP_PATH = "cmd/test_ui"
$script:FRONTEND_DIR = "frontend"
$script:BUILD_DIR = "bin"
$script:TMP_DIR = Join-Path $script:ROOT_DIR "bin/tmp"
$script:TOOLS_DIR = Join-Path $script:ROOT_DIR "bin/tools"
$script:FRONTEND_DEV_HOST = "0.0.0.0"
$script:FrontendDevPort = "5173"
$script:ModbusSharePort = "5020"
$script:BackendProcess = $null
$script:FrontendProcess = $null
$script:BackendLogFile = $null
$script:BackendErrorLogFile = $null
$script:FrontendLogFile = $null
$script:FrontendErrorLogFile = $null
$script:LogPositions = @{}
$script:VitePath = $null
$script:AirPath = $null
$script:OpenBrowser = $false
$script:RunLockDir = $null
$script:LibraryOnly = ($env:GATEWAY_START_PS1_LIBRARY_ONLY -eq "1")
. (Join-Path $script:ROOT_DIR "scripts/start-process-utils.ps1")
. (Join-Path $script:ROOT_DIR "scripts/start-port-utils.ps1")

function Write-Info { param([string]$Message) Write-Host "[INFO] $Message" -ForegroundColor Cyan }
function Write-Success { param([string]$Message) Write-Host "[OK] $Message" -ForegroundColor Green }
function Write-Warning { param([string]$Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Error { param([string]$Message) Write-Host "[ERR] $Message" -ForegroundColor Red; $script:ExitCode = 1 }
function Test-Command { param([Parameter(Mandatory = $true)][string]$Name) return $null -ne (Get-Command $Name -ErrorAction SilentlyContinue) }

function Get-ProcessEnvironmentValue {
    param([Parameter(Mandatory = $true)][string]$Name)
    $item = Get-Item -LiteralPath ("Env:{0}" -f $Name) -ErrorAction SilentlyContinue
    if ($null -eq $item) { return $null }
    return [string]$item.Value
}
function Set-ProcessEnvironmentValue {
    param([Parameter(Mandatory = $true)][string]$Name, [AllowNull()][object]$Value)
    if ($null -eq $Value) { Remove-Item -LiteralPath ("Env:{0}" -f $Name) -ErrorAction SilentlyContinue }
    else { Set-Item -LiteralPath ("Env:{0}" -f $Name) -Value ([string]$Value) }
}
function Import-DotEnvFile {
    param([Parameter(Mandatory = $true)][string]$Path)
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { return }
    foreach ($rawLine in (Get-Content -LiteralPath $Path)) {
        $line = [string]$rawLine
        if ($line -notmatch '^\s*(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*?)\s*$') { continue }
        $name = $Matches[1]; $value = $Matches[2]
        if ($null -ne (Get-ProcessEnvironmentValue -Name $name)) { continue }
        if ($value.Length -ge 2) {
            $first = $value.Substring(0, 1); $last = $value.Substring($value.Length - 1, 1)
            if (($first -eq '"' -and $last -eq '"') -or ($first -eq "'" -and $last -eq "'")) { $value = $value.Substring(1, $value.Length - 2) }
        }
        Set-ProcessEnvironmentValue -Name $name -Value $value
    }
}

function Test-IsWindows { return [Environment]::OSVersion.Platform -eq [PlatformID]::Win32NT }
function Request-OpenBrowser {
    $script:OpenBrowser = $false
    try { $answer = Read-Host "是否開啟瀏覽器？[y/N]" }
    catch { return }
    if ($answer -match '^(?i:y|yes)$') { $script:OpenBrowser = $true }
}
function Open-DevBrowser {
    if (-not $script:OpenBrowser) { return }
    $url = "http://127.0.0.1:$($script:FrontendDevPort)/studio/v2"
    try { Start-Process -FilePath $url }
    catch { Write-Warning "瀏覽器開啟失敗，請手動開啟：$url" }
}
function Convert-ToPort {
    param([Parameter(Mandatory = $true)][string]$Name, [AllowNull()][object]$Value)
    $text = if ($null -eq $Value) { "" } else { ([string]$Value).Trim() }
    if ($text -notmatch '^\d+$') { throw "$Name 必須是 1 到 65535 的整數" }
    $number = [int]$text
    if ($number -lt 1 -or $number -gt 65535) { throw "$Name 必須是 1 到 65535 的整數" }
    return $number
}
function Resolve-PortConfiguration {
    param([bool]$ExplicitPort)
    try {
        if ($ExplicitPort) { $script:Port = Convert-ToPort -Name "-Port" -Value $Port }
        else {
            $envPort = Get-ProcessEnvironmentValue -Name "PORT"
            $script:Port = if ([string]::IsNullOrWhiteSpace($envPort)) { 8080 } else { Convert-ToPort -Name "PORT" -Value $envPort }
        }
        $frontendPort = Convert-ToPort -Name "VITE_DEV_PORT/FRONTEND_DEV_PORT" -Value $script:FrontendDevPort
        $modbusPort = Convert-ToPort -Name "MODBUS_SHARE_PORT" -Value $script:ModbusSharePort
        $script:FrontendDevPort = [string]$frontendPort; $script:ModbusSharePort = [string]$modbusPort
        if ($script:Port -eq $frontendPort) { throw "後端 PORT 與前端 VITE_DEV_PORT 必須不同" }
        if ($script:Port -eq $modbusPort -or $frontendPort -eq $modbusPort) { throw "後端、前端與 MODBUS_SHARE_PORT 必須使用不同端口" }
        return $true
    }
    catch { Write-Error $_.Exception.Message; return $false }
}
function Get-HostGoSetting {
    param([Parameter(Mandatory = $true)][string]$Name)
    if (-not (Test-Command -Name "go")) { return $null }
    $value = & go env $Name 2>$null | Select-Object -First 1
    if ($LASTEXITCODE -ne 0 -or $null -eq $value) { return $null }
    return ([string]$value).Trim()
}
function Invoke-WithEnvironment {
    param([Parameter(Mandatory = $true)][hashtable]$Values, [Parameter(Mandatory = $true)][scriptblock]$Action)
    $saved = @{}
    foreach ($name in $Values.Keys) { $saved[$name] = Get-ProcessEnvironmentValue -Name $name }
    try {
        foreach ($name in $Values.Keys) { Set-ProcessEnvironmentValue -Name $name -Value ([string]$Values[$name]) }
        & $Action
    }
    finally { foreach ($name in $Values.Keys) { Set-ProcessEnvironmentValue -Name $name -Value $saved[$name] } }
}

function Install-LocalAir {
    $airPath = Join-Path $script:TOOLS_DIR "air.exe"
    if (Test-Path -LiteralPath $airPath -PathType Leaf) { return $airPath }
    if (-not (Test-Command -Name "go")) { Write-Error "找不到 Go，無法安裝本地 Air"; return $null }
    New-Item -ItemType Directory -Path $script:TOOLS_DIR -Force | Out-Null
    $hostOS = Get-HostGoSetting -Name "GOHOSTOS"; $hostArch = Get-HostGoSetting -Name "GOHOSTARCH"
    if ([string]::IsNullOrWhiteSpace($hostOS) -or [string]::IsNullOrWhiteSpace($hostArch)) { Write-Error "無法取得 Go host 平台，無法安裝 Air"; return $null }
    try {
        Write-Info "安裝固定版本 Air v1.63.0 到 $script:TOOLS_DIR"
        Invoke-WithEnvironment -Values @{ GOBIN = $script:TOOLS_DIR; GOOS = $hostOS; GOARCH = $hostArch; CGO_ENABLED = "0"; GOTOOLCHAIN = "local" } -Action {
            & go install "github.com/air-verse/air@v1.63.0" | Out-Host
            if ($LASTEXITCODE -ne 0) { throw "Air 安裝失敗（退出碼 $LASTEXITCODE）" }
        }
    }
    catch { Write-Error $_.Exception.Message; return $null }
    if (-not (Test-Path -LiteralPath $airPath -PathType Leaf)) { Write-Error "Air 安裝完成但找不到 $airPath"; return $null }
    return $airPath
}
function New-ProcessLogPath {
    param([Parameter(Mandatory = $true)][string]$Name)
    New-Item -ItemType Directory -Path $script:TMP_DIR -Force | Out-Null
    return Join-Path $script:TMP_DIR ("{0}-{1}.log" -f $Name, $PID)
}
function Start-BackendProcess {
    $script:BackendLogFile = New-ProcessLogPath -Name "backend"; $script:BackendErrorLogFile = New-ProcessLogPath -Name "backend-error"
    $startArgs = @{ FilePath = $script:AirPath; WorkingDirectory = $script:ROOT_DIR; PassThru = $true; RedirectStandardOutput = $script:BackendLogFile; RedirectStandardError = $script:BackendErrorLogFile; ArgumentList = @("-c", ".air.toml") }
    if (Test-IsWindows) { $startArgs["WindowStyle"] = "Hidden" }
    $hostOS = Get-HostGoSetting -Name "GOHOSTOS"; $hostArch = Get-HostGoSetting -Name "GOHOSTARCH"
    if ([string]::IsNullOrWhiteSpace($hostOS) -or [string]::IsNullOrWhiteSpace($hostArch)) { Write-Error "無法取得 Go host 平台，無法啟動 Air"; return $false }
    try {
        $script:BackendProcess = Invoke-WithEnvironment -Values @{ PORT = [string]$script:Port; GOOS = $hostOS; GOARCH = $hostArch; CGO_ENABLED = "0"; AUTO_OPEN_BROWSER = "false" } -Action { Start-Process @startArgs }
        return $null -ne $script:BackendProcess
    }
    catch { Write-Error "啟動 Air 失敗：$($_.Exception.Message)"; return $false }
}

function Ensure-FrontendVite {
    $frontendRoot = Join-Path $script:ROOT_DIR $script:FRONTEND_DIR
    $packagePath = Join-Path $frontendRoot "package.json"; $vitePath = Join-Path $frontendRoot "node_modules/vite/bin/vite.js"
    if (-not (Test-Path -LiteralPath $packagePath -PathType Leaf)) { Write-Error "找不到前端 package.json：$packagePath"; return $null }
    if (-not (Test-Path -LiteralPath $vitePath -PathType Leaf)) {
        if (-not (Test-Command -Name "npm")) { Write-Error "找不到 npm，無法執行 npm ci"; return $null }
        Write-Info "找不到本地 Vite，執行 npm ci"
        Push-Location $frontendRoot
        try {
            & npm ci | Out-Host
            if ($LASTEXITCODE -ne 0) { Write-Error "npm ci 失敗（退出碼 $LASTEXITCODE）"; return $null }
        }
        catch { Write-Error "npm ci 失敗：$($_.Exception.Message)"; return $null }
        finally { Pop-Location }
    }
    if (-not (Test-Path -LiteralPath $vitePath -PathType Leaf)) { Write-Error "npm ci 後仍找不到本地 Vite CLI：$vitePath"; return $null }
    return $vitePath
}
function Start-FrontendDevServer {
    $vitePath = if ($script:VitePath) { $script:VitePath } else { Ensure-FrontendVite }
    if ($null -eq $vitePath) { return $null }
    if (-not (Clear-PortForService -Port ([int]$script:FrontendDevPort))) { return $null }
    $node = (Get-Command node -ErrorAction SilentlyContinue).Source; if ([string]::IsNullOrWhiteSpace($node)) { $node = "node" }
    $proxy = Get-ProcessEnvironmentValue -Name "VITE_API_PROXY_TARGET"; if ([string]::IsNullOrWhiteSpace($proxy)) { $proxy = "http://127.0.0.1:$script:Port" }
    $script:FrontendLogFile = New-ProcessLogPath -Name "frontend"; $script:FrontendErrorLogFile = New-ProcessLogPath -Name "frontend-error"
    $startArgs = @{ FilePath = $node; ArgumentList = @("node_modules/vite/bin/vite.js", "--host", $script:FRONTEND_DEV_HOST, "--port", [string]$script:FrontendDevPort, "--strictPort"); WorkingDirectory = (Join-Path $script:ROOT_DIR $script:FRONTEND_DIR); PassThru = $true; RedirectStandardOutput = $script:FrontendLogFile; RedirectStandardError = $script:FrontendErrorLogFile }
    if (Test-IsWindows) { $startArgs["WindowStyle"] = "Hidden" }
    try {
        $script:FrontendProcess = Invoke-WithEnvironment -Values @{ PORT = [string]$script:Port; VITE_API_PROXY_TARGET = $proxy; VITE_DEV_PORT = [string]$script:FrontendDevPort; BROWSERSLIST_IGNORE_OLD_DATA = "true" } -Action { Start-Process @startArgs }
        if ($null -eq $script:FrontendProcess) { throw "無法建立前端進程" }
        if (-not (Wait-PortReady -TargetPort ([int]$script:FrontendDevPort) -Label "前端端口" -Attempts 300 -Process $script:FrontendProcess)) { return $null }
        Write-Success "Vite 已啟動（port $script:FrontendDevPort）"; return $script:FrontendProcess
    }
    catch { Write-Error "啟動 Vite 失敗：$($_.Exception.Message)"; return $null }
}
function Write-ProcessLogs {
    foreach ($path in @($script:BackendLogFile, $script:BackendErrorLogFile, $script:FrontendLogFile, $script:FrontendErrorLogFile)) {
        if ([string]::IsNullOrWhiteSpace($path) -or -not (Test-Path -LiteralPath $path)) { continue }
        $stream = $null; $reader = $null
        try {
            $offset = if ($script:LogPositions.ContainsKey($path)) { [int64]$script:LogPositions[$path] } else { [int64]0 }
            $stream = [System.IO.File]::Open($path, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, [System.IO.FileShare]::ReadWrite)
            if ($offset -gt $stream.Length) { $offset = 0 }
            if ($offset -eq 0 -and $stream.Length -ge 3) {
                $preamble = New-Object byte[] 3
                [void]$stream.Read($preamble, 0, 3)
                if ($preamble[0] -eq 0xEF -and $preamble[1] -eq 0xBB -and $preamble[2] -eq 0xBF) { $offset = 3 }
            }
            [void]$stream.Seek($offset, [System.IO.SeekOrigin]::Begin)
            $reader = New-Object System.IO.StreamReader -ArgumentList @($stream, [System.Text.Encoding]::UTF8, $true, 4096, $true)
            $chunk = $reader.ReadToEnd(); $cursor = 0
            foreach ($match in [regex]::Matches($chunk, '(?s).*?(?:\r\n|\n|\r)')) {
                Write-Host ($match.Value -replace '(\r\n|\n|\r)$', '')
                $cursor += $match.Length
            }
            $script:LogPositions[$path] = $offset
            if ($cursor -gt 0) { $script:LogPositions[$path] = $offset + $reader.CurrentEncoding.GetByteCount($chunk.Substring(0, $cursor)) }
        }
        catch { }
        finally {
            if ($reader) { $reader.Dispose() }
            if ($stream) { $stream.Dispose() }
        }
    }
}
function Sync-EmbeddedFrontendIfRequested {
    $index = Join-Path $script:ROOT_DIR "cmd/test_ui/static/index.html"
    if ((Test-Path -LiteralPath $index -PathType Leaf) -and (Get-Item -LiteralPath $index).Length -gt 0) { Write-Info "使用既有 embedded 前端快照，Vite 負責即時前端"; return $true }
    if (-not (Test-Command -Name "node")) { Write-Error "找不到 Node.js，無法準備 embedded 前端快照"; return $false }
    Write-Info "找不到 embedded 前端快照，執行前端準備建置"
    Push-Location $script:ROOT_DIR
    try {
        & node "scripts/build-single.mjs" "--prepare-dev" | Out-Host
        $code = [int]$LASTEXITCODE
        if ($code -ne 0) { Write-Error "embedded 前端準備建置失敗（退出碼 $code）"; return $false }
    }
    catch { Write-Error "embedded 前端準備建置失敗：$($_.Exception.Message)"; return $false }
    finally { Pop-Location }
    if (-not (Test-Path -LiteralPath $index -PathType Leaf) -or (Get-Item -LiteralPath $index).Length -eq 0) { Write-Error "前端準備建置完成但找不到有效的 embedded index.html"; return $false }
    return $true
}
function Wait-BackendExit {
    while ($true) {
        Write-ProcessLogs
        if ($script:FrontendProcess) {
            try {
                if ($script:FrontendProcess.HasExited) {
                    Write-ProcessLogs; $code = [int]$script:FrontendProcess.ExitCode
                    if ($code -eq 0) { Write-Error "Vite 進程意外結束"; return 1 }
                    return $code
                }
            }
            catch { return 1 }
        }
        if ($script:BackendProcess) {
            try { if ($script:BackendProcess.HasExited) { Write-ProcessLogs; return [int]$script:BackendProcess.ExitCode } }
            catch { return 1 }
        }
        Start-Sleep -Milliseconds 200
    }
}
function Assert-DevPreflight {
    if (-not (Test-Command -Name "go") -or -not (Test-Command -Name "node") -or -not (Test-Command -Name "npm")) { Write-Error "Dev 需要 go、node 與 npm"; return $false }
    foreach ($port in @([int]$script:Port, [int]$script:FrontendDevPort)) {
        if (Test-PortInUse -Port $port) {
            $process = Get-ProcessByPort -Port $port; $owner = if ($process) { "$($process.ProcessName) (PID: $($process.Id))" } else { "未知進程" }
            Write-Error "端口 $port 已被 $owner 佔用；請先手動釋放端口"; return $false
        }
    }
    $script:VitePath = Ensure-FrontendVite; if ($null -eq $script:VitePath) { return $false }
    $script:AirPath = Install-LocalAir
    return $null -ne $script:AirPath
}
function Cleanup-All {
    try {
        if ($script:FrontendProcess) { Stop-ProcessTree -ProcId $script:FrontendProcess.Id -Label "前端開發伺服器"; $script:FrontendProcess = $null }
        if ($script:BackendProcess) { Stop-ProcessTree -ProcId $script:BackendProcess.Id -Label "Air 後端"; $script:BackendProcess = $null }
        foreach ($path in @($script:BackendLogFile, $script:BackendErrorLogFile, $script:FrontendLogFile, $script:FrontendErrorLogFile)) { if ($path -and (Test-Path -LiteralPath $path)) { Remove-Item -LiteralPath $path -Force -ErrorAction SilentlyContinue } }
        $script:BackendLogFile = $null; $script:BackendErrorLogFile = $null; $script:FrontendLogFile = $null; $script:FrontendErrorLogFile = $null; $script:LogPositions = @{}; $script:VitePath = $null
    }
    finally { Release-RunLock }
}
function Start-DevMode {
    if (-not (Test-IsWindows) -and -not $script:LibraryOnly) { Write-Error "PowerShell 熱重載模式目前只支援 Windows；請在 Windows 執行 -Dev"; return }
    if (-not (Resolve-PortConfiguration -ExplicitPort:$script:PortWasExplicit)) { return }
    Request-OpenBrowser
    if (-not (Acquire-RunLock)) { return }
    try {
        $syncResult = @(Sync-EmbeddedFrontendIfRequested); if ($syncResult.Count -gt 0 -and -not $syncResult[-1]) { return }
        if (-not $script:LibraryOnly -and -not (Assert-DevPreflight)) { return }
        if (-not (Clear-PortForService -Port $script:Port)) { return }
        $backendResult = @(Start-BackendProcess); if ($backendResult.Count -gt 0 -and -not $backendResult[-1]) { return }
        if (-not (Wait-PortReady -TargetPort $script:Port -Label "後端端口" -Attempts 1200 -Process $script:BackendProcess)) { return }
        $frontend = Start-FrontendDevServer; if ($null -eq $frontend) { return }
        Open-DevBrowser
        $exitCode = Wait-BackendExit
        if ($exitCode -ne 0) { Write-Error "開發進程已退出（退出碼 $exitCode）"; $script:ExitCode = $exitCode }
    }
    catch { Write-Error "開發模式失敗：$($_.Exception.Message)" }
    finally { Cleanup-All }
}

function Build-SingleBinary {
    param([string]$TargetOS = "", [string]$TargetArch = "")
    if (-not (Test-Command -Name "node")) { Write-Error "找不到 Node.js，無法建立單一執行檔"; return $false }
    $arguments = @("scripts/build-single.mjs")
    if (-not [string]::IsNullOrWhiteSpace($TargetOS)) { $arguments += @("--os", $TargetOS) }
    if (-not [string]::IsNullOrWhiteSpace($TargetArch)) { $arguments += @("--arch", $TargetArch) }
    Push-Location $script:ROOT_DIR
    try { & node @arguments | Out-Host; $code = [int]$LASTEXITCODE }
    catch { Write-Error "單一執行檔建置失敗：$($_.Exception.Message)"; return $false }
    finally { Pop-Location }
    if ($code -ne 0) { Write-Error "單一執行檔建置失敗（退出碼 $code）"; return $false }
    Write-Success "單一執行檔建置完成"; return $true
}
function Read-BuildTarget {
    Write-Host "建置目標：1 Windows  2 macOS  3 Linux"
    $osChoice = (Read-Host "選擇 OS").Trim().ToLowerInvariant()
    $osMap = @{ "1" = "windows"; "2" = "darwin"; "3" = "linux"; "windows" = "windows"; "mac" = "darwin"; "darwin" = "darwin"; "linux" = "linux" }
    if (-not $osMap.ContainsKey($osChoice)) { throw "無效的 OS 選擇：$osChoice" }
    Write-Host "建置架構：1 amd64  2 arm64"
    $archChoice = (Read-Host "選擇架構").Trim().ToLowerInvariant()
    $archMap = @{ "1" = "amd64"; "2" = "arm64"; "amd64" = "amd64"; "arm64" = "arm64" }
    if (-not $archMap.ContainsKey($archChoice)) { throw "無效的架構選擇：$archChoice" }
    return @{ OS = $osMap[$osChoice]; Arch = $archMap[$archChoice] }
}
function Show-Help {
    @"
用法：.\start.ps1 [-Dev] [-BuildSingle] [-OS windows|mac|darwin|linux] [-Arch amd64|arm64] [-Port 8080]

  -Dev (-DevMode/-AirMode)  啟動後端 Air 與前端 Vite 熱重載（Windows）
  -BuildSingle               建立包含前後端的單一執行檔
  -OS / -Arch                建置目標；省略時由共用 Node builder 使用主機值
  -Port                     後端端口；CLI 值優先於外部環境與 .env
  -Help                     顯示說明
"@ | Write-Host
}

Import-DotEnvFile -Path (Join-Path $script:ROOT_DIR ".env")
$script:PortWasExplicit = $PSBoundParameters.ContainsKey("Port")
$script:FRONTEND_DEV_HOST = Get-ProcessEnvironmentValue -Name "FRONTEND_DEV_HOST"; if ([string]::IsNullOrWhiteSpace($script:FRONTEND_DEV_HOST)) { $script:FRONTEND_DEV_HOST = "0.0.0.0" }
$frontendEnv = Get-ProcessEnvironmentValue -Name "FRONTEND_DEV_PORT"; if ([string]::IsNullOrWhiteSpace($frontendEnv)) { $frontendEnv = Get-ProcessEnvironmentValue -Name "VITE_DEV_PORT" }; if (-not [string]::IsNullOrWhiteSpace($frontendEnv)) { $script:FrontendDevPort = $frontendEnv }
$modbusEnv = Get-ProcessEnvironmentValue -Name "MODBUS_SHARE_PORT"; if (-not [string]::IsNullOrWhiteSpace($modbusEnv)) { $script:ModbusSharePort = $modbusEnv }

if ($script:LibraryOnly) { return }
if ($Help) { Show-Help; exit 0 }
if ($Dev -and $BuildSingle) { Write-Error "-Dev 與 -BuildSingle 只能擇一"; exit 1 }
if ($Dev) {
    if (-not [string]::IsNullOrWhiteSpace($OS) -or -not [string]::IsNullOrWhiteSpace($Arch)) { Write-Error "-OS 與 -Arch 不可搭配 -Dev"; exit 1 }
    Start-DevMode; exit $script:ExitCode
}
if ($BuildSingle) { if (Build-SingleBinary -TargetOS $OS -TargetArch $Arch) { exit 0 }; exit 1 }
if (-not [string]::IsNullOrWhiteSpace($OS) -or -not [string]::IsNullOrWhiteSpace($Arch)) { Write-Error "-OS 與 -Arch 必須搭配 -BuildSingle"; exit 1 }

while ($true) {
    Write-Host "`nGo Gateway 啟動選單" -ForegroundColor Cyan
    Write-Host "1. 開發模式（Air + Vite 熱重載）"; Write-Host "2. 建置單一執行檔（前後端）"; Write-Host "0. 離開"
    switch ((Read-Host "選擇").Trim().ToLowerInvariant()) {
        "1" { Start-DevMode; exit $script:ExitCode }
        "2" { try { $target = Read-BuildTarget; if (Build-SingleBinary -TargetOS $target.OS -TargetArch $target.Arch) { exit 0 }; exit 1 } catch { Write-Error $_.Exception.Message; exit 1 } }
        "0" { exit 0 }
        default { Write-Warning "請選擇 1、2 或 0" }
    }
}
