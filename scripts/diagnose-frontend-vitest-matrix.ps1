[CmdletBinding()]
param(
    [ValidateSet("Environment", "NormalFull", "FreshIsolated", "DefaultGroup", "SingleWorkerGroup", "OneVariable")]
    [string]$Mode = "Environment",
    [ValidateRange(1, 2147483647)]
    [int]$Repeat = 1,
    [string]$OutputPath,
    [string]$SummaryPath,
    [string[]]$GroupFiles,
    [ValidateRange(1, 2147483647)]
    [int]$WatchdogSeconds = 600
)
$helperPath = Join-Path $PSScriptRoot "lib\FrontendVitestMatrix.psm1"
if (-not (Test-Path -LiteralPath $helperPath -PathType Leaf)) {
    throw "NormalFull helper module is missing"
}
Import-Module -Name $helperPath -Force -Global -ErrorAction Stop
$freshHelperPath = Join-Path $PSScriptRoot "lib\FrontendVitestFreshIsolated.psm1"
if (-not (Test-Path -LiteralPath $freshHelperPath -PathType Leaf)) { throw "FreshIsolated helper module is missing" }
Import-Module -Name $freshHelperPath -Force -ErrorAction Stop
$groupsHelperPath = Join-Path $PSScriptRoot "lib\FrontendVitestGroups.psm1"
if (-not (Test-Path -LiteralPath $groupsHelperPath -PathType Leaf)) { throw "group helper module is missing" }
Import-Module -Name $groupsHelperPath -Force -Global -ErrorAction Stop
function Resolve-FrontendVitestOutputDirectory {
    param([string]$RequestedPath)
    $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
    $workspaceRoot = [IO.Path]::GetFullPath((Split-Path -Parent $repoRoot))
    $isDefault = [string]::IsNullOrWhiteSpace($RequestedPath)
    if ($isDefault) {
        $base = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
        do {
            $candidate = Join-Path $base ("frontend-vitest-matrix-" + [Guid]::NewGuid().ToString("N"))
        } while (Test-Path -LiteralPath $candidate)
        New-Item -ItemType Directory -Path $candidate -Force | Out-Null
        return $candidate
    }
    if (-not [IO.Path]::IsPathRooted($RequestedPath)) { throw "OutputPath must be an absolute path" }
    if ($RequestedPath -match "(^|[\\/])\.\.([\\/]|$)") { throw "OutputPath cannot contain path traversal" }
    try { $resolved = [IO.Path]::GetFullPath($RequestedPath) } catch { throw "OutputPath is not resolvable" }
    $canonical = $resolved.TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar)
    $forbidden = @($repoRoot, $workspaceRoot, [IO.Path]::GetPathRoot($resolved)) | ForEach-Object {
        ([IO.Path]::GetFullPath($_)).TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar)
    }
    if ($forbidden -contains $canonical) { throw "OutputPath cannot be a workspace, repository, or disk root" }
    $parent = [IO.Directory]::GetParent($resolved)
    if ($null -eq $parent -or -not (Test-Path -LiteralPath $parent.FullName -PathType Container)) {
        throw "OutputPath parent does not exist"
    }
    if (Test-Path -LiteralPath $resolved) {
        if (-not (Test-Path -LiteralPath $resolved -PathType Container)) { throw "OutputPath is not a directory" }
        if (@(Get-ChildItem -LiteralPath $resolved -Force).Count -gt 0) { throw "OutputPath must be empty when it already exists" }
    } else {
        New-Item -ItemType Directory -Path $resolved -Force | Out-Null
    }
    return $resolved
}
function New-FrontendVitestEnvironmentRecord {
    param(
        [Parameter(Mandatory = $true)][int]$RepeatIndex,
        [Parameter(Mandatory = $true)][string]$OutputDirectory,
        [Parameter(Mandatory = $true)][int]$WatchdogSeconds,
        [Parameter(Mandatory = $true)]$EnvironmentEvidence
    )
    $started = [DateTime]::UtcNow
    $discovered = Get-FrontendVitestValue -Object $EnvironmentEvidence -Name "discovered"
    if ($null -eq $discovered) { $discovered = $EnvironmentEvidence }
    $identityStatement = Get-FrontendVitestValue -Object $discovered -Name "identityStatement"
    if ([string]::IsNullOrWhiteSpace([string]$identityStatement)) { $identityStatement = "This is a dynamic current discovery inventory; it is not a fixed timeout identity or permanent failure set." }
    $files = @(Get-FrontendVitestValue -Object $discovered -Name "files")
    $stdout = "Current filesystem discovery completed; no Vitest child process was started."
    $stderr = ""
    [ordered]@{
        schemaVersion = 1
        runId = [Guid]::NewGuid().ToString("N")
        mode = "Environment"
        repeatIndex = $RepeatIndex
        command = "Current filesystem discovery; no Vitest child process was started"
        settings = [ordered]@{
            isolate = "not-applicable/no-child"
            pool = "not-applicable/no-child"
            workers = "not-applicable/no-child"
            order = "not-applicable/no-child"
            groupFiles = @()
            processLifetime = "not-applicable/no-child"
            externalWatchdogSeconds = "not-applicable/no-child"
            testTimeoutMs = "not-applicable/no-child"
            hookTimeoutMs = "not-applicable/no-child"
            commandArguments = @()
            unchanged = $true
        }
        startedAtUtc = $started.ToString("o")
        completedAtUtc = [DateTime]::UtcNow.ToString("o")
        durationMs = [long]0
        discovered = [ordered]@{
            files = $files
            fileCount = $files.Count
            tests = if ($null -ne (Get-FrontendVitestValue -Object $discovered -Name "tests")) { Get-FrontendVitestValue -Object $discovered -Name "tests" } else { "unavailable" }
            passed = if ($null -ne (Get-FrontendVitestValue -Object $discovered -Name "passed")) { Get-FrontendVitestValue -Object $discovered -Name "passed" } else { "unavailable" }
            failed = if ($null -ne (Get-FrontendVitestValue -Object $discovered -Name "failed")) { Get-FrontendVitestValue -Object $discovered -Name "failed" } else { "unavailable" }
            timedOut = if ($null -ne (Get-FrontendVitestValue -Object $discovered -Name "timedOut")) { Get-FrontendVitestValue -Object $discovered -Name "timedOut" } else { "unavailable" }
            identityStatement = $identityStatement
        }
        failures = @()
        dynamicFailureSet = [ordered]@{ status = "not-run"; observedFiles = @(); fixedIdentity = $false; identityStatement = $identityStatement }
        resources = [ordered]@{
            runnerPid = "unavailable"
            childPids = @()
            workerCount = "unavailable"
            peakWorkingSetBytes = "unavailable"
            heapUsedBytes = "unavailable"
            processExitCode = "not-applicable/no-child"
            externalWatchdogTriggered = $false
            cleanupOutcome = "not-applicable/no-child"
            cleanupError = ""
        }
        result = "completed"
        classification = "unavailable"
        artifactPaths = @()
        stdout = $stdout
        stderr = $stderr
        outputDirectory = $OutputDirectory
    }
}
function Invoke-FrontendVitestMatrix {
    if ($Mode -in @("DefaultGroup", "SingleWorkerGroup")) {
        if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "$Mode derives files from NormalFull evidence; GroupFiles is not accepted" }
        $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
        $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
        Invoke-FrontendVitestGroup -RepoRoot $repoRoot -OutputDirectory $outputDirectory -Repeat $Repeat -WatchdogSeconds $WatchdogSeconds -Mode $Mode -SummaryPath $SummaryPath
        return
    }
    if ($Mode -eq "FreshIsolated") {
        if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "FreshIsolated derives files from NormalFull evidence; GroupFiles is not accepted" }
        $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
        $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
        Invoke-FrontendVitestFreshIsolated -RepoRoot $repoRoot -OutputDirectory $outputDirectory -Repeat $Repeat -WatchdogSeconds $WatchdogSeconds -SummaryPath $SummaryPath
        return
    }
    if ($Mode -eq "NormalFull") {
        if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "NormalFull mode does not accept GroupFiles" }
        if (-not [string]::IsNullOrWhiteSpace($SummaryPath)) { throw "NormalFull mode does not accept SummaryPath; use the owned output summaryPath" }
        $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
        $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
        $discoveredFiles = @(Get-FrontendVitestValue -Object (Get-FrontendVitestCurrentDiscovery -RepoRoot $repoRoot) -Name "files")
        Invoke-FrontendVitestNormalFull -RepoRoot $repoRoot -OutputDirectory $outputDirectory -Repeat $Repeat -WatchdogSeconds $WatchdogSeconds -DiscoveredFiles $discoveredFiles
        return
    }
    if ($Mode -eq "OneVariable") {
        if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "OneVariable mode does not accept GroupFiles" }
        if (-not [string]::IsNullOrWhiteSpace($SummaryPath)) { throw "OneVariable mode does not accept SummaryPath; no setting was changed" }
        $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
        $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
        Invoke-FrontendVitestOneVariable -OutputDirectory $outputDirectory -WatchdogSeconds $WatchdogSeconds
        return
    }
    if ($Mode -ne "Environment") { throw "Mode '$Mode' is not implemented; supported modes are Environment, NormalFull, FreshIsolated, DefaultGroup, SingleWorkerGroup, and OneVariable" }
    if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "Environment mode does not accept GroupFiles" }
    if (-not [string]::IsNullOrWhiteSpace($SummaryPath)) { throw "Environment mode does not accept SummaryPath; current filesystem discovery is used" }
    $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
    $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
    $environmentEvidence = Get-FrontendVitestCurrentDiscovery -RepoRoot $repoRoot
    $records = [System.Collections.Generic.List[object]]::new()
    for ($index = 1; $index -le $Repeat; $index++) {
        $records.Add((New-FrontendVitestEnvironmentRecord -RepeatIndex $index -OutputDirectory $outputDirectory -WatchdogSeconds $WatchdogSeconds -EnvironmentEvidence $environmentEvidence))
    }

    $artifactPaths = [System.Collections.Generic.List[string]]::new()
    for ($index = 1; $index -le $records.Count; $index++) {
        $suffix = if ($index -eq 1) { "" } else { "-$index" }
        $jsonName = "environment$suffix.json"
        $stdoutName = "stdout$suffix.txt"
        $stderrName = "stderr$suffix.txt"
        $artifactPaths.Add($jsonName)
        $artifactPaths.Add($stdoutName)
        $artifactPaths.Add($stderrName)
        $records[$index - 1]["artifactPaths"] = @($jsonName, $stdoutName, $stderrName)
        $records[$index - 1] | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath (Join-Path $outputDirectory $jsonName) -Encoding utf8
        Set-Content -LiteralPath (Join-Path $outputDirectory $stdoutName) -Value $records[$index - 1].stdout -Encoding utf8
        Set-Content -LiteralPath (Join-Path $outputDirectory $stderrName) -Value $records[$index - 1].stderr -Encoding utf8
    }
    $records[0]["artifactPaths"] = @($artifactPaths)
    $records[0] | ConvertTo-Json -Depth 30 | Set-Content -LiteralPath (Join-Path $outputDirectory "environment.json") -Encoding utf8
    Write-Output ($records[0] | ConvertTo-Json -Depth 30 -Compress)
}

if ($MyInvocation.InvocationName -ne ".") {
    try {
        Invoke-FrontendVitestMatrix
        exit 0
    } catch {
        Write-Error $_
        exit 1
    }
}
