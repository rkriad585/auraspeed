# AuraSpeed Architecture

A high-performance CLI tool for network speed testing and system diagnostics, built with Go. AuraSpeed combines real-time terminal UI, comprehensive network diagnostics, and system information gathering into a single extensible binary.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Layer (cmd/)                        │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌─────────┐  │
│  │  main.go  │→│ root.go   │→│commands.go│  │ version │  │
│  └───────────┘  └───────────┘  └───────────┘  └─────────┘  │
│         ↓                ↓                ↓                  │
└─────────────────────────────────────────────────────────────┘
         │                │                │
         ↓                ↓                ↓
┌─────────────────────────────────────────────────────────────┐
│                   Internal Packages                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  config  │  │ speedtest│  │   info   │  │ logging  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│  ┌──────────┐  ┌──────────┐                                 │
│  │ network  │  │    ui    │                                 │
│  └──────────┘  └──────────┘                                 │
└─────────────────────────────────────────────────────────────┘
         │                │                │
         ↓                ↓                ↓
┌─────────────────────────────────────────────────────────────┐
│                  External Dependencies                       │
│  cobra │ viper │ tview │ speedtest-go │ gopsutil │ zerolog  │
└─────────────────────────────────────────────────────────────┘
```

## Package Structure

### `cmd/` — CLI Entry Point

| File | Responsibility |
|------|----------------|
| `main.go` | Binary entry point, calls `root.Execute()` |
| `root/root.go` | Root Cobra command, persistent flags, CLI setup |
| `root/commands.go` | All subcommands: `speedtest`, `info`, `network`, `tui`, `history` |
| `root/version.go` | Version command, ldflags injection |
| `root/servers.json` | Cached speedtest server list |

> [!NOTE]
> Version information is injected at build time via `-ldflags` in `.goreleaser.yaml`.

### `internal/config` — Configuration Management

Uses **Viper** for hierarchical configuration with TOML files.

| File | Responsibility |
|------|----------------|
| `config.go` | Load/save config, defaults, validation |
| `config_test.go` | Unit tests for config operations |

**Config directory:** `~/.config/neostore/auraspeed/`  
**Directory structure:**
```
~/.config/neostore/auraspeed/
├── config.toml          # Main configuration
├── data/
│   ├── history.json     # Speed test history
│   └── servers.json     # Cached server list
└── logs/                # Application logs
~/Downloads/neostore/auraspeed/  # Exported output files
```

### `internal/theme` — Color Theme System

Provides 13 built-in themes with semantic color slots (Background, Surface, Primary, Secondary, Accent, Foreground, Success, Warning, Error). Exposes a `ThemeManager` singleton for runtime theme switching.

| File | Responsibility |
|------|----------------|
| `theme.go` | Theme palettes, color definitions |
| `manager.go` | ThemeManager singleton, listeners |

### `internal/speedtest` — Core Speed Test Engine

The heart of AuraSpeed. Wraps `showwin/speedtest-go` with enhanced UI and data persistence.

| File | Responsibility |
|------|----------------|
| `tui.go` | TUI implementation using `tview` — real-time test display |
| `graph.go` | Custom ASCII graph renderer (not `asciigraph`) |
| `history.go` | Test result persistence, JSON/TSV export |
| `history_test.go` | History storage tests |
| `mock.go` | Mock speedtest server for testing |
| `server_cache.go` | Caches server list to `servers.json` |

**Key types:**
- `TestResult` — Single speed test outcome (download, upload, ping, jitter)
- `HistoryStore` — Manages test result persistence

### `internal/info` — System Information

Cross-platform system metrics using `shirou/gopsutil`.

| Capability | gopsutil Source |
|------------|-----------------|
| CPU info | `cpu.Info()` |
| Memory usage | `mem.VirtualMemory()` |
| Disk usage | `disk.Usage()` |
| Network interfaces | `net.Interfaces()` |

### `internal/logging` — Structured Logging

Zero-alloc logging with **zerolog**.

| Feature | Implementation |
|---------|----------------|
| Log levels | `debug`, `info`, `warn`, `error` |
| Colored output | Optional, controlled by `--no-color` |
| Output format | Structured JSON in debug mode |

### `internal/network` — Network Diagnostics

| Function | Description |
|----------|-------------|
| `Ping(host, count)` | ICMP/TCP ping with statistics |
| `Traceroute(host)` | Route tracing |
| `LookupDNS(host)` | DNS A/AAAA/MX/NS lookup |

### `internal/ui` — UI Utilities

| File | Responsibility |
|------|----------------|
| `command.go` | Wraps command execution with UI feedback |

## Data Flow

### Speed Test Flow

```
auraspeed speedtest [flags]
        │
        ↓
