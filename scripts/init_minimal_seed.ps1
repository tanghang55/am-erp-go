param(
    [string]$EnvFile = ".env",
    [string]$SeedFile = "baseline\minimal\minimal_seed.sql",
    [string]$CredentialOut = "admin_credentials.txt",
    [string]$AdminUsername = "admin",
    [string]$AdminRealName = "系统管理员",
    [int]$PasswordLength = 20
)

$ErrorActionPreference = "Stop"
$projectRoot = Split-Path -Parent $PSScriptRoot

if (-not [System.IO.Path]::IsPathRooted($EnvFile)) {
    $EnvFile = Join-Path $projectRoot $EnvFile
}
if (-not [System.IO.Path]::IsPathRooted($SeedFile)) {
    $SeedFile = Join-Path $projectRoot $SeedFile
}
if (-not [System.IO.Path]::IsPathRooted($CredentialOut)) {
    $CredentialOut = Join-Path $projectRoot $CredentialOut
}

Write-Host "Initializing minimal seed..."
go run .\cmd\initseed\main.go `
    -env-file $EnvFile `
    -seed-file $SeedFile `
    -credential-out $CredentialOut `
    -admin-username $AdminUsername `
    -admin-real-name $AdminRealName `
    -password-length $PasswordLength

if ($LASTEXITCODE -ne 0) {
    throw "initseed command failed"
}

Write-Host "Minimal seed initialized: $CredentialOut"
