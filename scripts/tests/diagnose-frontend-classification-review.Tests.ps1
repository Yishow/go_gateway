$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$evidenceModule = Join-Path $repoRoot "scripts\lib\FrontendVitestEvidence.psm1"
Import-Module $evidenceModule -Force -Global
$reviewPath = Resolve-FrontendVitestEvidencePath -RepoRoot $repoRoot -FileName "classification-review.json"
$evidenceRoot = Split-Path -Parent $reviewPath

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) { throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')" }
}

function Assert-HashArtifact([string]$Path, [string]$ExpectedHash, [string]$Message) {
    Assert-True (Test-Path -LiteralPath $Path -PathType Leaf) "$Message exists"
    Assert-Equal $ExpectedHash (Get-FileHash -Algorithm SHA256 -LiteralPath $Path).Hash "$Message hash"
}

function Resolve-ReviewEvidencePath([string]$Reference, [string]$ExpectedFileName) {
    return Resolve-FrontendVitestLinkedEvidencePath -AnchorEvidenceRoot $evidenceRoot -Reference $Reference -ExpectedFileName $ExpectedFileName
}

function Get-ReviewSourcePath([string]$Name) {
    return Resolve-ReviewEvidencePath -Reference ([string]$sourceEvidenceByName[$Name].path) -ExpectedFileName $expectedSources[$Name]
}

