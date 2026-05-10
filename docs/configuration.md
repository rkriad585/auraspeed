# AuraSpeed Configuration Reference

Complete reference for all configuration options in AuraSpeed.

## Config File Location

The configuration file is stored in TOML format:

| OS | Path |
|----|------|
| Windows | `%USERPROFILE%\.config\neostore\auraspeed\config.toml` |
| Linux | `~/.config/neostore/auraspeed/config.toml` |
| macOS | `~/.config/neostore/auraspeed/config.toml` |

---

## Default Configuration

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

---

## Global Settings

Settings in the `[global]` section affect overall application behavior.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `loglevel` | string | `"info"` | Log level: `debug`, `info`, `warn`, `error` |
| `nocolor` | bool | `false` | Disable colored output |
| `autoupdate` | bool | `true` | Enable automatic updates (if implemented) |
| `confirmexit` | bool | `false` | Confirm before exiting application |

### Examples

**Enable debug logging:**
```bash
auraspeed config set global.loglevel debug
```

**Disable colored output:**
```bash
auraspeed config set global.nocolor true
```

---

## Speed Test Settings

Settings in the `[speedtest]` section control speed test behavior.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `timeout` | int | `30` | Test timeout in seconds |
| `defaultserverid` | int | `0` | Default server ID (`0` = auto-select) |
| `paralleldownloads` | int | `4` | Number of parallel download connections |
| `paralleluploads` | int | `2` | Number of parallel upload connections |

### Examples

**Increase test timeout to 60 seconds:**
```bash
auraspeed config set speedtest.timeout 60
```

**Set a specific default server:**
```bash
auraspeed config set speedtest.defaultserverid 1234
```

**Use more parallel connections:**
```bash
auraspeed config set speedtest.paralleldownloads 8
```

---

## UI Settings

Settings in the `[ui]` section control the terminal UI appearance and behavior.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `theme` | string | `"default"` | UI theme name |
| `graphheight` | int | `8` | Height of the throughput graph in rows |
| `historylimit` | int | `100` | Maximum number of history entries to keep |
| `autorefresh` | bool | `false` | Auto-refresh TUI display |
| `refreshrate` | int | `5` | Refresh rate in seconds (if autorefresh enabled) |
| `savehistory` | bool | `true` | Save test results to history file |

### Examples

**Set dark theme:**
```bash
auraspeed config set ui.theme dark
```

**Increase graph height:**
```bash
auraspeed config set ui.graphheight 12
```

**Disable history saving:**
```bash
auraspeed config set ui.savehistory false
```

---

## Command Aliases

The `[aliases]` section defines shortcuts for common commands.

| Alias | Expands To | Description |
|-------|-------------|-------------|
| `st` | `speedtest` | Run speed test |
| `si` | `info` | Show system information |
| `net` | `network` | Network diagnostics |
| `hist` | `history` | View test history |

### Custom Aliases

Add your own aliases by editing the config file:

```toml
[aliases]
  st = "speedtest"
  si = "info"
  net = "network"
  hist = "history"
  # Custom aliases
  test = "speedtest --json"
  myinfo = "info --verbose"
```

---

## Managing Configuration

### View All Settings

```bash
auraspeed config view
```

### View Specific Section

```bash
auraspeed config view speedtest
auraspeed config view ui
```

### Set a Value

```bash
auraspeed config set <key> <value>
```

Examples:
```bash
auraspeed config set speedtest.timeout 60
auraspeed config set ui.theme dark
auraspeed config set global.loglevel debug
```

### Reset to Defaults

```bash
auraspeed config reset
```

---

## Environment Variables

Currently, AuraSpeed uses config file only. Environment variables are not supported yet.

---

## Troubleshooting

### Config File Not Found

If the config file doesn't exist, AuraSpeed creates one with defaults on first run.

### Invalid TOML Syntax

If you manually edit the config file and introduce TOML syntax errors:

```bash
# Validate by running any command
auraspeed config view
# Error message will show the issue
```

### Permission Issues (Windows)

If you can't write to the config file:

```powershell
# Check permissions
Get-Acl "$env:USERPROFILE\.config\neostore\auraspeed\config.toml"

# Fix permissions
icacls "$env:USERPROFILE\.config\neostore\auraspeed\config.toml" /grant "$env:USERNAME:F"
```

---

## Links

- [GitHub Repository](https://github.com/rkriad585/auraspeed)
- [Getting Started Guide](getting-started.md)
- [API Reference](api-reference.md)
- [Examples](examples.md)
