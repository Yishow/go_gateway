$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evidenceRoot = Join-Path $repoRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence"
$reviewPath = Join-Path $evidenceRoot "process-safety-review.json"

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) {
        throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')"
    }
}

function Read-Json([string]$Path) {
    Assert-True (Test-Path -LiteralPath $Path -PathType Leaf) "process-safety review exists"
    return Get-Content -Raw -LiteralPath $Path | ConvertFrom-Json
}

function New-OwnedSleepProcess {
    $startInfo = [Diagnostics.ProcessStartInfo]::new()
    $pwsh = Get-Command pwsh -CommandType Application
    $startInfo.FileName = if ($pwsh.Source) { $pwsh.Source } else { $pwsh.Path }
    [void]$startInfo.ArgumentList.Add("-NoProfile")
    [void]$startInfo.ArgumentList.Add("-Command")
    [void]$startInfo.ArgumentList.Add("Start-Sleep -Seconds 3")
    $startInfo.UseShellExecute = $false
    $startInfo.CreateNoWindow = $true
    $process = [Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    [void]$process.Start()
    return $process
}

function Invoke-OwnedProcessProbe {
    $started = [DateTime]::UtcNow
    $a = $null
    $b = $null
    $aExited = $false
    $bAliveAfterAKill = $false
    $aCleanup = "not-started"
    $bCleanup = "not-started"
    try {
        $a = New-OwnedSleepProcess
        $b = New-OwnedSleepProcess
        Start-Sleep -Milliseconds 150
        if (-not $a.HasExited) { $a.Kill($true); $aCleanup = "owned-tree-kill-requested" }
        $aExited = $a.WaitForExit(2000)
        $bAliveAfterAKill = -not $b.HasExited
    }
    finally {
        if ($a -and -not $a.HasExited) { $a.Kill($true); $aCleanup = "owned-tree-kill-requested" }
        if ($b -and -not $b.HasExited) { $b.Kill($true); $bCleanup = "owned-tree-kill-requested" }
        if ($a) { [void]$a.WaitForExit(2000) }
        if ($b) { [void]$b.WaitForExit(2000) }
    }
    [ordered]@{
        durationMs = [long]([DateTime]::UtcNow - $started).TotalMilliseconds
        aExitedAfterOwnedTreeKill = $aExited
        bAliveAfterAKill = $bAliveAfterAKill
        aCleanupOutcome = $aCleanup
        bCleanupOutcome = $bCleanup
        pidsRecorded = ($null -ne $a -and $null -ne $b)
    }
}

function Invoke-TempOwnershipProbe {
    $root = Join-Path ([IO.Path]::GetTempPath()) ("frontend-vitest-process-safety-" + [Guid]::NewGuid().ToString("N"))
    $owned = Join-Path $root "owned"
    $external = Join-Path $root "external-sentinel.txt"
    New-Item -ItemType Directory -Path $owned -Force | Out-Null
    [IO.File]::WriteAllText($external, "sentinel")
    $preserved = $false
    try {
        [IO.Directory]::Delete($owned, $true)
        $preserved = [IO.File]::Exists($external)
    }
    finally {
        if ([IO.File]::Exists($external)) { [IO.File]::Delete($external) }
        if ([IO.Directory]::Exists($root)) { [IO.Directory]::Delete($root, $true) }
    }
    [ordered]@{ ownedDirectoryRemoved = -not [IO.Directory]::Exists($owned); externalSentinelPreserved = $preserved; externalSentinelRemoved = -not [IO.File]::Exists($external) }
}

try {
    $review = Read-Json $reviewPath
    Assert-Equal 1 $review.schemaVersion "process-safety schema version"
    Assert-Equal "ProcessSafetyReview" $review.mode "process-safety mode"
    Assert-Equal "pass/no-code" $review.result "process-safety result"

    $sourcePaths = @(
        "scripts/diagnose-frontend-vitest-matrix.ps1",
        "scripts/lib/FrontendVitestMatrix.psm1",
        "scripts/lib/FrontendVitestFreshIsolated.psm1",
        "scripts/lib/FrontendVitestGroups.psm1"
    )
    foreach ($relative in $sourcePaths) {
        $path = Join-Path $repoRoot ($relative -replace '/', '\')
        Assert-True (Test-Path -LiteralPath $path -PathType Leaf) "static review source exists: $relative"
        $source = Get-Content -Raw -LiteralPath $path
        Assert-True ($source -notmatch "Invoke-Expression|taskkill(?:\.exe)?|Stop-Process") "no dangerous process API: $relative"
    }
    Assert-True $review.staticReview.fixedCommandsOnly "fixed command contract"
    Assert-True $review.staticReview.watchdogAndTestTimeoutSeparated "watchdog/test timeout separated"
    Assert-True $review.staticReview.uniqueOutputOwnershipAndPathValidation "output ownership/path validation linked"

    $actualProcess = Invoke-OwnedProcessProbe
    Assert-True ($actualProcess.durationMs -lt 10000) "owned process probe is bounded"
    Assert-True $actualProcess.aExitedAfterOwnedTreeKill "owned A exits after tree kill"
    Assert-True $actualProcess.bAliveAfterAKill "unrelated B sentinel survives A kill"
    Assert-True ($actualProcess.aCleanupOutcome -eq "owned-tree-kill-requested" -and $actualProcess.bCleanupOutcome -eq "owned-tree-kill-requested") "both test-owned processes cleaned"
    Assert-True $review.actualProcessProbe.aExitedAfterOwnedTreeKill "recorded A cleanup"
    Assert-True $review.actualProcessProbe.bAliveAfterAKill "recorded B sentinel survival"
    Assert-True $review.actualProcessProbe.killEntireProcessTree "recorded tree-scoped kill"
    Assert-True ($review.actualProcessProbe.commandLineRecorded -eq $false -and $review.actualProcessProbe.environmentRecorded -eq $false) "no command line or environment recorded"

    $actualTemp = Invoke-TempOwnershipProbe
    Assert-True $actualTemp.ownedDirectoryRemoved "owned temp directory removed"
    Assert-True $actualTemp.externalSentinelPreserved "external sentinel untouched"
    Assert-True $actualTemp.externalSentinelRemoved "test-owned sentinel cleaned"
    Assert-True $review.tempOwnershipProbe.externalSentinelPreserved "recorded external sentinel preservation"

    Assert-Equal "refuse/no-kill" $review.unsafeOwnershipScenario.action "unsafe ownership refuses kill"
    Assert-Equal "unknown" $review.unsafeOwnershipScenario.classification "unsafe ownership is unknown"
    Assert-Equal "blocked" $review.unsafeOwnershipScenario.result "unsafe ownership blocks"
    Assert-True (-not $review.unsafeOwnershipScenario.killAttempted) "unsafe ownership never kills"

    Assert-True $review.rollbackRehearsal.dryRun "rollback is dry-run only"
    Assert-True (-not $review.rollbackRehearsal.deletedAnyPath -and @($review.rollbackRehearsal.plannedPaths).Count -gt 0) "rollback lists paths without deletion"
    Assert-True $review.linkedEvidence.rawEvidenceRetained "raw evidence retained"
    Assert-True $review.linkedEvidence.normalRunsWatchdogFalse "normal watchdog evidence remains false"

    Write-Output ("GREEN process-safety checks passed; probeMs={0}" -f $actualProcess.durationMs)
    exit 0
}
catch {
    Write-Error $_
    exit 1
}