┌──────────────────────────────────┐
│  Cobra command RunE triggered    │
│  (commands.go: NewSpeedtestCmd)  │
└──────────────────────────────────┘
        │
        ↓
┌──────────────────────────────────┐
│  Load config (viper)             │
│  Setup logging (zerolog)         │
└──────────────────────────────────┘
        │
        ↓
┌──────────────────────────────────┐
│  speedtest-go:                    │
│  1. FetchUserInfo()              │
│  2. FetchServers() [cached?]     │
│  3. FindServer(id or nearest)    │
│  4. DownloadTest()               │
│  5. UploadTest()                 │
└──────────────────────────────────┘
        │
        ↓
┌──────────────────────────────────┐
│  Output results:                 │
│  - Terminal table (default)      │
│  - JSON (--json flag)            │
│  - Save to history               │
└──────────────────────────────────┘
```

### TUI Mode Flow

```
auraspeed tui
        │
        ↓
┌──────────────────────────────────┐
│  tview.Application created       │
│  - Header bar                    │
│  - Results table                 │
│  - ASCII graph (graph.go)        │
│  - Status bar                    │
└──────────────────────────────────┘
        │
        ↓
┌──────────────────────────────────┐
│  goroutine: runAdvancedTest()    │
│  ┌────────────────────────────┐  │
│  │ speedtest-go operations    │  │
│  │ ↓                          │  │
│  │ app.QueueUpdateDraw(func() │  │
│  │   { update widgets }       │  │
│  │ )                          │  │
│  └────────────────────────────┘  │
└──────────────────────────────────┘
        │
        ↓
┌──────────────────────────────────┐
│  refreshGraph() called per       │
│  interval → ASCII graph updated  │
└──────────────────────────────────┘
```

## Dependencies

| Library | Purpose | Import Path |
|---------|---------|-------------|
| `spf13/cobra` | CLI framework | `github.com/spf13/cobra` |
| `spf13/viper` | Config management | `github.com/spf13/viper` |
| `rivo/tview` | Terminal UI widgets | `github.com/rivo/tview` |
| `gdamore/tcell/v2` | Terminal input handling | `github.com/gdamore/tcell/v2` |
| `showwin/speedtest-go` | Speed test engine | `github.com/showwin/speedtest-go` |
| `shirou/gopsutil` | System information | `github.com/shirou/gopsutil/v3` |
| `rs/zerolog` | Structured logging | `github.com/rs/zerolog` |
| `atotto/clipboard` | Clipboard access | `github.com/atotto/clipboard` |

## Build System

```
┌───────────────────────────────────────────┐
│              Build Targets                │
├───────────────────────────────────────────┤
│  make build        → Local binary         │
│  make build-all    → All platforms        │
│  make test         → Run all tests        │
│  make lint         → golangci-lint + vet  │
│  make format       → go fmt               │
│  make release      → Release build        │
│  make clean        → Remove artifacts     │
│  cmake -B build    → CMake configure      │
│  cmake --build build → CMake build        │
│  ./build.sh        → Unix release build   │
│  ./build.ps1       → Windows build        │
│  docker build      → Docker image         │
└───────────────────────────────────────────┘
```

**GitHub Actions** (`.github/workflows/`) produces release binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

## Concurrency Model

```
┌─────────────────────────────────────────────────────┐
│  Main goroutine:                                    │
│  - Cobra command execution                          │
│  - tview event loop (in TUI mode)                  │
│                                                     │
│  Worker goroutines:                                 │
│  - runAdvancedTest() — speed test in background     │
│  - Graph updates via QueueUpdateDraw (thread-safe)  │
└─────────────────────────────────────────────────────┘
```

> [!IMPORTANT]
> All tview widget updates must happen on the main goroutine. Use `app.QueueUpdateDraw()` to safely update UI from background goroutines.

## Error Handling Strategy

| Layer | Strategy |
|-------|----------|
| CLI | Cobra `RunE` returns errors → formatted by Cobra |
| Speed test | Errors wrapped with context, logged via zerolog |
| Config | Defaults used if file missing/corrupt |
| Network | Graceful degradation, partial results returned |
| TUI | Errors displayed in status bar, test continues |

## Configuration Hierarchy

```
Defaults (code)
    ↓
~/.config/neostore/auraspeed/config.toml
    ↓
Environment variables (AURASPEED_*)
    ↓
CLI flags (highest priority)
```
