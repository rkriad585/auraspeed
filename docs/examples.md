# AuraSpeed Examples

Real-world usage examples for common scenarios.

## Basic Speed Test

Run a quick speed test from the command line:

```bash
auraspeed speedtest
```

**Expected output:**
```
Running speed test...
Using cached server list...
Server: New York, NY
ISP:    Comcast (203.0.113.1)

Results:
--------
Download: 85.42 Mbps
Upload:   42.18 Mbps
Ping:     12 ms
```

---

## JSON Output for Scripting

Get results as JSON for automation:

```bash
auraspeed speedtest --json
```

**Output:**
```json
{
  "download": 85.42,
  "upload": 42.18,
  "ping": 12,
  "isp": "Comcast",
  "server": "New York, NY"
}
```

---

## Interactive TUI Mode

Launch the full terminal interface:

```bash
auraspeed tui
```

**Keyboard shortcuts:**
- `R` — Restart speed test
- `C` — Copy results to clipboard
- `H` — View test history
- `Esc` — Close popups
- `Ctrl+C` — Exit

---

## System Information

Display detailed system information:

```bash
auraspeed info
```

**Example output:**
```
System Information
------------------
OS:       Windows 10 22H2
CPU:      Intel(R) Core(TM) i7-10700K (16 cores)
Memory:   32 GB total, 16 GB used (50%)
Disk:     512 GB total, 256 GB used (50%)
Hostname: DESKTOP-ABC123
```

---

## Network Diagnostics

### Ping a Host

```bash
auraspeed network ping google.com
```

### DNS Lookup

```bash
auraspeed network dns google.com
```

### Traceroute

```bash
auraspeed network traceroute google.com
```

---

## Test History

### View Recent Tests

```bash
auraspeed history
```

### Limit Results

```bash
auraspeed history --limit 5
```

### Export to JSON

```bash
auraspeed history --export results.json
```

### Export to TSV

```bash
auraspeed history --export results.tsv
```

---

## Configuration Examples

### View Current Config

```bash
auraspeed config view
```

### Change Theme

```bash
auraspeed config set ui.theme dark
```

### Set Custom Timeout

```bash
auraspeed config set speedtest.timeout 60
```

### Reset to Defaults

```bash
auraspeed config reset
```

---

## Using Aliases

Quick shortcuts for common commands:

```bash
auraspeed st        # Same as: auraspeed speedtest
auraspeed si        # Same as: auraspeed info
auraspeed net       # Same as: auraspeed network
auraspeed hist      # Same as: auraspeed history
```

---

## Build from Source

### Windows (PowerShell)

```powershell
.\build.ps1
```

Outputs to `dist/`:
- `auraspeed-windows-amd64.exe`

### Linux/macOS (Bash/Zsh/Fish)

```bash
./build.sh
```

Outputs to `dist/`:
- `auraspeed-linux-amd64`
- `auraspeed-darwin-amd64`
- `auraspeed-darwin-arm64`

---

## CI/CD with GoReleaser

Create a version tag to trigger the release workflow:

```bash
git tag v1.0.1
git push origin v1.0.1
```

This triggers the GitHub Actions workflow which:
1. Runs tests
2. Builds for all platforms
3. Creates a GitHub Release with binaries

---

## Automation Example

### Bash Script for Periodic Tests

```bash
#!/bin/bash
# Run speed test every hour and log results

while true; do
    echo "=== $(date) ==="
    auraspeed speedtest --json >> speed_log.json
    sleep 3600
done
```

### PowerShell Script for Monitoring

```powershell
# Monitor connection quality
while ($true) {
    Write-Host "=== $(Get-Date) ==="
    auraspeed speedtest --json | Tee-Object -FilePath speed_log.json -Append
    Start-Sleep -Seconds 3600
}
```

---

## Advanced: Custom Server Selection

```bash
# List servers (check cache or fetch new)
auraspeed speedtest

# Use a specific server by ID
auraspeed speedtest --server-id 1234
```

---

## Links

- [GitHub Repository](https://github.com/rkriad585/auraspeed)
- [Releases](https://github.com/rkriad585/auraspeed/releases)
- [Issue Tracker](https://github.com/rkriad585/auraspeed/issues)
- [Getting Started Guide](getting-started.md)
- [API Reference](api-reference.md)
