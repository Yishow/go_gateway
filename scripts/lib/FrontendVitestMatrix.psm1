Set-StrictMode -Version Latest

Import-Module (Join-Path $PSScriptRoot "FrontendVitestEvidence.psm1") -Force -Global
Import-Module (Join-Path $PSScriptRoot "FrontendVitestProcess.psm1") -Force -Global

$script:NormalVitestCommand = "npm --prefix frontend run test -- --run"
$script:NormalVitestArguments = @("--prefix", "frontend", "run", "test", "--", "--run")

function Get-FrontendVitestValue {
    param([AllowNull()][object]$Object, [Parameter(Mandatory = $true)][string]$Name)
    if ($null -eq $Object) { return $null }
    if ($Object -is [System.Collections.IDictionary]) {
        if ($Object.Contains($Name)) { return $Object[$Name] }
        return $null
    }
    $property = $Object.PSObject.Properties[$Name]
    if ($null -ne $property) { return $property.Value }
    return $null
}

function Get-FrontendVitestCurrentDiscovery {
    param([Parameter(Mandatory = $true)][string]$RepoRoot)
    $root = [IO.Path]::GetFullPath($RepoRoot)
    $testsRoot = Join-Path $root "frontend\tests"
    if (-not (Test-Path -LiteralPath $testsRoot -PathType Container)) { throw "frontend tests directory is missing" }
    $files = @(
        Get-ChildItem -LiteralPath $testsRoot -Recurse -File -ErrorAction Stop |
            Where-Object { $_.Name -match '\.test\.tsx?$' } |
            ForEach-Object { [IO.Path]::GetRelativePath($root, $_.FullName).Replace('\', '/') } |
            Sort-Object -Unique
    )
    if ($files.Count -eq 0) { throw "current frontend test inventory is empty" }
    [ordered]@{
        files = $files
        fileCount = $files.Count
        tests = "unavailable"
        passed = "unavailable"
        failed = "unavailable"
        timedOut = "unavailable"
        identityStatement = "Current filesystem discovery under frontend/tests; this inventory is not a fixed timeout identity or archived evidence."
        source = "current-filesystem"
    }
}

function ConvertTo-MatrixRelativePath {
    param([string]$Path, [string]$RepoRoot)
    if ([string]::IsNullOrWhiteSpace($Path)) { return "unavailable" }
    $normalized = $Path.Trim().Replace('\', '/')
    if ($normalized -match '^(?:[A-Za-z]:/|/|\\)') {
        try { $normalized = [IO.Path]::GetRelativePath($RepoRoot, [IO.Path]::GetFullPath($Path)).Replace('\', '/') } catch { return "unavailable" }
    } elseif ($normalized -match '^tests/') { $normalized = "frontend/$normalized" }
    if ($normalized -notmatch '^frontend/tests/.+\.test\.tsx?$') { return "unavailable" }
    return $normalized
}

function Get-VitestSummaryCount {
    param([string]$Text, [string]$Label, [string]$State)
    $line = [regex]::Match($Text, "(?im)^\s*$([regex]::Escape($Label))\s+(?<body>.+)$")
    if (-not $line.Success) { return "unavailable" }
    $match = [regex]::Match($line.Groups["body"].Value, "(?<value>\d+)\s+$([regex]::Escape($State))")
    if ($match.Success) { return [int]$match.Groups["value"].Value }
    if ($line.Groups["body"].Value -match '\(\d+\)' -or $line.Groups["body"].Value -match '\d+') { return 0 }
    return "unavailable"
}

function Get-VitestTotalCount {
    param([string]$Text, [string]$Label)
    $line = [regex]::Match($Text, "(?im)^\s*$([regex]::Escape($Label))\s+(?<body>.+)$")
    if (-not $line.Success) { return "unavailable" }
    $match = [regex]::Match($line.Groups["body"].Value, '\((?<value>\d+)\)')
    if ($match.Success) { return [int]$match.Groups["value"].Value }
    $numbers = [regex]::Matches($line.Groups["body"].Value, '(?<value>\d+)')
    if ($numbers.Count -gt 0) { return [int]$numbers[$numbers.Count - 1].Groups["value"].Value }
    return "unavailable"
}

function ConvertFrom-FrontendVitestOutput {
    param(
        [AllowNull()][string]$Stdout, [AllowNull()][string]$Stderr, [AllowNull()][object]$ExitCode,
        [Parameter(Mandatory = $true)][int]$DurationMs, [Parameter(Mandatory = $true)][string]$Phase,
        [Parameter(Mandatory = $true)][string]$RunId, [Parameter(Mandatory = $true)][int]$RepeatIndex, [Parameter(Mandatory = $true)][string[]]$DiscoveredFiles,
        [Parameter(Mandatory = $true)][bool]$ExternalWatchdogTriggered,
        [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [string]$RepoRoot = (Split-Path -Parent (Split-Path -Parent $PSScriptRoot))
    )
    $raw = (([string]$Stdout) + "`n" + ([string]$Stderr)) -replace "`e\[[0-9;]*m", ""
    $records = [System.Collections.Generic.List[object]]::new(); $indexes = @{}
    $currentFile = $null; $currentTitle = $null; $timeoutValues = [System.Collections.Generic.List[int]]::new()
    foreach ($line in ($raw -split "`r?`n")) {
        $pathMatch = [regex]::Match($line, '(?<path>(?:frontend[\\/]?)?tests[\\/][^\s>]+\.test\.tsx?)')
        if ($pathMatch.Success) {
            $previousFile = $currentFile
            $currentFile = ConvertTo-MatrixRelativePath -Path $pathMatch.Groups["path"].Value -RepoRoot $RepoRoot
            if ($currentFile -ne $previousFile) { $currentTitle = $null }
            $tail = $line.Substring($pathMatch.Index + $pathMatch.Length)
            $titleMatch = [regex]::Match($tail, '^\s*>\s*(?<title>.+)$')
            if ($titleMatch.Success -and $line -match '^\s*(?:FAIL|[❯×✗])') {
                $currentTitle = [regex]::Replace($titleMatch.Groups["title"].Value.Trim(), '\s+\d+\s*ms$', '')
                $key = "$currentFile|$currentTitle"
                if (-not $indexes.ContainsKey($key)) { $indexes[$key] = $records.Count; $records.Add([ordered]@{ file = $currentFile; title = $currentTitle; kind = "assertion"; message = "unavailable"; timeoutMs = "unavailable"; phase = $Phase }) }
            }
        }
        $markerMatch = [regex]::Match($line, '^\s*[×✗]\s+(?<title>.+)$')
        if ($markerMatch.Success -and $null -ne $currentFile) {
            $currentTitle = [regex]::Replace($markerMatch.Groups["title"].Value.Trim(), '\s+\d+\s*ms$', '')
            $key = "$currentFile|$currentTitle"
            if (-not $indexes.ContainsKey($key)) { $indexes[$key] = $records.Count; $records.Add([ordered]@{ file = $currentFile; title = $currentTitle; kind = "assertion"; message = "unavailable"; timeoutMs = "unavailable"; phase = $Phase }) }
        }
        $timeoutMatch = [regex]::Match($line, '(?i)timed\s+out\s+in\s+(?<ms>\d+)\s*ms')
        if ($timeoutMatch.Success) {
            $timeoutValues.Add([int]$timeoutMatch.Groups["ms"].Value); $target = $null
            if ($null -ne $currentFile -and $null -ne $currentTitle -and $indexes.ContainsKey("$currentFile|$currentTitle")) { $target = $records[$indexes["$currentFile|$currentTitle"]] }
            if ($null -eq $target -and $null -ne $currentFile) { $sameFile = @($records | Where-Object { $_["file"] -eq $currentFile }); if ($sameFile.Count -eq 1) { $target = $sameFile[0] } }
            if ($null -ne $target) { $target["kind"] = "test-timeout"; $target["timeoutMs"] = [int]$timeoutMatch.Groups["ms"].Value; $target["message"] = $line.Trim() }
            elseif ($null -ne $currentFile) { $key = "$currentFile|unavailable"; if (-not $indexes.ContainsKey($key)) { $indexes[$key] = $records.Count; $records.Add([ordered]@{ file = $currentFile; title = if ($null -ne $currentTitle) { $currentTitle } else { "unavailable" }; kind = "test-timeout"; message = $line.Trim(); timeoutMs = [int]$timeoutMatch.Groups["ms"].Value; phase = $Phase }) } }
        }
        if ($records.Count -gt 0 -and $line -match '^\s*(?:→|Error:|AssertionError|Expected)') { $last = $records[$records.Count - 1]; if ([string]$last["message"] -eq "unavailable") { $last["message"] = $line.Trim() } }
    }
    $fileTotal = Get-VitestTotalCount $raw "Test Files"; $testTotal = Get-VitestTotalCount $raw "Tests"
    $filePassed = Get-VitestSummaryCount $raw "Test Files" "passed"; $fileFailed = Get-VitestSummaryCount $raw "Test Files" "failed"
    $testPassed = Get-VitestSummaryCount $raw "Tests" "passed"; $testFailed = Get-VitestSummaryCount $raw "Tests" "failed"
    $uniqueTimeoutCount = @($records | Where-Object { $_["kind"] -eq "test-timeout" }).Count
    $categories = [System.Collections.Generic.List[string]]::new()
    if ($ExternalWatchdogTriggered) { $categories.Add("worker-hang") }
    if ($uniqueTimeoutCount -gt 0) { $categories.Add("test-timeout") }
    if ($raw -match '(?i)spawn\s+EPERM|environment\s+boundary') { $categories.Add("environment") }
    $assertionObserved = @($records | Where-Object { $_["kind"] -eq "assertion" }).Count -gt 0
    $failedSummary = (($testFailed -is [int]) -and $testFailed -gt 0) -or (($fileFailed -is [int]) -and $fileFailed -gt 0)
    if ($assertionObserved -or ($failedSummary -and $uniqueTimeoutCount -eq 0 -and -not $ExternalWatchdogTriggered)) { $categories.Add("assertion") }
    $classification = "unknown"
    if ($categories.Count -eq 1) { $classification = $categories[0] }
    elseif ($categories.Count -eq 0 -and $ExitCode -eq 0 -and $fileTotal -ne "unavailable" -and $testTotal -ne "unavailable") { $classification = "pass" }
    if ($classification -eq "worker-hang" -and $records.Count -eq 0) { $records.Add([ordered]@{ file = "unavailable"; title = "external watchdog"; kind = "worker-hang"; message = "external watchdog boundary reached"; timeoutMs = "unavailable"; phase = $Phase }) }
    if ($classification -eq "assertion" -and $records.Count -eq 0) { $records.Add([ordered]@{ file = "unavailable"; title = "unavailable"; kind = "assertion"; message = "Vitest summary reported failed tests"; timeoutMs = "unavailable"; phase = $Phase }) }
    if ($classification -eq "environment" -and $records.Count -eq 0) { $records.Add([ordered]@{ file = "unavailable"; title = "environment"; kind = "environment"; message = "environment boundary prevented a comparable run"; timeoutMs = "unavailable"; phase = $Phase }) }
    if ($classification -eq "unknown" -and $records.Count -eq 0) { $records.Add([ordered]@{ file = "unavailable"; title = "unavailable"; kind = "unknown"; message = "Vitest output was insufficient to classify"; timeoutMs = "unavailable"; phase = $Phase }) }
    $timeoutMs = if ($timeoutValues.Count -gt 0 -and @($timeoutValues | Select-Object -Unique).Count -eq 1) { $timeoutValues[0] } else { "unavailable" }
    $result = if ($classification -eq "pass") { "pass" } elseif ($classification -eq "unknown") { "blocked" } else { "failed" }
    [ordered]@{
        schemaVersion = 1; runId = $RunId; mode = $Phase; repeatIndex = $RepeatIndex; command = $script:NormalVitestCommand; startedAtUtc = "unavailable"; durationMs = $DurationMs
        discovered = [ordered]@{ files = @($DiscoveredFiles); fileCount = @($DiscoveredFiles).Count; tests = $testTotal; passed = $testPassed; failed = $testFailed; timedOut = $uniqueTimeoutCount }
        discoveredFileSummary = [ordered]@{ total = $fileTotal; passed = $filePassed; failed = $fileFailed }; failures = @($records)
        resources = [ordered]@{ runnerPid = "unavailable"; childPids = "unavailable"; workerCount = "unavailable"; peakWorkingSetBytes = "unavailable"; heapUsedBytes = "unavailable"; processExitCode = $ExitCode; externalWatchdogTriggered = $ExternalWatchdogTriggered }
        result = $result; classification = $classification; externalWatchdogTriggered = $ExternalWatchdogTriggered
        testTimeout = [ordered]@{ observed = ($uniqueTimeoutCount -gt 0); count = $uniqueTimeoutCount; markerCount = $timeoutValues.Count; timeoutMs = $timeoutMs }
        stdout = [string]$Stdout; stderr = [string]$Stderr; artifactPaths = @(); watchdogSeconds = $WatchdogSeconds
    }
}

function Resolve-FrontendVitestClassification {
    param([Parameter(Mandatory = $true)][AllowNull()][object]$Evidence)
    $signals = Get-FrontendVitestValue $Evidence "signals"; $aliases = [ordered]@{ assertion = @("assertion", "assertionFailure"); "test-timeout" = @("testTimeout", "test-timeout", "testTimeoutDetected"); "worker-hang" = @("workerHang", "worker-hang", "externalWatchdog"); "resource leak" = @("resourceLeak", "resource-leak", "resource leak"); environment = @("environment", "environmentEffect") }
    $active = [System.Collections.Generic.List[string]]::new()
    foreach ($kind in $aliases.Keys) { foreach ($alias in $aliases[$kind]) { $value = Get-FrontendVitestValue $signals $alias; if ($value -is [bool] -and $value) { $active.Add($kind); break } } }
    $pass = Get-FrontendVitestValue $Evidence "pass"; $classification = if ($active.Count -eq 1) { $active[0] } elseif ($active.Count -eq 0 -and $pass -eq $true) { "pass" } else { "unknown" }
    [pscustomobject][ordered]@{ classification = $classification; result = if ($classification -eq "pass") { "pass" } elseif ($classification -eq "unknown") { "blocked" } else { "failed" }; blocked = ($classification -eq "unknown"); activeSignals = @($active) }
}

function New-FrontendVitestFailureRecord {
    param([Parameter(Mandatory = $true)][AllowNull()][object]$Evidence)
    $classification = Resolve-FrontendVitestClassification $Evidence; $resources = Get-FrontendVitestValue $Evidence "resources"
    if ($null -eq $resources) { $resources = [ordered]@{ runnerPid = "unavailable"; childPids = "unavailable"; workerCount = "unavailable"; peakWorkingSetBytes = "unavailable"; heapUsedBytes = "unavailable" } }
    $exitCode = Get-FrontendVitestValue $Evidence "processExitCode"; if ($null -eq $exitCode) { $exitCode = Get-FrontendVitestValue $Evidence "exitCode" }; if ($null -eq $exitCode) { $exitCode = "unavailable" }
    [pscustomobject][ordered]@{ file = if ($null -ne (Get-FrontendVitestValue $Evidence "file")) { Get-FrontendVitestValue $Evidence "file" } else { "unavailable" }; title = if ($null -ne (Get-FrontendVitestValue $Evidence "title")) { Get-FrontendVitestValue $Evidence "title" } else { "unavailable" }; kind = $classification.classification; classification = $classification.classification; message = if ($null -ne (Get-FrontendVitestValue $Evidence "message")) { Get-FrontendVitestValue $Evidence "message" } else { "unavailable" }; timeoutMs = if ($null -ne (Get-FrontendVitestValue $Evidence "timeoutMs")) { Get-FrontendVitestValue $Evidence "timeoutMs" } else { "unavailable" }; phase = if ($null -ne (Get-FrontendVitestValue $Evidence "phase")) { Get-FrontendVitestValue $Evidence "phase" } else { "unavailable" }; stdout = if ($null -ne (Get-FrontendVitestValue $Evidence "stdout")) { Get-FrontendVitestValue $Evidence "stdout" } else { "" }; stderr = if ($null -ne (Get-FrontendVitestValue $Evidence "stderr")) { Get-FrontendVitestValue $Evidence "stderr" } else { "" }; exitCode = $exitCode; resources = $resources; result = $classification.result; blocked = $classification.blocked }
}

function Get-FrontendVitestNormalSettings {
    param([int]$WatchdogSeconds)
    [ordered]@{ isolate = "unmodified-default"; pool = "unmodified-default"; workers = "unmodified-default"; order = "unmodified-default"; groupFiles = @(); processLifetime = "one-owned-process-per-run"; externalWatchdogSeconds = $WatchdogSeconds; testTimeoutMs = 5000; hookTimeoutMs = 10000; commandArguments = @($script:NormalVitestArguments); unchanged = $true }
}

function Invoke-FrontendVitestChildRun {
    param([Parameter(Mandatory = $true)][string]$RepoRoot, [Parameter(Mandatory = $true)][int]$WatchdogSeconds)
    $result = Invoke-FrontendVitestOwnedProcess -RepoRoot $RepoRoot -ArgumentList $script:NormalVitestArguments -WatchdogSeconds $WatchdogSeconds; $result.command = $script:NormalVitestCommand; return $result
}

function Test-FrontendVitestRecordFailure {
    param([object]$Record, [object]$Child)
    $startError = [string](Get-FrontendVitestChildResultValue -Child $Child -Name "startError" -Default "")
    $runtimeError = [string](Get-FrontendVitestChildResultValue -Child $Child -Name "runtimeError" -Default "")
    $cleanupOutcome = [string](Get-FrontendVitestChildResultValue -Child $Child -Name "cleanupOutcome" -Default "not-required")
    return ($Record.result -ne "pass" -or $Record.classification -ne "pass" -or -not [string]::IsNullOrWhiteSpace($startError) -or -not [string]::IsNullOrWhiteSpace($runtimeError) -or $cleanupOutcome -eq "unsafe-termination-refused")
}

function Invoke-FrontendVitestNormalFull {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot, [Parameter(Mandatory = $true)][string]$OutputDirectory,
        [Parameter(Mandatory = $true)][int]$Repeat, [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [Parameter(Mandatory = $true)][string[]]$DiscoveredFiles, [string]$SummaryPath, [scriptblock]$ChildInvoker
    )
    if (-not [string]::IsNullOrWhiteSpace($SummaryPath)) { throw "NormalFull mode does not accept SummaryPath; use the owned output summaryPath" }
    if (-not (Test-Path -LiteralPath $OutputDirectory -PathType Container) -or @(Get-ChildItem -LiteralPath $OutputDirectory -Force).Count -gt 0) { throw "NormalFull output directory must be an existing empty directory" }
    $settings = Get-FrontendVitestNormalSettings $WatchdogSeconds; $records = [System.Collections.Generic.List[object]]::new(); $hardFailure = $false; $blockedFailure = $false
    for ($index = 1; $index -le $Repeat; $index++) {
        $child = if ($null -ne $ChildInvoker) { & $ChildInvoker -RepoRoot $RepoRoot -WatchdogSeconds $WatchdogSeconds } else { Invoke-FrontendVitestChildRun -RepoRoot $RepoRoot -WatchdogSeconds $WatchdogSeconds }
        $record = ConvertFrom-FrontendVitestOutput -Stdout $child.stdout -Stderr $child.stderr -ExitCode $child.exitCode -DurationMs ([int]$child.durationMs) -Phase "NormalFull" -RunId ([Guid]::NewGuid().ToString("N")) -RepeatIndex $index -DiscoveredFiles $DiscoveredFiles -ExternalWatchdogTriggered ([bool]$child.externalWatchdogTriggered) -WatchdogSeconds $WatchdogSeconds -RepoRoot $RepoRoot
        $record.command = $script:NormalVitestCommand; $record.startedAtUtc = $child.startedAtUtc; $record.completedAtUtc = $child.completedAtUtc; $record.settings = $settings
        $childCleanupOutcome = [string](Get-FrontendVitestChildResultValue -Child $child -Name "cleanupOutcome" -Default "not-required")
        $childCleanupError = [string](Get-FrontendVitestChildResultValue -Child $child -Name "cleanupError" -Default "")
        $record.resources = [ordered]@{ runnerPid = $child.runnerPid; childPids = $child.childPids; workerCount = "unavailable"; peakWorkingSetBytes = $child.peakWorkingSetBytes; heapUsedBytes = "unavailable"; processExitCode = $child.exitCode; externalWatchdogTriggered = $child.externalWatchdogTriggered; cleanupOutcome = $childCleanupOutcome; cleanupError = $childCleanupError }
        $record.failureSet = @($record.failures | ForEach-Object { [ordered]@{ file = $_.file; title = $_.title; kind = $_.kind } }); $suffix = "-$index"; $jsonName = "normal-full$suffix.json"; $stdoutName = "stdout$suffix.txt"; $stderrName = "stderr$suffix.txt"
        Set-Content -LiteralPath (Join-Path $OutputDirectory $stdoutName) -Value $child.stdout -Encoding utf8; Set-Content -LiteralPath (Join-Path $OutputDirectory $stderrName) -Value $child.stderr -Encoding utf8
        $record.stdout = "artifact:$stdoutName"; $record.stderr = "artifact:$stderrName"; $record.artifactPaths = @($jsonName, $stdoutName, $stderrName); $record.artifactHashes = [ordered]@{ stdout = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $OutputDirectory $stdoutName)).Hash; stderr = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $OutputDirectory $stderrName)).Hash }
        $record | ConvertTo-Json -Depth 50 | Set-Content -LiteralPath (Join-Path $OutputDirectory $jsonName) -Encoding utf8; [void]$records.Add($record)
        $startError = [string](Get-FrontendVitestChildResultValue -Child $child -Name "startError" -Default "")
        $runtimeError = [string](Get-FrontendVitestChildResultValue -Child $child -Name "runtimeError" -Default "")
        if (Test-FrontendVitestRecordFailure $record $child) { $hardFailure = $true; if ($record.result -eq "blocked" -or -not [string]::IsNullOrWhiteSpace($startError) -or -not [string]::IsNullOrWhiteSpace($runtimeError) -or $childCleanupOutcome -eq "unsafe-termination-refused") { $blockedFailure = $true } }
    }
    $summary = [ordered]@{ schemaVersion = 1; runId = [Guid]::NewGuid().ToString("N"); mode = "NormalFull"; repeat = $Repeat; command = $script:NormalVitestCommand; commandArguments = @($script:NormalVitestArguments); settings = $settings; watchdogSeconds = $WatchdogSeconds; sequential = $true; outputDirectory = $OutputDirectory; records = @($records.ToArray()); result = if (-not $hardFailure) { "completed" } elseif ($blockedFailure) { "blocked" } else { "failed" }; artifactPaths = @("normal-full-summary.json") }
    $summaryJson = $summary | ConvertTo-Json -Depth 60; $summaryOutputPath = Join-Path $OutputDirectory "normal-full-summary.json"; Set-Content -LiteralPath $summaryOutputPath -Value $summaryJson -Encoding utf8
    $pointer = [ordered]@{ schemaVersion = 1; mode = "NormalFull"; repeat = $Repeat; outputDirectory = $OutputDirectory; summaryPath = [IO.Path]::GetFullPath($summaryOutputPath); runIds = @($records | ForEach-Object { $_.runId }); result = $summary.result; command = $script:NormalVitestCommand }
    Write-Output ($pointer | ConvertTo-Json -Depth 20 -Compress)
    if ($hardFailure) { throw "NormalFull evidence is $($summary.result); summary preserved at $summaryOutputPath" }
}

function Invoke-FrontendVitestOneVariable {
    param([Parameter(Mandatory = $true)][string]$OutputDirectory, [Parameter(Mandatory = $true)][int]$WatchdogSeconds)
    if (-not (Test-Path -LiteralPath $OutputDirectory -PathType Container) -or @(Get-ChildItem -LiteralPath $OutputDirectory -Force).Count -gt 0) { throw "OneVariable output directory must be an existing empty directory" }
    $reason = "No approved falsifiable hypothesis exists in current evidence; no setting was changed."
    $started = [DateTime]::UtcNow
    $resources = [ordered]@{
        runnerPid = "unavailable"; childPids = @(); workerCount = "unavailable"; peakWorkingSetBytes = "unavailable"; heapUsedBytes = "unavailable"
        processExitCode = "unavailable"; externalWatchdogTriggered = $false; cleanupOutcome = "not-applicable/no-child"; cleanupError = ""
    }
    $record = [ordered]@{
        schemaVersion = 1; runId = [Guid]::NewGuid().ToString("N"); mode = "OneVariable"; repeatIndex = 1; command = $script:NormalVitestCommand
        settings = [ordered]@{ isolate = "unmodified-default"; pool = "unmodified-default"; workers = "unmodified-default"; order = "unmodified-default"; groupFiles = @(); processLifetime = "not-applicable/no-child"; externalWatchdogSeconds = "not-applicable/no-child"; testTimeoutMs = "not-applicable/no-child"; hookTimeoutMs = "not-applicable/no-child"; commandArguments = @($script:NormalVitestArguments); unchanged = $true; variable = "none" }
        startedAtUtc = $started.ToString("o"); completedAtUtc = $started.ToString("o"); durationMs = [long]0
        discovered = [ordered]@{ files = @(); fileCount = 0; tests = "unavailable"; passed = "unavailable"; failed = "unavailable"; timedOut = "unavailable"; identityStatement = $reason }
        failures = @([ordered]@{ file = "unavailable"; title = "unavailable"; kind = "unknown"; classification = "unknown"; message = $reason; timeoutMs = "not-applicable/no-child"; phase = "OneVariable"; stdout = $reason; stderr = ""; exitCode = "unavailable"; resources = $resources; result = "blocked"; blocked = $true })
        resources = $resources; classification = "unknown"; result = "blocked"; reason = $reason; testTimeout = [ordered]@{ observed = $false; count = 0; markerCount = 0; timeoutMs = "not-applicable/no-child" }
        artifactPaths = @("one-variable.json", "stdout.txt", "stderr.txt"); outputDirectory = [IO.Path]::GetFullPath($OutputDirectory); stdout = $reason; stderr = ""
    }
    Set-Content -LiteralPath (Join-Path $OutputDirectory "stdout.txt") -Value $reason -Encoding utf8; Set-Content -LiteralPath (Join-Path $OutputDirectory "stderr.txt") -Value "" -Encoding utf8; $record | ConvertTo-Json -Depth 40 | Set-Content -LiteralPath (Join-Path $OutputDirectory "one-variable.json") -Encoding utf8
    $pointer = [ordered]@{ schemaVersion = 1; mode = "OneVariable"; result = "blocked"; outputDirectory = [IO.Path]::GetFullPath($OutputDirectory); summaryPath = [IO.Path]::GetFullPath((Join-Path $OutputDirectory "one-variable.json")); artifactPaths = @("one-variable.json", "stdout.txt", "stderr.txt") }
    Write-Output ($pointer | ConvertTo-Json -Depth 20 -Compress); throw "OneVariable evidence is blocked; no setting was changed"
}

Export-ModuleMember -Function Get-FrontendVitestValue, Get-FrontendVitestCurrentDiscovery, ConvertFrom-FrontendVitestOutput, Resolve-FrontendVitestClassification, New-FrontendVitestFailureRecord, Invoke-FrontendVitestNormalFull, Invoke-FrontendVitestOneVariable
