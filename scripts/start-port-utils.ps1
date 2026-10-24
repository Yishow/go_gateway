# Port helpers for start.ps1. An occupied port is an actionable error; the
# launcher never kills an unrelated process to make room for its services.

function Get-ManagedPorts {
    $ports = New-Object 'System.Collections.Generic.List[int]'
    $candidates = @($script:Port, $script:FrontendDevPort, $script:ModbusSharePort) + (8080..8090)
    foreach ($candidate in $candidates) {
        if ($null -eq $candidate) { continue }
        $text = ([string]$candidate).Trim()
        if ($text -notmatch '^\d+$') { continue }
        $number = [int]$text
        if ($number -lt 1 -or $number -gt 65535) { continue }
        if (-not $ports.Contains($number)) { [void]$ports.Add($number) }
    }
    return $ports
}

function Get-ProcessExecutablePath {
    param([AllowNull()][System.Diagnostics.Process]$Process)
    if ($null -eq $Process) { return $null }
    try {
        if ($Process.Path) { return $Process.Path }
    }
    catch { }
    try {
        return (Get-CimInstance Win32_Process -Filter "ProcessId = $($Process.Id)" -ErrorAction Stop).ExecutablePath
    }
    catch { return $null }
}

function Test-ManagedProcess {
    param([AllowNull()][System.Diagnostics.Process]$Process)
    if ($null -eq $Process) { return $false }
    $path = Get-ProcessExecutablePath -Process $Process
    if ([string]::IsNullOrWhiteSpace($path)) { return $false }
    $root = [regex]::Escape([string]$script:ROOT_DIR)
    return $path -match "^$root[\\/].*(?:gateway|test[_-]?ui|node(?:\.exe)?)"
}

function Stop-ProcessByPort {
    param(
        [Parameter(Mandatory = $true)][int]$Port,
        [switch]$Force
    )
    $process = Get-ProcessByPort -Port $Port
    if ($process) {
        Write-Error "拒絕終止端口 $Port 的外部進程 $($process.ProcessName) (PID: $($process.Id))"
        return $false
    }
    return $true
}

function Clear-PortForService {
    param(
        [Parameter(Mandatory = $true)][int]$Port,
        [switch]$AutoKill
    )
    if (-not (Test-PortInUse -Port $Port)) {
        Write-Info "端口 $Port 可用"
        return $true
    }

    $process = Get-ProcessByPort -Port $Port
    if ($process) {
        Write-Error "端口 $Port 已被 $($process.ProcessName) (PID: $($process.Id)) 佔用；為避免終止無關進程，請先手動釋放端口"
    }
    else {
        Write-Error "端口 $Port 已被佔用；無法確認進程，請先手動釋放端口"
    }
    return $false
}
