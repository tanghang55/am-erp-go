param(
    [Parameter(Mandatory = $true)]
    [string]$BackupFile,
    [string]$EnvFile = ".env",
    [switch]$Force
)

$ErrorActionPreference = "Stop"
$projectRoot = Split-Path -Parent $PSScriptRoot

function Get-EnvMap {
    param([string]$Path)

    if (-not (Test-Path $Path)) {
        throw "Env file not found: $Path"
    }

    $map = @{}
    foreach ($line in Get-Content $Path) {
        $trimmed = $line.Trim()
        if ([string]::IsNullOrWhiteSpace($trimmed) -or $trimmed.StartsWith("#")) {
            continue
        }
        $index = $trimmed.IndexOf("=")
        if ($index -lt 1) {
            continue
        }
        $key = $trimmed.Substring(0, $index).Trim()
        $value = $trimmed.Substring($index + 1).Trim().Trim('"')
        $map[$key] = $value
    }
    return $map
}

function Resolve-MySQLExe {
    $candidates = @(
        "C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe",
        "C:\Program Files\MySQL\MySQL Server 8.4\bin\mysql.exe"
    )

    foreach ($candidate in $candidates) {
        if (Test-Path $candidate) {
            return $candidate
        }
    }

    $cmd = Get-Command mysql.exe -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Source
    }

    throw "mysql.exe not found"
}

if (-not [System.IO.Path]::IsPathRooted($EnvFile)) {
    $EnvFile = Join-Path $projectRoot $EnvFile
}

if (-not [System.IO.Path]::IsPathRooted($BackupFile)) {
    $BackupFile = Join-Path $projectRoot $BackupFile
}

if (-not (Test-Path $BackupFile)) {
    throw "Backup file not found: $BackupFile"
}

if (-not $Force) {
    throw "Restore is destructive. Re-run with -Force after confirming the target database."
}

$envMap = Get-EnvMap -Path $EnvFile
$hostName = $envMap["DB_HOST"]
$port = $envMap["DB_PORT"]
$user = $envMap["DB_USER"]
$password = $envMap["DB_PASSWORD"]
$database = $envMap["DB_DATABASE"]

if ([string]::IsNullOrWhiteSpace($hostName) -or
    [string]::IsNullOrWhiteSpace($port) -or
    [string]::IsNullOrWhiteSpace($user) -or
    [string]::IsNullOrWhiteSpace($database)) {
    throw "Missing DB_HOST/DB_PORT/DB_USER/DB_DATABASE in $EnvFile"
}

$mysqlExe = Resolve-MySQLExe
$arguments = @(
    "--default-character-set=utf8mb4",
    "--host=$hostName",
    "--port=$port",
    "--user=$user",
    $database
)

if (-not [string]::IsNullOrWhiteSpace($password)) {
    $arguments += "--password=$password"
}

Write-Host "Restoring $database from $BackupFile"
Get-Content $BackupFile | & $mysqlExe @arguments
if ($LASTEXITCODE -ne 0) {
    throw "mysql restore failed"
}

Write-Host "Restore completed: $database <= $BackupFile"
