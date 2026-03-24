param(
    [string]$GoProjectDir = ".",
    [string]$VueProjectDir = "..\am-erp-vue"
)

$ErrorActionPreference = "Stop"
$projectBase = Split-Path -Parent $PSScriptRoot
if (-not [System.IO.Path]::IsPathRooted($GoProjectDir)) {
    $GoProjectDir = Join-Path $projectBase $GoProjectDir
}
if (-not [System.IO.Path]::IsPathRooted($VueProjectDir)) {
    $VueProjectDir = Join-Path $projectBase $VueProjectDir
}

$projectRoot = Resolve-Path $GoProjectDir
$vueRoot = Resolve-Path $VueProjectDir

function Invoke-Step {
    param(
        [string]$Title,
        [scriptblock]$Action
    )

    Write-Host ""
    Write-Host "==> $Title"
    & $Action
}

Invoke-Step "Go tests" {
    Push-Location $projectRoot
    try {
        & 'C:\Program Files\PowerShell\7\pwsh.exe' -Command "go test ./... -count=1"
        if ($LASTEXITCODE -ne 0) {
            throw "go test failed"
        }
    }
    finally {
        Pop-Location
    }
}

Invoke-Step "Vue build" {
    Push-Location $vueRoot
    try {
        & 'C:\Program Files\PowerShell\7\pwsh.exe' -Command "npm run build-only"
        if ($LASTEXITCODE -ne 0) {
            throw "npm run build-only failed"
        }
    }
    finally {
        Pop-Location
    }
}

Invoke-Step "Migration status" {
    Push-Location $projectRoot
    try {
        $output = go run .\cmd\migrate\main.go
        if ($LASTEXITCODE -ne 0) {
            throw "migration check failed"
        }
        $output | ForEach-Object { Write-Host $_ }
    }
    finally {
        Pop-Location
    }
}

Write-Host ""
Write-Host "Preflight check completed."
