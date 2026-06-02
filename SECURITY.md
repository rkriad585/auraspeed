# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 3.x     | :white_check_mark: |
| < 3.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in AuraSpeed, please follow these steps:

1. **Do not** open a public GitHub issue.
2. Email your findings to the maintainers at the address listed in [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
3. Include a detailed description, steps to reproduce, and any relevant code or configuration.

You can expect:
- **Acknowledgment** within 48 hours of your report.
- **Update** on the status of your report within 5 business days.
- **Fix timeline** communicated based on severity.

## Security Best Practices

When using AuraSpeed:

- **Config file permissions**: Ensure your config file at `~/.config/neostore/auraspeed/config.toml` is readable only by you. AuraSpeed sets `0600` permissions automatically.
- **Web server**: When using `auraspeed web`, bind to `127.0.0.1` (default) unless you need network access. Use a firewall to restrict access to the web port (`59733` by default).
- **Auto-update**: The auto-update feature fetches version info from GitHub over HTTPS. Binaries should be verified against checksums provided in release notes.
- **Third-party dependencies**: AuraSpeed relies on several open-source Go libraries. We regularly update dependencies to incorporate security patches.

## Known Security Considerations

- AuraSpeed's web server does not implement authentication by default. If exposed to a network, use a reverse proxy (e.g., nginx, Caddy) with authentication.
- Speed test data is stored locally in `history.json`. No data is sent to third parties other than the Speedtest.net servers used for speed measurement.
- CLI arguments may be visible in process listings on multi-user systems. Avoid passing sensitive information as arguments.
