# Wait-PortReady ownership contract: direct owner, descendant, foreign, and dead listeners.
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
. (Join-Path $repoRoot "scripts/start-process-utils.ps1")

$script:BackendProcess = $null
$script:listenerPid = 101
$script:parentByPid = @{}
$script:errors = @()

function Write-Error { param([string]$Message) $script:errors += $Message }
function Write-ProcessLogs { }
function Test-PortInUse { param([int]$Port) return $true }
function Get-ProcessByPort { param([int]$Port) return [PSCustomObject]@{ Id = $script:listenerPid } }
function Get-CimInstance {
    param([string]$ClassName, [string]$Filter, [string]$ErrorAction)
    if ($Filter -match 'ProcessId = (\d+)') {
        $processId = [int]$Matches[1]
        if ($script:parentByPid.ContainsKey($processId)) {
            return [PSCustomObject]@{ ParentProcessId = $script:parentByPid[$processId] }
        }
    }
    return $null
}

function New-TestProcess {
    param([int]$Id, [bool]$HasExited)
    return [PSCustomObject]@{ Id = $Id; HasExited = $HasExited }
}

function Assert-Ready {
    param(
        [Parameter(Mandatory = $true)][object]$Process,
        [Parameter(Mandatory = $true)][int]$Listener,
        [Parameter(Mandatory = $true)][bool]$Expected
    )
    $script:listenerPid = $Listener
    $result = Wait-PortReady -TargetPort 4545 -Attempts 1 -DelayMilliseconds 0 -Process $Process
    if ($result -ne $Expected) {
        throw "unexpected readiness result owner=$($Process.Id) listener=$Listener result=$result"
    }
}

Assert-Ready -Process (New-TestProcess 101 $false) -Listener 101 -Expected $true
$script:parentByPid = @{ 202 = 101 }
Assert-Ready -Process (New-TestProcess 101 $false) -Listener 202 -Expected $true
$script:parentByPid = @{ 303 = 999 }
Assert-Ready -Process (New-TestProcess 101 $false) -Listener 303 -Expected $false
Assert-Ready -Process (New-TestProcess 404 $true) -Listener 404 -Expected $false

# A helper caller without a process keeps the existing wait-only contract.
$script:BackendProcess = $null
$script:listenerPid = 303
if (-not (Wait-PortReady -TargetPort 4545 -Attempts 1 -DelayMilliseconds 0)) {
    throw "wait-only helper should accept a listening port without a process"
}

Write-Host "PASS: start-owned-readiness.ps1"
