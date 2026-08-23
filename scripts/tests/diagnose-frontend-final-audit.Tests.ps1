$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$handoffRelative = "openspec/changes/diagnose-frontend-full-suite-timeouts/evidence/final-audit-handoff.json"
$handoffPath = Join-Path $repoRoot ($handoffRelative -replace '/', '\')

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) {
        throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')"
    }
}

function Get-RelativePath([string]$RelativePath) {
    return Join-Path $repoRoot ($RelativePath -replace '/', '\')
}

function Get-Sha256([string]$RelativePath) {
    return (Get-FileHash -Algorithm SHA256 -LiteralPath (Get-RelativePath $RelativePath)).Hash
}

function Get-GitInventory {
    @(& rtk git ls-files --others --exclude-standard 2>$null | ForEach-Object { ([string]$_).Trim() } | Where-Object { $_ } | Sort-Object -Unique)
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
    Assert-Equal ($actualInventory -join '|') ($recordedInventory -join '|') "final audit inventory matches Git"
    foreach ($path in $actualInventory) {
        Assert-True ($path -notmatch "^(frontend/|.*\.go$|.*package(-lock)?\.json$|.*\.config\.|\.line-limit-ignore$)") "path is outside forbidden scope: $path"
    }

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
