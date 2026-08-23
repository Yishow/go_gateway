Set-StrictMode -Version Latest

$script:NormalVitestCommand = "npm --prefix frontend run test -- --run"

function ConvertTo-MatrixRelativePath {
    param([string]$Path, [string]$RepoRoot)

    if ([string]::IsNullOrWhiteSpace($Path)) { return "unavailable" }
    $normalized = $Path.Trim().Replace('\', '/')
    if ($normalized -match '^(?:[A-Za-z]:/|/|\\)') {
        try { $normalized = [IO.Path]::GetRelativePath($RepoRoot, [IO.Path]::GetFullPath($Path)).Replace('\', '/') } catch { return "unavailable" }
    } elseif ($normalized -match '^frontend/tests/') {
        return $normalized
    } elseif ($normalized -match '^tests/') {
        return "frontend/$normalized"
    } else {
        return "unavailable"
    }
    if ($normalized -notmatch '^frontend/tests/') { return "unavailable" }
    return $normalized
}

function Get-VitestSummaryCount {
    param([string]$Text, [string]$Label, [string]$State)

    $escapedLabel = [regex]::Escape($Label)
    $line = [regex]::Match($Text, "(?im)^\s*$escapedLabel\s+(?<body>.+)$")
    if (-not $line.Success) { return "unavailable" }
    $escapedState = [regex]::Escape($State)
    $stateMatch = [regex]::Match($line.Groups["body"].Value, "(?<value>\d+)\s+$escapedState")
    if ($stateMatch.Success) { return [int]$stateMatch.Groups["value"].Value }
    if ($line.Groups["body"].Value -match '\(\d+\)' -or $line.Groups["body"].Value -match '\d+') { return 0 }
    return "unavailable"
}

function Get-VitestTotalCount {
    param([string]$Text, [string]$Label)

    $escapedLabel = [regex]::Escape($Label)
    $line = [regex]::Match($Text, "(?im)^\s*$escapedLabel\s+(?<body>.+)$")
    if (-not $line.Success) { return "unavailable" }
    $totalMatch = [regex]::Match($line.Groups["body"].Value, "\((?<value>\d+)\)")
    if ($totalMatch.Success) { return [int]$totalMatch.Groups["value"].Value }
    $numbers = [regex]::Matches($line.Groups["body"].Value, "(?<value>\d+)")
    if ($numbers.Count -gt 0) { return [int]$numbers[$numbers.Count - 1].Groups["value"].Value }
    return "unavailable"
}

function ConvertFrom-FrontendVitestOutput {
    param(
        [AllowNull()][string]$Stdout,
        [AllowNull()][string]$Stderr,
        [AllowNull()][object]$ExitCode,
        [Parameter(Mandatory = $true)][int]$DurationMs,
        [Parameter(Mandatory = $true)][string]$Phase,
        [Parameter(Mandatory = $true)][string]$RunId,
        [Parameter(Mandatory = $true)][int]$RepeatIndex,
        [Parameter(Mandatory = $true)][string[]]$DiscoveredFiles,
        [Parameter(Mandatory = $true)][bool]$ExternalWatchdogTriggered,
        [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [string]$RepoRoot = (Split-Path -Parent (Split-Path -Parent $PSScriptRoot))
    )

    $raw = (([string]$Stdout) + "`n" + ([string]$Stderr)) -replace "`e\[[0-9;]*m", ""
    $records = [System.Collections.Generic.List[object]]::new()
    $indexes = @{}
    $currentFile = $null
    $currentTitle = $null
    $timeoutValues = [System.Collections.Generic.List[int]]::new()
    $lines = $raw -split "`r?`n"
    foreach ($line in $lines) {
        $pathMatch = [regex]::Match($line, '(?<path>(?:frontend[\\/]?)?tests[\\/][^\s>]+\.test\.tsx?)')
        if ($pathMatch.Success) {
            $previousFile = $currentFile
            $currentFile = ConvertTo-MatrixRelativePath -Path $pathMatch.Groups["path"].Value -RepoRoot $RepoRoot
            if ($currentFile -ne $previousFile) { $currentTitle = $null }
            $tail = $line.Substring($pathMatch.Index + $pathMatch.Length)
            $titleMatch = [regex]::Match($tail, '^\s*>\s*(?<title>.+)$')
            if ($titleMatch.Success -and $line -match '^\s*(?:FAIL|[❯×✗])') {
                $title = [regex]::Replace($titleMatch.Groups["title"].Value.Trim(), '\s+\d+\s*ms$', '')
                $currentTitle = $title
                $key = "$currentFile|$title"
                if (-not $indexes.ContainsKey($key)) {
                    $indexes[$key] = $records.Count
                    $records.Add([ordered]@{ file = $currentFile; title = $title; kind = "assertion"; message = "unavailable"; timeoutMs = "unavailable"; phase = $Phase })
                }
            }
        }
        $markerMatch = [regex]::Match($line, '^\s*[×✗]\s+(?<title>.+)$')
        if ($markerMatch.Success -and $null -ne $currentFile) {
            $title = [regex]::Replace($markerMatch.Groups["title"].Value.Trim(), '\s+\d+\s*ms$', '')
            $currentTitle = $title
            $key = "$currentFile|$title"
            if (-not $indexes.ContainsKey($key)) {
                $indexes[$key] = $records.Count
                $records.Add([ordered]@{ file = $currentFile; title = $title; kind = "assertion"; message = "unavailable"; timeoutMs = "unavailable"; phase = $Phase })
            }
        }
        $timeoutMatch = [regex]::Match($line, '(?i)timed\s+out\s+in\s+(?<ms>\d+)\s*ms')
        if ($timeoutMatch.Success) {
            $timeoutValues.Add([int]$timeoutMatch.Groups["ms"].Value)
            $target = $null
            if ($null -ne $currentFile -and $null -ne $currentTitle) {
                $key = "$currentFile|$currentTitle"
                if ($indexes.ContainsKey($key)) { $target = $records[$indexes[$key]] }
            }
            if ($null -eq $target -and $null -ne $currentFile) {
                $sameFile = @($records | Where-Object { $_["file"] -eq $currentFile })
                if ($sameFile.Count -eq 1) { $target = $sameFile[0] }
            }
            if ($null -ne $target) {
                $target["kind"] = "test-timeout"
                $target["timeoutMs"] = [int]$timeoutMatch.Groups["ms"].Value
                $target["message"] = $line.Trim()
            } elseif ($null -ne $currentFile) {
                $key = "$currentFile|$currentTitle"
                if ($indexes.ContainsKey($key)) { $target = $records[$indexes[$key]] }
                else {
                    $key = "$currentFile|unavailable"
                    $indexes[$key] = $records.Count
                    $records.Add([ordered]@{ file = $currentFile; title = if ($null -ne $currentTitle) { $currentTitle } else { "unavailable" }; kind = "test-timeout"; message = $line.Trim(); timeoutMs = [int]$timeoutMatch.Groups["ms"].Value; phase = $Phase })
                }
            } else {
                $currentTitle = $null
            }
        }
        if ($records.Count -gt 0 -and $line -match '^\s*(?:→|Error:|AssertionError|Expected)') {
            $last = $records[$records.Count - 1]
            if ([string]$last["message"] -eq "unavailable") { $last["message"] = $line.Trim() }
        }
    }

    $fileTotal = Get-VitestTotalCount -Text $raw -Label "Test Files"
    $testTotal = Get-VitestTotalCount -Text $raw -Label "Tests"
    $filePassed = Get-VitestSummaryCount -Text $raw -Label "Test Files" -State "passed"
    $fileFailed = Get-VitestSummaryCount -Text $raw -Label "Test Files" -State "failed"
    $testPassed = Get-VitestSummaryCount -Text $raw -Label "Tests" -State "passed"
    $testFailed = Get-VitestSummaryCount -Text $raw -Label "Tests" -State "failed"
    $uniqueTimeoutCount = @($records | Where-Object { $_["kind"] -eq "test-timeout" }).Count
    $testTimeoutObserved = $uniqueTimeoutCount -gt 0
    $assertionObserved = @($records | Where-Object { $_["kind"] -eq "assertion" }).Count -gt 0
    $failedSummaryObserved = (($testFailed -is [int]) -and $testFailed -gt 0) -or (($fileFailed -is [int]) -and $fileFailed -gt 0)
    $categories = @()
    if ($ExternalWatchdogTriggered) { $categories += "worker-hang" }
    if ($testTimeoutObserved) { $categories += "test-timeout" }
    if ($assertionObserved -or ($failedSummaryObserved -and -not $testTimeoutObserved -and -not $ExternalWatchdogTriggered)) { $categories += "assertion" }
    $classification = "unknown"
    if ($categories.Count -eq 1) {
        $classification = $categories[0]
    } elseif ($categories.Count -eq 0 -and $ExitCode -eq 0 -and $fileTotal -ne "unavailable" -and $testTotal -ne "unavailable") {
        $classification = "pass"
    }
    if ($classification -eq "worker-hang" -and $records.Count -eq 0) {
        $records.Add([ordered]@{ file = "unavailable"; title = "external watchdog"; kind = "worker-hang"; message = "external watchdog boundary reached"; timeoutMs = "unavailable"; phase = $Phase })
    }
    if ($classification -eq "assertion" -and $records.Count -eq 0) {
        $records.Add([ordered]@{ file = "unavailable"; title = "unavailable"; kind = "assertion"; message = "Vitest summary reported failed tests"; timeoutMs = "unavailable"; phase = $Phase })
    }
    if ($classification -eq "unknown" -and $records.Count -eq 0) {
        $records.Add([ordered]@{ file = "unavailable"; title = "unavailable"; kind = "unknown"; message = "Vitest output was insufficient to classify"; timeoutMs = "unavailable"; phase = $Phase })
    }
    $timeoutMs = if ($timeoutValues.Count -eq 1) { $timeoutValues[0] } elseif ($timeoutValues.Count -gt 1 -and (@($timeoutValues | Select-Object -Unique).Count -eq 1)) { $timeoutValues[0] } else { "unavailable" }
    $result = if ($classification -eq "pass") { "pass" } elseif ($classification -eq "unknown") { "blocked" } else { "failed" }
    [ordered]@{
        schemaVersion = 1
        runId = $RunId
        mode = $Phase
        repeatIndex = $RepeatIndex
        command = $script:NormalVitestCommand
        startedAtUtc = "unavailable"
        durationMs = $DurationMs
        discovered = [ordered]@{ files = @($DiscoveredFiles); fileCount = @($DiscoveredFiles).Count; tests = $testTotal; passed = $testPassed; failed = $testFailed; timedOut = $uniqueTimeoutCount }
        discoveredFileSummary = [ordered]@{ total = $fileTotal; passed = $filePassed; failed = $fileFailed }
        failures = @($records)
        resources = [ordered]@{ runnerPid = "unavailable"; childPids = "unavailable"; workerCount = "unavailable"; peakWorkingSetBytes = "unavailable"; heapUsedBytes = "unavailable"; processExitCode = $ExitCode; externalWatchdogTriggered = $ExternalWatchdogTriggered }
        result = $result
        classification = $classification
        externalWatchdogTriggered = $ExternalWatchdogTriggered
        testTimeout = [ordered]@{ observed = $testTimeoutObserved; count = $uniqueTimeoutCount; markerCount = $timeoutValues.Count; timeoutMs = $timeoutMs }
        stdout = [string]$Stdout
        stderr = [string]$Stderr
        artifactPaths = @()
        watchdogSeconds = $WatchdogSeconds
    }
}

function Invoke-FrontendVitestChildRun {
    param([Parameter(Mandatory = $true)][string]$RepoRoot, [Parameter(Mandatory = $true)][int]$WatchdogSeconds)

    $started = [DateTime]::UtcNow
    $stopwatch = [Diagnostics.Stopwatch]::StartNew()
    $psi = [Diagnostics.ProcessStartInfo]::new()
    $npmCommand = Get-Command -Name "npm.cmd" -CommandType Application -ErrorAction Stop
    $psi.FileName = if (-not [string]::IsNullOrWhiteSpace($npmCommand.Source)) { $npmCommand.Source } else { $npmCommand.Path }
    $psi.Arguments = "--prefix frontend run test -- --run"
    $psi.WorkingDirectory = $RepoRoot
    $psi.UseShellExecute = $false
    $psi.CreateNoWindow = $true
    $psi.RedirectStandardOutput = $true
    $psi.RedirectStandardError = $true
    $process = [Diagnostics.Process]::new()
    $process.StartInfo = $psi
    try {
        if (-not $process.Start()) { throw "npm process did not start" }
    } catch {
        return [ordered]@{ startedAtUtc = $started.ToString("o"); completedAtUtc = [DateTime]::UtcNow.ToString("o"); durationMs = [long]$stopwatch.ElapsedMilliseconds; stdout = ""; stderr = $_.Exception.Message; exitCode = "unavailable"; runnerPid = "unavailable"; childPids = "unavailable"; peakWorkingSetBytes = "unavailable"; externalWatchdogTriggered = $false; cleanupOutcome = "start-failed"; startError = $_.Exception.Message }
    }
    $runnerPid = $process.Id
    $stdoutTask = $process.StandardOutput.ReadToEndAsync()
    $stderrTask = $process.StandardError.ReadToEndAsync()
    $watchdog = $false
    while (-not $process.WaitForExit(250)) {
        if ($stopwatch.Elapsed.TotalSeconds -ge $WatchdogSeconds) { $watchdog = $true; break }
    }
    $cleanup = "not-required"
    $cleanupError = ""
    if ($watchdog) {
        try {
            if (-not $process.HasExited) { $process.Kill($true); $cleanup = "owned-tree-terminated" } else { $cleanup = "already-exited" }
        } catch {
            $cleanup = "unsafe-termination-refused"
            $cleanupError = $_.Exception.Message
        }
    }
    $stdout = if ($stdoutTask.Wait(5000)) { $stdoutTask.Result } else { "unavailable" }
    $stderr = if ($stderrTask.Wait(5000)) { $stderrTask.Result } else { "unavailable" }
    $exitCode = if ($process.HasExited) { $process.ExitCode } else { "unavailable" }
    $peakWorkingSet = try { if ($process.HasExited) { $process.PeakWorkingSet64 } else { "unavailable" } } catch { "unavailable" }
    $stopwatch.Stop()
    $completed = [DateTime]::UtcNow
    $process.Dispose()
    [ordered]@{ startedAtUtc = $started.ToString("o"); completedAtUtc = $completed.ToString("o"); durationMs = [long]$stopwatch.ElapsedMilliseconds; stdout = $stdout; stderr = $stderr; exitCode = $exitCode; runnerPid = $runnerPid; childPids = "unavailable"; peakWorkingSetBytes = $peakWorkingSet; externalWatchdogTriggered = $watchdog; cleanupOutcome = $cleanup; cleanupError = $cleanupError; startError = "" }
}

function Invoke-FrontendVitestNormalFull {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string]$OutputDirectory,
        [Parameter(Mandatory = $true)][int]$Repeat,
        [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [Parameter(Mandatory = $true)][string[]]$DiscoveredFiles
    )

    $records = [System.Collections.Generic.List[object]]::new()
    $hardFailure = $false
    for ($index = 1; $index -le $Repeat; $index++) {
        $child = Invoke-FrontendVitestChildRun -RepoRoot $RepoRoot -WatchdogSeconds $WatchdogSeconds
        $record = ConvertFrom-FrontendVitestOutput -Stdout $child.stdout -Stderr $child.stderr -ExitCode $child.exitCode -DurationMs $child.durationMs -Phase "NormalFull" -RunId ([Guid]::NewGuid().ToString("N")) -RepeatIndex $index -DiscoveredFiles $DiscoveredFiles -ExternalWatchdogTriggered $child.externalWatchdogTriggered -WatchdogSeconds $WatchdogSeconds -RepoRoot $RepoRoot
        $record.startedAtUtc = $child.startedAtUtc
        $record.completedAtUtc = $child.completedAtUtc
        $record.resources = [ordered]@{ runnerPid = $child.runnerPid; childPids = $child.childPids; workerCount = "unavailable"; peakWorkingSetBytes = $child.peakWorkingSetBytes; heapUsedBytes = "unavailable"; processExitCode = $child.exitCode; externalWatchdogTriggered = $child.externalWatchdogTriggered; cleanupOutcome = $child.cleanupOutcome; cleanupError = $child.cleanupError }
        $suffix = "-$index"
        $jsonName = "normal-full$suffix.json"
        $stdoutName = "stdout$suffix.txt"
        $stderrName = "stderr$suffix.txt"
        Set-Content -LiteralPath (Join-Path $OutputDirectory $stdoutName) -Value $child.stdout -Encoding utf8
        Set-Content -LiteralPath (Join-Path $OutputDirectory $stderrName) -Value $child.stderr -Encoding utf8
        $record.stdout = "artifact:$stdoutName"
        $record.stderr = "artifact:$stderrName"
        $record.artifactPaths = @($jsonName, $stdoutName, $stderrName)
        $record.artifactHashes = [ordered]@{ stdout = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $OutputDirectory $stdoutName)).Hash; stderr = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $OutputDirectory $stderrName)).Hash }
        $record | ConvertTo-Json -Depth 40 | Set-Content -LiteralPath (Join-Path $OutputDirectory $jsonName) -Encoding utf8
        $records.Add($record)
        if ($child.startError -ne "" -or $child.cleanupOutcome -eq "unsafe-termination-refused" -or $record.classification -eq "unknown") { $hardFailure = $true }
    }
    $summary = [ordered]@{
        schemaVersion = 1
        runId = [Guid]::NewGuid().ToString("N")
        mode = "NormalFull"
        repeat = $Repeat
        command = $script:NormalVitestCommand
        sequential = $true
        outputDirectory = $OutputDirectory
        records = @($records)
        result = if ($hardFailure) { "blocked" } else { "completed" }
        artifactPaths = @("normal-full-summary.json")
    }
    $summary | ConvertTo-Json -Depth 50 | Set-Content -LiteralPath (Join-Path $OutputDirectory "normal-full-summary.json") -Encoding utf8
    $pointer = [ordered]@{ schemaVersion = 1; mode = "NormalFull"; repeat = $Repeat; outputDirectory = $OutputDirectory; summaryPath = "normal-full-summary.json"; runIds = @($records | ForEach-Object { $_.runId }); result = $summary.result; command = $script:NormalVitestCommand }
    Write-Output ($pointer | ConvertTo-Json -Depth 10 -Compress)
    if ($hardFailure) { throw "NormalFull evidence is blocked or unsafe; summary preserved at $OutputDirectory\normal-full-summary.json" }
}

Export-ModuleMember -Function ConvertFrom-FrontendVitestOutput, Invoke-FrontendVitestNormalFull
