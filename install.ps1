# AuraSpeed Installer for Windows
# Downloads and installs AuraSpeed to ~/.config/neostore/auraspeed/bin/

$ErrorActionPreference = "Stop"

$VERSION_URL = "https://raw.githubusercontent.com/rkriad585/auraspeed/main/.version"
$RELEASE_URL = "https://github.com/rkriad585/auraspeed/releases/download"

# Get latest version from GitHub
try {
    $VERSION = (Invoke-WebRequest -Uri $VERSION_URL -UseBasicParsing).Content.Trim()
    if (-not $VERSION) { $VERSION = "v3.0.1" }
} catch {
    $VERSION = "v3.0.1"
}

Write-Host "AuraSpeed Installer $VERSION" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan

# Detect architecture
$BIN_NAME = "auraspeed.exe"

# Set install directory
$BIN_DIR = "$env:USERPROFILE\.config\neostore\auraspeed\bin"

# Create bin directory
if (-not (Test-Path $BIN_DIR)) {
    Write-Host "Creating installation directory: $BIN_DIR"
    New-Item -ItemType Directory -Path $BIN_DIR -Force | Out-Null
}

# Download binary
$downloadUrl = "$RELEASE_URL/$VERSION/auraspeed-windows-amd64.exe"
$outputPath = Join-Path $BIN_DIR $BIN_NAME

Write-Host "Downloading AuraSpeed v$VERSION from: $downloadUrl"
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath -UseBasicParsing
    Write-Host "Download complete!" -ForegroundColor Green
} catch {
    Write-Host "Error downloading: $_" -ForegroundColor Red
    exit 1
}

# Add to PATH
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$pathEntry = $BIN_DIR

if ($userPath -notlike "*$pathEntry*") {
    Write-Host "Adding to user PATH..."
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$pathEntry", "User")
    Write-Host "PATH updated. You may need to restart your terminal." -ForegroundColor Yellow
} else {
    Write-Host "Already in PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Run 'auraspeed --help' to get started"
Write-Host "To uninstall, run: auraspeed --selfuninstall"