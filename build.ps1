# Build script for AuraSpeed
# Builds cross-platform binaries and saves to bin/

$ErrorActionPreference = "Stop"

$projectName = "auraspeed"
$outputDir = "bin"

# Get version from .version file
$version = "dev"
$versionFile = ".version"
if (Test-Path $versionFile) {
    $version = Get-Content $versionFile -Raw
    $version = $version.Trim()
    if (-not $version) { $version = "dev" }
}

# Get commit hash
try {
    $commit = git rev-parse --short HEAD 2>$null
    if (-not $commit) { $commit = "unknown" }
} catch {
    $commit = "unknown"
}

$buildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

Write-Host "Building $projectName..." -ForegroundColor Cyan
Write-Host "Version: $version" -ForegroundColor Cyan
Write-Host "Commit: $commit" -ForegroundColor Cyan
Write-Host ""

# Create output directory
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

# Build targets
$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Output = "$projectName-windows-amd64.exe" },
    @{ GOOS = "windows"; GOARCH = "arm64"; Output = "$projectName-windows-arm64.exe" },
    @{ GOOS = "linux";   GOARCH = "amd64"; Output = "$projectName-linux-amd64" },
    @{ GOOS = "linux";   GOARCH = "arm64"; Output = "$projectName-linux-arm64" },
    @{ GOOS = "darwin";  GOARCH = "amd64"; Output = "$projectName-darwin-amd64" },
    @{ GOOS = "darwin";  GOARCH = "arm64"; Output = "$projectName-darwin-arm64" }
)

# LDFLAGS for version info
$ldflags = "-s -w -X 'auraspeed/cmd/root.Version=$version' -X 'auraspeed/cmd/root.Commit=$commit' -X 'auraspeed/cmd/root.BuildTime=$buildTime'"

foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    $env:CGO_ENABLED = "0"

    $outputPath = Join-Path $outputDir $target.Output

    Write-Host "Building for $($target.GOOS)/$($target.GOARCH)..." -ForegroundColor Yellow

    try {
        go build -ldflags $ldflags -o $outputPath ./cmd/main.go

        if (Test-Path $outputPath) {
            $fileSize = (Get-Item $outputPath).Length / 1KB
            Write-Host "  ✓ Created: $outputPath ($($fileSize.ToString('0.0')) KB)" -ForegroundColor Green
        } else {
            Write-Host "  ✗ Failed: Output file not found" -ForegroundColor Red
        }
    } catch {
        Write-Host "  ✗ Build failed: $_" -ForegroundColor Red
    }
}

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""
$env:CGO_ENABLED = ""

Write-Host ""
Write-Host "Build complete! Output directory: $outputDir" -ForegroundColor Green
Get-ChildItem $outputDir | Format-Table Name, @{Name="Size (KB)";Expression={$_.Length / 1KB}} -AutoSize
