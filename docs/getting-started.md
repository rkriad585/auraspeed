# Getting Started with AuraSpeed

AuraSpeed is a cross-platform terminal tool for network diagnostics, built with Go. This beginner-friendly guide walks you through installing, configuring, and using AuraSpeed in ~10 minutes.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Option 1: Pre-built Binaries](#option-1-pre-built-binaries)
  - [Option 2: Build from Source](#option-2-build-from-source)
- [First Run & Verification](#first-run--verification)
- [Core CLI Commands](#core-cli-commands)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

---

## Prerequisites
### If Building from Source
- Go 1.21+ installed ([download Go](https://go.dev/dl/))
### OS-Specific Requirements
| OS | Required Shell |
|----|----------------|
| Windows | PowerShell or CMD |
| Linux | bash/zsh/fish terminal |
| macOS | bash/zsh/fish terminal |

---

## Installation
### Option 1: Pre-built Binaries
Fastest method, no Go required. Download the correct binary for your OS from [AuraSpeed Releases](https://github.com/rkriad585/auraspeed/releases):

| OS | Architecture | Binary Filename |
|----|--------------|-----------------|
| Windows | x86-64 | `auraspeed-windows-amd64.exe` |
| Linux | x86-64 | `auraspeed-linux-amd64` |
| macOS (Intel) | x86-64 | `auraspeed-darwin-amd64` |
| macOS (Apple Silicon) | ARM64 | `auraspeed-darwin-arm64` |

#### Windows Steps
1. Download `auraspeed-windows-amd64.exe` from the latest release
2. (Optional) Rename to `auraspeed.exe` for easier use
3. Add the binary's folder to your system `PATH` to run from any terminal

#### Linux/macOS Steps
1. Download the matching binary for your OS/architecture
2. Open your terminal and navigate to the download folder
3. Make the binary executable:
   ```bash
   chmod +x auraspeed-<your-os>-<arch>
   ```
4. (Optional) Move to a PATH directory like `/usr/local/bin`:
   ```bash
   sudo mv auraspeed-<your-os>-<arch> /usr/local/bin/auraspeed
   ```

### Option 2: Build from Source
Requires Go 1.21+.

#### Windows
Run the provided build script in PowerShell:
```powershell
.\build.ps1
```
Output binary: `dist/auraspeed-windows-amd64.exe`

#### Linux/macOS
Run the provided build script in your terminal:
```bash
./build.sh
```
Output binaries: `dist/` folder with linux, darwin (amd64/arm64) binaries.

---

## First Run & Verification
Verify installation by running the help command:
```bash
auraspeed --help
```

Expected output (abbreviated):
```
AuraSpeed - Cross-platform terminal tool for network diagnostics

Usage:
  auraspeed [command]

Available Commands:
  help        Help about any command
  info        Show system and network info
  network     Run network diagnostics
  speedtest   Run a quick speed test
  tui         Launch interactive TUI mode
```

> [!TIP]
> If you see "command not found", check that the binary is in your system PATH.

---

## Core CLI Commands
All commands work on Windows (use `.exe` suffix if not in PATH), Linux, and macOS.

### 1. Launch Interactive TUI Mode
The TUI (Terminal User Interface) provides a visual dashboard for all features:
```bash
auraspeed tui
```
> [!NOTE]
> Use arrow keys to navigate, Enter to select, Esc/Q to exit.

### 2. Quick Speed Test
Run a full download/upload speed test with auto-selected servers:
```bash
auraspeed speedtest
```
Expected output:
```
Running speed test...
Download: 150.2 Mbps
Upload: 45.8 Mbps
Latency: 12 ms
```

### 3. System & Network Info
View your device's network configuration and system details:
```bash
auraspeed info
```
Expected output includes OS version, CPU, RAM, network interfaces, and public IP.

### 4. Network Diagnostics
Run ping, traceroute, and other diagnostics against a target:
```bash
auraspeed network ping google.com
```
Expected output:
```
Pinging google.com (142.250.80.46)...
Reply from 142.250.80.46: bytes=32 time=14ms TTL=117
```

---

## Configuration
AuraSpeed uses a TOML config file, created automatically on first run.

### Config File Location
| OS | Path |
|----|------|
| Windows | `%USERPROFILE%\.auraspeed\config.toml` |
| Linux | `~/.auraspeed/config.toml` |
| macOS | `~/.auraspeed/config.toml` |

### Default Config Options
Edit the config file to customize behavior:
```toml
[global]
  loglevel = "info"  # Options: debug, info, warn, error
  nocolor = false
  autoupdate = true

[speedtest]
  timeout = 30
  defaultserverid = 0  # 0 = auto-select closest server
  paralleldownloads = 4
  paralleluploads = 2

[ui]
  theme = "default"
  savehistory = true
```

| Section | Option | Type | Default | Description |
|---------|--------|------|---------|-------------|
| global | loglevel | string | info | Logging verbosity level |
| global | nocolor | boolean | false | Disable colored terminal output |
| global | autoupdate | boolean | true | Auto-check for updates on startup |
| speedtest | timeout | int | 30 | Speed test timeout in seconds |
| speedtest | defaultserverid | int | 0 | Fixed speed test server ID (0 = auto) |
| speedtest | paralleldownloads | int | 4 | Concurrent download connections |
| speedtest | paralleluploads | int | 2 | Concurrent upload connections |
| ui | theme | string | default | TUI theme (default, dark, light) |
| ui | savehistory | boolean | true | Save TUI command history |

---

## Troubleshooting
### "auraspeed: command not found"
- Windows: Ensure the binary folder is added to your System PATH, or run with full path `C:\path\to\auraspeed-windows-amd64.exe`
- Linux/macOS: Verify the binary is executable (`chmod +x`) and in a PATH directory

### Build script fails (Windows)
> [!WARNING]
> Ensure you have Go 1.21+ installed and added to PATH. Check with `go version`.

### Permission denied (Linux/macOS)
Run the binary with sudo if needed for network diagnostics:
```bash
sudo auraspeed network ping google.com
```

### Config file not found
Run any AuraSpeed command once to auto-generate the default config file at the location listed above.

---

## Next Steps
- Explore more commands with `auraspeed [command] --help`
- Check the [GitHub Repository](https://github.com/rkriad585/auraspeed) for advanced usage
- Report issues at [GitHub Issues](https://github.com/rkriad585/auraspeed/issues)
