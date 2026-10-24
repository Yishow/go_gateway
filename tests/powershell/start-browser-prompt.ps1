# Browser consent defaults to no; only open the Vite URL after readiness.
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$env:GATEWAY_START_PS1_LIBRARY_ONLY = "1"
. (Join-Path $repoRoot "start.ps1")
Remove-Item Env:GATEWAY_START_PS1_LIBRARY_ONLY
$script:FrontendDevPort = "5175"
$script:Opened = @()
function Read-Host { param([string]$Prompt) return $script:Answer }
function Start-Process { param([string]$FilePath) $script:Opened += $FilePath }

foreach ($answer in @("", "n", "N", "y", "Y", "yes")) {
    $script:Answer = $answer
    $script:OpenBrowser = $true
    $script:Opened = @()
    Request-OpenBrowser
    Open-DevBrowser
    $expected = $answer -match '^(?i:y|yes)$'
    if ($script:OpenBrowser -ne $expected) { throw "Unexpected browser choice for '$answer'" }
    if ($expected) {
        if ($script:Opened.Count -ne 1 -or $script:Opened[0] -ne "http://127.0.0.1:5175/studio/v2") { throw "Expected Vite URL" }
    }
    elseif ($script:Opened.Count -ne 0) { throw "Browser must remain closed by default" }
}
function Read-Host { throw "No interactive input" }
Request-OpenBrowser
if ($script:OpenBrowser) { throw "Unavailable input must default to no" }

# Backend reloads must never inherit auto-open, even when externally enabled.
$savedAutoOpen = Get-ProcessEnvironmentValue -Name "AUTO_OPEN_BROWSER"
function Get-HostGoSetting { param([string]$Name) if ($Name -eq "GOHOSTOS") { return "windows" }; return "amd64" }
function New-ProcessLogPath { param([string]$Name) return "mock-$Name.log" }
function Start-Process {
    $script:ObservedAutoOpen = $env:AUTO_OPEN_BROWSER
    return [PSCustomObject]@{ Id = 202 }
}
try {
    $env:AUTO_OPEN_BROWSER = "true"
    $script:AirPath = "mock-air.exe"
    if (-not (Start-BackendProcess)) { throw "Backend start failed" }
    if ($script:ObservedAutoOpen -ne "false") { throw "Backend must not open a browser" }
    if ($env:AUTO_OPEN_BROWSER -ne "true") { throw "Parent environment must be preserved" }
}
finally { Set-ProcessEnvironmentValue -Name "AUTO_OPEN_BROWSER" -Value $savedAutoOpen }

# Exercise the dev orchestration; no browser before frontend readiness or on failure.
$script:Order = @()
function Resolve-PortConfiguration { param([bool]$ExplicitPort) return $true }
function Request-OpenBrowser { $script:Order += "prompt"; $script:OpenBrowser = $true }
function Acquire-RunLock { return $true }
function Sync-EmbeddedFrontendIfRequested { return $true }
function Clear-PortForService { param([int]$Port) return $true }
function Start-BackendProcess { return $true }
function Wait-PortReady { return $true }
function Start-FrontendDevServer {
    $script:Order += "frontend-ready"
    if ($script:FailFrontend) { return $null }
    return [PSCustomObject]@{ Id = 202 }
}
function Start-Process { param([string]$FilePath) $script:Order += "browser:$FilePath" }
function Wait-BackendExit { $script:Order += "monitor"; return 0 }
function Cleanup-All { $script:Order += "cleanup" }
$script:FailFrontend = $false
Start-DevMode
$expected = "prompt>frontend-ready>browser:http://127.0.0.1:5175/studio/v2>monitor>cleanup"
if (($script:Order -join ">") -ne $expected) { throw "Incorrect browser readiness order" }
$script:FailFrontend = $true
$script:Order = @()
Start-DevMode
if (($script:Order -join ">") -ne "prompt>frontend-ready>cleanup") { throw "Failed startup must not open browser" }
Write-Host "PASS: browser prompt defaults and readiness"
