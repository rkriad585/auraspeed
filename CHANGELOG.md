# Changelog

All notable changes to AuraSpeed are documented here.

## [Unreleased]

### Added
-

### Changed
-

### Fixed
-

## [1.0.1] - 2026-05-06

### Fixed
- Resolve TUI panic and display issues (nil pointer in speedtest-go)
- Add missing main path in GoReleaser config
- Update goreleaser config to version 2 syntax
- Update goreleaser config to version 1 syntax
- Correct GitHub Actions release condition syntax
- Add go.sum to fix GitHub Actions
- Update GitHub username and .gitignore

### Added
- Unix build script (build.sh) for bash/zsh/fish
- Windows build script (build.ps1) for PowerShell
- Ignore dist/ folder in .gitignore

### Documentation
- Update CONTRIBUTING.md with correct GitHub username
- Add comprehensive documentation in docs/ folder

## [1.0.0] - 2026-05-06

### Added
- Interactive TUI mode with real-time throughput graphs
- Speed test with download/upload/latency measurements
- System information display (CPU, memory, disk)
- Network diagnostics (ping, traceroute, DNS)
- Test history tracking with JSON/TSV export
- TOML-based configuration with viper
- Clipboard integration for copying results
- Cross-platform support (Windows, Linux, macOS Intel & Apple Silicon)
- Command aliases (st, si, net, hist)
- GoReleaser configuration for automated releases
- GitHub Actions CI/CD pipeline

### Built With
- cobra (CLI framework)
- tview (terminal UI)
- speedtest-go (speed test library)
- gopsutil (system information)
- zerolog (logging)
