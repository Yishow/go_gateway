$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evidenceRoot = Join-Path $repoRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence"
$handoffRelative = "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/stability-handoff.json"
$handoffPath = Join-Path $repoRoot ($handoffRelative -replace '/', '\')

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) {
        throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')"
    }
}

function Read-Json([string]$RelativePath) {
    $path = Join-Path $repoRoot ($RelativePath -replace '/', '\')
    Assert-True (Test-Path -LiteralPath $path -PathType Leaf) "evidence exists: $RelativePath"
    return Get-Content -Raw -LiteralPath $path | ConvertFrom-Json
}

function Get-Sha256([string]$RelativePath) {
    $path = Join-Path $repoRoot ($RelativePath -replace '/', '\')
    return (Get-FileHash -Algorithm SHA256 -LiteralPath $path).Hash
}

try {
    Assert-True (Test-Path -LiteralPath $handoffPath -PathType Leaf) "stability handoff exists"
    $handoff = Get-Content -Raw -LiteralPath $handoffPath | ConvertFrom-Json
    Assert-Equal 1 $handoff.schemaVersion "stability schema version"
    Assert-Equal "StabilityHandoff" $handoff.mode "stability mode"
    Assert-Equal "blocked/no-code/N/A" $handoff.result "stability result"
    Assert-True $handoff.noStabilityClaim "blocked branch makes no stability claim"

    $reviewRelative = "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/classification-review.json"
    $review = Read-Json $reviewRelative
    Assert-Equal (Get-Sha256 $reviewRelative) $handoff.links.classificationReview.sha256 "classification hash"
    Assert-Equal "locked" $review.repairGate.status "classification repair gate"
    Assert-True (-not $review.repairGate.sourceChangeAllowed) "source change remains locked"

    $repairRelative = "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/repair-branch-handoff.json"
    $repair = Read-Json $repairRelative
    Assert-Equal (Get-Sha256 $repairRelative) $handoff.links.repairBranch.sha256 "repair handoff hash"
    Assert-Equal "blocked/no-code/N/A" $repair.result "repair branch remains blocked"

    $normalRelative = "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/normal-full-summary.json"
    $normal = Read-Json $normalRelative
    Assert-Equal (Get-Sha256 $normalRelative) $handoff.links.normalFull.sha256 "normal baseline hash"
    Assert-Equal "NormalFull" $normal.mode "normal baseline mode"
    Assert-Equal 3 $normal.repeat "normal baseline repeat"
    $expected = @(
        @{ failed = 7; timedOut = 7; classification = "test-timeout" },
        @{ failed = 11; timedOut = 11; classification = "test-timeout" },
        @{ failed = 38; timedOut = 35; classification = "unknown" }
    )
    $records = @($normal.records)
    Assert-Equal 3 $records.Count "three existing normal records"
    for ($i = 0; $i -lt $expected.Count; $i++) {
        Assert-Equal $expected[$i].failed $records[$i].discovered.failed "baseline failed total $($i + 1)"
        Assert-Equal $expected[$i].timedOut $records[$i].discovered.timedOut "baseline timeout total $($i + 1)"
        Assert-Equal $expected[$i].classification $records[$i].classification "baseline classification $($i + 1)"
    }

    Assert-True $handoff.noExtraStabilityFullRuns "no extra three-run stability suite was requested"
    Assert-Equal "not-applicable/no-repair" $handoff.notRun.focusedRepeat "focused repeat is N/A"
    Assert-Equal "not-applicable/no-repair" $handoff.notRun.defaultGroup "default group is N/A"
    Assert-Equal "existing-baseline-only" $handoff.notRun.normalFull "normal full uses existing baseline only"
    Assert-Equal "not-applicable/no-repair" $handoff.notRun.lint "lint is N/A"
    Assert-Equal "not-applicable/no-repair" $handoff.notRun.build "build is N/A"
    Assert-True ($handoff.sourceChangeAllowed -eq $false -and @($handoff.modifiedImpactPaths).Count -eq 0) "no source change"
    Assert-True $handoff.linkedEvidence.rawEvidenceRetained "raw matrix evidence retained"

    Write-Output "GREEN stability-handoff checks passed"
    exit 0
}
catch {
    Write-Error $_
    exit 1
}
