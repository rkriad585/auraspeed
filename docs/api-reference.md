# AuraSpeed API Reference

> Complete reference for AuraSpeed CLI commands and internal Go packages.

## Table of Contents

- [CLI Commands](#cli-commands)
  - [Global Flags](#global-flags)
  - [`auraspeed speedtest`](#auraspeed-speedtest)
  - [`auraspeed info`](#auraspeed-info)
  - [`auraspeed network`](#auraspeed-network)
  - [`auraspeed history`](#auraspeed-history)
  - [`auraspeed config`](#auraspeed-config)
  - [`auraspeed tui`](#auraspeed-tui)
- [Internal Packages](#internal-packages)
  - [internal/speedtest/tui.go](#internalspeedtesttuigo)
  - [internal/speedtest/history.go](#internalspeedtesthistorygo)
  - [internal/speedtest/server_cache.go](#internalspeedtestserver_cachego)
  - [internal/speedtest/graph.go](#internalspeedtestgraphgo)
  - [internal/config/config.go](#internalconfigconfiggo)
  - [internal/info/info.go](#internalinfoinfo)
  - [internal/logging/logging.go](#internallogginglogginggo)

---

## CLI Commands

### Global Flags

These flags work with all `auraspeed` commands.

| Flag | Type | Description |
|------|------|-------------|
| `--log-level` | `string` | Set log level (`debug`, `info`, `warn`, `error`) |
| `--no-color` | — | Disable colored output |
| `-v`, `--verbose` | — | Enable verbose logging |

**Examples:**

```bash
auraspeed --log-level debug speedtest
auraspeed --no-color info
auraspeed -v network ping google.com
```

---

### `auraspeed speedtest`

Run a network speed test.

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--server-id` | `string` | Specific server ID to test against |
| `--json` | — | Output results as JSON |

**Examples:**

```bash
auraspeed speedtest
auraspeed speedtest --server-id 1234
auraspeed speedtest --json
```

---

### `auraspeed info`

Display system information (CPU, memory, disk, hostname).

**Examples:**

```bash
auraspeed info
```

---

### `auraspeed network`

Network diagnostics subcommands.

**Subcommands:**

| Subcommand | Syntax | Description |
|------------|--------|-------------|
| `ping` | `ping <host>` | Ping a host |
| `traceroute` | `traceroute <host>` | Run traceroute |
| `dns` | `dns <host>` | DNS lookup |

**Examples:**

```bash
auraspeed network ping google.com
auraspeed network traceroute google.com
auraspeed network dns google.com
```

---

### `auraspeed history`

View test history.

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | `int` | `10` | Limit number of results |
| `--export` | `string` | — | Export to file (auto-detects JSON/TSV from extension) |

**Examples:**

```bash
auraspeed history
auraspeed history --limit 5
auraspeed history --export results.json
auraspeed history --export results.tsv
```

---

### `auraspeed config`

Configuration management.

**Subcommands:**

| Subcommand | Syntax | Description |
|------------|--------|-------------|
| `view` | `view [section]` | View all or specific config section |
| `set` | `set <key> <value>` | Set config value |
| `reset` | `reset` | Reset to defaults |
| `theme` | `theme [name]` | Interactive theme selector or set by name |
| `toggle-dark` | `toggle-dark` | Toggle dark mode override |

**Examples:**

```bash
auraspeed config view
auraspeed config view speedtest
auraspeed config set speedtest.timeout 60
auraspeed config theme
auraspeed config theme dark
auraspeed config theme --list
auraspeed config toggle-dark
auraspeed config reset
```

---

### `auraspeed tui`

Launch the interactive TUI.

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `R` | Restart speed test |
| `C` | Copy results to clipboard |
| `H` | View history |
| `Esc` | Close popups |
| `Ctrl+C` | Exit |

**Example:**

```bash
auraspeed tui
```

---

## Internal Packages

### internal/speedtest/tui.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `RunTUI` | `RunTUI() error` | Launch the interactive TUI |
| `StopTUI` | `StopTUI()` | Stop the TUI if running |
| `IsTUIRunning` | `IsTUIRunning() bool` | Check if TUI is active |
| `updateStatus` | `updateStatus(msg string)` | Update status text in TUI |

**Examples:**

```go
// Launch the TUI
if err := speedtest.RunTUI(); err != nil {
    log.Fatal(err)
}

// Stop the TUI programmatically
speedtest.StopTUI()

// Check if TUI is running
if speedtest.IsTUIRunning() {
    // TUI is active
}

// Update status message inside TUI
speedtest.updateStatus("Running speed test...")
```

---

### internal/speedtest/history.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `saveToHistory` | `saveToHistory(dl, ul float64, ping int64, isp string)` | Save test result to history |
| `showHistory` | `showHistory()` | Display history in TUI |

**Examples:**

```go
// Save a completed test result
speedtest.saveToHistory(100.5, 50.2, 20, "MyISP")

// Display history view inside TUI
speedtest.showHistory()
```

---

### internal/speedtest/server_cache.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `GetCachedServers` | `GetCachedServers() (speedtest.Servers, bool)` | Get cached servers |
| `SaveServersToCache` | `SaveServersToCache(servers speedtest.Servers) error` | Cache servers to disk |

**Examples:**

```go
// Retrieve cached servers
servers, found := speedtest.GetCachedServers()
if found {
    // Use cached servers
    fmt.Printf("Found %d cached servers\n", len(servers))
}

// Save servers to cache
servers, _ := fetchServersFromAPI()
if err := speedtest.SaveServersToCache(servers); err != nil {
    log.Fatal(err)
}
```

---

### internal/speedtest/graph.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `refreshGraph` | `refreshGraph()` | Redraw the speed graph (throttled to 500ms) |
| `generateSimpleGraph` | `generateSimpleGraph(data []float64) string` | Generate ASCII bar graph from data |

**Examples:**

```go
// Trigger a graph redraw (throttled internally)
speedtest.refreshGraph()

// Generate an ASCII graph from speed data
data := []float64{10.5, 20.3, 30.1, 45.7, 60.2}
graph := speedtest.generateSimpleGraph(data)
fmt.Println(graph)
```

---

### internal/config/config.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `Init` | `Init(appName string) error` | Initialize configuration |
| `Get` | `Get() *Config` | Get current config struct |
| `Set` | `Set(key string, value interface{}) error` | Set a config value by key path |
| `Reset` | `Reset()` | Reset configuration to defaults |
| `EnsureSensitiveFilePermissions` | `EnsureSensitiveFilePermissions() error` | Fix file permissions on sensitive files |

**Examples:**

```go
// Initialize config for the application
if err := config.Init("auraspeed"); err != nil {
    log.Fatal(err)
}

// Retrieve and use current config
cfg := config.Get()
fmt.Printf("Timeout: %ds\n", cfg.Speedtest.Timeout)

// Set a config value
if err := config.Set("speedtest.timeout", 60); err != nil {
    log.Fatal(err)
}

// Reset to defaults
config.Reset()

// Ensure config file has secure permissions
if err := config.EnsureSensitiveFilePermissions(); err != nil {
    log.Fatal(err)
}
```

---

### internal/info/info.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `GetSystemInfo` | `GetSystemInfo() (*SystemInfo, error)` | Get CPU, RAM, disk, and hostname info |

**Examples:**

```go
sysInfo, err := info.GetSystemInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("CPU:    %s\n", sysInfo.CPU)
fmt.Printf("Memory: %d MB\n", sysInfo.Memory)
fmt.Printf("Disk:   %d GB\n", sysInfo.Disk)
fmt.Printf("Host:   %s\n", sysInfo.Hostname)
```

---

### internal/logging/logging.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `Get` | `Get() *Logger` | Get the logger instance |
| `SetLevel` | `SetLevel(level string) error` | Set the log level |
| `SetNoColor` | `SetNoColor(noColor bool)` | Toggle colored output |

**Examples:**

```go
// Get logger and write a log message
logger := logging.Get()
logger.Info("Application started")
logger.Debug("Debug details here")

// Change log level at runtime
if err := logging.SetLevel("debug"); err != nil {
    log.Fatal(err)
}

// Disable colored output (useful for logs to file)
logging.SetNoColor(true)
```
