param(
    [Parameter(Mandatory = $true)]
    [string]$BundleDir,
    [string]$EnvFile = ".env",
    [string]$CredentialOut = "",
    [switch]$Force
)

$ErrorActionPreference = "Stop"
$projectRoot = Split-Path -Parent $PSScriptRoot

if (-not [System.IO.Path]::IsPathRooted($BundleDir)) {
    $BundleDir = Join-Path $projectRoot $BundleDir
}

if (-not [System.IO.Path]::IsPathRooted($EnvFile)) {
    $EnvFile = Join-Path $projectRoot $EnvFile
}

$structureFile = Join-Path $BundleDir "structure.sql"
$seedFile = Join-Path $BundleDir "minimal_seed.sql"
$versionFile = Join-Path $BundleDir "schema_migration_versions.txt"

if (-not (Test-Path $BundleDir)) {
    throw "Baseline bundle not found: $BundleDir"
}
if (-not (Test-Path $structureFile)) {
    throw "Structure SQL not found: $structureFile"
}
if (-not (Test-Path $seedFile)) {
    throw "Minimal seed SQL not found: $seedFile"
}
if (-not (Test-Path $versionFile)) {
    throw "schema_migration version file not found: $versionFile"
}
if (-not $Force) {
    throw "Restore is destructive. Re-run with -Force after confirming the target database."
}

if ([string]::IsNullOrWhiteSpace($CredentialOut)) {
    $CredentialOut = Join-Path $BundleDir "admin_credentials.txt"
} elseif (-not [System.IO.Path]::IsPathRooted($CredentialOut)) {
    $CredentialOut = Join-Path $projectRoot $CredentialOut
}

Write-Host "Restoring structure SQL..."
& "$projectRoot\scripts\restore_db.ps1" -BackupFile $structureFile -EnvFile $EnvFile -Force

Write-Host "Rebuilding schema_migration from version manifest..."
go run .\cmd\migrate\main.go -env-file $EnvFile -baseline-file $versionFile
if ($LASTEXITCODE -ne 0) {
    throw "baseline-file migrate command failed"
}

Write-Host "Applying minimal seed..."
& "$projectRoot\scripts\init_minimal_seed.ps1" -EnvFile $EnvFile -SeedFile $seedFile -CredentialOut $CredentialOut
if ($LASTEXITCODE -ne 0) {
    throw "minimal seed initialization failed"
}

Write-Host "Applying incremental migrations after baseline..."
go run .\cmd\migrate\main.go -env-file $EnvFile
if ($LASTEXITCODE -ne 0) {
    throw "incremental migrate command failed"
}

Write-Host "Baseline bundle restore completed: $BundleDir"
