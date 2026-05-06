# Contributing to AuraSpeed

Thank you for your interest in contributing to AuraSpeed!

## How to Contribute

### Reporting Bugs
- Use the GitHub issue tracker
- Include steps to reproduce
- Include your OS, Go version, and AuraSpeed version

### Suggesting Features
- Open a GitHub issue with the "enhancement" label
- Describe the feature and its use case

### Pull Requests
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup
```bash
git clone https://github.com/rkriad585/auraspeed.git
cd auraspeed
go mod download
make build
```

### Coding Standards
- Run `go fmt ./...` before committing
- Run `golangci-lint run ./...` to check code quality
- Write tests for new functionality
- Update documentation as needed

### Testing
```bash
make test
# or
go test ./... -v -cover
```
