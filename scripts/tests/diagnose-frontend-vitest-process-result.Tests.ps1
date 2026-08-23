$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$ownedOutputs = [System.Collections.Generic.List[string]]::new()
Import-Module (Join-Path $repoRoot "scripts\lib\FrontendVitestProcess.psm1") -Force -Global
Import-Module (Join-Path $repoRoot "scripts\lib\FrontendVitestMatrix.psm1") -Force -Global
Import-Module (Join-Path $repoRoot "scripts\lib\FrontendVitestFreshIsolated.psm1") -Force -Global
Import-Module (Join-Path $repoRoot "scripts\lib\FrontendVitestGroups.psm1") -Force -Global

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) { throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')" }
}

function New-OwnedOutput {
    $path = Join-Path ([IO.Path]::GetTempPath()) ("frontend-process-result-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $path -Force | Out-Null
    [void]$ownedOutputs.Add($path)
    return $path
}

function New-AssertionChild {
    [ordered]@{
        startedAtUtc = [DateTime]::UtcNow.ToString("o")
        completedAtUtc = [DateTime]::UtcNow.ToString("o")
        durationMs = 1
        stdout = " FAIL  tests/unit/app-routing-lazy-load.test.tsx > synthetic > assertion`n AssertionError: synthetic failure`n Test Files  1 failed (1)`n Tests  1 failed (1)"
        stderr = ""
        exitCode = 1
        runnerPid = "synthetic"
        childPids = "unavailable"
        peakWorkingSetBytes = "unavailable"
        externalWatchdogTriggered = $false
        cleanupOutcome = "not-required"
        cleanupError = ""
        startError = ""
        runtimeError = "synthetic runtime failure"
    }
}

try {
    $orderedProbe = [ordered]@{ runtimeError = "ordered runtime failure"; cleanupOutcome = $null }
    Assert-Equal "ordered runtime failure" (Get-FrontendVitestChildResultValue -Child $orderedProbe -Name "runtimeError" -Default "") "OrderedDictionary runtimeError getter"
    Assert-Equal "default" (Get-FrontendVitestChildResultValue -Child $orderedProbe -Name "missing" -Default "default") "child getter default"

    $summaryPath = Join-Path (New-OwnedOutput) "normal-full-summary.json"
    [ordered]@{
        schemaVersion = 1
        mode = "NormalFull"
        records = @([ordered]@{ failureSet = @([ordered]@{ file = "frontend/tests/unit/app-routing-lazy-load.test.tsx"; title = "synthetic"; kind = "assertion" }) })
    } | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $summaryPath -Encoding utf8

    $normalOutput = New-OwnedOutput
    $normalThrown = $false
    try { Invoke-FrontendVitestNormalFull -RepoRoot $repoRoot -OutputDirectory $normalOutput -Repeat 1 -WatchdogSeconds 1 -DiscoveredFiles @("frontend/tests/unit/app-routing-lazy-load.test.tsx") -ChildInvoker { param($RepoRoot, $WatchdogSeconds) New-AssertionChild } | Out-Null } catch { $normalThrown = $true }
    $normalSummary = Get-Content -Raw (Join-Path $normalOutput "normal-full-summary.json") | ConvertFrom-Json
    Assert-True $normalThrown "NormalFull runtimeError gates nonzero"
    Assert-Equal "blocked" $normalSummary.result "NormalFull OrderedDictionary runtimeError is blocked"

    $freshOutput = New-OwnedOutput
    $freshThrown = $false
    $freshError = ""
    try { Invoke-FrontendVitestFreshIsolated -RepoRoot $repoRoot -OutputDirectory $freshOutput -Repeat 1 -WatchdogSeconds 1 -FailureFiles @("frontend/tests/unit/app-routing-lazy-load.test.tsx") -SummaryPath $summaryPath -ChildInvoker { param($RepoRoot, $FilePath, $RunId, $RepeatIndex, $WatchdogSeconds) New-AssertionChild } | Out-Null } catch { $freshThrown = $true; $freshError = "$(($_.Exception.Message)) at $(($_.ScriptStackTrace))" }
    if (-not (Test-Path -LiteralPath (Join-Path $freshOutput "fresh-isolated-summary.json") -PathType Leaf)) { throw "Fresh synthetic invocation failed before summary: $freshError" }
    $freshSummary = Get-Content -Raw (Join-Path $freshOutput "fresh-isolated-summary.json") | ConvertFrom-Json
    Assert-True $freshThrown "FreshIsolated runtimeError gates nonzero"
    Assert-Equal "blocked" $freshSummary.result "FreshIsolated OrderedDictionary runtimeError is blocked"

    $groupOutput = New-OwnedOutput
    $groupThrown = $false
    try { Invoke-FrontendVitestGroup -RepoRoot $repoRoot -OutputDirectory $groupOutput -Repeat 1 -WatchdogSeconds 1 -Mode DefaultGroup -SummaryPath $summaryPath -ChildInvoker { param($RepoRoot, $Mode, $Files, $RunId, $RepeatIndex, $WatchdogSeconds) New-AssertionChild } | Out-Null } catch { $groupThrown = $true }
    $groupSummary = Get-Content -Raw (Join-Path $groupOutput "group-summary.json") | ConvertFrom-Json
    Assert-True $groupThrown "Group runtimeError gates nonzero"
    Assert-Equal "blocked" $groupSummary.result "Group OrderedDictionary runtimeError is blocked"

    Write-Output "GREEN process-result contracts passed"
    exit 0
}
finally {
    foreach ($path in $ownedOutputs | Sort-Object -Unique) {
        if (Test-Path -LiteralPath $path -PathType Container) { Remove-Item -LiteralPath $path -Recurse -Force -ErrorAction SilentlyContinue }
    }
}
