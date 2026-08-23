$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$runnerPath = Join-Path $repoRoot "scripts\diagnose-frontend-vitest-matrix.ps1"
$ownedOutputs = [System.Collections.Generic.List[string]]::new()
function Assert-True {
    param(
        [Parameter(Mandatory = $true)][bool]$Condition,
        [Parameter(Mandatory = $true)][string]$Message
    )
    if (-not $Condition) {
        throw "ASSERTION FAILED: $Message"
    }
}
function Assert-Equal {
    param(
        [Parameter(Mandatory = $true)]$Expected,
        [Parameter(Mandatory = $true)]$Actual,
        [Parameter(Mandatory = $true)][string]$Message
    )
    if ($Expected -ne $Actual) {
        throw "ASSERTION FAILED: $Message (expected '$Expected', actual '$Actual')"
    }
}
function Invoke-Runner {
    param(
        [string]$Mode = "Environment",
        [int]$Repeat = 1,
        [string]$OutputPath,
        [string]$SummaryPath,
        [string[]]$GroupFiles,
        [int]$WatchdogSeconds = 600
    )
    $arguments = @("-NoProfile", "-File", $runnerPath, "-Mode", $Mode, "-Repeat", $Repeat, "-WatchdogSeconds", $WatchdogSeconds)
    if ($PSBoundParameters.ContainsKey("OutputPath")) {
        $arguments += @("-OutputPath", $OutputPath)
    }
    if ($PSBoundParameters.ContainsKey("SummaryPath")) {
        $arguments += @("-SummaryPath", $SummaryPath)
    }
    if ($PSBoundParameters.ContainsKey("GroupFiles")) {
        $arguments += @("-GroupFiles") + $GroupFiles
    }
    $output = & pwsh @arguments 2>&1
    [pscustomobject]@{
        ExitCode = $LASTEXITCODE
        Output = (@($output) -join [Environment]::NewLine)
    }
}
function New-TestOutputPath {
    $path = Join-Path ([IO.Path]::GetTempPath()) ("frontend-vitest-contract-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $path -Force | Out-Null
    $ownedOutputs.Add($path)
    return $path
}
try {
    Assert-True (Test-Path -LiteralPath $runnerPath -PathType Leaf) "runner script exists"
    $defaultOne = Invoke-Runner
    Assert-Equal 0 $defaultOne.ExitCode "default Environment mode succeeds"
    $defaultRecordOne = $defaultOne.Output | ConvertFrom-Json
    $ownedOutputs.Add([string]$defaultRecordOne.outputDirectory)
    Assert-True (Test-Path -LiteralPath $defaultRecordOne.outputDirectory -PathType Container) "default output directory exists"
    $defaultTwo = Invoke-Runner
    Assert-Equal 0 $defaultTwo.ExitCode "second default Environment mode succeeds"
    $defaultRecordTwo = $defaultTwo.Output | ConvertFrom-Json
    $ownedOutputs.Add([string]$defaultRecordTwo.outputDirectory)
    Assert-True ($defaultRecordOne.outputDirectory -ne $defaultRecordTwo.outputDirectory) "default output directories are unique"
    $explicitPath = New-TestOutputPath
    $explicit = Invoke-Runner -OutputPath $explicitPath
    Assert-Equal 0 $explicit.ExitCode "explicit safe output succeeds"
    $record = $explicit.Output | ConvertFrom-Json
    Assert-Equal $explicitPath $record.outputDirectory "explicit output path is honored"
    Assert-True (Test-Path -LiteralPath (Join-Path $explicitPath "environment.json") -PathType Leaf) "environment evidence file exists"
    $requiredFields = @(
        "schemaVersion", "runId", "mode", "repeatIndex", "command", "settings",
        "startedAtUtc", "durationMs", "discovered", "failures", "resources", "result",
        "artifactPaths", "stdout", "stderr", "outputDirectory"
    )
    foreach ($field in $requiredFields) {
        Assert-True ($null -ne $record.PSObject.Properties[$field]) "schema field '$field' exists"
    }
    Assert-Equal 1 $record.schemaVersion "schema version is one"
    Assert-Equal "Environment" $record.mode "mode is Environment"
    Assert-True ([Guid]::TryParseExact([string]$record.runId, "N", [ref]([Guid]::Empty))) "runId is a GUID"
    Assert-True ($record.discovered.files.Count -gt 0) "dynamic discovery has current candidate files"; Assert-True ($record.dynamicFailureSet.fixedIdentity -eq $false) "failure identity is not fixed"; Assert-True ([string]$record.discovered.identityStatement -match "not a fixed") "identity statement rejects fixed timeout identity"
    Assert-Equal "not-applicable/no-child" $record.settings.processLifetime "Environment has no child process lifetime"; Assert-Equal "not-applicable/no-child" $record.settings.externalWatchdogSeconds "Environment does not execute a watchdog"; Assert-Equal "not-applicable/no-child" $record.settings.testTimeoutMs "Environment does not execute Vitest timeouts"; Assert-Equal "not-applicable/no-child" $record.settings.hookTimeoutMs "Environment does not execute hooks"
    Assert-Equal "unavailable" $record.resources.workerCount "unsupported worker metric is unavailable"; Assert-Equal "unavailable" $record.resources.heapUsedBytes "unsupported heap metric is unavailable"; Assert-Equal "not-applicable/no-child" $record.resources.processExitCode "Environment has no child exit code"; Assert-Equal "not-applicable/no-child" $record.resources.cleanupOutcome "Environment cleanup is not applicable"
    $currentFiles = @(Get-ChildItem -LiteralPath (Join-Path $repoRoot "frontend\tests") -Recurse -File | Where-Object { $_.Name -match '\.test\.tsx?$' } | ForEach-Object { [IO.Path]::GetRelativePath($repoRoot, $_.FullName).Replace('\', '/') } | Sort-Object -Unique)
    Assert-Equal ($currentFiles -join '|') (@($record.discovered.files) -join '|') "Environment inventory equals current filesystem inventory"; Assert-True ($record.artifactPaths.Count -ge 3) "stdout/stderr and JSON artifacts are listed"
    $invalidRepeat = Invoke-Runner -Repeat 0
    Assert-True ($invalidRepeat.ExitCode -ne 0) "zero Repeat fails loudly"
    $invalidWatchdog = Invoke-Runner -WatchdogSeconds 0
    Assert-True ($invalidWatchdog.ExitCode -ne 0) "zero WatchdogSeconds fails loudly"
    $invalidMode = Invoke-Runner -Mode InvalidMode
    Assert-True ($invalidMode.ExitCode -ne 0) "invalid Mode fails loudly"
    $normalSummaryInput = Join-Path ([IO.Path]::GetTempPath()) "frontend-vitest-normal-summary-input.json"
    $normalSummaryRejected = Invoke-Runner -Mode NormalFull -SummaryPath $normalSummaryInput
    Assert-True ($normalSummaryRejected.ExitCode -ne 0) "NormalFull rejects SummaryPath handoff input"
    $environmentSummaryRejected = Invoke-Runner -Mode Environment -SummaryPath $normalSummaryInput
    Assert-True ($environmentSummaryRejected.ExitCode -ne 0) "Environment rejects SummaryPath handoff input"
    $oneVariableOutput = New-TestOutputPath
    $oneVariable = Invoke-Runner -Mode OneVariable -OutputPath $oneVariableOutput
    Assert-True ($oneVariable.ExitCode -ne 0) "OneVariable remains a non-zero blocked branch"
    $oneVariableEvidencePath = Join-Path $oneVariableOutput "one-variable.json"
    Assert-True (Test-Path -LiteralPath $oneVariableEvidencePath -PathType Leaf) "OneVariable blocked JSON is persisted"
    $oneVariableEvidence = Get-Content -Raw -LiteralPath $oneVariableEvidencePath | ConvertFrom-Json
    $oneVariableFields = @("schemaVersion", "runId", "mode", "repeatIndex", "command", "settings", "startedAtUtc", "completedAtUtc", "durationMs", "discovered", "failures", "resources", "classification", "result", "stdout", "stderr", "artifactPaths", "testTimeout", "outputDirectory"); foreach ($field in $oneVariableFields) { Assert-True ($null -ne $oneVariableEvidence.PSObject.Properties[$field]) "OneVariable schema field '$field' exists" }
    Assert-Equal 1 $oneVariableEvidence.schemaVersion "OneVariable schema version"; Assert-Equal "OneVariable" $oneVariableEvidence.mode "OneVariable mode"; Assert-Equal "unknown" $oneVariableEvidence.classification "OneVariable is not classified without a hypothesis"; Assert-Equal "blocked" $oneVariableEvidence.result "OneVariable is fail-closed"
    Assert-True ($oneVariableEvidence.settings.unchanged -eq $true) "OneVariable settings remain unchanged"; Assert-Equal "none" $oneVariableEvidence.settings.variable "OneVariable does not invent a setting"; Assert-Equal $oneVariableOutput $oneVariableEvidence.outputDirectory "OneVariable records owned output directory"; Assert-True (@($oneVariableEvidence.artifactPaths).Count -ge 3) "OneVariable records raw artifacts"
    Assert-Equal "unavailable" $oneVariableEvidence.failures[0].file "OneVariable failure file is unavailable"; Assert-Equal "unavailable" $oneVariableEvidence.failures[0].title "OneVariable failure title is unavailable"; Assert-Equal "unknown" $oneVariableEvidence.failures[0].kind "OneVariable failure is unknown"; Assert-Equal "not-applicable/no-child" $oneVariableEvidence.resources.cleanupOutcome "OneVariable records no-child cleanup"
    $oneVariableHandoffOutput = New-TestOutputPath
    $oneVariableHandoff = Invoke-Runner -Mode OneVariable -OutputPath $oneVariableHandoffOutput -SummaryPath $normalSummaryInput
    Assert-True ($oneVariableHandoff.ExitCode -ne 0) "OneVariable rejects SummaryPath handoff"
    Assert-Equal 0 @(Get-ChildItem -LiteralPath $oneVariableHandoffOutput -Force).Count "rejected OneVariable handoff writes no evidence"
    $groupInEnvironment = Invoke-Runner -GroupFiles "frontend/tests/example.test.ts"
    Assert-True ($groupInEnvironment.ExitCode -ne 0) "Environment rejects unrelated GroupFiles"
    $repoRootOutput = Invoke-Runner -OutputPath $repoRoot
    Assert-True ($repoRootOutput.ExitCode -ne 0) "repo root output is rejected"
    $diskRoot = [IO.Path]::GetPathRoot($repoRoot)
    $diskRootOutput = Invoke-Runner -OutputPath $diskRoot
    Assert-True ($diskRootOutput.ExitCode -ne 0) "disk root output is rejected"
    $traversalOutput = Invoke-Runner -OutputPath (Join-Path ([IO.Path]::GetTempPath()) "frontend-vitest-contract-..\escape")
    Assert-True ($traversalOutput.ExitCode -ne 0) "path traversal output is rejected"
    $nonEmptyPath = New-TestOutputPath
    New-Item -ItemType File -Path (Join-Path $nonEmptyPath "unrelated.txt") -Force | Out-Null
    $nonEmptyOutput = Invoke-Runner -OutputPath $nonEmptyPath
    Assert-True ($nonEmptyOutput.ExitCode -ne 0) "existing non-empty output is rejected"
    . $runnerPath
    $successOutput = @'
 Test Files  2 passed (2)
      Tests  3 passed (3)
   Duration  123ms
'@
    $successParsed = ConvertFrom-FrontendVitestOutput -Stdout $successOutput -Stderr "" -ExitCode 0 -DurationMs 123 -Phase "NormalFull" -RunId "synthetic-pass" -RepeatIndex 1 -DiscoveredFiles @(
        "frontend/tests/unit/synthetic-a.test.ts",
        "frontend/tests/unit/synthetic-b.test.ts"
    ) -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    Assert-Equal "pass" $successParsed.result "synthetic successful output is pass"
    Assert-Equal 2 $successParsed.discovered.files.Count "synthetic success retains discovered file inventory"
    Assert-Equal 3 $successParsed.discovered.tests "synthetic success parses test count"
    Assert-Equal 0 $successParsed.discovered.failed "synthetic success parses failed test count"
    Assert-True ($successParsed.testTimeout.observed -eq $false) "synthetic success has no test timeout"
    Assert-True ($successParsed.externalWatchdogTriggered -eq $false) "synthetic success has no external watchdog"
    $timeoutOutput = @'
 FAIL  tests/unit/synthetic-timeout.test.tsx > timeout suite > waits
 ❯ tests/unit/synthetic-timeout.test.tsx > timeout suite > waits
     Error: Test timed out in 5000ms.
     Error: Test timed out in 5000ms.
 Test Files  1 failed (1)
      Tests  1 failed (1)
'@
    $timeoutParsed = ConvertFrom-FrontendVitestOutput -Stdout $timeoutOutput -Stderr "" -ExitCode 1 -DurationMs 5100 -Phase "NormalFull" -RunId "synthetic-timeout" -RepeatIndex 1 -DiscoveredFiles @("frontend/tests/unit/synthetic-timeout.test.tsx") -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    Assert-Equal "test-timeout" $timeoutParsed.classification "synthetic timeout classification"
    Assert-True ($timeoutParsed.testTimeout.observed -eq $true) "synthetic timeout is separate from watchdog"
    Assert-True ($timeoutParsed.externalWatchdogTriggered -eq $false) "synthetic timeout has no external watchdog"
    Assert-Equal "frontend/tests/unit/synthetic-timeout.test.tsx" $timeoutParsed.failures[0].file "timeout failure path is project-root-relative"
    Assert-Equal 5000 $timeoutParsed.failures[0].timeoutMs "timeout milliseconds are parsed"
    Assert-Equal 1 $timeoutParsed.discovered.timedOut "duplicate timeout markers count one timeout test"
    Assert-Equal 1 $timeoutParsed.testTimeout.count "duplicate timeout markers count one timeout test"
    Assert-Equal 2 $timeoutParsed.testTimeout.markerCount "raw timeout marker count is preserved"
    $crossFileOutput = @'
 FAIL  tests/unit/title-reset-a.test.tsx > suite > inherited-title
     Error: Test timed out in 5000ms.
 FAIL  tests/unit/title-reset-b.test.tsx
     Error: Test timed out in 5000ms.
 Test Files  2 failed (2)
      Tests  2 failed (2)
'@
    $crossFileParsed = ConvertFrom-FrontendVitestOutput -Stdout $crossFileOutput -Stderr "" -ExitCode 1 -DurationMs 5100 -Phase "NormalFull" -RunId "synthetic-title-reset" -RepeatIndex 1 -DiscoveredFiles @("frontend/tests/unit/title-reset-a.test.tsx", "frontend/tests/unit/title-reset-b.test.tsx") -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    $crossFileB = @($crossFileParsed.failures | Where-Object file -eq "frontend/tests/unit/title-reset-b.test.tsx")
    Assert-Equal 1 $crossFileB.Count "cross-file timeout creates a second failure"
    Assert-Equal "unavailable" $crossFileB[0].title "cross-file timeout does not inherit previous title"
    $driftA = ConvertFrom-FrontendVitestOutput -Stdout " FAIL  tests/unit/drift-a.test.ts > suite > title-a`n Test Files  1 failed (1)`n Tests  1 failed (1)" -Stderr "" -ExitCode 1 -DurationMs 10 -Phase "NormalFull" -RunId "synthetic-drift-a" -RepeatIndex 1 -DiscoveredFiles @("frontend/tests/unit/drift-a.test.ts") -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    $driftB = ConvertFrom-FrontendVitestOutput -Stdout " FAIL  tests/unit/drift-b.test.ts > suite > title-b`n Test Files  1 failed (1)`n Tests  1 failed (1)" -Stderr "" -ExitCode 1 -DurationMs 11 -Phase "NormalFull" -RunId "synthetic-drift-b" -RepeatIndex 2 -DiscoveredFiles @("frontend/tests/unit/drift-b.test.ts") -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    Assert-True ($driftA.failures[0].file -ne $driftB.failures[0].file) "failure identity drift is retained"
    Assert-True ($driftA.runId -ne $driftB.runId) "drift records retain independent run ids"
    $watchdogParsed = ConvertFrom-FrontendVitestOutput -Stdout "runner output before watchdog" -Stderr "runner stderr" -ExitCode "unavailable" -DurationMs 600000 -Phase "NormalFull" -RunId "synthetic-watchdog" -RepeatIndex 1 -DiscoveredFiles @("frontend/tests/unit/synthetic-hang.test.ts") -ExternalWatchdogTriggered:$true -WatchdogSeconds 600
    Assert-Equal "worker-hang" $watchdogParsed.classification "watchdog classification is worker-hang"
    Assert-True ($watchdogParsed.externalWatchdogTriggered -eq $true) "watchdog signal is preserved"
    Assert-True ($watchdogParsed.testTimeout.observed -eq $false) "watchdog is separate from test timeout"
    Assert-Equal "runner output before watchdog" $watchdogParsed.stdout "watchdog stdout is preserved"
    Assert-Equal "runner stderr" $watchdogParsed.stderr "watchdog stderr is preserved"
    $malformedParsed = ConvertFrom-FrontendVitestOutput -Stdout "not a Vitest summary" -Stderr "malformed stderr" -ExitCode 1 -DurationMs 12 -Phase "NormalFull" -RunId "synthetic-malformed" -RepeatIndex 1 -DiscoveredFiles @("frontend/tests/unit/synthetic-malformed.test.ts") -ExternalWatchdogTriggered:$false -WatchdogSeconds 600
    Assert-Equal "unknown" $malformedParsed.classification "malformed output is unknown"
    Assert-Equal "blocked" $malformedParsed.result "malformed output is blocked"
    Assert-Equal "malformed stderr" $malformedParsed.stderr "malformed stderr is preserved"
    $normalFailureOutput = @'
 FAIL  tests/unit/synthetic-assertion.test.tsx > assertion suite > fails
     AssertionError: expected true to be false
 Test Files  1 failed (1)
      Tests  1 failed (1)
'@
    $syntheticNormalChild = {
        [ordered]@{
            startedAtUtc = [DateTime]::UtcNow.ToString("o")
            completedAtUtc = [DateTime]::UtcNow.ToString("o")
            durationMs = 8
            stdout = $normalFailureOutput
            stderr = ""
            exitCode = 1
            runnerPid = "synthetic"
            childPids = "synthetic"
            peakWorkingSetBytes = "unavailable"
            externalWatchdogTriggered = $false
            cleanupOutcome = "not-required"
            cleanupError = ""
            startError = ""
        }
    }
    $normalOutputPath = New-TestOutputPath
    $normalFailureThrown = $false
    try {
        Invoke-FrontendVitestNormalFull -RepoRoot $repoRoot -OutputDirectory $normalOutputPath -Repeat 1 -WatchdogSeconds 600 -DiscoveredFiles @("frontend/tests/unit/synthetic-assertion.test.tsx") -ChildInvoker $syntheticNormalChild | Out-Null
    } catch {
        $normalFailureThrown = $true
    }
    Assert-True $normalFailureThrown "classified NormalFull failure is a non-zero gate"
    $normalSummaryPath = Join-Path $normalOutputPath "normal-full-summary.json"
    Assert-True (Test-Path -LiteralPath $normalSummaryPath -PathType Leaf) "NormalFull summary survives classified failure"
    $normalSummary = Get-Content -Raw -LiteralPath $normalSummaryPath | ConvertFrom-Json
    Assert-Equal "failed" $normalSummary.result "classified NormalFull failure is failed"
    Assert-Equal "assertion" $normalSummary.records[0].classification "NormalFull preserves assertion classification"
    Assert-Equal "failed" $normalSummary.records[0].result "NormalFull record is non-pass"
    Assert-Equal "unmodified-default" $normalSummary.settings.isolate "NormalFull isolate setting is unchanged"
    Assert-Equal "unmodified-default" $normalSummary.settings.pool "NormalFull pool setting is unchanged"
    Assert-Equal "unmodified-default" $normalSummary.settings.workers "NormalFull worker setting is unchanged"
    Assert-Equal "unmodified-default" $normalSummary.settings.order "NormalFull order setting is unchanged"
    Assert-Equal 0 $normalSummary.settings.groupFiles.Count "NormalFull has no groups"
    Assert-Equal "one-owned-process-per-run" $normalSummary.settings.processLifetime "NormalFull process lifetime is explicit"
    Assert-Equal 5000 $normalSummary.settings.testTimeoutMs "NormalFull test timeout is recorded"
    Assert-Equal 10000 $normalSummary.settings.hookTimeoutMs "NormalFull hook timeout is recorded"
    Assert-True ($normalSummary.settings.unchanged -eq $true) "NormalFull settings record no configuration change"
    Assert-True ($normalSummary.records[0].artifactPaths.Count -ge 3) "NormalFull classified failure retains artifacts"
    $runtimeChild = { param($RepoRoot, $WatchdogSeconds) $child = & $syntheticNormalChild -RepoRoot $RepoRoot -WatchdogSeconds $WatchdogSeconds; $child.runtimeError = "synthetic runtime failure"; $child }
    $runtimeOutputPath = New-TestOutputPath; $runtimeFailureThrown = $false
    try { Invoke-FrontendVitestNormalFull -RepoRoot $repoRoot -OutputDirectory $runtimeOutputPath -Repeat 1 -WatchdogSeconds 600 -DiscoveredFiles @("frontend/tests/unit/synthetic-assertion.test.tsx") -ChildInvoker $runtimeChild | Out-Null } catch { $runtimeFailureThrown = $true }
    $runtimeSummary = Get-Content -Raw (Join-Path $runtimeOutputPath "normal-full-summary.json") | ConvertFrom-Json
    Assert-True ($runtimeFailureThrown -and $runtimeSummary.records[0].classification -eq "assertion") "runtimeError keeps assertion evidence and gates nonzero"
    Assert-Equal "blocked" $runtimeSummary.result "NormalFull runtimeError is blocked"
    $runnerParameters = (Get-Command $runnerPath).Parameters.Keys
    Assert-True ($runnerParameters -notcontains "Command") "CLI does not expose arbitrary command input"
    $helperPath = Join-Path $repoRoot "scripts\lib\FrontendVitestMatrix.psm1"
    Assert-True ((Get-Content -Raw $helperPath) -match "npm --prefix frontend run test -- --run") "NormalFull uses the fixed normal command"
    $classificationCases = @(
        @{ Name = "assertion"; Signals = @{ assertion = $true } },
        @{ Name = "test-timeout"; Signals = @{ testTimeout = $true } },
        @{ Name = "worker-hang"; Signals = @{ workerHang = $true } },
        @{ Name = "resource leak"; Signals = @{ resourceLeak = $true } },
        @{ Name = "environment"; Signals = @{ environment = $true } }
    )
    foreach ($case in $classificationCases) {
        $synthetic = @{
            file = "frontend/tests/synthetic.test.ts"
            title = "synthetic $($case.Name)"
            message = "synthetic evidence"
            phase = "synthetic"
            processExitCode = 1
            stdout = "synthetic stdout"
            stderr = "synthetic stderr"
            signals = $case.Signals
        }
        $classified = New-FrontendVitestFailureRecord -Evidence $synthetic
        Assert-Equal $case.Name $classified.kind "synthetic $($case.Name) classification"
        Assert-Equal "failed" $classified.result "synthetic $($case.Name) is non-pass"
        Assert-True ($classified.blocked -eq $false) "synthetic $($case.Name) has one classification"
        Assert-Equal 1 @($classified.kind).Count "synthetic $($case.Name) has exactly one kind"
        Assert-Equal "synthetic stdout" $classified.stdout "synthetic stdout is preserved"
        Assert-Equal "synthetic stderr" $classified.stderr "synthetic stderr is preserved"
        Assert-Equal 1 $classified.exitCode "synthetic exit code is preserved"
        Assert-Equal "synthetic" $classified.phase "synthetic phase is preserved"
    }
    $unknown = New-FrontendVitestFailureRecord -Evidence @{ file = "frontend/tests/synthetic.test.ts"; title = "unknown"; signals = @{} }
    Assert-Equal "unknown" $unknown.kind "insufficient evidence is unknown"; Assert-Equal "blocked" $unknown.result "unknown result is blocked"; Assert-True ($unknown.blocked -eq $true) "unknown classification is fail-closed"
    $conflict = New-FrontendVitestFailureRecord -Evidence @{ signals = @{ assertion = $true; testTimeout = $true } }
    Assert-Equal "unknown" $conflict.kind "conflicting signals are unknown"; Assert-Equal "blocked" $conflict.result "conflicting signals are blocked"
    $pass = New-FrontendVitestFailureRecord -Evidence @{ pass = $true; signals = @{} }
    Assert-Equal "pass" $pass.kind "explicit pass is pass"; Assert-Equal "pass" $pass.result "explicit pass result is pass"
    Assert-True ($pass.blocked -eq $false) "explicit pass is not blocked"
    $orderedAssertion = New-FrontendVitestFailureRecord -Evidence ([ordered]@{ signals = [ordered]@{ assertion = $true }; processExitCode = 1 })
    $orderedPass = New-FrontendVitestFailureRecord -Evidence ([ordered]@{ pass = $true; signals = [ordered]@{} })
    $orderedConflict = New-FrontendVitestFailureRecord -Evidence ([ordered]@{ signals = [ordered]@{ assertion = $true; testTimeout = $true } })
    Assert-Equal "assertion" $orderedAssertion.kind "OrderedDictionary assertion is classified"; Assert-Equal "pass" $orderedPass.kind "OrderedDictionary pass is classified"; Assert-Equal "unknown" $orderedConflict.kind "OrderedDictionary conflict is blocked"
    $runnerText = Get-Content -Raw $runnerPath
    Assert-True ($runnerText -notmatch "(?i)Invoke-Expression|\bStop-Process\b|\btaskkill(?:\.exe)?\b|\bkill\.exe\b") "runner has no command injection or global kill primitive"
    Write-Output "GREEN contract checks passed"
    exit 0
}
finally {
    foreach ($path in $ownedOutputs | Sort-Object -Unique) {
        if ([string]::IsNullOrWhiteSpace($path)) {
            continue
        }
        if (Test-Path -LiteralPath $path -PathType Container) {
            Remove-Item -LiteralPath $path -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}
