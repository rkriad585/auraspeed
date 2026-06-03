# AuraSpeed

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)](https://go.dev/)
[![MIT License](https://img.shields.io/github/license/rkriad585/auraspeed?color=blue)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/rkriad585/auraspeed)](https://github.com/rkriad585/auraspeed/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/rkriad585/auraspeed)](https://goreportcard.com/report/github.com/rkriad585/auraspeed)

Cross-platform terminal tool for network diagnostics, system monitoring, and performance optimization. Built with Go, it provides an interactive TUI with real-time graphs, comprehensive speed testing, system information, and a web dashboard.

## Features

- **Interactive TUI** — Real-time throughput graphs, live metrics, and keyboard navigation
- **Speed Test** — Download/upload speed and latency measurement via Speedtest.net
- **System Information** — CPU, memory, disk, and hostname display
- **Network Diagnostics** — Ping, traceroute, DNS lookup, and host lookup
- **Web Dashboard** — HTTP API server with HTML UI for remote access
- **13 Color Themes** — Built-in themes with dark/light mode, interactive theme selector
- **Test History** — Track, view, and export past speed test results
- **Clipboard Integration** — Copy results with a single keypress in TUI mode
- **Command Aliases** — Shortcuts: `st`, `si`, `net`, `hist`
- **Docker Support** — Multi-stage production-ready image
- **Cross-Platform** — Windows, Linux, macOS (Intel & Apple Silicon)

## Installation

### Quick Install

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/rkriad585/auraspeed/main/install.ps1 | iex
```

**Linux/macOS:**
```bash
curl -sL https://raw.githubusercontent.com/rkriad585/auraspeed/main/install.sh | bash
```

### From Source

```bash
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed
go build -o auraspeed ./cmd/main.go
```

Or use the build scripts:
```bash
./build.sh          # Linux/macOS
.\build.ps1          # Windows
```

### Using Docker

```bash
docker build -t auraspeed .
docker run -d --name auraspeed -p 59733:59733 auraspeed
```

Or with Docker Compose:
```bash
docker-compose up -d
```

See [Docker Usage](#docker-usage) for details.

## Usage

### Interactive TUI

```bash
auraspeed tui
```

**Keyboard shortcuts:** `R` Restart | `C` Copy | `H` History | `?` Help | `Esc` Close | `Ctrl+C` Exit

### Speed Test

```bash
auraspeed speedtest
auraspeed speedtest --json          # JSON output for scripting
auraspeed speedtest --server-id 1234 # Specific server
```

### System Information

```bash
auraspeed info
auraspeed info --json
```

### Network Diagnostics

```bash
auraspeed network ping google.com
auraspeed network dns google.com
auraspeed network traceroute google.com
```

### Web Server

```bash
auraspeed web
# http://localhost:59733
```

### Test History

```bash
auraspeed history                  # View history
auraspeed history --limit 10       # Last 10 results
auraspeed history --export results.json
```

### Configuration Management

```bash
auraspeed config                   # View all settings
auraspeed config theme             # Interactive theme selector
auraspeed config theme dark        # Set theme directly
auraspeed config theme --list      # List available themes
auraspeed config toggle-dark       # Toggle dark mode override
auraspeed config set speedtest.timeout 60
auraspeed config reset
```

### Command Aliases

```
auraspeed st    → auraspeed speedtest
auraspeed si    → auraspeed info
auraspeed net   → auraspeed network
auraspeed hist  → auraspeed history
```

## Configuration

AuraSpeed uses TOML configuration at:
- **Linux/macOS:** `~/.config/neostore/auraspeed/config.toml`
- **Windows:** `%USERPROFILE%\.config\neostore\auraspeed\config.toml`

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
  theme = "sunny-beach"
  darkmode = false
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

### Project Directories

| Directory | Path | Purpose |
|-----------|------|---------|
| Config | `~/.config/neostore/auraspeed/` | Config file, aliases |
| Data | `~/.config/neostore/auraspeed/data/` | History, server cache |
| Logs | `~/.config/neostore/auraspeed/logs/` | Application logs |
| Downloads | `~/Downloads/neostore/auraspeed/` | Exported output files |

### Configuration API

The `internal/config` package provides these helpers:

```go
import "auraspeed/internal/config"

config.GetConfigDir()     // ~/.config/neostore/auraspeed/
config.ConfigFile("x")   // ~/.config/neostore/auraspeed/x
config.GetLogsDir()      // ~/.config/neostore/auraspeed/logs/
config.GetDownloadsDir() // ~/Downloads/neostore/auraspeed/
config.EnsureConfigDir() // create dirs if missing
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
| ui | theme | string | "sunny-beach" | UI theme name |
| ui | darkmode | bool | false | Force dark theme |
| ui | graphheight | int | 8 | Height of graph in rows |
| ui | historylimit | int | 100 | Max history entries |
| ui | autorefresh | bool | false | Auto-refresh TUI display |
| ui | refreshrate | int | 5 | Refresh rate in seconds |
| ui | savehistory | bool | true | Save test results |

### Themes

AuraSpeed includes 13 built-in color themes. Run `auraspeed config theme` for an interactive selector.

| Theme | Type |
|-------|------|
| Dark Theme | Dark |
| Light Theme | Light |
| Sunny Beach Day | Dark |
| Olive Garden Feast | Dark |
| Summer Ocean Breeze | Dark |
| Refreshing Summer Fun | Dark |
| Black & Gold Elegance | Dark |
| Vibrant Color Fiesta | Dark |
| Light Steel | Light |
| Golden Twilight | Dark |
| Deep Sea | Dark |
| Bright Green | Dark |
| Vivid Nightfall | Dark |

## Build Instructions

### Using Go Directly

```bash
go build -o auraspeed ./cmd/main.go
```

### Using Make

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make test           # Run tests
make lint           # Run linters
make format         # Format code
make clean          # Remove build artifacts
make release        # Build release binaries
make run            # Build and run
```

### Using CMake

```bash
cmake -B build
cmake --build build --target build
cmake --build build --target test
cmake --build build --target clean
```

### Build Scripts

```bash
./build.sh          # Linux/macOS (all platforms)
.\build.ps1         # Windows (all platforms)
```

All binaries include version info embedded via ldflags and output to `bin/`.

## Docker Usage

### Build Image

```bash
docker build -t auraspeed .
```

Or with version metadata:

```bash
docker build \
  --build-arg VERSION=v3.1.8 \
  --build-arg COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t auraspeed .
```

### Run Container

```bash
docker run -d \
  --name auraspeed \
  -p 59733:59733 \
  -v auraspeed-config:/root/.config/neostore/auraspeed \
  -e AURASPEED_LOGLEVEL=info \
  -e AURASPEED_AUTOUPDATE=false \
  auraspeed
```

### Docker Compose

```bash
docker-compose up -d
```

The `docker-compose.yml` includes health checks, restart policy, and persistent volume for configuration.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AURASPEED_LOGLEVEL` | `info` | Log level override |
| `AURASPEED_AUTOUPDATE` | `false` | Disable auto-update in container |
| `AS_CONFIG_PATH` | `/app/config.toml` | Custom config file path |
| `AS_PORT` | `59733` | Web server port |

## Examples

### Run Speed Test with JSON Output

```bash
auraspeed speedtest --json
```

Output:
```json
{
  "download": 85.42,
  "upload": 42.18,
  "ping": 12,
  "isp": "Comcast",
  "server": "New York, NY"
}
```

### Monitor Connection with a Script

```bash
#!/bin/bash
while true; do
  echo "=== $(date) ==="
  auraspeed speedtest --json >> speed_log.json
  sleep 3600
done
```

### Change Theme Non-Interactively

```bash
auraspeed config theme dark
```

### Export History to JSON

```bash
auraspeed history --export results.json
```

## Project Structure

```
auraspeed/
├── bin/                    # Cross-platform build output
├── cmd/
│   ├── main.go             # Entry point
│   └── root/               # CLI commands (cobra)
│       ├── root.go         # Root command
│       ├── commands.go     # All subcommands
│       ├── web.go          # Web server
│       ├── update.go       # Update command
│       └── install.go      # Install command
├── docs/                   # Documentation
├── internal/
│   ├── config/             # Configuration management (viper)
│   ├── http/               # HTTP client utilities
│   ├── info/               # System information (gopsutil)
│   ├── logging/            # Logging (zerolog)
│   ├── network/            # Network diagnostics
│   ├── speedtest/          # Speed test & TUI (tview)
│   ├── theme/              # Color theme system
│   └── ui/                 # UI command wrapper
├── .dockerignore
├── .version                # Current version
├── build.ps1               # Windows build script
├── build.sh                # Unix build script
├── CMakeLists.txt          # CMake build configuration
├── docker-compose.yml      # Docker Compose
├── Dockerfile              # Docker image
├── go.mod / go.sum
└── Makefile                # Build automation
```

## Development Workflow

1. **Fork and clone** the repository
2. **Install dependencies:** `go mod download`
3. **Make changes** following Go best practices
4. **Format:** `go fmt ./...`
5. **Lint:** `go vet ./...` or `make lint`
6. **Test:** `go test ./... -v -cover` or `make test`
7. **Build:** `make build`
8. **Commit** with [Conventional Commits](https://www.conventionalcommits.org/)

Branch naming: `feat/description`, `fix/description`, `docs/description`

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for full guidelines.

Quick summary:
- Fork the repo and create a feature branch
- Write tests for new functionality
- Ensure `go test ./...` passes
- Open a Pull Request

## License

MIT License — see [LICENSE](LICENSE) for details.

Copyright (c) 2025 AuraSpeed Contributors
