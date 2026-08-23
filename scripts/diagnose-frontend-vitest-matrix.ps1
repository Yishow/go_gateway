[CmdletBinding()]
param(
    [ValidateSet("Environment", "NormalFull", "FreshIsolated", "DefaultGroup", "SingleWorkerGroup", "OneVariable")]
    [string]$Mode = "Environment",
    [ValidateRange(1, 2147483647)]
    [int]$Repeat = 1,
    [string]$OutputPath,
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
function Get-FrontendVitestValue {
    param(
        [AllowNull()][object]$Object,
        [Parameter(Mandatory = $true)][string]$Name
    )
    if ($null -eq $Object) {
        return $null
    }
    if ($Object -is [hashtable]) {
        return $Object[$Name]
    }
    $property = $Object.PSObject.Properties[$Name]
    if ($null -ne $property) {
        return $property.Value
    }
    return $null
}
function ConvertTo-FrontendVitestBoolean {
    param([AllowNull()][object]$Value)
    if ($Value -is [bool]) {
        return $Value
    }
    if ($Value -is [string]) {
        if ($Value -ieq "true") { return $true }
        if ($Value -ieq "false") { return $false }
    }
    if ($Value -is [int] -or $Value -is [long] -or $Value -is [double]) {
        if ($Value -eq 1) { return $true }
        if ($Value -eq 0) { return $false }
    }
    return $false
}
function Resolve-FrontendVitestClassification {
    param(
        [Parameter(Mandatory = $true)][AllowNull()][object]$Evidence
    )
    $signals = Get-FrontendVitestValue -Object $Evidence -Name "signals"
    $signalAliases = [ordered]@{
        "assertion" = @("assertion", "assertionFailure")
        "test-timeout" = @("testTimeout", "test-timeout", "testTimeoutDetected")
        "worker-hang" = @("workerHang", "worker-hang", "externalWatchdog")
        "resource leak" = @("resourceLeak", "resource-leak", "resource leak")
        "environment" = @("environment", "environmentEffect")
    }
    $active = [System.Collections.Generic.List[string]]::new()
    foreach ($classification in $signalAliases.Keys) {
        foreach ($alias in $signalAliases[$classification]) {
            if (ConvertTo-FrontendVitestBoolean (Get-FrontendVitestValue -Object $signals -Name $alias)) {
                $active.Add($classification)
                break
            }
        }
    }
    $kind = "unknown"
    if ($active.Count -eq 1) {
        $kind = $active[0]
    } elseif ($active.Count -eq 0 -and (ConvertTo-FrontendVitestBoolean (Get-FrontendVitestValue -Object $Evidence -Name "pass"))) {
        $kind = "pass"
    }
    $result = if ($kind -eq "pass") { "pass" } elseif ($kind -eq "unknown") { "blocked" } else { "failed" }
    [pscustomobject][ordered]@{
        classification = $kind
        result = $result
        blocked = ($kind -eq "unknown")
        activeSignals = @($active)
    }
}
function New-FrontendVitestFailureRecord {
    param(
        [Parameter(Mandatory = $true)][AllowNull()][object]$Evidence
    )
    $classification = Resolve-FrontendVitestClassification -Evidence $Evidence
    $resources = Get-FrontendVitestValue -Object $Evidence -Name "resources"
    if ($null -eq $resources) {
        $resources = [ordered]@{
            runnerPid = "unavailable"
            childPids = "unavailable"
            workerCount = "unavailable"
            peakWorkingSetBytes = "unavailable"
            heapUsedBytes = "unavailable"
        }
    }
    $exitCode = Get-FrontendVitestValue -Object $Evidence -Name "processExitCode"
    if ($null -eq $exitCode) { $exitCode = Get-FrontendVitestValue -Object $Evidence -Name "exitCode" }
    if ($null -eq $exitCode) { $exitCode = "unavailable" }
    [pscustomobject][ordered]@{
        file = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "file")) { Get-FrontendVitestValue -Object $Evidence -Name "file" } else { "unavailable" }
        title = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "title")) { Get-FrontendVitestValue -Object $Evidence -Name "title" } else { "unavailable" }
        kind = $classification.classification
        classification = $classification.classification
        message = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "message")) { Get-FrontendVitestValue -Object $Evidence -Name "message" } else { "unavailable" }
        timeoutMs = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "timeoutMs")) { Get-FrontendVitestValue -Object $Evidence -Name "timeoutMs" } else { "unavailable" }
        phase = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "phase")) { Get-FrontendVitestValue -Object $Evidence -Name "phase" } else { "unavailable" }
        stdout = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "stdout")) { Get-FrontendVitestValue -Object $Evidence -Name "stdout" } else { "" }
        stderr = if ($null -ne (Get-FrontendVitestValue -Object $Evidence -Name "stderr")) { Get-FrontendVitestValue -Object $Evidence -Name "stderr" } else { "" }
        exitCode = $exitCode
        resources = $resources
        result = $classification.result
        blocked = $classification.blocked
    }
}
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
function Read-FrontendVitestEnvironmentEvidence {
    $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
    $evidencePath = Join-Path $repoRoot "openspec\changes\diagnose-frontend-full-suite-timeouts\evidence\environment-capture.json"
    if (Test-Path -LiteralPath $evidencePath -PathType Leaf) {
        try { return (Get-Content -Raw -LiteralPath $evidencePath | ConvertFrom-Json) } catch { throw "task 1 environment evidence is malformed" }
    }
    $testsRoot = Join-Path $repoRoot "frontend\tests"
    if (-not (Test-Path -LiteralPath $testsRoot -PathType Container)) { throw "frontend test directory is missing" }
    $files = @(Get-ChildItem -LiteralPath $testsRoot -Recurse -File | Where-Object { $_.Name -match '\.test\.tsx?$' } | ForEach-Object {
        [IO.Path]::GetRelativePath($repoRoot, $_.FullName).Replace('\', '/')
    } | Sort-Object)
    return [pscustomobject]@{
        discovered = [pscustomobject]@{ files = $files; fileCount = $files.Count; tests = "unavailable"; passed = "unavailable"; failed = "unavailable"; timedOut = "unavailable"; identityStatement = "This is a dynamic current discovery inventory; it is not a fixed timeout identity or permanent failure set." }
    }
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
    if ($null -eq $discovered) { throw "environment evidence has no discovered section" }
    $identityStatement = Get-FrontendVitestValue -Object $discovered -Name "identityStatement"
    if ([string]::IsNullOrWhiteSpace([string]$identityStatement)) { $identityStatement = "This is a dynamic current discovery inventory; it is not a fixed timeout identity or permanent failure set." }
    $files = @(Get-FrontendVitestValue -Object $discovered -Name "files")
    $stdout = "Environment capture completed; no Vitest child process was started."
    $stderr = ""
    [ordered]@{
        schemaVersion = 1
        runId = [Guid]::NewGuid().ToString("N")
        mode = "Environment"
        repeatIndex = $RepeatIndex
        command = "Environment capture using task 1 evidence; no Vitest execution"
        settings = [ordered]@{
            isolate = "unavailable"
            pool = "unavailable"
            workers = "unavailable"
            order = "unavailable"
            groupFiles = @()
            externalWatchdogSeconds = $WatchdogSeconds
            testTimeoutMs = 5000
            hookTimeoutMs = 10000
        }
        startedAtUtc = $started.ToString("o")
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
            childPids = "unavailable"
            workerCount = "unavailable"
            peakWorkingSetBytes = "unavailable"
            heapUsedBytes = "unavailable"
            processExitCode = 0
            externalWatchdogTriggered = $false
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
        Invoke-FrontendVitestGroup -RepoRoot $repoRoot -OutputDirectory $outputDirectory -Repeat $Repeat -WatchdogSeconds $WatchdogSeconds -Mode $Mode
        return
    }
    if ($Mode -eq "FreshIsolated") {
        if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "FreshIsolated derives files from NormalFull evidence; GroupFiles is not accepted" }
        $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
        $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
        Invoke-FrontendVitestFreshIsolated -RepoRoot $repoRoot -OutputDirectory $outputDirectory -Repeat $Repeat -WatchdogSeconds $WatchdogSeconds
        return
    }
    if ($Mode -eq "NormalFull") {
        if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "NormalFull mode does not accept GroupFiles" }
        $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
        $environmentEvidence = Read-FrontendVitestEnvironmentEvidence
        $discovered = Get-FrontendVitestValue -Object $environmentEvidence -Name "discovered"
        if ($null -eq $discovered) { throw "task 1 environment evidence has no discovered section" }
        $discoveredFiles = @(Get-FrontendVitestValue -Object $discovered -Name "files")
        if ($discoveredFiles.Count -eq 0) { throw "task 1 environment evidence has no discovered files" }
        $repoRoot = [IO.Path]::GetFullPath((Split-Path -Parent $PSScriptRoot))
        Invoke-FrontendVitestNormalFull -RepoRoot $repoRoot -OutputDirectory $outputDirectory -Repeat $Repeat -WatchdogSeconds $WatchdogSeconds -DiscoveredFiles $discoveredFiles
        return
    }
    if ($Mode -ne "Environment") { throw "Mode '$Mode' is not implemented; only Environment, NormalFull, and FreshIsolated are available" }
    if ($null -ne $GroupFiles -and @($GroupFiles).Count -gt 0) { throw "Environment mode does not accept GroupFiles" }
    $outputDirectory = Resolve-FrontendVitestOutputDirectory -RequestedPath $OutputPath
    $environmentEvidence = Read-FrontendVitestEnvironmentEvidence
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
