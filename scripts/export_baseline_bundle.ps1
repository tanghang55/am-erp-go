param(
    [string]$EnvFile = ".env",
    [string]$OutputDir = "backups\baseline"
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

if (-not [System.IO.Path]::IsPathRooted($OutputDir)) {
    $OutputDir = Join-Path $projectRoot $OutputDir
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

$bundleDir = Join-Path $OutputDir (Get-Date -Format "yyyyMMdd-HHmmss")
New-Item -ItemType Directory -Path $bundleDir -Force | Out-Null

$dumpExe = Resolve-MySQLDump
$mysqlExe = Resolve-MySQLExe
$structureFile = Join-Path $bundleDir "structure.sql"
$seedFile = Join-Path $bundleDir "minimal_seed.sql"
$versionFile = Join-Path $bundleDir "schema_migration_versions.txt"

$dumpArgs = @(
    "--default-character-set=utf8mb4",
    "--single-transaction",
    "--no-data",
    "--routines",
    "--triggers",
    "--host=$hostName",
    "--port=$port",
    "--user=$user"
)
if (-not [string]::IsNullOrWhiteSpace($password)) {
    $dumpArgs += "--password=$password"
}
$dumpArgs += $database

Write-Host "Exporting structure SQL to $structureFile"
$warningFile = [System.IO.Path]::GetTempFileName()
& $dumpExe @dumpArgs 2> $warningFile | Set-Content -Path $structureFile -Encoding UTF8
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

Write-Host "Exporting minimal seed SQL to $seedFile"
& "$projectRoot\scripts\export_minimal_seed.ps1" -EnvFile $EnvFile -OutputFile $seedFile
if ($LASTEXITCODE -ne 0) {
    throw "minimal seed export failed"
}

$query = "SELECT version FROM schema_migration ORDER BY version"
$mysqlArgs = @(
    "--batch",
    "--skip-column-names",
    "--host=$hostName",
    "--port=$port",
    "--user=$user"
)
if (-not [string]::IsNullOrWhiteSpace($password)) {
    $mysqlArgs += "--password=$password"
}
$mysqlArgs += @($database, "-e", $query)

Write-Host "Exporting schema_migration versions to $versionFile"
& $mysqlExe @mysqlArgs | Set-Content -Path $versionFile -Encoding UTF8
if ($LASTEXITCODE -ne 0) {
    throw "mysql query failed while exporting schema_migration versions"
}

Write-Host "Baseline bundle completed: $bundleDir"
