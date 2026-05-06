# Contributing to AuraSpeed

Thank you for your interest in contributing to AuraSpeed! This guide will help you get started with the development workflow and coding standards.

## Development Setup

Clone the repository and install dependencies:

```bash
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed
go mod download
```

### Build

```bash
# Quick build
go build -o auraspeed ./cmd/main.go

# Or use build scripts
./build.sh          # Unix
.\build.ps1          # Windows
```

## Project Structure

AuraSpeed follows standard Go project layout conventions:

| Directory | Purpose |
|-----------|---------|
| `cmd/` | Application entry points (main package, CLI setup) |
| `internal/` | Private application code (not importable by external projects) |
| `docs/` | Project documentation |
| `.github/` | GitHub workflows and community files |

> [!TIP]
> Code in `internal/` is enforced by the Go compiler — external projects cannot import it. This keeps the public API clean.

## Branch Naming

Use descriptive branch names following these patterns:

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feat/description` | `feat/add-dark-mode` |
| Bug fix | `fix/description` | `fix/speedtest-panic` |
| Docs | `docs/description` | `docs/api-reference` |
| Refactor | `refactor/description` | `refactor/tui-cleanup` |

## Commit Messages

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

## Pull Request Process

1. Create a branch from `main` using the [branch naming convention](#branch-naming)
2. Add tests for any new functionality
3. Run the test suite: `go test ./... -v -cover`
4. Run static analysis: `go vet ./...`
5. Run formatting: `go fmt ./...`
6. Update documentation if behavior changes
7. Open a PR with a clear description of the changes

> [!IMPORTANT]
> All PRs must pass CI checks and include appropriate test coverage before merging.

## Coding Standards

- Run `go fmt ./...` before committing
- Follow [Go best practices](https://golang.org/doc/effective_go.html)
- Document all exported functions using [godoc](https://godoc.org) style
- Keep functions small and focused on a single responsibility
- Use meaningful variable and function names

### Example Godoc Comment

```go
// CalculateSpeed computes the download speed in Mbps.
// It takes the total bytes downloaded and the duration in seconds.
func CalculateSpeed(bytes int64, duration float64) float64 {
    // ...
}
```

## Testing

Run the full test suite with:

```bash
# Verbose output
go test ./... -v

# With coverage report
go test ./... -cover
```

> [!TIP]
> Aim for meaningful test coverage on new code. Tests go in `_test.go` files next to the source.

## Build Commands

| Command | Description |
|---------|-------------|
| `go build -o auraspeed ./cmd/main.go` | Quick local build |
| `.\build.ps1` | Windows: build for all platforms |
| `./build.sh` | Unix: build for all platforms |
| `goreleaser build --snapshot --clean` | GoReleaser snapshot build |

## Getting Help

- Open a [GitHub issue](https://github.com/rkriad585/auraspeed/issues) for bugs or feature requests
- Check existing issues before creating a new one
- Include your OS, Go version, and AuraSpeed version when reporting bugs

Thank you for contributing to AuraSpeed! 🚀

Thank you for contributing to AuraSpeed! 🚀
