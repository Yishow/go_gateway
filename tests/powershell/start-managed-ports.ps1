# Get-ManagedPorts 合約測試（含 range 展開回歸：逗號優先於 .. 的陷阱）
# 執行：pwsh -NoProfile -File tests/powershell/start-managed-ports.ps1
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
. (Join-Path $repoRoot "scripts/start-port-utils.ps1")

$exitCode = 0
try {
    # 情境 1：預設 Port=8080、FrontendDevPort=5173、ModbusSharePort=5020 → 8080..8090 + 5173 + 5020 去重
    $Port = 8080
    $script:FrontendDevPort = "5173"
    $script:ModbusSharePort = "5020"
    $ports = @(Get-ManagedPorts | ForEach-Object { [int]$_ })
    if ($ports.Count -ne 13) { throw "預期 13 個端口（8080..8090 共 11 個 + 5173 + 5020），got $($ports.Count): $($ports -join ',')" }
    if ($ports -notcontains 5173) { throw "應包含前端端口 5173" }
    if ($ports -notcontains 5020) { throw "應包含 Modbus share 端口 5020" }
    if ($ports -notcontains 8080 -or $ports -notcontains 8090) { throw "應包含 8080 與 8090" }
    if (@($ports | Group-Object | Where-Object { $_.Count -gt 1 }).Count -gt 0) { throw "不得有重複端口" }

    # 情境 2：自訂 Port 不在 8080..8090 時也要納入
    $Port = 3333
    $script:FrontendDevPort = "4173"
    $script:ModbusSharePort = "15020"
    $ports2 = @(Get-ManagedPorts | ForEach-Object { [int]$_ })
    if ($ports2.Count -ne 14) { throw "預期 14 個端口（8080..8090 共 11 個 + 3333 + 4173 + 15020），got $($ports2.Count)" }
    if ($ports2 -notcontains 3333 -or $ports2 -notcontains 4173 -or $ports2 -notcontains 15020) { throw "應包含自訂 Port、前端端口與 Modbus share 端口" }

    # 情境 3：非法輸入（空字串/null）被略過不炸
    $Port = 8080
    $script:FrontendDevPort = $null
    $script:ModbusSharePort = $null
    $ports3 = @(Get-ManagedPorts | ForEach-Object { [int]$_ })
    if ($ports3 -notcontains 8080) { throw "null 前端/Modbus 端口不應影響後端端口" }

    Write-Host "PASS: start-managed-ports.ps1"
}
catch {
    Write-Host "FAIL: $($_.Exception.Message)" -ForegroundColor Red
    $exitCode = 1
}
exit $exitCode
