# AuraSpeed FAQ

Frequently asked questions about AuraSpeed.

---

## General Questions

### What is AuraSpeed?

AuraSpeed is a cross-platform terminal tool for network diagnostics, system monitoring, and performance optimization. It provides:
- Interactive TUI with real-time graphs
- Speed test with download/upload measurements
- System information display (CPU, RAM, disk)
- Network diagnostics (ping, traceroute, DNS)

---

### Is AuraSpeed free?

Yes! AuraSpeed is open-source software licensed under the MIT License.

---

### Which platforms are supported?

AuraSpeed supports:
- **Windows** (amd64)
- **Linux** (amd64)
- **macOS** (Intel amd64 & Apple Silicon arm64)

---

### What terminal should I use?

**Windows:**
- Windows Terminal (recommended)
- PowerShell
- Command Prompt

**Linux/macOS:**
- Any modern terminal (GNOME Terminal, iTerm2, Alacritty, etc.)

---

## Installation

### How do I install AuraSpeed?

**Option 1: Download pre-built binary**
1. Go to [Releases](https://github.com/rkriad585/auraspeed/releases)
2. Download the binary for your platform
3. Add to PATH or run directly

**Option 2: Build from source**
```bash
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed
# Linux/macOS
./build.sh
# Windows
.\build.ps1
```

---

### Do I need Go installed?

Only if building from source. Pre-built binaries don't require Go.

---

### Where is the config file?

| OS | Path |
|----|------|
| Windows | `%USERPROFILE%\.config\neostore\auraspeed\config.toml` |
| Linux | `~/.config/neostore/auraspeed/config.toml` |
| macOS | `~/.config/neostore/auraspeed/config.toml` |

### Where are logs and exported files stored?

| Item | Path |
|------|------|
| History | `~/.config/neostore/auraspeed/data/history.json` |
| Logs | `~/.config/neostore/auraspeed/logs/` |
| Exported files | `~/Downloads/neostore/auraspeed/` |

---

## Speed Test

### How accurate is the speed test?

AuraSpeed uses the Speedtest.net infrastructure via the speedtest-go library. Results are comparable to the official Speedtest.net website or app.

---

### Why is my speed test result different from other tools?

- Different servers may be selected
- Network conditions change over time
- Some tools use different measurement methods

For consistent results, use the same server:
```bash
auraspeed speedtest --server-id 1234
```

---

### Can I test a specific server?

Yes! First, run a test to cache servers, then use a specific ID:
```bash
auraspeed speedtest  # Lists available servers
auraspeed speedtest --server-id 1234
```

---

### Why is upload speed slower than download?

This is normal. Most ISPs provide asymmetric speeds (faster download than upload).

---

## TUI Mode

### How do I launch the TUI?

```bash
auraspeed tui
```

---

### What do the keyboard shortcuts do?

| Key | Action |
|-----|--------|
| `R` | Restart speed test |
| `C` | Copy results to clipboard |
| `H` | View test history |
| `Esc` | Close popups |
| `Ctrl+C` | Exit |

---

### Can I disable the graph?

Yes, set graph height to 0 in config:
```toml
[ui]
  graphheight = 0
```

---

## Configuration

### How do I change settings?

**Via CLI:**
```bash
auraspeed config set speedtest.timeout 60
```

**Via config file:**
Edit `~/.config/neostore/auraspeed/config.toml` (or Windows equivalent)

---

### How do I reset to defaults?

```bash
auraspeed config reset
```

---

### Can I use environment variables?

Currently, AuraSpeed uses TOML config file only. Environment variables are not supported yet.

---

## Network Diagnostics

### What network tools are available?

```bash
auraspeed network ping <host>      # Ping a host
auraspeed network traceroute <host> # Traceroute
auraspeed network dns <host>        # DNS lookup
auraspeed network host <host>       # Host lookup
```

---

### Why does traceroute take so long?

Traceroute can take time as it probes each hop. Typical time is 10-30 seconds.

---

## History

### Where is history stored?

History is stored in `history.json` in the same directory as the config file.

---

### How do I export history?

```bash
# Export to JSON
auraspeed history --export results.json

# Export to TSV
auraspeed history --export results.tsv
```

---

### Can I limit the number of history entries?

Yes, set `historylimit` in config:
```toml
[ui]
  historylimit = 50  # Keep only last 50 entries
```

---

## Development

### How do I contribute?

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

### How do I run tests?

```bash
go test ./... -v -cover
```

---

### How do I build for all platforms?

```bash
# Linux/macOS
./build.sh

# Windows
.\build.ps1
```

---

## Troubleshooting

### AuraSpeed crashes on startup

1. Check Go installation: `go version`
2. Rebuild from source
3. Check for config file corruption

---

### Speed test hangs or times out

1. Check internet connection
2. Increase timeout: `auraspeed config set speedtest.timeout 60`
3. Try a different server

---

### More help?

- Check [Troubleshooting Guide](troubleshooting.md)
- Open an [issue](https://github.com/rkriad585/auraspeed/issues)

---

## Links

- [GitHub Repository](https://github.com/rkriad585/auraspeed)
- [Getting Started Guide](getting-started.md)
- [API Reference](api-reference.md)
- [Configuration Reference](configuration.md)
- [Troubleshooting Guide](troubleshooting.md)
