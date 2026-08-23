$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evidenceModule = Join-Path $repoRoot "scripts\lib\FrontendVitestEvidence.psm1"
Import-Module $evidenceModule -Force -Global
$handoffRelative = "final-audit-handoff.json"
$handoffPath = Resolve-FrontendVitestEvidencePath -RepoRoot $repoRoot -FileName $handoffRelative
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
    if ($Expected -ne $Actual) {
        throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')"
    }
}

function Assert-Throws([scriptblock]$Action, [string]$Message) {
    $thrown = $false
    try { & $Action } catch { $thrown = $true }
    Assert-True $thrown $Message
}

function Test-ExactCurrentInventory([string[]]$Actual, [string[]]$Allowed) {
    $actualSet = @($Actual | ForEach-Object { ([string]$_).Trim().Replace('\', '/') } | Where-Object { $_ } | Sort-Object -Unique)
    $allowedSet = @($Allowed | ForEach-Object { ([string]$_).Trim().Replace('\', '/') } | Where-Object { $_ } | Sort-Object -Unique)
    if ($actualSet.Count -eq 0) { return $true }
    return $actualSet.Count -eq $allowedSet.Count -and (($actualSet -join "`n") -eq ($allowedSet -join "`n"))
}

function Get-RelativePath([string]$RelativePath) {
    if ($RelativePath -match '(?i)evidence[\\/](?<file>[A-Za-z0-9][A-Za-z0-9._-]*\.json)$') {
        return Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $handoffEvidenceRoot -Reference $RelativePath -ExpectedFileName $Matches.file
    }
    return Join-Path $repoRoot ($RelativePath -replace '/', '\')
}

function Get-Sha256([string]$RelativePath) {
    return (Get-FileHash -Algorithm SHA256 -LiteralPath (Get-RelativePath $RelativePath)).Hash
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
    Assert-True (Test-Path -LiteralPath $handoffPath -PathType Leaf) "final audit handoff exists"
    $handoff = Get-Content -Raw -LiteralPath $handoffPath | ConvertFrom-Json
    Assert-Equal 1 $handoff.schemaVersion "final audit schema version"
    Assert-Equal "FinalAuditHandoff" $handoff.mode "final audit mode"
    Assert-Equal "blocked/no-code/N/A" $handoff.result "final audit no-code result"
    Assert-True $handoff.noStabilityClaim "final audit makes no stability claim"

    foreach ($link in @($handoff.linkedEvidence.items)) {
        $path = [string]$link.path
        Assert-True (Test-Path -LiteralPath (Get-RelativePath $path) -PathType Leaf) "linked evidence exists: $path"
        Assert-Equal (Get-Sha256 $path) $link.sha256 "linked evidence hash: $path"
    }
    Assert-True $handoff.linkedEvidence.rawEvidenceRetained "raw evidence retained"
    Assert-True (-not $handoff.linkedEvidence.staleEnvironmentMarkdownPresent) "stale environment markdown is absent"
    Assert-Throws { Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $handoffEvidenceRoot -Reference "openspec/changes/other-change/evidence/classification-review.json" -ExpectedFileName "classification-review.json" } "other change reference is rejected"
    Assert-Throws { Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $handoffEvidenceRoot -Reference "openspec/changes/archive/2026-08-22-other-change/evidence/classification-review.json" -ExpectedFileName "classification-review.json" } "other archive reference is rejected"
    Assert-Throws { Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $handoffEvidenceRoot -Reference "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/classification-review.json" -ExpectedFileName "repair-branch-handoff.json" } "filename mismatch is rejected"
    Assert-Throws { Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $handoffEvidenceRoot -Reference "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/../classification-review.json" -ExpectedFileName "classification-review.json" } "reference traversal is rejected"

    Assert-Equal 0 $handoff.spectraAnalyze.critical "Spectra analyze critical findings"
    Assert-Equal 0 $handoff.spectraAnalyze.warning "Spectra analyze warning findings"
    Assert-Equal 13 $handoff.spectraAnalyze.suggestions "Spectra analyze suggestions are non-blocking"
    Assert-True $handoff.spectraValidate.valid "Spectra validate passes"
    Assert-Equal "Passed" $handoff.lineGate.status "line gate passes"
    Assert-Equal "clean" $handoff.diffCheck.status "diff check is clean"
    Assert-Equal 0 @($handoff.diffCheck.trackedDiffPaths).Count "no tracked diff paths"
    Assert-Equal 0 $handoff.scopeAudit.forbiddenCount "forbidden scope count"
    Assert-Equal 0 $handoff.secretScan.matchCount "secret scan match count"

    $actualInventory = Get-GitInventory
    $recordedInventory = @($handoff.scopeAudit.untrackedPaths | ForEach-Object { ([string]$_).Replace('\', '/') } | Sort-Object -Unique)
    Assert-True (@($recordedInventory).Count -gt 0) "archived final audit retains its recorded inventory"
    $inventorySource = Get-Content -Raw -LiteralPath $PSCommandPath
    Assert-True ($inventorySource -match '\$trackedExitCode\s*=\s*\$LASTEXITCODE' -and $inventorySource -match '\$untrackedExitCode\s*=\s*\$LASTEXITCODE') "Git inventory checks both command exit codes"
    Assert-Equal 17 @($allowedCurrentPaths).Count "current repair allowlist is exact"
    if (@($actualInventory).Count -gt 0) { Assert-True (Test-ExactCurrentInventory -Actual $actualInventory -Allowed $allowedCurrentPaths) "non-empty current inventory exactly equals owned repair set" }
    foreach ($path in $actualInventory) {
        Assert-True ($path -notmatch "^(frontend/|.*\.go$|.*package(-lock)?\.json$|.*\.config\.|\.line-limit-ignore$)") "path is outside forbidden scope: $path"
        Assert-True ($allowedCurrentPaths -contains $path) "current diff path is exactly owned: $path"
    }
    $syntheticAllowed = @($allowedCurrentPaths | Sort-Object)
    Assert-True (Test-ExactCurrentInventory -Actual $syntheticAllowed -Allowed $allowedCurrentPaths) "exact inventory helper accepts complete set"
    Assert-True (-not (Test-ExactCurrentInventory -Actual @($syntheticAllowed | Select-Object -Skip 1) -Allowed $allowedCurrentPaths)) "exact inventory helper rejects missing path"
    Assert-True (-not (Test-ExactCurrentInventory -Actual (@($syntheticAllowed | Select-Object -Skip 1) + @("scripts/tests/arbitrary-unowned.ps1")) -Allowed $allowedCurrentPaths)) "exact inventory helper rejects extra path"

    Assert-True $handoff.auditDiscipline.fixedCommandsOnly "fixed command discipline"
    Assert-True $handoff.auditDiscipline.failClosedClassification "fail-closed classification"
    Assert-True $handoff.auditDiscipline.noGlobalKill "no global process kill"
    Assert-True $handoff.auditDiscipline.pathValidation "path validation remains enabled"
    Assert-True $handoff.auditDiscipline.noCommandInjection "no command injection surface"
    Assert-True (-not $handoff.sourceChangeAllowed -and @($handoff.modifiedImpactPaths).Count -eq 0) "no source change"
    Assert-Equal "13/14" $handoff.tasksProgress.beforeFinalMark "pre-final task progress"
    Assert-Equal "14/14" $handoff.tasksProgress.expectedAfterFinalMark "expected final task progress"
    Assert-True $handoff.tasksProgress.tasksFileIsSourceOfTruth "tasks file remains source of truth"
    Assert-Equal "not-run/no-repair" $handoff.notRun.npmVitest "npm/Vitest not run in no-code branch"
    Assert-Equal "not-run/no-repair" $handoff.notRun.frontendLintBuild "frontend lint/build not run"

    Write-Output "GREEN final-audit checks passed"
    exit 0
}
catch {
    Write-Error $_
    exit 1
}
