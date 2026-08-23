$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$handoffPath = Join-Path $repoRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence\repair-branch-handoff.json"
$reviewRelative = "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/classification-review.json"

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) { throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')" }
}

function Get-GitInventory {
    @(& rtk git ls-files --others --exclude-standard 2>$null | ForEach-Object { ([string]$_).Trim() } | Where-Object { $_ } | Sort-Object -Unique)
}

try {
    Assert-True (Test-Path -LiteralPath $handoffPath -PathType Leaf) "repair branch handoff exists"
    $handoff = Get-Content -Raw -LiteralPath $handoffPath | ConvertFrom-Json
    Assert-Equal 1 $handoff.schemaVersion "handoff schema version"
    Assert-Equal "RepairBranchHandoff" $handoff.mode "handoff mode"
    Assert-Equal "blocked/no-code/N/A" $handoff.result "handoff result"

    $reviewPath = Join-Path $repoRoot $reviewRelative
    Assert-True (Test-Path -LiteralPath $reviewPath -PathType Leaf) "classification review exists"
    Assert-Equal $reviewRelative $handoff.classificationReview.path "classification review path"
    Assert-Equal (Get-FileHash -Algorithm SHA256 -LiteralPath $reviewPath).Hash $handoff.classificationReview.sha256 "classification review hash"
    $review = Get-Content -Raw -LiteralPath $reviewPath | ConvertFrom-Json
    Assert-Equal "locked" $review.repairGate.status "linked repair gate is locked"
    Assert-True (-not $review.repairGate.sourceChangeAllowed -and -not $review.repairGate.testLocalRootCauseProven) "linked repair gate disallows source change"

    Assert-Equal "N/A/not-created" $handoff.taskDecisions.'5.1'.characterizationRed "5.1 does not create RED"
    Assert-Equal "N/A/no-source-change" $handoff.taskDecisions.'5.2'.greenRepair "5.2 has no source change"
    Assert-Equal "blocked/no-code" $handoff.taskDecisions.'5.3'.branch "5.3 closes blocked no-code"
    Assert-True ($handoff.modifiedImpactPaths.Count -eq 0) "no Impact path was modified"
    Assert-True (-not $handoff.workaroundProof.timeoutChanged -and -not $handoff.workaroundProof.workerChanged -and -not $handoff.workaroundProof.poolChanged) "no timeout/worker/pool workaround"
    Assert-True (-not $handoff.workaroundProof.configChanged -and -not $handoff.workaroundProof.dependencyChanged -and -not $handoff.workaroundProof.discoveryChanged) "no config/dependency/discovery workaround"
    Assert-True (-not $handoff.workaroundProof.skipRetryExcludeChanged -and -not $handoff.workaroundProof.coverageChanged) "no skip/retry/exclude/coverage workaround"
    Assert-True $handoff.workaroundProof.singleWorkerDiagnosisOnly "single worker remains diagnosis-only"

    $actualInventory = Get-GitInventory
    $handoffInventory = @($handoff.sourceDiffProof.actualUntrackedPaths | ForEach-Object { ([string]$_).Replace('\', '/') } | Sort-Object -Unique)
    Assert-Equal ($actualInventory -join "|") ($handoffInventory -join "|") "handoff matches actual untracked Git inventory"
    $allowed = @(
        "openspec/changes/diagnose-frontend-full-suite-timeouts/",
        "scripts/diagnose-frontend-vitest-matrix.ps1",
        "scripts/lib/FrontendVitest",
        "scripts/tests/diagnose-frontend-"
    )
    foreach ($path in $actualInventory) {
        $ok = $path.StartsWith($allowed[0], [StringComparison]::OrdinalIgnoreCase) -or $path -eq $allowed[1] -or $path.StartsWith($allowed[2], [StringComparison]::OrdinalIgnoreCase) -or $path.StartsWith($allowed[3], [StringComparison]::OrdinalIgnoreCase)
        Assert-True $ok "path is inside diagnosis/change allowlist: $path"
        Assert-True ($path -notmatch "^(frontend/|.*\.go$|.*package(-lock)?\.json$|.*\.config\.|\.line-limit-ignore$)") "path is outside forbidden scope: $path"
    }
    Assert-Equal 0 @($handoff.sourceDiffProof.forbiddenPathsFound).Count "no forbidden path is recorded"
    Assert-True $handoff.sourceDiffProof.scopePass "scope proof passes"

    Assert-True (-not $handoff.ingestDecision.required) "no spectra-ingest needed for no-edit branch"
    Assert-True ($handoff.ingestDecision.futureOutOfImpactRequiresIngest) "future out-of-Impact edits require ingest"
    Assert-True ($handoff.linkedEvidence.rawEvidenceRetained -and $handoff.linkedEvidence.noStabilityClaim) "raw evidence retained without stability claim"
    Assert-Equal "blocked/no-code/N/A" $handoff.branchClosure.result "branch closure result"
    Assert-True $handoff.branchClosure.noSourceRepair "branch closure records no source repair"
    Assert-True $handoff.branchClosure.tasks51to53Closed "branch closure covers 5.1-5.3"

    Write-Output "GREEN repair-gate checks passed"
    exit 0
}
catch {
    Write-Error $_
    exit 1
}
