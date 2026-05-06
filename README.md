# AuraSpeed

AuraSpeed is a cross-platform terminal tool for network diagnostics, system monitoring, and performance optimization. Built with Go, it provides an interactive TUI (Terminal User Interface) mode with real-time graphs, comprehensive speed testing, and system information display.

## Features

- **Interactive TUI Mode** - Full-featured terminal interface with real-time throughput graphs
- **Speed Test** - Measure download/upload speeds with configurable parallel connections
- **System Information** - View CPU, memory, disk usage, and hostname
- **Network Diagnostics** - Ping, traceroute, and DNS lookup utilities
- **Test History** - Track and export speed test results (JSON/TSV formats)
- **Configurable** - TOML-based configuration with sensible defaults
- **Command Aliases** - Quick shortcuts (`st`, `si`, `net`, `hist`)
- **Clipboard Integration** - Copy results with a keypress
- **Cross-Platform** - Works on Windows, Linux, and macOS (Intel & Apple Silicon)

## Installation

### Prerequisites

- Go 1.21 or higher (for building from source)

### From Source

```bash
git clone https://github.com/yourusername/auraspeed.git
cd auraspeed
make build
```

For Windows:
```bash
make build-windows
```

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/yourusername/auraspeed/releases):

- **Windows**: `auraspeed-windows-amd64.exe`
- **Linux**: `auraspeed-linux-amd64`
- **macOS (Intel)**: `auraspeed-darwin-amd64`
- **macOS (Apple Silicon)**: `auraspeed-darwin-arm64`

### Build All Platforms

```bash
make build-all
```

## Usage

### TUI Mode (Interactive)

Launch the full interactive terminal interface:

```bash
auraspeed tui
```

**TUI Keyboard Shortcuts:**
- `R` - Restart speed test
- `C` - Copy results to clipboard
- `H` - View test history
- `Esc` - Close popups
- `Ctrl+C` - Exit

### Speed Test (CLI)

Run a quick speed test from the command line:

```bash
auraspeed speedtest
```

**With specific server:**
```bash
auraspeed speedtest --server-id 1234
```

**JSON output (for scripting):**
```bash
auraspeed speedtest --json
```

### System Information

Display system details:

```bash
auraspeed info
```

Output example:
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

**Ping a host:**
```bash
auraspeed network ping google.com
```

**Traceroute:**
```bash
auraspeed network traceroute google.com
```

**DNS lookup:**
```bash
auraspeed network dns google.com
```

### View Test History

**Show all history:**
```bash
auraspeed history
```

**Limit to last N results:**
```bash
auraspeed history --limit 10
```

**Export to file (auto-detects format from extension):**
```bash
auraspeed history --export results.json
auraspeed history --export results.tsv
```

**Clear all history:**
```bash
auraspeed history --clear
```

### Configuration

**View all settings:**
```bash
auraspeed config view
```

**View specific section:**
```bash
auraspeed config view speedtest
auraspeed config view ui
auraspeed config view global
```

**Set a value:**
```bash
auraspeed config set speedtest.timeout 60
auraspeed config set ui.theme dark
auraspeed config set ui.savehistory true
```

**Reset to defaults:**
```bash
auraspeed config reset
```

### Command Aliases

Quick shortcuts are available:
- `auraspeed st` = `auraspeed speedtest`
- `auraspeed si` = `auraspeed info`
- `auraspeed net` = `auraspeed network`
- `auraspeed hist` = `auraspeed history`

### Global Flags

```bash
auraspeed --log-level debug    # Set log level (debug, info, warn, error)
auraspeed --no-color          # Disable colored output
auraspeed -v                  # Enable verbose logging
```

## Configuration

Config file location: `~/.auraspeed/config.toml` (Windows: `%USERPROFILE%\.auraspeed\config.toml`)

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

## TUI Mode Screenshot

```
─── AURASPEED PRO: ADVANCED NETWORK ANALYZER ───

ISP: Comcast (203.0.113.1)  Server: New York, NY

Running Download Test...

   80.0 ┤     ╭╮
   70.0 ┤    ╭╯╰╮
   60.0 ┤   ╭╯  ╰╮
   50.0 ┤  ╭╯    ╰╮
   40.0 ┤ ╭╯      ╰╮
   30.0 ┼╭╯        ╰╮
   20.0 ┤╯           ╰─
   10.0 ┤
    0.0 ┼──────────────
        Download (Mbps)

┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ DOWNLOAD     │ │  UPLOAD      │ │  LATENCY     │ │   JITTER     │
│   (Mbps)     │ │    (Mbps)    │ │    (ms)      │ │    (ms)      │
│              │ │              │ │              │ │              │
│    85.42     │ │    42.18     │ │      12      │ │      3       │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘

Press C Copy Link | R Restart | H History | Esc Close Popups | Ctrl+C Exit
```

## Project Structure

```
auraspeed/
├── cmd/
│   ├── main.go           # Entry point
│   └── root/
│       ├── root.go       # Root command & CLI setup
│       └── commands.go   # All subcommands
├── internal/
│   ├── config/          # Configuration management
│   ├── info/            # System information
│   ├── logging/         # Logging utilities
│   ├── network/         # Network diagnostics
│   ├── speedtest/       # Speed test & TUI
│   └── ui/              # UI command wrapper
├── Makefile             # Build targets
├── go.mod               # Go module definition
└── README.md            # This file
```

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Building for Development

```bash
go build -o auraspeed ./cmd/main.go
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [tview](https://github.com/rivo/tview) - Terminal UI framework
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [asciigraph](https://github.com/guptarohit/asciigraph) - ASCII graph generation
- [speedtest-go](https://github.com/showwin/speedtest-go) - Speed test library
- [gopsutil](https://github.com/shirou/gopsutil) - System information
- [zerolog](https://github.com/rs/zerolog) - Logging
- [clipboard](https://github.com/atotto/clipboard) - Clipboard access

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Guidelines

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed
- Ensure `make test` and `make lint` pass

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by the need for a fast, terminal-based network diagnostics tool
- Built with excellent Go libraries from the open-source community

---

**Note**: Speed test functionality uses the Speedtest.net infrastructure via the speedtest-go library.