try {
    Assert-True (Test-Path -LiteralPath $reviewPath -PathType Leaf) "classification review artifact exists"
    $review = Get-Content -Raw -LiteralPath $reviewPath | ConvertFrom-Json
    Assert-Equal 1 $review.schemaVersion "review schema version"
    Assert-Equal "ClassificationReview" $review.mode "review mode"
    Assert-Equal "locked" $review.repairGate.status "repair gate is locked"
    Assert-True ($review.repairGate.sourceChangeAllowed -eq $false) "source repair is not allowed"
    Assert-Equal "stopped" $review.oneVariableBranch.decision "one-variable branch is stopped"
    Assert-True ($review.identityPolicy.dynamicFailureSet -and -not $review.identityPolicy.fixedTimeoutIdentity) "failure identity remains dynamic"

    $expectedSources = [ordered]@{
        environment = "environment-capture.json"
        normalFull = "normal-full-summary.json"
        freshIsolated = "fresh-isolated-summary.json"
        groupComparison = "group-comparison-summary.json"
    }
    $sourceEvidenceByName = @{}
    foreach ($name in $expectedSources.Keys) {
        $source = @($review.sourceEvidence | Where-Object name -eq $name)
        Assert-Equal 1 $source.Count "source '$name' is linked once"
        $relative = ([string]$source[0].path).Replace('\', '/')
        Assert-Equal $expectedSources[$name] ([IO.Path]::GetFileName($relative)) "source '$name' path"
        $sourceEvidenceByName[$name] = $source[0]
        $full = Resolve-ReviewEvidencePath -Reference ([string]$source[0].path) -ExpectedFileName $expectedSources[$name]
        Assert-True (Test-Path -LiteralPath $full -PathType Leaf) "source '$name' exists"
        Assert-Equal (Get-FileHash -Algorithm SHA256 -LiteralPath $full).Hash $source[0].sha256 "source '$name' hash"
    }
    Assert-True ($review.hashValidation.allPresent -and $review.hashValidation.allSha256Match) "source hashes are validated"
    Assert-Equal 4 $review.hashValidation.sourceCount "all four source summaries are hashed"

    Assert-Equal 181 $review.evidenceFacts.environment.discoveredFiles "environment inventory count"
    Assert-True ($review.evidenceFacts.environment.notFixedIdentity -and -not $review.evidenceFacts.environment.suiteExecuted) "environment capture is inventory-only"
    Assert-Equal 23 $review.evidenceFacts.freshIsolated.fileCount "fresh dynamic file count"
    Assert-Equal 230 $review.evidenceFacts.freshIsolated.recordCount "fresh record count"
    Assert-Equal 230 $review.evidenceFacts.freshIsolated.passCount "fresh pass count"
    Assert-Equal 4 $review.evidenceFacts.groupComparison.recordCount "group record count"
    Assert-True ($review.evidenceFacts.groupComparison.filesEqual -and $review.evidenceFacts.groupComparison.nonWorkerSettingsDiff.Count -eq 0) "group settings are comparable"
    Assert-True ($review.evidenceFacts.groupComparison.singleWorkerDiagnosisOnly) "single worker is diagnosis-only"

    Assert-Equal 9 $review.rawValidation.normalFull.hashChecks "NormalFull raw hash checks"
    Assert-Equal 690 $review.rawValidation.freshIsolated.artifactChecks "FreshIsolated artifact checks"
    Assert-Equal 12 $review.rawValidation.groupComparison.artifactChecks "group artifact checks"
    Assert-True ($review.rawValidation.normalFull.allPresent -and $review.rawValidation.normalFull.allHashesMatch) "NormalFull raw evidence integrity"
    Assert-True ($review.rawValidation.freshIsolated.allPresent -and $review.rawValidation.freshIsolated.allHashesMatch) "FreshIsolated raw evidence integrity"
    Assert-True ($review.rawValidation.groupComparison.allPresent -and $review.rawValidation.groupComparison.allHashesMatch) "group raw evidence integrity"

    $normal = Get-Content -Raw (Get-ReviewSourcePath "normalFull") | ConvertFrom-Json
    foreach ($run in @($normal.records)) { foreach ($artifact in @($run.rawArtifacts)) { Assert-HashArtifact -Path $artifact.path -ExpectedHash $artifact.sha256 -Message "NormalFull raw artifact" } }
    $fresh = Get-Content -Raw (Get-ReviewSourcePath "freshIsolated") | ConvertFrom-Json
    $freshRecords = @(Get-ChildItem -LiteralPath $fresh.currentRawLocation -Recurse -Filter "record.json" -File)
    Assert-Equal 230 $freshRecords.Count "FreshIsolated raw record count"
    foreach ($recordFile in $freshRecords) {
        $record = Get-Content -Raw $recordFile.FullName | ConvertFrom-Json
        $runDir = Split-Path -Parent $recordFile.FullName
        foreach ($kind in @("stdout", "stderr")) { Assert-HashArtifact -Path (Join-Path $runDir ($kind + ".txt")) -ExpectedHash $record.artifactHashes.$kind -Message "FreshIsolated $kind artifact" }
    }
    $group = Get-Content -Raw (Get-ReviewSourcePath "groupComparison") | ConvertFrom-Json
    foreach ($phase in @($group.phases)) { foreach ($i in 1..2) { $runDir = Join-Path $phase.rawLocation ("run-{0:D3}" -f $i); $record = Get-Content -Raw (Join-Path $runDir "record.json") | ConvertFrom-Json; foreach ($kind in @("stdout", "stderr")) { Assert-HashArtifact -Path (Join-Path $runDir ($kind + ".txt")) -ExpectedHash $record.artifactHashes.$kind -Message "$($phase.mode) $kind artifact" } } }

    $expectedClassification = @{
        "NormalFull#1" = @{ classification = "test-timeout"; result = "failed" }
        "NormalFull#2" = @{ classification = "test-timeout"; result = "failed" }
        "NormalFull#3" = @{ classification = "unknown"; result = "blocked" }
        "DefaultGroup#1" = @{ classification = "unknown"; result = "blocked" }
        "DefaultGroup#2" = @{ classification = "unknown"; result = "blocked" }
        "FreshIsolated#all" = @{ classification = "pass"; result = "completed" }
        "SingleWorkerGroup#1" = @{ classification = "pass"; result = "pass" }
        "SingleWorkerGroup#2" = @{ classification = "pass"; result = "pass" }
    }
    foreach ($key in $expectedClassification.Keys) {
        $case = @($review.classificationTable | Where-Object identity -eq $key)
        Assert-Equal 1 $case.Count "classification case '$key' exists"
        Assert-Equal $expectedClassification[$key].classification $case[0].classification "classification '$key'"
        Assert-Equal $expectedClassification[$key].result $case[0].result "result '$key'"
    }
    Assert-True (-not $review.classificationFacts.resourceLeakObserved) "resource-leak classification is not claimed"
    Assert-True (-not $review.classificationFacts.workerHangObserved) "worker-hang classification is not claimed"
    Assert-True ($review.classificationFacts.unknownAlwaysBlocked) "unknown is fail-closed"

    $expectedHypotheses = @{
        "worker-resource-saturation" = "strongly-supported"
        "test-local-lifecycle" = "not-supported-for-current-union"
        "order-isolation-coupling" = "insufficient-evidence"
        "deterministic-assertion-defect" = "not-supported"
        "environment-effect" = "separate-boundary-only"
    }
    foreach ($id in $expectedHypotheses.Keys) {
        $hypothesis = @($review.hypotheses | Where-Object id -eq $id)
        Assert-Equal 1 $hypothesis.Count "hypothesis '$id' exists"
        Assert-Equal $expectedHypotheses[$id] $hypothesis[0].status "hypothesis '$id' status"
        Assert-True (@($hypothesis[0].evidenceLinks).Count -ge 2) "hypothesis '$id' has comparable evidence links"
        Assert-True (-not [string]::IsNullOrWhiteSpace([string]$hypothesis[0].falsificationTest)) "hypothesis '$id' has falsification test"
        Assert-True (-not [string]::IsNullOrWhiteSpace([string]$hypothesis[0].consequence)) "hypothesis '$id' has consequence"
    }
    Assert-Equal "unavailable" $review.hypotheses[0].limits.directResourceCorrelation "resource correlation is unavailable"
    Assert-True (-not $review.hypotheses[0].limits.productionRootCauseProven) "production root cause is unproven"
    foreach ($branch in @($review.oneVariableBranch.branches)) { Assert-Equal "stopped" $branch.status "branch '$($branch.name)' is stopped" }
    Assert-True ($review.repairGate.decision -match "blocked/no-code/N/A") "repair branch is blocked no-code N/A"
    Assert-True (-not $review.secretsPresent) "review contains no secrets"

    Write-Output "GREEN classification review checks passed"
    exit 0
}
catch {
    Write-Error $_
    exit 1
}
