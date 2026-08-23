$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$runnerPath = Join-Path $repoRoot "scripts\diagnose-frontend-vitest-matrix.ps1"
$summaryPath = Join-Path $repoRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence\normal-full-summary.json"
$ownedOutputs = [System.Collections.Generic.List[string]]::new()

function Assert-True([bool]$Condition, [string]$Message) {
    if (-not $Condition) { throw "ASSERTION FAILED: $Message" }
}

function Assert-Equal($Expected, $Actual, [string]$Message) {
    if ($Expected -ne $Actual) { throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')" }
}

function New-OwnedOutput {
    $path = Join-Path ([IO.Path]::GetTempPath()) ("frontend-groups-contract-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $path -Force | Out-Null
    [void]$ownedOutputs.Add($path)
    $path
}

try {
    . $runnerPath
    $sourcePath = Join-Path (New-OwnedOutput) "synthetic-normal.json"
    $fileA = "frontend/tests/unit/app-routing-lazy-load.test.tsx"
    $fileB = "frontend/tests/unit/components/ui/spinner.test.tsx"
    [ordered]@{
        schemaVersion = 1
        mode = "NormalFull"
        records = @(
            [ordered]@{ failureSet = @(
                [ordered]@{ file = $fileA; title = "a"; kind = "test-timeout" },
                [ordered]@{ file = "frontend\tests\unit\app-routing-lazy-load.test.tsx"; title = "a"; kind = "test-timeout" },
                [ordered]@{ file = $fileB; title = "b"; kind = "assertion" }
            ) },
            [ordered]@{ failureSet = @([ordered]@{ file = $fileA; title = "a"; kind = "test-timeout" }) }
        )
    } | ConvertTo-Json -Depth 20 | Set-Content -LiteralPath $sourcePath -Encoding utf8

    $source = Get-FrontendVitestGroupSource -SummaryPath $sourcePath -RepoRoot $repoRoot
    Assert-Equal 2 @($source.files).Count "dynamic group union is unique"
    Assert-Equal $fileA $source.files[0] "first observed selector order is retained"
    Assert-Equal $fileB $source.files[1] "second observed selector order is retained"
    Assert-True ($source.derivation -match "ordered") "source derivation records ordered dynamic union"
    Assert-True (-not [string]::IsNullOrWhiteSpace($source.summarySha256)) "source summary hash is recorded"

    $defaultArgs = @(Get-FrontendVitestGroupCommandArguments -Mode DefaultGroup -Files $source.files)
    $singleArgs = @(Get-FrontendVitestGroupCommandArguments -Mode SingleWorkerGroup -Files $source.files)
    Assert-Equal "--prefix|frontend|run|test|--|--run|tests/unit/app-routing-lazy-load.test.tsx|tests/unit/components/ui/spinner.test.tsx" ($defaultArgs -join "|") "default command keeps frontend-relative selectors"
    Assert-Equal "--prefix|frontend|run|test|--|--run|tests/unit/app-routing-lazy-load.test.tsx|tests/unit/components/ui/spinner.test.tsx|--no-file-parallelism|--maxWorkers=1|--minWorkers=1" ($singleArgs -join "|") "single command changes only worker flags"
    $defaultSettings = Get-FrontendVitestGroupSettings -Mode DefaultGroup -Files $source.files
    $singleSettings = Get-FrontendVitestGroupSettings -Mode SingleWorkerGroup -Files $source.files
    foreach ($key in @("isolate", "pool", "order", "groupFiles", "processLifetime", "testTimeoutMs", "hookTimeoutMs")) {
        Assert-Equal ([string]$defaultSettings[$key]) ([string]$singleSettings[$key]) "group setting '$key' is comparable"
    }
    Assert-Equal "default" $defaultSettings.workerMode "default worker mode is explicit"
    Assert-Equal "single" $singleSettings.workerMode "single worker mode is diagnosis-only"

    foreach ($invalid in @("C:\outside.test.ts", "frontend/tests/../outside.test.ts", "frontend/tests/unit/not-a-test.ts", "frontend/tests/unit/missing-group.test.tsx", "frontend/tests/unit/evil`;Write-Output.test.tsx")) {
        $rejected = $false
        try { Get-FrontendVitestGroupCommandArguments -Mode DefaultGroup -Files @($invalid) | Out-Null } catch { $rejected = $true }
        Assert-True $rejected "unsafe or missing selector '$invalid' is rejected"
    }

    $syntheticChild = {
        param($RepoRoot, $Mode, $Files, $RunId, $RepeatIndex, $WatchdogSeconds)
        $start = [DateTime]::UtcNow
        [ordered]@{
            startedAtUtc = $start.ToString("o")
            completedAtUtc = [DateTime]::UtcNow.ToString("o")
            durationMs = 2
            stdout = " Test Files  2 passed (2)`n Tests  3 passed (3)"
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
    $runOutput = New-OwnedOutput
    $pointer = Invoke-FrontendVitestGroup -RepoRoot $repoRoot -OutputDirectory $runOutput -Repeat 2 -WatchdogSeconds 1 -Mode DefaultGroup -SummaryPath $sourcePath -ChildInvoker $syntheticChild | ConvertFrom-Json
    Assert-Equal 2 $pointer.recordCount "synthetic group repeats are recorded"
    Assert-Equal 2 @($pointer.runIds | Sort-Object -Unique).Count "synthetic group run ids are unique"
    Assert-True ($pointer.sequential -eq $true) "synthetic group records are sequential"
    Assert-True ($pointer.artifactIntegrity.allPresent -and $pointer.artifactIntegrity.allHashesMatch) "synthetic group artifacts are present and hashed"
    $record = Get-Content -Raw (Join-Path $runOutput "run-001\record.json") | ConvertFrom-Json
    Assert-Equal "pass" $record.classification "synthetic pass parser is retained"
    Assert-True ($record.stdout -match "artifact:") "stdout is captured as an artifact"
    Assert-True ($record.stderr -match "artifact:") "stderr is captured as an artifact"

    $timeout = ConvertFrom-FrontendVitestOutput -Stdout " FAIL  tests/unit/timeout.test.tsx > waits`n Error: Test timed out in 5000ms.`n Test Files  1 failed (1)`n Tests  1 failed (1)" -Stderr "" -ExitCode 1 -DurationMs 5000 -Phase "DefaultGroup" -RunId "timeout" -RepeatIndex 1 -DiscoveredFiles @($fileA) -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    Assert-Equal "test-timeout" $timeout.classification "timeout parser remains separate"
    $unknown = ConvertFrom-FrontendVitestOutput -Stdout "malformed" -Stderr "" -ExitCode 1 -DurationMs 1 -Phase "DefaultGroup" -RunId "unknown" -RepeatIndex 1 -DiscoveredFiles @($fileA) -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    Assert-Equal "blocked" $unknown.result "unknown parser is blocked"
    Assert-True (-not $timeout.externalWatchdogTriggered) "test timeout has no watchdog signal"

    $override = & pwsh -NoProfile -File $runnerPath -Mode DefaultGroup -GroupFiles $fileA 2>&1
    Assert-True ($LASTEXITCODE -ne 0) "DefaultGroup rejects caller GroupFiles override"
    Assert-True ((Get-Command $runnerPath).Parameters.Keys -notcontains "Command") "CLI does not expose arbitrary command input"
    Write-Output "GREEN groups contract checks passed"
    exit 0
}
finally {
    foreach ($path in $ownedOutputs | Sort-Object -Unique) {
        if (Test-Path -LiteralPath $path -PathType Container) { Remove-Item -LiteralPath $path -Recurse -Force -ErrorAction SilentlyContinue }
    }
}
