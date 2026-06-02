# Contributing to AuraSpeed

Thank you for your interest in contributing to AuraSpeed! This guide will help you get started with the development workflow and coding standards.

## 📋 Before You Start

- Read the [Code of Conduct](CODE_OF_CONDUCT.md)
- Check [existing issues](https://github.com/rkriad585/auraspeed/issues) before opening a new one
- For big changes — open an issue first to discuss

## 🛠️ Development Setup

```bash
# 1. Fork and clone
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed

# 2. Install dependencies
go mod download

# 3. Build the project
go build -o auraspeed ./cmd/main.go

# Or use build scripts
./build.sh          # Unix (bash/zsh/fish)
.\build.ps1          # Windows (PowerShell)

# 4. Run tests to verify setup
go test ./... -v

# 5. Start development
# TUI mode: auraspeed tui
# CLI mode: auraspeed speedtest
```

## 🌿 Branch Naming

Use descriptive branch names following these patterns:

| Type      | Pattern                  | Example                    |
|-----------|--------------------------|----------------------------|
| Feature   | `feat/short-description` | `feat/add-dark-mode`       |
| Bug fix   | `fix/short-description`  | `fix/speedtest-panic`      |
| Docs      | `docs/short-description` | `docs/api-reference`       |
| Refactor  | `refactor/description`   | `refactor/tui-cleanup`      |

## 📝 Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/) for clear, automated changelog generation:

| Type | Description |
|------|-------------|
| `feat:` | New features |
| `fix:` | Bug fixes |
| `docs:` | Documentation changes |
| `refactor:` | Code refactoring (no behavior change) |
| `test:` | Adding or updating tests |
| `chore:` | Build/tooling changes |

### Examples

```bash
feat: add dark mode toggle to TUI
fix: resolve panic when speedtest server is unreachable
docs: update API reference for config package
refactor: simplify TUI rendering logic
```

## 🔄 Pull Request Process

1. Create a branch from `main` using the [branch naming convention](#branch-naming)
2. Make your changes with tests
3. Run the test suite: `go test ./... -v -cover`
4. Run static analysis: `go vet ./...`
5. Run formatting: `go fmt ./...`
6. Update documentation if behavior changes
7. Open a PR with a clear description

> [!IMPORTANT]
> All PRs must pass CI checks and include appropriate test coverage before merging.

## ✅ PR Checklist

- [ ] Tests added for new functionality
- [ ] All existing tests pass
- [ ] Docs updated if behavior changed
- [ ] Commit messages follow convention
- [ ] No merge conflicts with main

## 🧪 Running Tests

```bash
# Verbose output
go test ./... -v

# With coverage report
go test ./... -cover

# Run specific test
go test ./internal/speedtest/... -v
```

> [!TIP]
> Aim for meaningful test coverage on new code. Tests go in `_test.go` files next to the source.

## Coding Standards

- Run `go fmt ./...` before committing
- Follow [Go best practices](https://golang.org/doc/effective_go.html)
- Document all exported functions using [godoc](https://godoc.org) style
- Keep functions small and focused on a single responsibility
- Use meaningful variable and function names
- Avoid global state when possible

### Example Godoc Comment

```go
// CalculateSpeed computes the download speed in Mbps.
// It takes the total bytes downloaded and the duration in seconds.
func CalculateSpeed(bytes int64, duration float64) float64 {
    // ...
}
```

## 🏗️ Project Structure

AuraSpeed follows standard Go project layout conventions:

| Directory | Purpose |
|-----------|---------|
| `cmd/` | Application entry points (main package, CLI setup) |
| `internal/` | Private application code (not importable by external projects) |
| `docs/` | Project documentation |
| `.github/` | GitHub workflows and community files |

> [!TIP]
> Code in `internal/` is enforced by the Go compiler — external projects cannot import it. This keeps the public API clean.

## 📦 Build Commands

| Command | Description |
|---------|-------------|
| `go build -o auraspeed ./cmd/main.go` | Quick local build |
| `make build` | Build using Makefile |
| `make test` | Run all tests |
| `make lint` | Run linters |
| `make format` | Format source code |
| `make release` | Build for all platforms |
| `make clean` | Remove build artifacts |
| `cmake -B build && cmake --build build` | Build via CMake |
| `.\build.ps1` | Windows: build for all platforms |
| `./build.sh` | Unix: build for all platforms |

## 🆘 Getting Help

- Open a [GitHub issue](https://github.com/rkriad585/auraspeed/issues) for bugs or feature requests
- Check existing issues before creating a new one
- Include your OS, Go version, and AuraSpeed version when reporting bugs

Thank you for contributing to AuraSpeed! 🚀
