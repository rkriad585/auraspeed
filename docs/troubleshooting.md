# AuraSpeed Troubleshooting Guide

Common issues and their solutions.

## Build Issues

### Error: `go: module not found`

**Cause:** Dependencies not downloaded.

**Fix:**
```bash
go mod download
go mod tidy
```

---

### Error: `build failed: build for auraspeed does not contain a main function`

**Cause:** GoReleaser can't find the entry point.

**Fix:** Ensure `.goreleaser.yaml` has the correct `main` field:
```yaml
builds:
  - main: ./cmd/main.go
    # ...
```

---

### Error: `goreleaser: command not found`

**Cause:** GoReleaser not installed.

**Fix (local testing):**
```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Run snapshot build
goreleaser build --snapshot --clean
```

---

## Speed Test Issues

### Error: `Cannot fetch server list`

**Cause:** Network connectivity or Speedtest.net API issue.

**Fix:**
1. Check internet connection
2. Retry after a few minutes
3. Try with a specific server ID:
   ```bash
   auraspeed speedtest --server-id 1234
   ```

---

### Error: `Download Test Failed`

**Cause:** Speedtest-go library bug or network issue.

**Fix:**
1. Update dependencies:
   ```bash
   go get github.com/showwin/speedtest-go@latest
   go mod tidy
   ```
2. Check if server context is initialized (bug fix in v1.0.1+)

---

### Panic: `runtime error: invalid memory address or nil pointer dereference`

**Cause:** Server `Context` field is nil in speedtest-go library.

**Fix:** This was fixed in commit `945420e`. Update to latest version:
```bash
git pull origin main
go build -o auraspeed ./cmd/main.go
```

---

## TUI Issues

### Issue: Garbled Output / Corrupted Display

**Cause:** Terminal doesn't support tview or asciigraph escape sequences.

**Fix:**
1. Use a modern terminal (Windows Terminal, iTerm2, GNOME Terminal)
2. Disable graph rendering (edit config):
   ```toml
   [ui]
     graphheight = 0
   ```
3. Use CLI mode instead: `auraspeed speedtest`

---

### Issue: TUI Doesn't Launch

**Cause:** Terminal not compatible with tcell.

**Fix:**
```bash
# Set TERM environment variable (Unix)
export TERM=xterm-256color
auraspeed tui

# Windows: Use Windows Terminal
```

---

## Configuration Issues

### Error: `Config file not found`

**Cause:** First run or config file deleted.

**Fix:** AuraSpeed creates a default config on first run. If issues persist:
```bash
# Windows
notepad %USERPROFILE%\.config\neostore\auraspeed\config.toml

# Linux/macOS
nano ~/.config/neostore/auraspeed/config.toml
```

---

### Error: `Invalid TOML syntax`

**Cause:** Manual edit introduced syntax error.

**Fix:** Validate TOML syntax:
```bash
# Use a TOML validator or check with:
auraspeed config view
# Error message will indicate the issue
```

---

## Network Diagnostics Issues

### Error: `ping: command not found` (Windows)

**Cause:** Ping not in PATH or using unsupported terminal.

**Fix:** Use Windows Command Prompt or PowerShell:
```powershell
auraspeed.exe network ping google.com
```

---

### Error: `traceroute: command not found` (Windows)

**Cause:** Windows uses `tracert` instead of `traceroute`.

**Fix:** The `network` command should handle this. If not:
```powershell
tracert google.com
```

---

## History Issues

### Error: `History file not found`

**Cause:** No tests run yet or history disabled.

**Fix:**
1. Run a speed test first: `auraspeed speedtest`
2. Enable history in config:
   ```toml
   [ui]
     savehistory = true
   ```

---

### Error: Export fails with permission denied

**Cause:** No write permission to target directory.

**Fix:**
```bash
# Linux/macOS
chmod u+w /path/to/export/

# Windows PowerShell
icacls "C:\path\to\export" /grant "$env:USERNAME:F"
```

---

## General Debugging

### Enable Debug Logging

```bash
auraspeed --log-level debug speedtest
```

### Check Version and Build Info

```bash
auraspeed --version
# or
auraspeed -v
```

### Verify Installation

```bash
# Check binary exists
which auraspeed
# or (Windows)
where auraspeed

# Check dependencies
go version
```

---

## Reporting Bugs

When opening a [GitHub issue](https://github.com/rkriad585/auraspeed/issues), include:

1. **AuraSpeed version:** `auraspeed --version`
2. **OS and architecture:** Windows 10, Ubuntu 22.04, macOS 14, etc.
3. **Go version:** `go version`
4. **Terminal type:** Windows Terminal, iTerm2, etc.
5. **Full error message** (copy/paste, no screenshots)
6. **Steps to reproduce** the issue

---

## Links

- [GitHub Issues](https://github.com/rkriad585/auraspeed/issues)
- [Getting Started Guide](getting-started.md)
- [API Reference](api-reference.md)
- [Configuration Reference](configuration.md)
