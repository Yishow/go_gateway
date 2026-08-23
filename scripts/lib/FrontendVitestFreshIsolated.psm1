Set-StrictMode -Version Latest

Import-Module (Join-Path $PSScriptRoot "FrontendVitestEvidence.psm1") -Force -Global
Import-Module (Join-Path $PSScriptRoot "FrontendVitestProcess.psm1") -Force -Global
Import-Module (Join-Path $PSScriptRoot "FrontendVitestMatrix.psm1") -Force -Global

function Test-FreshIsolatedFilePath {
    param([string]$FilePath, [string]$RepoRoot)
    if ([string]::IsNullOrWhiteSpace($FilePath) -or [IO.Path]::IsPathRooted($FilePath)) { return $false }
    $normalized = $FilePath.Trim().Replace('\', '/')
    if ($normalized -match '(^|/)\.\.(/|$)' -or $normalized -notmatch '^frontend/tests/.+\.test\.tsx?$' -or $normalized -match '["`;&|<>\x00-\x1F]') { return $false }
    try { $full = [IO.Path]::GetFullPath((Join-Path $RepoRoot $normalized)) } catch { return $false }
    $root = ([IO.Path]::GetFullPath($RepoRoot)).TrimEnd('\', '/') + '\'
    if (-not $full.StartsWith($root, [StringComparison]::OrdinalIgnoreCase)) { return $false }
    return Test-Path -LiteralPath $full -PathType Leaf
}

function Get-FreshIsolatedFailureSet {
    param([string]$SummaryPath, [Parameter(Mandatory = $true)][string]$RepoRoot)
    $SummaryPath = Resolve-FrontendVitestEvidencePath -RepoRoot $RepoRoot -FileName "normal-full-summary.json" -ExplicitPath $SummaryPath
    try { $summary = Get-Content -Raw -LiteralPath $SummaryPath | ConvertFrom-Json } catch { throw "NormalFull evidence is malformed" }
    if ($summary.schemaVersion -ne 1 -or $summary.mode -ne "NormalFull" -or @($summary.records).Count -eq 0) { throw "NormalFull evidence is not a validated schemaVersion 1 summary" }
    $files = [System.Collections.Generic.List[string]]::new()
    foreach ($run in @($summary.records)) {
        foreach ($failure in @($run.failureSet)) {
            $file = ([string]$failure.file).Trim().Replace('\', '/')
            if (-not (Test-FreshIsolatedFilePath -FilePath $file -RepoRoot $RepoRoot)) { throw "NormalFull failure path is unsafe or missing: $file" }
            if (-not $files.Contains($file)) { $files.Add($file) }
        }
    }
    if ($files.Count -eq 0) { throw "NormalFull evidence has no observed failure files" }
    $resolved = [IO.Path]::GetFullPath($SummaryPath)
    [ordered]@{ files = @($files | Sort-Object); summaryPath = $resolved; summarySha256 = (Get-FileHash -Algorithm SHA256 -LiteralPath $resolved).Hash; sourceRecordCount = @($summary.records).Count; derivation = "dynamic union of current NormalFull records.failureSet.file; not a fixed path list" }
}

function Get-FreshIsolatedCommandArguments {
    param([Parameter(Mandatory = $true)][string]$FilePath)
    $normalized = $FilePath.Trim().Replace('\', '/')
    if ($normalized -notmatch '^frontend/tests/.+\.test\.tsx?$' -or $normalized -match '(^|/)\.\.(/|$)' -or $normalized -match '["`;&|<>\x00-\x1F]' -or [IO.Path]::IsPathRooted($normalized)) { throw "FreshIsolated file selector is unsafe" }
    @("--prefix", "frontend", "run", "test", "--", "--run", $normalized.Substring("frontend/".Length))
}

function Get-FreshIsolatedCommandText {
    param([Parameter(Mandatory = $true)][string]$FilePath)
    $args = @(Get-FreshIsolatedCommandArguments -FilePath $FilePath)
    "npm $($args -join ' ')"
}

function Invoke-FreshIsolatedChildRun {
    param([string]$RepoRoot, [string]$FilePath, [string]$RunId, [int]$RepeatIndex, [int]$WatchdogSeconds)
    $result = Invoke-FrontendVitestOwnedProcess -RepoRoot $RepoRoot -ArgumentList (Get-FreshIsolatedCommandArguments -FilePath $FilePath) -WatchdogSeconds $WatchdogSeconds
    $result.command = Get-FreshIsolatedCommandText -FilePath $FilePath
    return $result
}

function Get-FreshDurationSummary {
    param([object[]]$Records)
    $values = @($Records | ForEach-Object { if ($_.durationMs -is [int] -or $_.durationMs -is [long] -or $_.durationMs -is [double]) { [double]$_.durationMs } } | Sort-Object)
    if ($values.Count -eq 0) { return [ordered]@{ min = "unavailable"; max = "unavailable"; median = "unavailable" } }
    $middle = [int][Math]::Floor($values.Count / 2); $median = if ($values.Count % 2) { $values[$middle] } else { ($values[$middle - 1] + $values[$middle]) / 2 }
    [ordered]@{ min = $values[0]; max = $values[$values.Count - 1]; median = $median }
}

function Invoke-FrontendVitestFreshIsolated {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string]$OutputDirectory,
        [Parameter(Mandatory = $true)][int]$Repeat,
        [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [string[]]$FailureFiles,
        [string]$SummaryPath,
        [scriptblock]$ChildInvoker
    )
    if (-not (Test-Path -LiteralPath $OutputDirectory -PathType Container) -or @(Get-ChildItem -LiteralPath $OutputDirectory -Force).Count -gt 0) { throw "FreshIsolated output directory must be an existing empty directory" }
    $source = Get-FreshIsolatedFailureSet -SummaryPath $SummaryPath -RepoRoot $RepoRoot
    $files = if ($null -ne $FailureFiles -and @($FailureFiles).Count -gt 0) { @($FailureFiles | ForEach-Object { $n = $_.Trim().Replace('\', '/'); if (-not (Test-FreshIsolatedFilePath -FilePath $n -RepoRoot $RepoRoot)) { throw "FreshIsolated file path is unsafe or missing: $n" }; $n } | Sort-Object -Unique) } else { @($source.files) }
    if (($files -join "|") -ne (@($source.files) -join "|")) { throw "FreshIsolated files must match the current NormalFull failure union" }
    $records = [System.Collections.Generic.List[object]]::new(); $artifactIndex = [System.Collections.Generic.List[object]]::new(); $recordPaths = [System.Collections.Generic.List[string]]::new(); $hardFailure = $false; $blockedFailure = $false; $fileIndex = 0
    foreach ($file in $files) {
        $fileIndex++
        for ($repeatIndex = 1; $repeatIndex -le $Repeat; $repeatIndex++) {
            $runId = [Guid]::NewGuid().ToString("N"); $runDirectory = Join-Path $OutputDirectory ("run-{0:D3}-{1:D2}-{2}" -f $fileIndex, $repeatIndex, $runId); New-Item -ItemType Directory -Path $runDirectory -Force | Out-Null
            try { $child = if ($null -ne $ChildInvoker) { & $ChildInvoker -RepoRoot $RepoRoot -FilePath $file -RunId $runId -RepeatIndex $repeatIndex -WatchdogSeconds $WatchdogSeconds } else { Invoke-FreshIsolatedChildRun -RepoRoot $RepoRoot -FilePath $file -RunId $runId -RepeatIndex $repeatIndex -WatchdogSeconds $WatchdogSeconds } } catch { $child = [ordered]@{ command = Get-FreshIsolatedCommandText -FilePath $file; startedAtUtc = [DateTime]::UtcNow.ToString("o"); completedAtUtc = [DateTime]::UtcNow.ToString("o"); durationMs = 0; stdout = ""; stderr = $_.Exception.Message; exitCode = "unavailable"; runnerPid = "unavailable"; childPids = "unavailable"; peakWorkingSetBytes = "unavailable"; externalWatchdogTriggered = $false; cleanupOutcome = "unsafe-termination-refused"; cleanupError = $_.Exception.Message; startError = $_.Exception.Message } }
            $childCommand = Get-FreshIsolatedCommandText -FilePath $file; $childCommandValue = Get-FrontendVitestChildResultValue -Child $child -Name "command" -Default $null; if ($null -ne $childCommandValue) { $childCommand = [string]$childCommandValue }
            $record = ConvertFrom-FrontendVitestOutput -Stdout $child.stdout -Stderr $child.stderr -ExitCode $child.exitCode -DurationMs ([int]$child.durationMs) -Phase "FreshIsolated" -RunId $runId -RepeatIndex $repeatIndex -DiscoveredFiles @($file) -ExternalWatchdogTriggered ([bool]$child.externalWatchdogTriggered) -WatchdogSeconds $WatchdogSeconds -RepoRoot $RepoRoot
            if ($record.classification -eq "unknown" -and ([string]$child.stderr -match "(?i)spawn EPERM|environment boundary")) { $record.classification = "environment"; $record.result = "failed"; foreach ($failure in @($record.failures)) { $failure.kind = "environment"; $failure.classification = "environment"; $failure.message = "environment boundary: spawn EPERM" } }
            $childStartError = [string](Get-FrontendVitestChildResultValue -Child $child -Name "startError" -Default "")
            $childRuntimeError = [string](Get-FrontendVitestChildResultValue -Child $child -Name "runtimeError" -Default "")
            $childCleanupOutcome = [string](Get-FrontendVitestChildResultValue -Child $child -Name "cleanupOutcome" -Default "not-required")
            $childCleanupError = [string](Get-FrontendVitestChildResultValue -Child $child -Name "cleanupError" -Default "")
            $record.command = $childCommand; $record.startedAtUtc = $child.startedAtUtc; $record.completedAtUtc = $child.completedAtUtc; $record.fileIndex = $fileIndex; $record.targetFile = $file; $record.sourceSummaryPath = $source.summaryPath; $record.sourceSummarySha256 = $source.summarySha256; $record.settings = [ordered]@{ isolate = "fresh-process"; pool = "unmodified-default"; workers = "unmodified-default"; order = "unmodified-default"; groupFiles = @($file); processLifetime = "one-owned-process-per-run"; externalWatchdogSeconds = $WatchdogSeconds; testTimeoutMs = 5000; hookTimeoutMs = 10000; commandArguments = @(Get-FreshIsolatedCommandArguments -FilePath $file); unchanged = $true }; $record.resources = [ordered]@{ runnerPid = $child.runnerPid; childPids = $child.childPids; workerCount = "unavailable"; peakWorkingSetBytes = $child.peakWorkingSetBytes; heapUsedBytes = "unavailable"; processExitCode = $child.exitCode; externalWatchdogTriggered = $child.externalWatchdogTriggered; cleanupOutcome = $childCleanupOutcome; cleanupError = $childCleanupError }
            $stdoutName = "stdout.txt"; $stderrName = "stderr.txt"; $jsonName = "record.json"; Set-Content -LiteralPath (Join-Path $runDirectory $stdoutName) -Value $child.stdout -Encoding utf8; Set-Content -LiteralPath (Join-Path $runDirectory $stderrName) -Value $child.stderr -Encoding utf8; $record.stdout = "artifact:$stdoutName"; $record.stderr = "artifact:$stderrName"; $record.artifactPaths = @($jsonName, $stdoutName, $stderrName); $record.artifactHashes = [ordered]@{ stdout = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $runDirectory $stdoutName)).Hash; stderr = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $runDirectory $stderrName)).Hash }; $jsonPath = Join-Path $runDirectory $jsonName; $record | ConvertTo-Json -Depth 50 | Set-Content -LiteralPath $jsonPath -Encoding utf8; [void]$records.Add($record); [void]$recordPaths.Add($jsonPath)
            foreach ($name in @($jsonName, $stdoutName, $stderrName)) { $path = Join-Path $runDirectory $name; [void]$artifactIndex.Add([ordered]@{ path = $path; sha256 = (Get-FileHash -Algorithm SHA256 -LiteralPath $path).Hash }) }
            if ($childStartError -ne "" -or $childRuntimeError -ne "" -or $childCleanupOutcome -eq "unsafe-termination-refused" -or $record.result -ne "pass" -or $record.classification -ne "pass") {
                $hardFailure = $true
                if ($childStartError -ne "" -or $childRuntimeError -ne "" -or $childCleanupOutcome -eq "unsafe-termination-refused" -or $record.result -eq "blocked") { $blockedFailure = $true }
            }
        }
    }
    $recordArray = @($records.ToArray()); $perFile = foreach ($group in @($recordArray | Group-Object { $_["targetFile"] })) { $items = @($group.Group); [ordered]@{ file = $group.Name; repeatCount = $items.Count; passCount = @($items | Where-Object result -eq "pass").Count; failedCount = @($items | Where-Object result -eq "failed").Count; blockedCount = @($items | Where-Object result -eq "blocked").Count; timeoutCount = @($items | Where-Object classification -eq "test-timeout").Count; assertionCount = @($items | Where-Object classification -eq "assertion").Count; workerHangCount = @($items | Where-Object classification -eq "worker-hang").Count; resourceLeakCount = @($items | Where-Object classification -eq "resource leak").Count; environmentCount = @($items | Where-Object classification -eq "environment").Count; unknownCount = @($items | Where-Object classification -eq "unknown").Count; durationMs = Get-FreshDurationSummary -Records $items } }
    $orderedRecords = @($recordArray); $sequential = $true; for ($i = 1; $i -lt $orderedRecords.Count; $i++) { if ([DateTime]::Parse($orderedRecords[$i].startedAtUtc) -lt [DateTime]::Parse($orderedRecords[$i - 1].completedAtUtc)) { $sequential = $false } }
    $allPresent = $true; $allHashesMatch = $true
    foreach ($artifact in @($artifactIndex)) {
        if (-not (Test-Path -LiteralPath $artifact.path -PathType Leaf)) { $allPresent = $false; continue }
        try { if ((Get-FileHash -Algorithm SHA256 -LiteralPath $artifact.path).Hash -ne $artifact.sha256) { $allHashesMatch = $false } } catch { $allHashesMatch = $false }
    }
    if (-not $allPresent -or -not $allHashesMatch) { $hardFailure = $true; $blockedFailure = $true }
    if (-not $sequential) { $hardFailure = $true; $blockedFailure = $true }
    $summary = [ordered]@{ schemaVersion = 1; runId = [Guid]::NewGuid().ToString("N"); mode = "FreshIsolated"; repeat = $Repeat; command = "npm --prefix frontend run test -- --run <validated-file>"; sourceSummaryPath = $source.summaryPath; sourceSummarySha256 = $source.summarySha256; sourceFileCount = @($files).Count; outputDirectory = $OutputDirectory; sequential = $sequential; recordCount = @($records).Count; perFile = @($perFile); artifactIntegrity = [ordered]@{ checked = @($artifactIndex).Count; allPresent = $allPresent; allHashesMatch = $allHashesMatch; rawLocation = $OutputDirectory }; records = @($recordArray | ForEach-Object { [ordered]@{ runId = $_.runId; file = $_.targetFile; repeatIndex = $_.repeatIndex; result = $_.result; classification = $_.classification; durationMs = $_.durationMs } }); result = if (-not $hardFailure) { "completed" } elseif ($blockedFailure) { "blocked" } else { "failed" } }
    $summaryPathOut = Join-Path $OutputDirectory "fresh-isolated-summary.json"; $summary | ConvertTo-Json -Depth 50 | Set-Content -LiteralPath $summaryPathOut -Encoding utf8
    $pointer = [ordered]@{ schemaVersion = 1; mode = "FreshIsolated"; repeat = $Repeat; fileCount = @($files).Count; recordCount = @($records).Count; runIds = @($recordArray | ForEach-Object runId); recordPaths = @($recordPaths); artifactPaths = @($recordPaths); summaryPath = $summaryPathOut; sourceSummaryPath = $source.summaryPath; sourceSummarySha256 = $source.summarySha256; sequential = $sequential; result = $summary.result; artifactIntegrity = $summary.artifactIntegrity }
    Write-Output ($pointer | ConvertTo-Json -Depth 20 -Compress)
    if ($hardFailure) { throw "FreshIsolated evidence is $($summary.result); summary preserved at $summaryPathOut" }
}

Export-ModuleMember -Function Get-FreshIsolatedFailureSet, Get-FreshIsolatedCommandArguments, Invoke-FrontendVitestFreshIsolated
