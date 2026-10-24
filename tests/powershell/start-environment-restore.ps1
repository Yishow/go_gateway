# Temporary Air environment overrides must restore existing and absent entries.
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$env:GATEWAY_START_PS1_LIBRARY_ONLY = "1"
. (Join-Path $repoRoot "start.ps1")
Remove-Item Env:\GATEWAY_START_PS1_LIBRARY_ONLY -ErrorAction SilentlyContinue

$absentName = "START_PS_TEST_ABSENT"
$existingName = "START_PS_TEST_EXISTING"
Remove-Item -LiteralPath ("Env:{0}" -f $absentName) -ErrorAction SilentlyContinue
Set-Item -LiteralPath ("Env:{0}" -f $existingName) -Value "original"
$exitCode = 0
try {
    Invoke-WithEnvironment -Values @{ $absentName = "temporary"; $existingName = "temporary" } -Action {
        if ($env:START_PS_TEST_ABSENT -ne "temporary" -or $env:START_PS_TEST_EXISTING -ne "temporary") { throw "temporary environment was not applied" }
    }
    if ($null -ne (Get-Item -LiteralPath ("Env:{0}" -f $absentName) -ErrorAction SilentlyContinue)) { throw "absent environment entry was not removed" }
    if ($env:START_PS_TEST_EXISTING -ne "original") { throw "existing environment entry was not restored" }
    Write-Host "PASS: start-environment-restore.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
finally {
    Remove-Item -LiteralPath ("Env:{0}" -f $absentName) -ErrorAction SilentlyContinue
    Remove-Item -LiteralPath ("Env:{0}" -f $existingName) -ErrorAction SilentlyContinue
}
exit $exitCode
