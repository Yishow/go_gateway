$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$modulePath = Join-Path $repoRoot "scripts\lib\FrontendVitestEvidence.psm1"
$roots = [System.Collections.Generic.List[string]]::new()
Import-Module -Name $modulePath -Force -Global

function Assert-True {
    param([Parameter(Mandatory = $true)][bool]$Condition, [Parameter(Mandatory = $true)][string]$Message)
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}
function Assert-Throws {
    param([Parameter(Mandatory = $true)][scriptblock]$Action, [Parameter(Mandatory = $true)][string]$Message)
    $thrown = $false
    try { & $Action } catch { $thrown = $true }
    Assert-True $thrown $Message
}
function New-SyntheticRepo {
    $root = Join-Path ([IO.Path]::GetTempPath()) ("frontend-evidence-resolver-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path (Join-Path $root "openspec\changes\archive") -Force | Out-Null
    [void]$roots.Add($root)
    return $root
}
function Write-SyntheticEvidence {
    param([Parameter(Mandatory = $true)][string]$Root, [Parameter(Mandatory = $true)][string]$ChangeDirectory, [Parameter(Mandatory = $true)][string]$Value = "ok")
    $path = Join-Path $Root ("openspec\changes\" + $ChangeDirectory + "\evidence\environment-capture.json")
    New-Item -ItemType Directory -Path (Split-Path -Parent $path) -Force | Out-Null
    Set-Content -LiteralPath $path -Value ([ordered]@{ source = $Value } | ConvertTo-Json) -Encoding utf8
    return $path
}
try {
    $activeRoot = New-SyntheticRepo
    $activePath = Write-SyntheticEvidence -Root $activeRoot -ChangeDirectory "diagnose-frontend-full-suite-timeouts" -Value "active"
    [void](Write-SyntheticEvidence -Root $activeRoot -ChangeDirectory "archive\2026-08-22-diagnose-frontend-full-suite-timeouts" -Value "old")
    Assert-True ((Resolve-FrontendVitestEvidencePath -RepoRoot $activeRoot -FileName "environment-capture.json") -eq [IO.Path]::GetFullPath($activePath)) "active evidence wins over archive"

    $missingRoot = New-SyntheticRepo
    New-Item -ItemType Directory -Path (Join-Path $missingRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence") -Force | Out-Null
    [void](Write-SyntheticEvidence -Root $missingRoot -ChangeDirectory "archive\2026-08-22-diagnose-frontend-full-suite-timeouts" -Value "archive")
    Assert-Throws { Resolve-FrontendVitestEvidencePath -RepoRoot $missingRoot -FileName "environment-capture.json" } "active missing evidence does not fall back to archive"

    $malformedRoot = New-SyntheticRepo
    $malformedPath = Join-Path $malformedRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence\environment-capture.json"
    New-Item -ItemType Directory -Path (Split-Path -Parent $malformedPath) -Force | Out-Null
    Set-Content -LiteralPath $malformedPath -Value "not-json" -Encoding utf8
    [void](Write-SyntheticEvidence -Root $malformedRoot -ChangeDirectory "archive\2026-08-22-diagnose-frontend-full-suite-timeouts" -Value "archive")
    Assert-Throws { Resolve-FrontendVitestEvidencePath -RepoRoot $malformedRoot -FileName "environment-capture.json" } "active malformed evidence does not fall back to archive"

    $archiveRoot = New-SyntheticRepo
    [void](Write-SyntheticEvidence -Root $archiveRoot -ChangeDirectory "archive\2026-08-21-diagnose-frontend-full-suite-timeouts" -Value "old")
    $latestPath = Write-SyntheticEvidence -Root $archiveRoot -ChangeDirectory "archive\2026-08-22-diagnose-frontend-full-suite-timeouts" -Value "latest"
    Assert-True ((Resolve-FrontendVitestEvidencePath -RepoRoot $archiveRoot -FileName "environment-capture.json") -eq [IO.Path]::GetFullPath($latestPath)) "latest unique archive is selected"

    $ambiguousRoot = New-SyntheticRepo
    [void](Write-SyntheticEvidence -Root $ambiguousRoot -ChangeDirectory "archive\2026-08-23-diagnose-frontend-full-suite-timeouts-a" -Value "a")
    [void](Write-SyntheticEvidence -Root $ambiguousRoot -ChangeDirectory "archive\2026-08-23-diagnose-frontend-full-suite-timeouts-b" -Value "b")
    Assert-Throws { Resolve-FrontendVitestEvidencePath -RepoRoot $ambiguousRoot -FileName "environment-capture.json" } "same-date archive ambiguity fails closed"

    Assert-Throws { Resolve-FrontendVitestEvidencePath -RepoRoot $archiveRoot -FileName "environment-capture.json" -ExplicitPath "relative.json" } "relative explicit evidence path is rejected"
    Write-Output "GREEN evidence resolver contract checks passed"
    exit 0
}
finally {
    foreach ($root in @($roots)) {
        if (Test-Path -LiteralPath $root -PathType Container) { Remove-Item -LiteralPath $root -Recurse -Force -ErrorAction SilentlyContinue }
    }
}
