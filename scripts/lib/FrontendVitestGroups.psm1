Set-StrictMode -Version Latest

function Get-FrontendVitestGroupRepoRoot {
    [IO.Path]::GetFullPath((Join-Path $PSScriptRoot "..\.."))
}

function Test-FrontendVitestGroupFilePath {
    param([string]$FilePath, [string]$RepoRoot)
    if ([string]::IsNullOrWhiteSpace($FilePath) -or [IO.Path]::IsPathRooted($FilePath)) { return $false }
    $normalized = $FilePath.Trim().Replace('\', '/')
    if ($normalized -notmatch '^frontend/tests/.+\.test\.tsx?$' -or $normalized -match '(^|/)\.\.(/|$)' -or $normalized -match '["`;&|<>\x00-\x1F]') { return $false }
    try { $full = [IO.Path]::GetFullPath((Join-Path $RepoRoot $normalized)) } catch { return $false }
    $root = [IO.Path]::GetFullPath($RepoRoot).TrimEnd('\', '/') + '\'
    if (-not $full.StartsWith($root, [StringComparison]::OrdinalIgnoreCase)) { return $false }
    Test-Path -LiteralPath $full -PathType Leaf
}

function Get-FrontendVitestGroupSource {
    param([Parameter(Mandatory = $true)][string]$SummaryPath, [Parameter(Mandatory = $true)][string]$RepoRoot)
    if (-not (Test-Path -LiteralPath $SummaryPath -PathType Leaf)) { throw "NormalFull evidence is missing" }
    try { $summary = Get-Content -Raw -LiteralPath $SummaryPath | ConvertFrom-Json } catch { throw "NormalFull evidence is malformed" }
    if ($summary.schemaVersion -ne 1 -or $summary.mode -ne "NormalFull" -or $null -eq $summary.PSObject.Properties["records"] -or @($summary.records).Count -eq 0) { throw "NormalFull evidence is not a validated schemaVersion 1 summary" }
    $files = [System.Collections.Generic.List[string]]::new()
    $seen = [System.Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
    foreach ($run in @($summary.records)) {
        if ($null -eq $run.PSObject.Properties["failureSet"]) { throw "NormalFull evidence has no failureSet" }
        foreach ($failure in @($run.failureSet)) {
            if ($null -eq $failure.PSObject.Properties["file"]) { throw "NormalFull failure has no file" }
            $file = ([string]$failure.file).Trim().Replace('\', '/')
            if (-not (Test-FrontendVitestGroupFilePath -FilePath $file -RepoRoot $RepoRoot)) { throw "NormalFull failure path is unsafe or missing: $file" }
            if ($seen.Add($file)) { [void]$files.Add($file) }
        }
    }
    if ($files.Count -eq 0) { throw "NormalFull evidence has no observed failure files" }
    $resolved = [IO.Path]::GetFullPath($SummaryPath)
    [ordered]@{
        files = $files.ToArray()
        summaryPath = $resolved
        summarySha256 = (Get-FileHash -Algorithm SHA256 -LiteralPath $resolved).Hash
        sourceRecordCount = @($summary.records).Count
        derivation = "ordered dynamic union of current NormalFull records.failureSet.file; not a fixed path list"
    }
}

function Get-FrontendVitestGroupCommandArguments {
    param([Parameter(Mandatory = $true)][ValidateSet("DefaultGroup", "SingleWorkerGroup")][string]$Mode, [Parameter(Mandatory = $true)][string[]]$Files, [string]$RepoRoot)
    if ([string]::IsNullOrWhiteSpace($RepoRoot)) { $RepoRoot = Get-FrontendVitestGroupRepoRoot }
    $selectors = [System.Collections.Generic.List[string]]::new()
    $seen = [System.Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
    foreach ($file in @($Files)) {
        $normalized = $file.Trim().Replace('\', '/')
        if (-not (Test-FrontendVitestGroupFilePath -FilePath $normalized -RepoRoot $RepoRoot)) { throw "group file selector is unsafe or missing: $normalized" }
        if (-not $seen.Add($normalized)) { throw "group file selector is duplicated: $normalized" }
        [void]$selectors.Add($normalized.Substring("frontend/".Length))
    }
    if ($selectors.Count -eq 0) { throw "group file selector list is empty" }
    $arguments = [System.Collections.Generic.List[string]]::new()
    foreach ($argument in @("--prefix", "frontend", "run", "test", "--", "--run")) { [void]$arguments.Add($argument) }
    foreach ($selector in $selectors) { [void]$arguments.Add($selector) }
    if ($Mode -eq "SingleWorkerGroup") {
        foreach ($argument in @("--no-file-parallelism", "--maxWorkers=1", "--minWorkers=1")) { [void]$arguments.Add($argument) }
    }
    $arguments.ToArray()
}

function Get-FrontendVitestGroupCommandText {
    param([Parameter(Mandatory = $true)][ValidateSet("DefaultGroup", "SingleWorkerGroup")][string]$Mode, [Parameter(Mandatory = $true)][string[]]$Files, [string]$RepoRoot)
    "npm $((Get-FrontendVitestGroupCommandArguments -Mode $Mode -Files $Files -RepoRoot $RepoRoot) -join ' ')"
}

function Get-FrontendVitestGroupSettings {
    param([Parameter(Mandatory = $true)][ValidateSet("DefaultGroup", "SingleWorkerGroup")][string]$Mode, [Parameter(Mandatory = $true)][string[]]$Files, [int]$WatchdogSeconds = 600, [string]$RepoRoot)
    $arguments = @(Get-FrontendVitestGroupCommandArguments -Mode $Mode -Files $Files -RepoRoot $RepoRoot)
    [ordered]@{
        isolate = "unmodified-default"
        pool = "unmodified-default"
        workers = if ($Mode -eq "SingleWorkerGroup") { "maxWorkers=1,minWorkers=1" } else { "unmodified-default" }
        workerMode = if ($Mode -eq "SingleWorkerGroup") { "single" } else { "default" }
        workerArgs = if ($Mode -eq "SingleWorkerGroup") { @("--no-file-parallelism", "--maxWorkers=1", "--minWorkers=1") } else { @() }
        order = "unmodified-default"
        groupFiles = @($Files)
        processLifetime = "one-owned-process-per-run"
        externalWatchdogSeconds = $WatchdogSeconds
        testTimeoutMs = 5000
        hookTimeoutMs = 10000
        commandArguments = $arguments
        diagnosisOnly = $true
    }
}

function Invoke-FrontendVitestGroupChildRun {
    param([string]$RepoRoot, [string]$Mode, [string[]]$Files, [int]$WatchdogSeconds)
    $started = [DateTime]::UtcNow; $timer = [Diagnostics.Stopwatch]::StartNew(); $process = [Diagnostics.Process]::new()
    try {
        $npm = Get-Command -Name "npm.cmd" -CommandType Application -ErrorAction Stop
        $psi = [Diagnostics.ProcessStartInfo]::new(); $psi.FileName = if ($npm.Source) { $npm.Source } else { $npm.Path }
        foreach ($argument in @(Get-FrontendVitestGroupCommandArguments -Mode $Mode -Files $Files -RepoRoot $RepoRoot)) { [void]$psi.ArgumentList.Add([string]$argument) }
        $psi.WorkingDirectory = $RepoRoot; $psi.UseShellExecute = $false; $psi.CreateNoWindow = $true; $psi.RedirectStandardOutput = $true; $psi.RedirectStandardError = $true; $process.StartInfo = $psi
        if (-not $process.Start()) { throw "npm process did not start" }
    } catch {
        $process.Dispose()
        return [ordered]@{ startedAtUtc = $started.ToString("o"); completedAtUtc = [DateTime]::UtcNow.ToString("o"); durationMs = [long]$timer.ElapsedMilliseconds; stdout = ""; stderr = $_.Exception.Message; exitCode = "unavailable"; runnerPid = "unavailable"; childPids = "unavailable"; peakWorkingSetBytes = "unavailable"; externalWatchdogTriggered = $false; cleanupOutcome = "start-failed"; cleanupError = ""; startError = $_.Exception.Message }
    }
    $runnerPid = $process.Id; $stdoutTask = $process.StandardOutput.ReadToEndAsync(); $stderrTask = $process.StandardError.ReadToEndAsync(); $watchdog = $false
    while (-not $process.WaitForExit(250)) { if ($timer.Elapsed.TotalSeconds -ge $WatchdogSeconds) { $watchdog = $true; break } }
    $cleanup = "not-required"; $cleanupError = ""
    if ($watchdog) { try { if (-not $process.HasExited) { $process.Kill($true); $cleanup = "owned-tree-terminated" } else { $cleanup = "already-exited" } } catch { $cleanup = "unsafe-termination-refused"; $cleanupError = $_.Exception.Message } }
    $stdout = if ($stdoutTask.Wait(5000)) { $stdoutTask.Result } else { "unavailable" }; $stderr = if ($stderrTask.Wait(5000)) { $stderrTask.Result } else { "unavailable" }; $exit = if ($process.HasExited) { $process.ExitCode } else { "unavailable" }; $peak = "unavailable"; try { if ($process.HasExited) { $peakValue = $process.PeakWorkingSet64; if ($null -ne $peakValue) { $peak = [long]$peakValue } } } catch { $peak = "unavailable" }; $timer.Stop(); $completed = [DateTime]::UtcNow; $process.Dispose()
    [ordered]@{ startedAtUtc = $started.ToString("o"); completedAtUtc = $completed.ToString("o"); durationMs = [long]$timer.ElapsedMilliseconds; stdout = $stdout; stderr = $stderr; exitCode = $exit; runnerPid = $runnerPid; childPids = "unavailable"; peakWorkingSetBytes = $peak; externalWatchdogTriggered = $watchdog; cleanupOutcome = $cleanup; cleanupError = $cleanupError; startError = "" }
}

function Get-FrontendVitestGroupDurationSummary {
    param([object[]]$Records)
    $values = @($Records | ForEach-Object { if ($_.durationMs -is [int] -or $_.durationMs -is [long] -or $_.durationMs -is [double]) { [double]$_.durationMs } } | Sort-Object)
    if ($values.Count -eq 0) { return [ordered]@{ min = "unavailable"; max = "unavailable"; median = "unavailable" } }
    $middle = [int][Math]::Floor($values.Count / 2); $median = if ($values.Count % 2) { $values[$middle] } else { ($values[$middle - 1] + $values[$middle]) / 2 }
    [ordered]@{ min = $values[0]; max = $values[$values.Count - 1]; median = $median }
}

function Invoke-FrontendVitestGroup {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string]$OutputDirectory,
        [Parameter(Mandatory = $true)][int]$Repeat,
        [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [Parameter(Mandatory = $true)][ValidateSet("DefaultGroup", "SingleWorkerGroup")][string]$Mode,
        [string]$SummaryPath,
        [scriptblock]$ChildInvoker
    )
    if (-not (Test-Path -LiteralPath $OutputDirectory -PathType Container) -or @(Get-ChildItem -LiteralPath $OutputDirectory -Force).Count -gt 0) { throw "group output directory must be an existing empty directory" }
    if ([string]::IsNullOrWhiteSpace($SummaryPath)) { $SummaryPath = Join-Path $RepoRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence\normal-full-summary.json" }
    $source = Get-FrontendVitestGroupSource -SummaryPath $SummaryPath -RepoRoot $RepoRoot; $files = @($source.files); $settings = Get-FrontendVitestGroupSettings -Mode $Mode -Files $files -WatchdogSeconds $WatchdogSeconds -RepoRoot $RepoRoot
    $records = [System.Collections.Generic.List[object]]::new(); $artifactIndex = [System.Collections.Generic.List[object]]::new(); $recordPaths = [System.Collections.Generic.List[string]]::new(); $hardFailure = $false
    for ($repeatIndex = 1; $repeatIndex -le $Repeat; $repeatIndex++) {
        $runId = [Guid]::NewGuid().ToString("N"); $runDirectory = Join-Path $OutputDirectory ("run-{0:D3}" -f $repeatIndex); New-Item -ItemType Directory -Path $runDirectory -Force | Out-Null
        try { $child = if ($null -ne $ChildInvoker) { & $ChildInvoker -RepoRoot $RepoRoot -Mode $Mode -Files $files -RunId $runId -RepeatIndex $repeatIndex -WatchdogSeconds $WatchdogSeconds } else { Invoke-FrontendVitestGroupChildRun -RepoRoot $RepoRoot -Mode $Mode -Files $files -WatchdogSeconds $WatchdogSeconds } } catch { $child = [ordered]@{ startedAtUtc = [DateTime]::UtcNow.ToString("o"); completedAtUtc = [DateTime]::UtcNow.ToString("o"); durationMs = 0; stdout = ""; stderr = $_.Exception.Message; exitCode = "unavailable"; runnerPid = "unavailable"; childPids = "unavailable"; peakWorkingSetBytes = "unavailable"; externalWatchdogTriggered = $false; cleanupOutcome = "unsafe-termination-refused"; cleanupError = $_.Exception.Message; startError = $_.Exception.Message } }
        $record = ConvertFrom-FrontendVitestOutput -Stdout $child.stdout -Stderr $child.stderr -ExitCode $child.exitCode -DurationMs ([int]$child.durationMs) -Phase $Mode -RunId $runId -RepeatIndex $repeatIndex -DiscoveredFiles $files -ExternalWatchdogTriggered ([bool]$child.externalWatchdogTriggered) -WatchdogSeconds $WatchdogSeconds -RepoRoot $RepoRoot
        $record.command = Get-FrontendVitestGroupCommandText -Mode $Mode -Files $files -RepoRoot $RepoRoot; $record.startedAtUtc = $child.startedAtUtc; $record.completedAtUtc = $child.completedAtUtc; $record.settings = $settings; $record.groupFiles = $files; $record.sourceSummaryPath = $source.summaryPath; $record.sourceSummarySha256 = $source.summarySha256; $record.resources = [ordered]@{ runnerPid = $child.runnerPid; childPids = $child.childPids; workerCount = "unavailable"; peakWorkingSetBytes = $child.peakWorkingSetBytes; heapUsedBytes = "unavailable"; processExitCode = $child.exitCode; externalWatchdogTriggered = $child.externalWatchdogTriggered; cleanupOutcome = $child.cleanupOutcome; cleanupError = $child.cleanupError }
        $stdoutName = "stdout.txt"; $stderrName = "stderr.txt"; $jsonName = "record.json"; Set-Content -LiteralPath (Join-Path $runDirectory $stdoutName) -Value $child.stdout -Encoding utf8; Set-Content -LiteralPath (Join-Path $runDirectory $stderrName) -Value $child.stderr -Encoding utf8; $record.stdout = "artifact:$stdoutName"; $record.stderr = "artifact:$stderrName"; $record.artifactPaths = @($jsonName, $stdoutName, $stderrName); $record.artifactHashes = [ordered]@{ stdout = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $runDirectory $stdoutName)).Hash; stderr = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $runDirectory $stderrName)).Hash }; $jsonPath = Join-Path $runDirectory $jsonName; $record | ConvertTo-Json -Depth 60 | Set-Content -LiteralPath $jsonPath -Encoding utf8; [void]$records.Add($record); [void]$recordPaths.Add($jsonPath)
        foreach ($name in @($jsonName, $stdoutName, $stderrName)) { $path = Join-Path $runDirectory $name; [void]$artifactIndex.Add([ordered]@{ path = $path; sha256 = (Get-FileHash -Algorithm SHA256 -LiteralPath $path).Hash }) }
        if ($child.startError -ne "" -or $child.cleanupOutcome -eq "unsafe-termination-refused" -or $record.classification -eq "unknown") { $hardFailure = $true }
    }
    $recordArray = $records.ToArray(); $sequential = $true; for ($i = 1; $i -lt $recordArray.Count; $i++) { if ([DateTime]::Parse($recordArray[$i].startedAtUtc) -lt [DateTime]::Parse($recordArray[$i - 1].completedAtUtc)) { $sequential = $false } }
    $allPresent = $true; $allHashesMatch = $true; foreach ($artifact in @($artifactIndex)) { if (-not (Test-Path -LiteralPath $artifact.path -PathType Leaf)) { $allPresent = $false; continue }; try { if ((Get-FileHash -Algorithm SHA256 -LiteralPath $artifact.path).Hash -ne $artifact.sha256) { $allHashesMatch = $false } } catch { $allHashesMatch = $false } }
    if (-not $allPresent -or -not $allHashesMatch -or -not $sequential) { $hardFailure = $true }
    $classificationTotals = [ordered]@{ pass = 0; assertion = 0; testTimeout = 0; workerHang = 0; resourceLeak = 0; environment = 0; unknown = 0 }
    foreach ($record in $recordArray) { switch ([string]$record.classification) { "pass" { $classificationTotals.pass++ } "assertion" { $classificationTotals.assertion++ } "test-timeout" { $classificationTotals.testTimeout++ } "worker-hang" { $classificationTotals.workerHang++ } "resource leak" { $classificationTotals.resourceLeak++ } "environment" { $classificationTotals.environment++ } default { $classificationTotals.unknown++ } } }
    $summary = [ordered]@{ schemaVersion = 1; runId = [Guid]::NewGuid().ToString("N"); mode = $Mode; repeat = $Repeat; command = Get-FrontendVitestGroupCommandText -Mode $Mode -Files $files -RepoRoot $RepoRoot; commandContract = "fixed npm command with validated selectors; no arbitrary command input"; sourceSummaryPath = $source.summaryPath; sourceSummarySha256 = $source.summarySha256; sourceFileCount = $files.Count; groupFiles = $files; settings = $settings; outputDirectory = $OutputDirectory; sequential = $sequential; recordCount = $recordArray.Count; classificationTotals = $classificationTotals; durationMs = Get-FrontendVitestGroupDurationSummary -Records $recordArray; records = @($recordArray | ForEach-Object { [ordered]@{ runId = $_.runId; repeatIndex = $_.repeatIndex; result = $_.result; classification = $_.classification; durationMs = $_.durationMs; failureSet = @($_.failures | ForEach-Object { [ordered]@{ file = $_.file; title = $_.title; kind = $_.kind } }) } }); artifactIntegrity = [ordered]@{ checked = $artifactIndex.Count; allPresent = $allPresent; allHashesMatch = $allHashesMatch; rawLocation = $OutputDirectory }; resources = [ordered]@{ childPids = "unavailable"; workerCount = "unavailable"; peakWorkingSetBytes = "unavailable"; heapUsedBytes = "unavailable" }; watchdog = [ordered]@{ externalWatchdogTriggered = @($recordArray | Where-Object externalWatchdogTriggered).Count -gt 0; seconds = $WatchdogSeconds }; result = if ($hardFailure) { "blocked" } else { "completed" }; diagnosisSignal = "SingleWorkerGroup is diagnosis evidence only; it is not a worker-configuration repair." }
    $summaryPathOut = Join-Path $OutputDirectory "group-summary.json"; $summary | ConvertTo-Json -Depth 80 | Set-Content -LiteralPath $summaryPathOut -Encoding utf8
    $pointer = [ordered]@{ schemaVersion = 1; mode = $Mode; repeat = $Repeat; fileCount = $files.Count; recordCount = $recordArray.Count; runIds = @($recordArray | ForEach-Object runId); recordPaths = @($recordPaths); summaryPath = $summaryPathOut; sourceSummaryPath = $source.summaryPath; sourceSummarySha256 = $source.summarySha256; sequential = $sequential; result = $summary.result; classificationTotals = $classificationTotals; artifactIntegrity = $summary.artifactIntegrity }
    Write-Output ($pointer | ConvertTo-Json -Depth 30 -Compress)
    if ($hardFailure) { throw "$Mode evidence is blocked; summary preserved at $summaryPathOut" }
}

Export-ModuleMember -Function Get-FrontendVitestGroupSource, Get-FrontendVitestGroupCommandArguments, Get-FrontendVitestGroupCommandText, Get-FrontendVitestGroupSettings, Invoke-FrontendVitestGroup
