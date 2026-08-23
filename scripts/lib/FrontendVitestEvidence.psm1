Set-StrictMode -Version Latest

$script:FrontendVitestChangeName = "diagnose-frontend-full-suite-timeouts"

function Get-FrontendVitestEvidenceRoot {
    param([Parameter(Mandatory = $true)][string]$RepoRoot)

    $root = [IO.Path]::GetFullPath($RepoRoot)
    $changesRoot = Join-Path $root "openspec\changes"
    if (-not (Test-Path -LiteralPath $changesRoot -PathType Container)) {
        throw "OpenSpec changes directory is missing"
    }
    return $changesRoot
}

function Test-FrontendVitestEvidenceJson {
    param([Parameter(Mandatory = $true)][string]$Path)

    try {
        Get-Content -Raw -LiteralPath $Path | ConvertFrom-Json | Out-Null
        return $true
    } catch {
        return $false
    }
}

function Resolve-FrontendVitestEvidencePath {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string]$FileName,
        [string]$ExplicitPath
    )

    if ($FileName -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]*\.json$') {
        throw "evidence file name is unsafe: $FileName"
    }

    if (-not [string]::IsNullOrWhiteSpace($ExplicitPath)) {
        if (-not [IO.Path]::IsPathRooted($ExplicitPath)) { throw "evidence path must be absolute" }
        try { $resolvedExplicit = [IO.Path]::GetFullPath($ExplicitPath) } catch { throw "evidence path is not resolvable" }
        if (-not (Test-Path -LiteralPath $resolvedExplicit -PathType Leaf)) { throw "explicit evidence file is missing: $resolvedExplicit" }
        if (-not (Test-FrontendVitestEvidenceJson -Path $resolvedExplicit)) { throw "explicit evidence file is malformed: $resolvedExplicit" }
        return $resolvedExplicit
    }

    $changesRoot = Get-FrontendVitestEvidenceRoot -RepoRoot $RepoRoot
    $activeDirectories = @(
        Get-ChildItem -LiteralPath $changesRoot -Directory -ErrorAction Stop |
            Where-Object {
                $_.Name -eq $script:FrontendVitestChangeName -or
                $_.Name -match ("^\d{4}-\d{2}-\d{2}-" + [regex]::Escape($script:FrontendVitestChangeName) + "$")
            }
    )
    if ($activeDirectories.Count -gt 1) { throw "multiple active evidence directories found for $FileName" }
    if ($activeDirectories.Count -eq 1) {
        $activePath = [IO.Path]::GetFullPath((Join-Path $activeDirectories[0].FullName (Join-Path "evidence" $FileName)))
        if (-not (Test-Path -LiteralPath $activePath -PathType Leaf)) { throw "active evidence file is missing: $activePath" }
        if (-not (Test-FrontendVitestEvidenceJson -Path $activePath)) { throw "active evidence file is malformed: $activePath" }
        return $activePath
    }

    $archiveRoot = Join-Path $changesRoot "archive"
    $archives = @()
    if (Test-Path -LiteralPath $archiveRoot -PathType Container) {
        $archives = @(
            Get-ChildItem -LiteralPath $archiveRoot -Directory -ErrorAction Stop |
                Where-Object { $_.Name -match ("^(?<date>\d{4}-\d{2}-\d{2})-" + [regex]::Escape($script:FrontendVitestChangeName) + "(?:-.+)?$") } |
                ForEach-Object {
                    $dateText = $_.Name.Substring(0, 10)
                    $date = [DateTime]::MinValue
                    if (-not [DateTime]::TryParseExact($dateText, "yyyy-MM-dd", [Globalization.CultureInfo]::InvariantCulture, [Globalization.DateTimeStyles]::None, [ref]$date)) { return }
                    $candidate = Join-Path $_.FullName (Join-Path "evidence" $FileName)
                    if (-not (Test-Path -LiteralPath $candidate -PathType Leaf)) { return }
                    if (-not (Test-FrontendVitestEvidenceJson -Path $candidate)) { throw "archive evidence file is malformed: $candidate" }
                    [pscustomobject]@{ Date = $date; Path = [IO.Path]::GetFullPath($candidate); Archive = $_.FullName }
                }
        )
    }
    if ($archives.Count -eq 0) { throw "no active or archived evidence file found for $FileName" }
    $latestDate = ($archives | Measure-Object -Property Date -Maximum).Maximum
    $latest = @($archives | Where-Object { $_.Date -eq $latestDate })
    if ($latest.Count -ne 1) { throw "ambiguous latest archived evidence for $FileName" }
    return $latest[0].Path
}

function Read-FrontendVitestEvidenceJson {
    param(
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string]$FileName,
        [string]$ExplicitPath
    )

    $path = Resolve-FrontendVitestEvidencePath -RepoRoot $RepoRoot -FileName $FileName -ExplicitPath $ExplicitPath
    try {
        [pscustomobject]@{ Path = $path; Evidence = (Get-Content -Raw -LiteralPath $path | ConvertFrom-Json) }
    } catch {
        throw "evidence file is malformed: $path"
    }
}

function Resolve-FrontendVitestLinkedEvidencePath {
    param(
        [Parameter(Mandatory = $true)][string]$AnchorEvidenceRoot,
        [Parameter(Mandatory = $true)][string]$Reference,
        [Parameter(Mandatory = $true)][string]$ExpectedFileName
    )
    if ($ExpectedFileName -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]*\.json$') { throw "linked evidence file name is unsafe: $ExpectedFileName" }
    $normalizedReference = $Reference.Replace('\', '/')
    $legacyReference = "openspec/changes/$script:FrontendVitestChangeName/evidence/$ExpectedFileName"
    if ($normalizedReference -cne $legacyReference) { throw "linked evidence reference is not the expected legacy path: $Reference" }
    try { $anchor = [IO.Path]::GetFullPath($AnchorEvidenceRoot) } catch { throw "linked evidence anchor is not resolvable" }
    $normalizedAnchor = $anchor.Replace('\', '/')
    if ($normalizedAnchor -notmatch ("(?i)/openspec/changes/(?:archive/)?(?:\d{4}-\d{2}-\d{2}-)?" + [regex]::Escape($script:FrontendVitestChangeName) + "/evidence$")) { throw "linked evidence anchor is outside the expected change" }
    if (-not (Test-Path -LiteralPath $anchor -PathType Container)) { throw "linked evidence anchor is missing: $anchor" }
    $candidate = [IO.Path]::GetFullPath((Join-Path $anchor $ExpectedFileName))
    $anchorPrefix = $anchor.TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar) + [IO.Path]::DirectorySeparatorChar
    if (-not $candidate.StartsWith($anchorPrefix, [StringComparison]::OrdinalIgnoreCase)) { throw "linked evidence path escapes anchor" }
    if (-not (Test-Path -LiteralPath $candidate -PathType Leaf)) { throw "linked evidence file is missing: $candidate" }
    return $candidate
}

Export-ModuleMember -Function Resolve-FrontendVitestEvidencePath, Read-FrontendVitestEvidenceJson, Resolve-FrontendVitestLinkedEvidencePath
