Set-StrictMode -Version Latest

function Get-FrontendVitestChildResultValue {
    param(
        [AllowNull()][object]$Child,
        [Parameter(Mandatory = $true)][string]$Name,
        [AllowNull()][object]$Default = $null
    )
    if ($null -eq $Child) { return $Default }
    if ($Child -is [System.Collections.IDictionary]) {
        if ($Child.Contains($Name)) {
            $value = $Child[$Name]
            if ($null -ne $value) { return $value }
            return $Default
        }
        return $Default
    }
    $property = $Child.PSObject.Properties[$Name]
    if ($null -ne $property -and $null -ne $property.Value) { return $property.Value }
    return $Default
}

function Invoke-FrontendVitestOwnedProcess {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string[]]$ArgumentList,
        [Parameter(Mandatory = $true)][ValidateRange(1, 2147483647)][int]$WatchdogSeconds
    )

    $started = [DateTime]::UtcNow
    $timer = [Diagnostics.Stopwatch]::StartNew()
    $process = $null
    $stdoutTask = $null
    $stderrTask = $null
    $runnerPid = "unavailable"
    $stdout = ""
    $stderr = ""
    $exitCode = "unavailable"
    $peakWorkingSet = "unavailable"
    $watchdog = $false
    $cleanupOutcome = "not-required"
    $cleanupError = ""
    $startError = ""
    $runtimeError = ""

    try {
        if (@($ArgumentList).Count -eq 0) { throw "npm argument list is empty" }
        foreach ($argument in @($ArgumentList)) {
            if ($null -eq $argument -or ([string]$argument).IndexOf([char]0) -ge 0) { throw "npm argument is invalid" }
        }
        $npm = Get-Command -Name "npm.cmd" -CommandType Application -ErrorAction Stop
        $psi = [Diagnostics.ProcessStartInfo]::new()
        $psi.FileName = if (-not [string]::IsNullOrWhiteSpace($npm.Source)) { $npm.Source } else { $npm.Path }
        foreach ($argument in @($ArgumentList)) { [void]$psi.ArgumentList.Add([string]$argument) }
        $psi.WorkingDirectory = [IO.Path]::GetFullPath($RepoRoot)
        $psi.UseShellExecute = $false
        $psi.CreateNoWindow = $true
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $process = [Diagnostics.Process]::new()
        $process.StartInfo = $psi
        if (-not $process.Start()) { throw "npm process did not start" }
        $runnerPid = $process.Id
        $stdoutTask = $process.StandardOutput.ReadToEndAsync()
        $stderrTask = $process.StandardError.ReadToEndAsync()
        while (-not $process.WaitForExit(250)) {
            if ($timer.Elapsed.TotalSeconds -ge $WatchdogSeconds) {
                $watchdog = $true
                break
            }
        }
        if ($watchdog) {
            try {
                if (-not $process.HasExited) {
                    $process.Kill($true)
                    if (-not $process.WaitForExit(5000)) { throw "owned process tree did not exit after watchdog termination" }
                    $cleanupOutcome = "owned-tree-terminated"
                } else {
                    $cleanupOutcome = "already-exited"
                }
            } catch {
                $cleanupOutcome = "unsafe-termination-refused"
                $cleanupError = $_.Exception.Message
            }
        }
        if ($null -ne $stdoutTask -and $stdoutTask.Wait(5000)) { $stdout = $stdoutTask.Result } else { $stdout = "unavailable" }
        if ($null -ne $stderrTask -and $stderrTask.Wait(5000)) { $stderr = $stderrTask.Result } else { $stderr = "unavailable" }
        if ($process.HasExited) { $exitCode = $process.ExitCode }
        try {
            if ($process.HasExited) { $peakWorkingSet = [long]$process.PeakWorkingSet64 }
        } catch { $peakWorkingSet = "unavailable" }
    } catch {
        if ($null -eq $process -or $runnerPid -eq "unavailable") { $startError = $_.Exception.Message } else { $runtimeError = $_.Exception.Message }
        if ([string]::IsNullOrWhiteSpace($stderr)) { $stderr = $_.Exception.Message }
        if ($null -ne $process -and $runnerPid -ne "unavailable") {
            try {
                if (-not $process.HasExited) {
                    $process.Kill($true)
                    if ($process.WaitForExit(5000)) { $cleanupOutcome = "owned-tree-terminated-after-error" }
                    else { $cleanupOutcome = "unsafe-termination-refused" }
                }
            } catch {
                $cleanupOutcome = "unsafe-termination-refused"
                $cleanupError = $_.Exception.Message
            }
        }
        if ($cleanupOutcome -eq "not-required" -and $watchdog) { $cleanupOutcome = "unsafe-termination-refused"; $cleanupError = $_.Exception.Message }
        if ([string]::IsNullOrWhiteSpace($cleanupError) -and $cleanupOutcome -eq "unsafe-termination-refused") { $cleanupError = $_.Exception.Message }
    } finally {
        $timer.Stop()
        if ($null -ne $process -and $runnerPid -ne "unavailable") {
            try {
                if (-not $process.HasExited) {
                    $process.Kill($true)
                    if ($process.WaitForExit(5000)) { $cleanupOutcome = "owned-tree-terminated-finally" }
                    else { $cleanupOutcome = "unsafe-termination-refused"; if ([string]::IsNullOrWhiteSpace($cleanupError)) { $cleanupError = "owned process tree remained alive during final cleanup" } }
                }
            } catch {
                $cleanupOutcome = "unsafe-termination-refused"
                if ([string]::IsNullOrWhiteSpace($cleanupError)) { $cleanupError = $_.Exception.Message }
            }
        }
        if ($null -ne $process) { $process.Dispose() }
    }

    [ordered]@{
        startedAtUtc = $started.ToString("o")
        completedAtUtc = [DateTime]::UtcNow.ToString("o")
        durationMs = [long]$timer.ElapsedMilliseconds
        stdout = $stdout
        stderr = $stderr
        exitCode = $exitCode
        runnerPid = $runnerPid
        childPids = "unavailable"
        peakWorkingSetBytes = $peakWorkingSet
        externalWatchdogTriggered = $watchdog
        cleanupOutcome = $cleanupOutcome
        cleanupError = $cleanupError
        startError = $startError
        runtimeError = $runtimeError
    }
}

Export-ModuleMember -Function Invoke-FrontendVitestOwnedProcess, Get-FrontendVitestChildResultValue
