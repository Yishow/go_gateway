$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evidenceModule = Join-Path $repoRoot "scripts\lib\FrontendVitestEvidence.psm1"
Import-Module $evidenceModule -Force -Global
$handoffPath = Resolve-FrontendVitestEvidencePath -RepoRoot $repoRoot -FileName "stability-handoff.json"
$evidenceRoot = Split-Path -Parent $handoffPath

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) {
        throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')"
    }
}

function Resolve-TestEvidencePath([string]$Reference, [string]$ExpectedFileName) {
    return Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $evidenceRoot -Reference $Reference -ExpectedFileName $ExpectedFileName
}

function Read-Json([string]$Reference, [string]$ExpectedFileName) {
    $path = Resolve-TestEvidencePath $Reference $ExpectedFileName
    Assert-True (Test-Path -LiteralPath $path -PathType Leaf) "evidence exists: $Reference"
    return Get-Content -Raw -LiteralPath $path | ConvertFrom-Json
}

function Get-Sha256([string]$Reference, [string]$ExpectedFileName) {
    return (Get-FileHash -Algorithm SHA256 -LiteralPath (Resolve-TestEvidencePath $Reference $ExpectedFileName)).Hash
}

try {
    Assert-True (Test-Path -LiteralPath $handoffPath -PathType Leaf) "stability handoff exists"
    $handoff = Get-Content -Raw -LiteralPath $handoffPath | ConvertFrom-Json
    Assert-Equal 1 $handoff.schemaVersion "stability schema version"
    Assert-Equal "StabilityHandoff" $handoff.mode "stability mode"
    Assert-Equal "blocked/no-code/N/A" $handoff.result "stability result"
    Assert-True $handoff.noStabilityClaim "blocked branch makes no stability claim"

    $reviewReference = [string]$handoff.links.classificationReview.path
    $reviewFile = [IO.Path]::GetFileName($reviewReference)
    $review = Read-Json $reviewReference "classification-review.json"
    Assert-Equal (Get-Sha256 $reviewReference "classification-review.json") $handoff.links.classificationReview.sha256 "classification hash"
    Assert-Equal "locked" $review.repairGate.status "classification repair gate"
    Assert-True (-not $review.repairGate.sourceChangeAllowed) "source change remains locked"

    $repairReference = [string]$handoff.links.repairBranch.path
    $repairFile = [IO.Path]::GetFileName($repairReference)
    $repair = Read-Json $repairReference "repair-branch-handoff.json"
    Assert-Equal (Get-Sha256 $repairReference "repair-branch-handoff.json") $handoff.links.repairBranch.sha256 "repair handoff hash"
    Assert-Equal "blocked/no-code/N/A" $repair.result "repair branch remains blocked"

    $normalReference = [string]$handoff.links.normalFull.path
    $normalFile = [IO.Path]::GetFileName($normalReference)
    $normal = Read-Json $normalReference "normal-full-summary.json"
    Assert-Equal (Get-Sha256 $normalReference "normal-full-summary.json") $handoff.links.normalFull.sha256 "normal baseline hash"
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
