param(
    [string]$EnvFile = ".env",
    [string]$OutputDir = "backups\db"
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

if (-not [System.IO.Path]::IsPathRooted($EnvFile)) {
    $EnvFile = Join-Path $projectRoot $EnvFile
}

if (-not [System.IO.Path]::IsPathRooted($OutputDir)) {
    $OutputDir = Join-Path $projectRoot $OutputDir
}

function Resolve-MySQLDump {
    $candidates = @(
        "C:\Program Files\MySQL\MySQL Server 8.0\bin\mysqldump.exe",
        "C:\Program Files\MySQL\MySQL Server 8.4\bin\mysqldump.exe"
    )

    foreach ($candidate in $candidates) {
        if (Test-Path $candidate) {
            return $candidate
        }
    }

    $cmd = Get-Command mysqldump.exe -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Source
    }

    throw "mysqldump.exe not found"
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

$dumpExe = Resolve-MySQLDump
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

$outputFile = Join-Path $OutputDir "$database-$timestamp.sql"

$arguments = @(
    "--default-character-set=utf8mb4",
    "--single-transaction",
    "--routines",
    "--triggers",
    "--host=$hostName",
    "--port=$port",
    "--user=$user"
)

if (-not [string]::IsNullOrWhiteSpace($password)) {
    $arguments += "--password=$password"
}

$arguments += $database

Write-Host "Backing up $database to $outputFile"
$warningFile = [System.IO.Path]::GetTempFileName()
& $dumpExe @arguments 2> $warningFile | Set-Content -Path $outputFile -Encoding UTF8
if ($LASTEXITCODE -ne 0) {
    $warningText = Get-Content $warningFile -Raw
    Remove-Item $warningFile -Force -ErrorAction SilentlyContinue
    throw "mysqldump failed: $warningText"
}

$warningText = Get-Content $warningFile -Raw
Remove-Item $warningFile -Force -ErrorAction SilentlyContinue
if (-not [string]::IsNullOrWhiteSpace($warningText)) {
    Write-Warning $warningText.Trim()
}

Write-Host "Backup completed: $outputFile"
