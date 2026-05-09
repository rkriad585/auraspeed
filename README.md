# AuraSpeed

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=flat&logo=go)](https://go.dev/)
[![MIT License](https://img.shields.io/github/license/rkriad585/auraspeed?color=blue)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/rkriad585/auraspeed)](https://github.com/rkriad585/auraspeed/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/rkriad585/auraspeed)](https://goreportcard.com/report/github.com/rkriad585/auraspeed)

Cross-platform terminal tool for network diagnostics, system monitoring, and performance optimization. Built with Go 1.25.0, it provides an interactive TUI with real-time graphs, comprehensive speed testing, and system information display.

## Quick Start

Get started in seconds with a pre-built binary:

```bash
# Download the latest release for your platform from https://github.com/rkriad585/auraspeed/releases
# Or build from source:
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed
# Linux/macOS
bash build.sh
# Windows
.\build.ps1
```

## Features

- **Interactive TUI** — Launch with `auraspeed tui` for real-time throughput graphs and live metrics
- **Speed Test** — Measure download/upload speeds and latency with `auraspeed speedtest`
- **System Information** — Display CPU, memory, and disk usage via `auraspeed info`
- **Network Diagnostics** — Run ping, DNS, and host lookups with `auraspeed network`
- **Web Server** — Start HTTP API server with `auraspeed web` for remote access
- **Test History** — Track and export past speed tests using `auraspeed history`
- **Flexible Configuration** — Manage settings via TOML config or `auraspeed config` commands
- **Clipboard Integration** — Copy results to clipboard with a single keypress in TUI mode
- **Cross-Platform** — Supports Windows, Linux, and macOS (Intel & Apple Silicon)
- **Command Aliases** — Quick shortcuts: `st` (speedtest), `si` (info), `net` (network), `hist` (history)

## Installation

### Quick Install (Recommended)

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/rkriad585/auraspeed/main/install.ps1 | iex
```

**Linux/macOS (Bash):**
```bash
curl -sL https://raw.githubusercontent.com/rkriad585/auraspeed/main/install.sh | bash
```

The installer will:
- Download the latest binary to `~/.config/neostore/auraspeed/bin/`
- Add it to your PATH
- Print next steps

To uninstall, run: `auraspeed --selfuninstall`

### Prerequisites

- Go 1.25.0+ (only required for building from source)

### From Source

Clone the repository and run the appropriate build script for your platform:

```bash
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed
```

**Linux/macOS (Bash/Zsh/Fish):**
```bash
bash build.sh
```

**Windows (PowerShell):**
```powershell
.\build.ps1
```

The compiled binary will be output to the project root.

### Pre-built Binaries

Download the latest pre-built binary for your platform from the [GitHub Releases page](https://github.com/rkriad585/auraspeed/releases):

| Platform | Binary Name |
|----------|-------------|
| Windows (amd64) | `auraspeed-windows-amd64.exe` |
| Linux (amd64) | `auraspeed-linux-amd64` |
| macOS (Intel) | `auraspeed-darwin-amd64` |
| macOS (Apple Silicon) | `auraspeed-darwin-arm64` | |

## Usage

### Web Server

Start a web server to access AuraSpeed features via HTTP:

```bash
auraspeed web
# Custom port
auraspeed web --port 3000
```

**Endpoints:**
- `GET /` — HTML UI with buttons to run tests
- `POST /api/speedtest` — Run speed test
- `GET /api/info` — Get system information
- `GET /health` — Health check

**Access:** http://localhost:8080

---

### Interactive TUI

Launch the full-featured terminal interface with real-time graphs:

```bash
auraspeed tui
```

**TUI Keyboard Shortcuts:**
- `R` — Restart speed test
- `C` — Copy results to clipboard
- `H` — View test history
- `Esc` — Close popups
- `Ctrl+C` — Exit

### Speed Test

Run a speed test from the CLI:

```bash
auraspeed speedtest
```

**With a specific server:**
```bash
auraspeed speedtest --server-id 1234
```

**JSON output (for scripting):**
```bash
auraspeed speedtest --json
```

### System Information

Display CPU, memory, and disk usage:

```bash
auraspeed info
```

Example output:
```
System Information
------------------
OS:       Windows 10 22H2
CPU:      Intel(R) Core(TM) i7-10700K (16 cores)
Memory:   32 GB total, 16 GB used (50%)
Disk:     512 GB total, 256 GB used (50%)
Hostname: DESKTOP-ABC123
```

### Network Diagnostics

Run network diagnostics with the `network` subcommand:

**Ping a host:**
```bash
auraspeed network ping google.com
```

**DNS lookup:**
```bash
auraspeed network dns google.com
```

**Host lookup:**
```bash
auraspeed network host google.com
```

### Test History

View and export past speed test results:

**List all history:**
```bash
auraspeed history
```

**Limit to last N results:**
```bash
auraspeed history --limit 10
```

**Export to JSON/TSV:**
```bash
auraspeed history --export results.json
auraspeed history --export results.tsv
```

### Configuration Management

Manage AuraSpeed settings via CLI:

**View all settings:**
```bash
auraspeed config view
```

**View a specific section:**
```bash
auraspeed config view speedtest
auraspeed config view ui
```

**Set a value:**
```bash
auraspeed config set speedtest.timeout 60
auraspeed config set ui.theme dark
```

**Reset to defaults:**
```bash
auraspeed config reset
```

### Command Aliases

Use these shortcuts for common commands:
- `auraspeed st` → `auraspeed speedtest`
- `auraspeed si` → `auraspeed info`
- `auraspeed net` → `auraspeed network`
- `auraspeed hist` → `auraspeed history`

### Global Flags

```bash
auraspeed --log-level debug    # Set log level (debug, info, warn, error)
auraspeed --no-color          # Disable colored output
auraspeed -v                  # Enable verbose logging
```

## Configuration

AuraSpeed uses TOML for configuration via [spf13/viper](https://github.com/spf13/viper). The config file is located at:
- Linux/macOS: `~/.auraspeed/config.toml`
- Windows: `%USERPROFILE%\.auraspeed\config.toml`

### Default Configuration

```toml
[global]
  loglevel = "info"
  nocolor = false
  autoupdate = true
  confirmexit = false

[speedtest]
  timeout = 30
  defaultserverid = 0
  paralleldownloads = 4
  paralleluploads = 2

[ui]
  theme = "default"
  graphheight = 8
  historylimit = 100
  autorefresh = false
  refreshrate = 5
  savehistory = true

[aliases]
  st = "speedtest"
  si = "info"
  net = "network"
  hist = "history"
```

### Configuration Options

| Section | Key | Type | Default | Description |
|---------|-----|------|---------|-------------|
| global | loglevel | string | "info" | Log level: debug, info, warn, error |
| global | nocolor | bool | false | Disable colored output |
| global | autoupdate | bool | true | Enable automatic updates |
| global | confirmexit | bool | false | Confirm before exiting |
| speedtest | timeout | int | 30 | Test timeout in seconds |
| speedtest | defaultserverid | int | 0 | Default server ID (0 = auto) |
| speedtest | paralleldownloads | int | 4 | Parallel download connections |
| speedtest | paralleluploads | int | 2 | Parallel upload connections |
| ui | theme | string | "default" | UI theme name |
| ui | graphheight | int | 8 | Height of graph in rows |
| ui | historylimit | int | 100 | Max history entries to keep |
| ui | autorefresh | bool | false | Auto-refresh TUI display |
| ui | refreshrate | int | 5 | Refresh rate in seconds |
| ui | savehistory | bool | true | Save test results to history |

## Project Structure

```
auraspeed/
├── cmd/
│   ├── main.go           # Entry point (package main)
│   └── root/
│       ├── root.go       # Root command & CLI setup (cobra)
│       ├── commands.go   # All subcommands
│       ├── web.go        # Web server command
│       ├── web.html      # Web UI
│       ├── update.go     # Update command
│       ├── install.go    # Install command
│       └── install.sh    # Linux/macOS installer script
├── internal/
│   ├── config/           # Configuration management (viper)
│   ├── info/             # System information (gopsutil)
│   ├── logging/          # Logging utilities (zerolog)
│   ├── network/          # Network diagnostics
│   ├── speedtest/        # Speed test & TUI (speedtest-go, tview)
│   └── ui/               # UI command wrapper
├── Dockerfile            # Docker container definition
├── docker-compose.yml    # Docker Compose configuration
├── auraspeed.service     # systemd service file
├── install.ps1           # Windows installer script
├── build.ps1             # PowerShell build script
├── build.sh              # Bash/Zsh/Fish build script
├── go.mod                # Go module definition (auraspeed)
└── README.md             # This file
```

## Dependencies

AuraSpeed is built with these excellent Go libraries:
- [spf13/cobra](https://github.com/spf13/cobra) — CLI framework
- [spf13/viper](https://github.com/spf13/viper) — Configuration management
- [rivo/tview](https://github.com/rivo/tview) — Terminal UI framework
- [gdamore/tcell/v2](https://github.com/gdamore/tcell) — Terminal handling
- [showwin/speedtest-go](https://github.com/showwin/speedtest-go) — Speed test library
- [shirou/gopsutil](https://github.com/shirou/gopsutil) — System information
- [rs/zerolog](https://github.com/rs/zerolog) — Logging
- [atotto/clipboard](https://github.com/atotto/clipboard) — Clipboard access
- [guptarohit/asciigraph](https://github.com/guptarohit/asciigraph) — ASCII graph generation

## Documentation

Full documentation is available in the [docs/](docs/) folder:

| Document | Description |
|----------|-------------|
| [Getting Started](docs/getting-started.md) | Beginner-friendly setup guide |
| [API Reference](docs/api-reference.md) | CLI commands and Go functions |
| [Architecture](docs/architecture.md) | System design and data flow |

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Quick start:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Ensure your changes pass `go test ./...` and adhere to Go best practices.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---
**Note**: Speed test functionality uses the Speedtest.net infrastructure via the speedtest-go library.
