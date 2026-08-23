$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evidenceModule = Join-Path $repoRoot "scripts\lib\FrontendVitestEvidence.psm1"
Import-Module $evidenceModule -Force -Global
$handoffPath = Resolve-FrontendVitestEvidencePath -RepoRoot $repoRoot -FileName "repair-branch-handoff.json"
$handoffEvidenceRoot = Split-Path -Parent $handoffPath
$allowedCurrentPaths = @(
    "scripts/diagnose-frontend-vitest-matrix.ps1",
    "scripts/lib/FrontendVitestEvidence.psm1", "scripts/lib/FrontendVitestFreshIsolated.psm1", "scripts/lib/FrontendVitestGroups.psm1", "scripts/lib/FrontendVitestMatrix.psm1", "scripts/lib/FrontendVitestProcess.psm1",
    "scripts/tests/diagnose-frontend-classification-review.Tests.ps1", "scripts/tests/diagnose-frontend-evidence-resolver.Tests.ps1", "scripts/tests/diagnose-frontend-final-audit.Tests.ps1", "scripts/tests/diagnose-frontend-fresh-isolated.Tests.ps1", "scripts/tests/diagnose-frontend-groups.Tests.ps1", "scripts/tests/diagnose-frontend-process-safety.Tests.ps1", "scripts/tests/diagnose-frontend-repair-gate.Tests.ps1", "scripts/tests/diagnose-frontend-stability-handoff.Tests.ps1", "scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1", "scripts/tests/diagnose-frontend-vitest-process-result.Tests.ps1",
    "tests/powershell/start-process-tree.ps1"
)

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) { throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')" }
}

function Test-ExactCurrentInventory([string[]]$Actual, [string[]]$Allowed) {
    $actualSet = @($Actual | ForEach-Object { ([string]$_).Trim().Replace('\', '/') } | Where-Object { $_ } | Sort-Object -Unique)
    $allowedSet = @($Allowed | ForEach-Object { ([string]$_).Trim().Replace('\', '/') } | Where-Object { $_ } | Sort-Object -Unique)
    if ($actualSet.Count -eq 0) { return $true }
    return $actualSet.Count -eq $allowedSet.Count -and (($actualSet -join "`n") -eq ($allowedSet -join "`n"))
}

function Get-GitInventory {
    $trackedOutput = @(& rtk git diff --name-only HEAD 2>$null)
    $trackedExitCode = $LASTEXITCODE
    if ($trackedExitCode -ne 0) { throw "git diff inventory command failed with exit code $trackedExitCode" }
    $untrackedOutput = @(& rtk git ls-files --others --exclude-standard 2>$null)
    $untrackedExitCode = $LASTEXITCODE
    if ($untrackedExitCode -ne 0) { throw "git ls-files inventory command failed with exit code $untrackedExitCode" }
    @($trackedOutput + $untrackedOutput) |
        ForEach-Object { ([string]$_).Trim().Replace('\', '/') } | Where-Object { $_ -and $_ -ne "Changes:" } | Sort-Object -Unique
}

try {
    Assert-True (Test-Path -LiteralPath $handoffPath -PathType Leaf) "repair branch handoff exists"
    $handoff = Get-Content -Raw -LiteralPath $handoffPath | ConvertFrom-Json
    Assert-Equal 1 $handoff.schemaVersion "handoff schema version"
    Assert-Equal "RepairBranchHandoff" $handoff.mode "handoff mode"
    Assert-Equal "blocked/no-code/N/A" $handoff.result "handoff result"

    $reviewPath = Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $handoffEvidenceRoot -Reference ([string]$handoff.classificationReview.path) -ExpectedFileName "classification-review.json"
    Assert-True (Test-Path -LiteralPath $reviewPath -PathType Leaf) "classification review exists"
    Assert-Equal "classification-review.json" ([IO.Path]::GetFileName(([string]$handoff.classificationReview.path))) "classification review path"
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
    Assert-True (@($handoffInventory).Count -gt 0) "archived handoff retains its recorded untracked inventory"
    $inventorySource = Get-Content -Raw -LiteralPath $PSCommandPath
    Assert-True ($inventorySource -match '\$trackedExitCode\s*=\s*\$LASTEXITCODE' -and $inventorySource -match '\$untrackedExitCode\s*=\s*\$LASTEXITCODE') "Git inventory checks both command exit codes"
    Assert-Equal 17 @($allowedCurrentPaths).Count "current repair allowlist is exact"
    if (@($actualInventory).Count -gt 0) { Assert-True (Test-ExactCurrentInventory -Actual $actualInventory -Allowed $allowedCurrentPaths) "non-empty current inventory exactly equals owned repair set" }
    foreach ($path in $actualInventory) {
        Assert-True ($allowedCurrentPaths -contains $path) "path is exactly inside diagnosis/change allowlist: $path"
        Assert-True ($path -notmatch "^(frontend/|.*\.go$|.*package(-lock)?\.json$|.*\.config\.|\.line-limit-ignore$)") "path is outside forbidden scope: $path"
    }
    $syntheticAllowed = @($allowedCurrentPaths | Sort-Object)
    Assert-True (Test-ExactCurrentInventory -Actual $syntheticAllowed -Allowed $allowedCurrentPaths) "exact inventory helper accepts complete set"
    Assert-True (-not (Test-ExactCurrentInventory -Actual @($syntheticAllowed | Select-Object -Skip 1) -Allowed $allowedCurrentPaths)) "exact inventory helper rejects missing path"
    Assert-True (-not (Test-ExactCurrentInventory -Actual (@($syntheticAllowed | Select-Object -Skip 1) + @("scripts/tests/arbitrary-unowned.ps1")) -Allowed $allowedCurrentPaths)) "exact inventory helper rejects extra path"
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
