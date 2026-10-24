# Occupied ports are refused; the launcher must not kill unrelated processes.
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
. (Join-Path $repoRoot "scripts/start-port-utils.ps1")

$errors = @()
function Write-Error { param([string]$Message) $script:errors += $Message }
function Write-Info { param([string]$Message) }
function Test-PortInUse { param([int]$Port) return $true }
function Get-ProcessByPort { param([int]$Port) return [PSCustomObject]@{ ProcessName = "unrelated"; Id = 45678 } }
function Stop-Process { throw "Clear-PortForService must not stop an unrelated process" }

$exitCode = 0
try {
    if (Clear-PortForService -Port 4545 -AutoKill) { throw "occupied port should be refused" }
    if (($errors -join "`n") -notmatch "4545|無關進程|手動釋放") { throw "refusal should explain the occupied port" }
    Write-Host "PASS: start-port-ownership.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
exit $exitCode
