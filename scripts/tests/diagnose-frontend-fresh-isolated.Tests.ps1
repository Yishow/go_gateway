$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$runnerPath = Join-Path $repoRoot "scripts\diagnose-frontend-vitest-matrix.ps1"
$ownedOutputs = [System.Collections.Generic.List[string]]::new()

function Assert-True {
    param([bool]$Condition, [string]$Message)
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal {
    param($Expected, $Actual, [string]$Message)
    if ($Expected -ne $Actual) { throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')" }
}

function New-OwnedOutput {
    $path = Join-Path ([IO.Path]::GetTempPath()) ("frontend-fresh-contract-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $path -Force | Out-Null
    $ownedOutputs.Add($path)
    return $path
}

try {
    . $runnerPath
    $summaryPath = Join-Path (New-OwnedOutput) "normal-full-summary.json"
    $syntheticSummary = [ordered]@{
        schemaVersion = 1
        mode = "NormalFull"
        records = @(
            [ordered]@{ failureSet = @(
                [ordered]@{ file = "frontend/tests/unit/app-routing-lazy-load.test.tsx"; title = "route"; kind = "test-timeout" },
                [ordered]@{ file = "frontend/tests/unit/app-routing-lazy-load.test.tsx"; title = "route"; kind = "test-timeout" },
                [ordered]@{ file = "frontend\tests\unit\app-routing-lazy-load.test.tsx"; title = "route"; kind = "test-timeout" },
                [ordered]@{ file = "frontend/tests/unit/components/ui/spinner.test.tsx"; title = "spinner"; kind = "assertion" }
            ) }
        )
    }
    $syntheticSummary | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $summaryPath -Encoding utf8

    $source = Get-FreshIsolatedFailureSet -SummaryPath $summaryPath -RepoRoot $repoRoot
    Assert-Equal 2 @($source.files).Count "dynamic failure union deduplicates normalized paths"
    Assert-True ($source.files -contains "frontend/tests/unit/app-routing-lazy-load.test.tsx") "union retains first observed path"
    Assert-True ($source.files -contains "frontend/tests/unit/components/ui/spinner.test.tsx") "union retains second observed path"
    Assert-True (-not [string]::IsNullOrWhiteSpace($source.summarySha256)) "source summary hash is recorded"
    Assert-True ($source.derivation -match "not a fixed path list") "failure union remains dynamic"

    $arguments = @(Get-FreshIsolatedCommandArguments -FilePath $source.files[0])
    Assert-Equal 7 $arguments.Count "fresh command has one validated file selector"
    $expectedArguments = @("--prefix", "frontend", "run", "test", "--", "--run", $source.files[0].Substring("frontend/".Length))
    Assert-Equal ($expectedArguments -join "|") ($arguments -join "|") "fresh command arguments are exact"
    Assert-Equal "--prefix" $arguments[0] "fixed npm prefix argument"
    Assert-Equal "frontend" $arguments[1] "fixed frontend working target"
    Assert-Equal "--run" $arguments[5] "fixed Vitest run argument"
    Assert-Equal ($source.files[0].Substring("frontend/".Length)) $arguments[6] "validated file selector is last"
    Assert-True ($arguments[6] -notmatch '^frontend/') "frontend package selector is relative to its working root"
    $injectionRejected = $false
    try { Get-FreshIsolatedCommandArguments -FilePath 'frontend/tests/unit/evil";Write-Output.test.tsx' | Out-Null } catch { $injectionRejected = $true }
    Assert-True $injectionRejected "command metacharacters are rejected"
    Assert-True ((Get-Content -Raw $runnerPath) -notmatch '(?im)^\s*\[string\]\$Command\b') "CLI does not expose arbitrary command input"

    foreach ($invalidPath in @(
        "C:\outside.test.ts",
        "frontend/tests/../outside.test.ts",
        "frontend/tests/unit/not-a-test.ts",
        "frontend/tests/unit/missing-fresh-file.test.ts"
    )) {
        $invalidSummaryPath = Join-Path (New-OwnedOutput) "invalid.json"
        [ordered]@{ schemaVersion = 1; mode = "NormalFull"; records = @([ordered]@{ failureSet = @([ordered]@{ file = $invalidPath; title = "invalid"; kind = "assertion" }) }) } |
            ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $invalidSummaryPath -Encoding utf8
        $rejected = $false
        try { Get-FreshIsolatedFailureSet -SummaryPath $invalidSummaryPath -RepoRoot $repoRoot | Out-Null } catch { $rejected = $true }
        Assert-True $rejected "invalid path '$invalidPath' is rejected"
    }

    $missingRejected = $false
    try { Get-FreshIsolatedFailureSet -SummaryPath (Join-Path (New-OwnedOutput) "missing.json") -RepoRoot $repoRoot | Out-Null } catch { $missingRejected = $true }
    Assert-True $missingRejected "missing NormalFull evidence fails loudly"

    $groupRejectedOutput = & pwsh -NoProfile -File $runnerPath -Mode FreshIsolated -GroupFiles "frontend/tests/unit/app-routing-lazy-load.test.tsx" 2>&1
    Assert-True ($LASTEXITCODE -ne 0) "FreshIsolated rejects unrelated GroupFiles"

    $runOutput = New-OwnedOutput
    $syntheticChild = {
        param($RepoRoot, $FilePath, $RunId, $RepeatIndex, $WatchdogSeconds)
        $started = [DateTime]::UtcNow
        $completed = [DateTime]::UtcNow
        [ordered]@{
            startedAtUtc = $started.ToString("o")
            completedAtUtc = $completed.ToString("o")
            durationMs = 1
            stdout = " Test Files  1 passed (1)`n Tests  1 passed (1)"
            stderr = ""
            exitCode = 0
            runnerPid = "synthetic"
            childPids = "unavailable"
            peakWorkingSetBytes = "unavailable"
            externalWatchdogTriggered = $false
            cleanupOutcome = "not-required"
            cleanupError = ""
            startError = ""
        }
    }
    $fresh = Invoke-FrontendVitestFreshIsolated -RepoRoot $repoRoot -OutputDirectory $runOutput -Repeat 2 -WatchdogSeconds 1 -FailureFiles $source.files -SummaryPath $summaryPath -ChildInvoker $syntheticChild | ConvertFrom-Json
    Assert-Equal 4 $fresh.recordCount "synthetic fresh matrix record count"
    Assert-Equal 4 @($fresh.runIds | Sort-Object -Unique).Count "synthetic fresh run ids are unique"
    Assert-True ($fresh.sequential -eq $true) "synthetic fresh runs are sequential"
    Assert-Equal 4 @($fresh.artifactPaths).Count "synthetic fresh artifacts are independently listed"
    Assert-Equal "completed" $fresh.result "synthetic fresh phase completes"
    $freshSummary = Get-Content -Raw -LiteralPath (Join-Path $runOutput "fresh-isolated-summary.json") | ConvertFrom-Json
    Assert-Equal 2 @($freshSummary.perFile).Count "synthetic fresh summary keeps per-file aggregates"
    Assert-True (@($freshSummary.perFile | ForEach-Object repeatCount) -contains 2) "synthetic fresh per-file repeat counts are isolated"

    Write-Output "GREEN fresh-isolated contract checks passed"
    exit 0
}
finally {
    foreach ($path in $ownedOutputs | Sort-Object -Unique) {
        if (Test-Path -LiteralPath $path -PathType Container) { Remove-Item -LiteralPath $path -Recurse -Force -ErrorAction SilentlyContinue }
    }
}
