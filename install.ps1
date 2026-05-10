# AuraSpeed Installer for Windows
# Downloads and installs AuraSpeed to ~/.config/neostore/auraspeed/bin/
# and adds it to the user PATH.

$ErrorActionPreference = "Stop"

$projectName = "auraspeed"
$versionUrl = "https://raw.githubusercontent.com/rkriad585/$projectName/main/.version"
$releaseUrl = "https://github.com/rkriad585/$projectName/releases/download"

# Get latest version from GitHub
try {
    $version = (Invoke-WebRequest -Uri $versionUrl -UseBasicParsing).Content.Trim()
    if (-not $version) { $version = "dev" }
} catch {
    $version = "dev"
}

Write-Host ">>> Installing $projectName $version..." -ForegroundColor Cyan
Write-Host ""

# Detect architecture
$arch = "amd64"
if ([Environment]::Is64BitOperatingSystem) {
    try {
        $cpuInfo = Get-WmiObject -Class Win32_Processor -ErrorAction SilentlyContinue
        if ($cpuInfo -and $cpuInfo.Architecture -eq 9) { $arch = "arm64" }
    } catch {
        # fallback to amd64
    }
}

$binName = "$projectName.exe"
$downloadName = "$projectName-windows-$arch.exe"

# Install directory
$binDir = "$env:USERPROFILE\.config\neostore\$projectName\bin"

# Create bin directory
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir -Force | Out-Null
}

# Download binary
$downloadUrl = "$releaseUrl/$version/$downloadName"
$outputPath = Join-Path $binDir $binName

Write-Host "Downloading $downloadUrl"
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath -UseBasicParsing
} catch {
    Write-Host "Error downloading: $_" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path $outputPath)) {
    Write-Host "Error: Download failed - output file not found" -ForegroundColor Red
    exit 1
}

$fileSize = (Get-Item $outputPath).Length / 1KB
Write-Host "OK   Installed to $outputPath ($($fileSize.ToString('0.0')) KB)" -ForegroundColor Green

# Add to user PATH
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$pathEntry = $binDir

if ($userPath -notlike "*$pathEntry*") {
    if ($userPath.EndsWith(";") -or [string]::IsNullOrEmpty($userPath)) {
        [Environment]::SetEnvironmentVariable("PATH", "$userPath$pathEntry", "User")
    } else {
        [Environment]::SetEnvironmentVariable("PATH", "$userPath;$pathEntry", "User")
    }
    Write-Host "Added to PATH. You may need to restart your terminal." -ForegroundColor Yellow
} else {
    Write-Host "Already in PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "OK   $projectName $version installed successfully!" -ForegroundColor Green
Write-Host "Run '$projectName --help' to get started"
