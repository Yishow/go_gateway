# BuildSingle must keep child output out of its boolean return value.
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$env:GATEWAY_START_PS1_LIBRARY_ONLY = "1"
. (Join-Path $repoRoot "start.ps1")
Remove-Item Env:\GATEWAY_START_PS1_LIBRARY_ONLY -ErrorAction SilentlyContinue

$testExitCode = 0
$mockRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("start-build-failure-" + [System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $mockRoot -Force | Out-Null
try {
    # Direct function call: mocked node writes output and exits nonzero.
    function node { Write-Output "mock builder output"; $global:LASTEXITCODE = 17 }
    if (Build-SingleBinary -TargetOS "linux" -TargetArch "amd64") { throw "Build-SingleBinary must return false on child failure" }
    Remove-Item Function:\node -ErrorAction SilentlyContinue

    # CLI call: the same mocked child must make start.ps1 exit 1.
    if (Test-IsWindows) {
        $nodePath = Join-Path $mockRoot "node.cmd"
        [System.IO.File]::WriteAllText($nodePath, "@echo mock builder output`r`n@exit /b 17`r`n", [System.Text.Encoding]::ASCII)
    }
    else {
        $nodePath = Join-Path $mockRoot "node"
        [System.IO.File]::WriteAllText($nodePath, "#!/bin/sh`necho mock builder output`nexit 17`n", [System.Text.Encoding]::ASCII)
        & chmod +x $nodePath
    }
    $pwsh = (Get-Command pwsh -ErrorAction SilentlyContinue).Source
    if ([string]::IsNullOrWhiteSpace($pwsh)) { $pwsh = (Get-Command powershell -ErrorAction Stop).Source }
    $oldPath = $env:PATH
    try {
        $env:PATH = $mockRoot + [System.IO.Path]::PathSeparator + $oldPath
        $childOutput = @(& $pwsh -NoProfile -File (Join-Path $repoRoot "start.ps1") -BuildSingle -OS linux -Arch amd64 2>&1)
        $childCode = [int]$LASTEXITCODE
    }
    finally { $env:PATH = $oldPath }
    if ($childCode -ne 1) { throw "CLI should exit 1 on builder failure, got ${childCode}: $($childOutput -join ' | ')" }
    if (($childOutput -join "`n") -notmatch "mock builder output") { throw "CLI should surface mocked builder output" }
    Write-Host "PASS: start-build-single-failure.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $testExitCode = 1
}
finally { Remove-Item -LiteralPath $mockRoot -Recurse -Force -ErrorAction SilentlyContinue }
exit $testExitCode
