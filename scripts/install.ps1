$ErrorActionPreference = "Stop"

# ==============================
# TRAK Windows Installer
# ==============================

$Version = "v1.1.0"

# Architecture detection (amd64 vs arm64)
$Arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64" -or $env:PROCESSOR_ARCHITEW6432 -eq "ARM64") {
    $Arch = "arm64"
}

$DownloadUrl = "https://github.com/ndk123-web/trak/releases/download/$Version/trak-windows-$Arch.exe"

# Install location:
# C:\Users\<User>\trak\bin\trak.exe
$InstallRoot = Join-Path $env:USERPROFILE "trak"
$BinDir = Join-Path $InstallRoot "bin"
$ExePath = Join-Path $BinDir "trak.exe"

Write-Host ""
Write-Host "TRAK Installer" -ForegroundColor Cyan
Write-Host "=============="
Write-Host ""

# --------------------------------
# 1. Detect architecture
# --------------------------------

if (-not [Environment]::Is64BitOperatingSystem) {
    throw "TRAK currently requires a 64-bit Windows system."
}

Write-Host "[1/4] Detected Windows ($Arch)" -ForegroundColor Green

# --------------------------------
# 2. Create install directory
# --------------------------------

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
}

Write-Host "[2/4] Install directory:" -ForegroundColor Green
Write-Host "      $BinDir"

# --------------------------------
# 3. Download TRAK executable
# --------------------------------

Write-Host ""
Write-Host "[3/4] Downloading TRAK $Version..." -ForegroundColor Yellow

Invoke-WebRequest `
    -Uri $DownloadUrl `
    -OutFile $ExePath

if (-not (Test-Path $ExePath)) {
    throw "TRAK executable was not downloaded."
}

Write-Host "      Downloaded successfully." -ForegroundColor Green

# --------------------------------
# 4. Add bin directory to USER PATH
# --------------------------------

Write-Host ""
Write-Host "[4/4] Configuring PATH..." -ForegroundColor Yellow

$currentUserPath = [Environment]::GetEnvironmentVariable(
    "Path",
    [EnvironmentVariableTarget]::User
)

$pathEntries = @()

if ($currentUserPath) {
    $pathEntries = $currentUserPath -split ";" |
        Where-Object { $_ -and $_.Trim() -ne "" }
}

$alreadyExists = $pathEntries | Where-Object {
    $_.TrimEnd("\") -ieq $BinDir.TrimEnd("\")
}

if (-not $alreadyExists) {

    if ($currentUserPath) {
        $newUserPath = "$currentUserPath;$BinDir"
    }
    else {
        $newUserPath = $BinDir
    }

    [Environment]::SetEnvironmentVariable(
        "Path",
        $newUserPath,
        [EnvironmentVariableTarget]::User
    )

    Write-Host "      Added to USER PATH." -ForegroundColor Green
}
else {
    Write-Host "      Already present in USER PATH." -ForegroundColor Green
}

# Update current active PowerShell session PATH immediately
if ($env:PATH -notlike "*$BinDir*") {
    $env:PATH = "$BinDir;$env:PATH"
}

# Optional: Ensure profile also has PATH configured
try {
    if ($PROFILE) {
        $profileDir = Split-Path -Parent $PROFILE
        if (-not (Test-Path $profileDir)) {
            New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
        }
        $profileLine = "`$env:PATH = `"$BinDir;`$env:PATH`""
        if (Test-Path $PROFILE) {
            $content = Get-Content $PROFILE -Raw -ErrorAction SilentlyContinue
            if ($content -notlike "*$BinDir*") {
                Add-Content -Path $PROFILE -Value "`n# TRAK CLI Path`n$profileLine"
            }
        }
    }
} catch {
    # Non-fatal if profile update fails
}

# --------------------------------
# Done
# --------------------------------

Write-Host ""
Write-Host "TRAK installed successfully! 🚀" -ForegroundColor Cyan
Write-Host ""
Write-Host "Installed at:"
Write-Host "  $ExePath"
Write-Host ""
Write-Host "💡 TIP FOR DEVELOPERS:" -ForegroundColor Yellow
Write-Host "If 'trak' is not recognized in other open tabs, please restart your terminal" -ForegroundColor Yellow
Write-Host "or open a new PowerShell window so the updated PATH is loaded." -ForegroundColor Yellow
Write-Host ""
Write-Host "Quick Test Commands:"
Write-Host "  trak list" -ForegroundColor Cyan
Write-Host "  trak init lang/go" -ForegroundColor Cyan
Write-Host ""
