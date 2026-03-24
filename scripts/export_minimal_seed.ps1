param(
    [string]$EnvFile = ".env",
    [string]$OutputFile = "baseline\minimal\minimal_seed.sql"
)

$ErrorActionPreference = "Stop"
$projectRoot = Split-Path -Parent $PSScriptRoot

if (-not [System.IO.Path]::IsPathRooted($EnvFile)) {
    $EnvFile = Join-Path $projectRoot $EnvFile
}
if (-not [System.IO.Path]::IsPathRooted($OutputFile)) {
    $OutputFile = Join-Path $projectRoot $OutputFile
}

Write-Host "Exporting minimal seed..."
go run .\cmd\exportseed\main.go -env-file $EnvFile -output-file $OutputFile

if ($LASTEXITCODE -ne 0) {
    throw "exportseed command failed"
}

Write-Host "Minimal seed exported: $OutputFile"
