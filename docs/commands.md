# Command Reference

AuraSpeed follows a sub-command structure. You can use `auraspeed [command] --help` for detailed flag information.

## Core Commands

### `tui`
Launches the interactive Terminal User Interface.
- **Flags:** `--fullscreen, -f`
- **Description:** Provides real-time throughput graphs and system metrics.

### `speedtest` (Alias: `st`)
Runs a non-interactive network speed test.
- **Flags:** 
  - `--server-id`: Specify a specific Ookla server ID.
  - `--json`: Output results in JSON format.
  - `--verbose`: Show detailed connection steps.
  - `--timeout`: Set test timeout in seconds.

### `info` (Alias: `si`)
Displays comprehensive system information.
- **Output:** OS version, CPU model/cores, Memory usage, and Disk statistics.

### `network` (Alias: `net`)
Suite of diagnostic tools.
- `auraspeed network ping [host]`: Standard ICMP ping.
- `auraspeed network traceroute [host]`: Trace network hops.
- `auraspeed network dns [host]`: Forward and reverse DNS lookups.

### `history` (Alias: `hist`)
Manage past test results.
- **Flags:**
  - `--limit`: Number of entries to show.
  - `--clear`: Wipe the history file.
  - `--save`: Export history to a file.

### `servers`
Manage speed test servers.
- `auraspeed servers list`: List nearby servers.
- `auraspeed servers favorites`: View your bookmarked servers.
- `auraspeed servers add [id]`: Add a server to favorites.

## System Commands

### `config`
Manage application settings.
- `auraspeed config view`: Show current configuration.
- `auraspeed config set [key] [value]`: Update a setting.
- `auraspeed config theme [name]`: Change UI theme (interactive if no name).
- `auraspeed config toggle-dark`: Toggle dark mode override.
- `auraspeed config reset`: Restore defaults.

### `self-update`
Checks for and replaces the running binary with the latest release from GitHub.
Supports Windows safe replacement (rename → replace → cleanup).
- **Flags:** `--proxy, -p` (Proxy URL for download).

### `web`
Starts the HTTP server for remote diagnostics.
- **Flags:** `--port, -p` (Default: 8080).

### `install`
(Linux only) Installs AuraSpeed as a `systemd` service.

written by Neorwc