# AuraSpeed Web Command

Documentation for the `auraspeed web` command.

## Overview

The `web` command starts an HTTP server that exposes AuraSpeed functionality via RESTful API endpoints. It also serves a simple HTML UI for browser-based access.

## Usage

### Basic Start

```bash
auraspeed web
```

Default port: `59733`
Base URL: `http://localhost:59733`

### Custom Port

```bash
auraspeed web --port 3000
# or
auraspeed web -p 3000
```

### Help

```bash
auraspeed web --help
```

---

## Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | HTML UI with test buttons |
| `/health` | GET | Health check |
| `/api/speedtest` | POST | Run speed test |
| `/api/info` | GET | Get system information |

---

## API Examples

### Health Check

```bash
curl http://localhost:59733/health
# Output: {"status": "ok", "service": "auraspeed"}
```

### Speed Test

```bash
curl -X POST http://localhost:59733/api/speedtest
```

**Response:**
```json
{
  "download": 85.42,
  "upload": 42.18,
  "ping": 12,
  "isp": "Comcast",
  "server": "New York, NY"
}
```

### System Information

```bash
curl http://localhost:59733/api/info
```

**Response:**
```json
{
  "os": "Windows 10 22H2",
  "cpu": "Intel(R) Core(TM) i7-10700K (16 cores)",
  "memory_total": "32 GB",
  "memory_used": "16 GB",
  "disk_total": "512 GB",
  "disk_used": "256 GB",
  "hostname": "DESKTOP-ABC123"
}
```

---

## Web UI

Open in browser: `http://localhost:59733`

The UI provides:
- **Run Speed Test** button — triggers speed test
- **Get System Info** button — fetches system information
- Results displayed directly in the browser

---

## Use Cases

### 1. Local Development

```bash
# Start web server
auraspeed web &

# Run tests from command line
curl -X POST http://localhost:59733/api/speedtest
```

### 2. Remote Monitoring

```bash
# On remote machine
auraspeed web --port 59733

# From local machine
curl http://remote-machine:59733/api/info
```

### 3. Integration Testing

```bash
# Start server
auraspeed web &
sleep 2

# Run test
response=$(curl -s -X POST http://localhost:59733/api/speedtest)
echo $response | jq '.download'
```

### 4. Automation Script

```bash
#!/bin/bash
# Periodic speed test logger

while true; do
    echo "=== $(date) ==="
    curl -s -X POST http://localhost:59733/api/speedtest | jq '.'
    sleep 3600  # Every hour
done
```

---

## Configuration

The web command uses the same configuration as the CLI:

| Config Key | Default | Description |
|------------|---------|-------------|
| `global.loglevel` | `"info"` | Log level for web server |
| `speedtest.timeout` | `30` | Speed test timeout |
| `speedtest.defaultserverid` | `0` | Default server (0 = auto) |

Config file location:
- Windows: `%USERPROFILE%\.config\neostore\auraspeed\config.toml`
- Linux/macOS: `~/.config/neostore/auraspeed/config.toml`

---

## Security Considerations

**⚠️ Important:**

1. **No Authentication** — The web server does not implement auth
2. **Local Use Only** — Designed for localhost/internal network
3. **No HTTPS** — Uses plain HTTP
4. **No Rate Limiting** — Be careful with frequent requests

**Do NOT expose to public internet without:**
- Adding authentication middleware
- Enabling HTTPS/TLS
- Implementing rate limiting
- Adding CORS restrictions

---

## Troubleshooting

### Port Already in Use

**Error:**
```
listen tcp :59733: bind: Only one usage of each socket address is normally permitted.
```

**Fix:**
```bash
# Use a different port
auraspeed web --port 8081
```

### Connection Refused

**Cause:** Web server not running.

**Fix:**
```bash
# Check if running
ps aux | grep auraspeed

# Start server
auraspeed web
```

### Speed Test Fails

**Cause:** Network issue or speedtest-go library error.

**Fix:**
```bash
# Check logs
auraspeed --log-level debug web

# Try with different server
# (Edit config: speedtest.defaultserverid = 1234)
```

---

## Development

### Code Location

- Command definition: `cmd/root/web.go`
- HTML template: `cmd/root/web.html`

### Running in Debug Mode

```bash
auraspeed --log-level debug web
```

### Testing the API

```bash
# Health check
curl http://localhost:59733/health

# Speed test
curl -X POST http://localhost:59733/api/speedtest | jq '.'

# System info
curl http://localhost:59733/api/info | jq '.'
```

---

## Links

- [API Documentation](api-web.md)
- [Getting Started Guide](getting-started.md)
- [Main README](../README.md)
- [GitHub Repository](https://github.com/rkriad585/auraspeed)
