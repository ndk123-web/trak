# ==============================================================================
# Trak CLI - Cross-Platform Multi-Architecture Build & Distribution Script
# Compiles binaries for Windows, Linux, and macOS (Intel & Apple Silicon ARM64)
# ==============================================================================

param (
    [string]$Version = "1.2.0",
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"

$rootPath = Split-Path -Parent $PSScriptRoot
Set-Location $rootPath

Write-Host ""
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "  Trak CLI - Multi-Platform Distribution Builder v$Version" -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host ""

# Clean or create output directories
$distPath = Join-Path $rootPath $OutputDir
$releasePath = Join-Path $distPath "v$Version"
$binariesPath = Join-Path $releasePath "binaries"

if (Test-Path $releasePath) {
    Write-Host "[CLEAN] Removing existing directory: $releasePath" -ForegroundColor Yellow
    Remove-Item -Recurse -Force $releasePath
}

New-Item -ItemType Directory -Path $binariesPath -Force | Out-Null

# Cross-compilation target matrix
$targets = @(
    # Windows
    @{ OS = "windows"; Arch = "amd64"; Ext = ".exe"; Label = "Windows (x86_64 64-bit)" },
    @{ OS = "windows"; Arch = "arm64"; Ext = ".exe"; Label = "Windows (ARM64)" },
    @{ OS = "windows"; Arch = "386";   Ext = ".exe"; Label = "Windows (32-bit)" },

    # Linux
    @{ OS = "linux";   Arch = "amd64"; Ext = "";     Label = "Linux (x86_64 64-bit)" },
    @{ OS = "linux";   Arch = "arm64"; Ext = "";     Label = "Linux (ARM64 / aarch64)" },
    @{ OS = "linux";   Arch = "386";   Ext = "";     Label = "Linux (32-bit)" },

    # macOS (Darwin)
    @{ OS = "darwin";  Arch = "amd64"; Ext = "";     Label = "macOS (Intel x86_64)" },
    @{ OS = "darwin";  Arch = "arm64"; Ext = "";     Label = "macOS (Apple Silicon M1/M2/M3/M4)" }
)

# Optimized build flags (strip debug symbols -s -w for compact production binaries)
$ldflags = "-s -w -X 'github.com/ndk123-web/trak/cmd.Version=v$Version'"

$builtFiles = @()
$stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

foreach ($target in $targets) {
    $os = $target.OS
    $arch = $target.Arch
    $ext = $target.Ext
    $label = $target.Label

    $binaryName = "trak-${os}-${arch}${ext}"
    $outputFilePath = Join-Path $binariesPath $binaryName

    Write-Host "[BUILD] Compiling for $label [$os/$arch]... " -ForegroundColor White -NoNewline

    $env:GOOS = $os
    $env:GOARCH = $arch
    $env:CGO_ENABLED = "0"

    # Execute Go Build
    go build -ldflags $ldflags -o $outputFilePath .

    if ($LASTEXITCODE -ne 0) {
        Write-Host "FAILED" -ForegroundColor Red
        Exit 1
    }

    $fileInfo = Get-Item $outputFilePath
    $sizeMB = [math]::Round($fileInfo.Length / 1MB, 2)
    Write-Host "OK ($sizeMB MB)" -ForegroundColor Green

    # Package into archive (.zip for windows, .tar.gz for unix)
    $archiveName = "trak_${Version}_${os}_${arch}"
    
    if ($os -eq "windows") {
        $zipPath = Join-Path $releasePath "$archiveName.zip"
        $tempDir = Join-Path $env:TEMP ([System.Guid]::NewGuid().ToString())
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        Copy-Item $outputFilePath (Join-Path $tempDir "trak.exe")
        Compress-Archive -Path (Join-Path $tempDir "*") -DestinationPath $zipPath -Force
        Remove-Item -Recurse -Force $tempDir
        $builtFiles += $zipPath
    } else {
        $tarPath = Join-Path $releasePath "$archiveName.tar.gz"
        $tempDir = Join-Path $env:TEMP ([System.Guid]::NewGuid().ToString())
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        Copy-Item $outputFilePath (Join-Path $tempDir "trak")
        
        tar -czf $tarPath -C $tempDir trak 2>$null
        if (-not (Test-Path $tarPath)) {
            $zipPath = Join-Path $releasePath "$archiveName.zip"
            Compress-Archive -Path (Join-Path $tempDir "*") -DestinationPath $zipPath -Force
            $builtFiles += $zipPath
        } else {
            $builtFiles += $tarPath
        }
        Remove-Item -Recurse -Force $tempDir
    }
}

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""
$env:CGO_ENABLED = ""

# Generate SHA-256 Checksums
Write-Host ""
Write-Host "[CHECKSUM] Generating SHA-256 Checksums (checksums.txt)..." -ForegroundColor Yellow
$checksumFile = Join-Path $releasePath "checksums.txt"
$checksumLines = @()

foreach ($file in $builtFiles) {
    if (Test-Path $file) {
        $hash = (Get-FileHash -Path $file -Algorithm SHA256).Hash.ToLower()
        $fileName = Split-Path -Leaf $file
        $checksumLines += "$hash  $fileName"
    }
}

$checksumLines | Out-File -FilePath $checksumFile -Encoding utf8
Write-Host "[OK] Checksums created at: $checksumFile" -ForegroundColor Green

$stopwatch.Stop()
$elapsedSec = [math]::Round($stopwatch.Elapsed.TotalSeconds, 1)

# Summary Display
Write-Host ""
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "  Build Completed Successfully in $elapsedSec seconds!" -ForegroundColor Green
Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "Output Directory: $releasePath" -ForegroundColor White
Write-Host ""

Get-ChildItem -Path $releasePath -Filter "trak_*" | ForEach-Object {
    $sizeKB = [math]::Round($_.Length / 1KB, 1)
    Write-Host "  - $($_.Name) ($sizeKB KB)" -ForegroundColor Cyan
}

Write-Host "  - checksums.txt" -ForegroundColor Yellow
Write-Host ""
Write-Host "Ready for GitHub Release v$Version!" -ForegroundColor Green
Write-Host ""
