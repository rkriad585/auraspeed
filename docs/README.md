# AuraSpeed Documentation

Welcome to the AuraSpeed documentation hub. This folder contains comprehensive guides for using, configuring, and contributing to AuraSpeed.

## 📚 Documentation Index

### Getting Started

| Document | Description |
|----------|-------------|
| [Getting Started](getting-started.md) | Beginner-friendly setup guide (installation, first run) |
| [FAQ](faq.md) | Frequently asked questions |

### Command Reference

| Document | Description |
|----------|-------------|
| [API Reference](api-reference.md) | CLI commands and Go functions reference |
| [Web Command](web-command.md) | Documentation for `auraspeed web` HTTP server |
| [Web API](api-web.md) | HTTP endpoints for web server (`/api/speedtest`, etc.) |

### Configuration & Usage

| Document | Description |
|----------|-------------|
| [Configuration Reference](configuration.md) | Complete TOML config options explained |
| [Examples](examples.md) | Real-world usage examples for all commands |

### System & Architecture

| Document | Description |
|----------|-------------|
| [Architecture](architecture.md) | System design, data flow, and dependencies |

### Support

| Document | Description |
|----------|-------------|
| [Troubleshooting](troubleshooting.md) | Common issues and solutions |

### Community

| Document | Description |
|----------|-------------|
| [Contributing](../CONTRIBUTING.md) | Guidelines for contributors |
| [Code of Conduct](../CODE_OF_CONDUCT.md) | Contributor Covenant v2.1 |
| [Changelog](../CHANGELOG.md) | Version history (Keep a Changelog format) |

---

## 🚀 Quick Links

### For Users

1. [Install AuraSpeed](getting-started.md#installation) (from source or pre-built binaries)
2. [Run your first speed test](getting-started.md#first-run) in under 5 minutes
3. [Configure settings](configuration.md) via TOML or CLI commands
4. [Explore examples](examples.md) for common use cases

### For Developers

1. [Set up development environment](getting-started.md#development-setup)
2. [Understand the architecture](architecture.md)
3. [Read contributing guidelines](../CONTRIBUTING.md)
4. [Review API reference](api-reference.md)

### For DevOps / CI/CD

1. [Build for all platforms](examples.md#build-from-source) with `build.sh` or `build.ps1`
2. [Set up web server](web-command.md) for remote access
3. [Use GoReleaser](getting-started.md#cicd-with-goreleaser) for automated releases

---

## 📦 Project Structure

```
auraspeed/
├── cmd/
│   ├── main.go              # Entry point (package main)
│   └── root/
│       ├── root.go          # Root cobra command & CLI setup
│       ├── commands.go      # All CLI subcommands
│       ├── web.go          # Web server command (NEW)
│       └── web.html        # HTML template for web UI
├── internal/
│   ├── config/           # Configuration management (viper)
│   ├── info/             # System information (gopsutil)
│   ├── logging/          # Logging utilities (zerolog)
│   ├── network/          # Network diagnostics
│   ├── speedtest/        # Speed test & TUI (speedtest-go, tview)
│   └── ui/               # UI command wrapper
├── docs/                   # ← YOU ARE HERE
│   ├── README.md          # This file (documentation index)
│   ├── getting-started.md # Beginner guide
│   ├── api-reference.md   # CLI & Go API reference
│   ├── web-command.md     # Web server documentation
│   ├── api-web.md         # Web API endpoints
│   ├── configuration.md    # Complete config reference
│   ├── examples.md         # Real-world usage examples
│   ├── architecture.md     # System design & data flow
│   ├── troubleshooting.md   # Common issues & fixes
│   ├── faq.md             # Frequently asked questions
│   └── ...
├── build.sh                 # Unix build script (bash/zsh/fish)
├── build.ps1                # Windows build script (PowerShell)
├── .goreleaser.yaml         # GoReleaser config for releases
├── README.md                # Project homepage
├── CONTRIBUTING.md          # Contributor guidelines
├── CODE_OF_CONDUCT.md        # Code of Conduct
├── CHANGELOG.md             # Version history
└── go.mod                   # Go module (auraspeed)
```

---

## 🔗 External Links

| Resource | Link |
|----------|------|
| GitHub Repository | https://github.com/rkriad585/auraspeed |
| Releases | https://github.com/rkriad585/auraspeed/releases |
| Issue Tracker | https://github.com/rkriad585/auraspeed/issues |
| GoReleaser | https://goreleaser.com/ |
| Cobra CLI | https://github.com/spf13/cobra |
| tview (TUI) | https://github.com/rivo/tview |
| speedtest-go | https://github.com/showwin/speedtest-go |

---

## 📝 Documentation Stats

| Metric | Count |
|--------|-------|
| Total Documents | 10 |
| Getting Started Guides | 1 |
| API References | 2 |
| Configuration Guides | 1 |
| Example Collections | 1 |
| Architecture Docs | 1 |
| Troubleshooting Guides | 1 |
| FAQ Documents | 1 |
| Community Docs (Contributing, CoC) | 2 |

---

## 🤝 Contributing to Documentation

Found a typo? Want to improve a guide?

1. Fork the repository: https://github.com/rkriad585/auraspeed/fork
2. Create a branch: `git checkout -b docs/update-guide`
3. Edit files in the `docs/` folder
4. Commit: `git commit -m "docs: improve getting-started guide"`
5. Push: `git push origin docs/update-guide`
6. Open a Pull Request!

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

---

## 📄 License

AuraSpeed is open-source software licensed under the [MIT License](../LICENSE).

---

**GitHub:** https://github.com/rkriad585/auraspeed
**Built with ❤️ using Go, Cobra, tview, and speedtest-go**
