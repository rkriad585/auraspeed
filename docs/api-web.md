# AuraSpeed API Documentation

Complete API reference for AuraSpeed HTTP endpoints (available via `auraspeed web` command).

## Base URL

When running the web server:
```
http://localhost:8080
```

Custom port:
```bash
auraspeed web --port 3000
# Base URL: http://localhost:3000
```

---

## Health Check

### `GET /health`

Check if the server is running.

**Response:**
```json
{"status": "ok", "service": "auraspeed"}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

## Speed Test

### `POST /api/speedtest`

Run a network speed test (download + upload + latency).

**Method:** `POST`

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

**Example:**
```bash
curl -X POST http://localhost:8080/api/speedtest
```

**Error Responses:**

| Code | Message |
|------|---------|
| 405 | Method not allowed (use POST) |
| 500 | Failed to fetch user info |
| 500 | Failed to fetch servers |
| 500 | No servers found |
| 500 | Download test failed |
| 500 | Upload test failed |

---

## System Information

### `GET /api/info`

Get system information (OS, CPU, memory, disk).

**Method:** `GET`

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

**Example:**
```bash
curl http://localhost:8080/api/info
```

**Error Responses:**

| Code | Message |
|------|---------|
| 500 | Failed to get system info |

---

## Web UI

### `GET /`

Simple HTML interface with buttons to run tests.

**Method:** `GET`

**Returns:** HTML page with:
- "Run Speed Test" button → calls `POST /api/speedtest`
- "Get System Info" button → calls `GET /api/info`

**Example:**
```bash
# Open in browser
start http://localhost:8080
# Windows
open http://localhost:8080
# macOS
xdg-open http://localhost:8080
# Linux
```

---

## Authentication

Currently, the web server does not implement authentication. It is intended for:
- Local development
- Internal network use
- Testing purposes

**⚠️ Warning:** Do not expose the web server to the public internet without adding authentication.

---

## Rate Limiting

Currently, no rate limiting is implemented. Be mindful of:
- Speed tests use bandwidth
- Frequent tests may slow down your connection

---

## JavaScript Examples

### Vanilla JS (Browser)

```javascript
// Run speed test
async function runSpeedTest() {
    const response = await fetch('/api/speedtest', { method: 'POST' });
    const data = await response.json();
    console.log('Download:', data.download, 'Mbps');
    console.log('Upload:', data.upload, 'Mbps');
    console.log('Ping:', data.ping, 'ms');
}

// Get system info
async function getSystemInfo() {
    const response = await fetch('/api/info');
    const data = await response.json();
    console.log('OS:', data.os);
    console.log('CPU:', data.cpu);
}
```

### Node.js (Server-side)

```javascript
const axios = require('axios');

// Speed test
async function runSpeedTest() {
    try {
        const response = await axios.post('http://localhost:8080/api/speedtest');
        console.log(response.data);
    } catch (error) {
        console.error('Error:', error.message);
    }
}

// System info
async function getSystemInfo() {
    try {
        const response = await axios.get('http://localhost:8080/api/info');
        console.log(response.data);
    } catch (error) {
        console.error('Error:', error.message);
    }
}
```

### Python

```python
import requests

# Speed test
def run_speed_test():
    try:
        response = requests.post('http://localhost:8080/api/speedtest')
        data = response.json()
        print(f"Download: {data['download']} Mbps")
        print(f"Upload: {data['upload']} Mbps")
        print(f"Ping: {data['ping']} ms")
    except Exception as e:
        print(f"Error: {e}")

# System info
def get_system_info():
    try:
        response = requests.get('http://localhost:8080/api/info')
        data = response.json()
        print(f"OS: {data['os']}")
        print(f"CPU: {data['cpu']}")
    except Exception as e:
        print(f"Error: {e}")
```

---

## Error Handling

All endpoints return JSON error messages:

```json
{"error": "Failed to fetch user info: connection timeout"}
```

**Common HTTP Status Codes:**

| Code | Meaning |
|------|---------|
| 200 | Success |
| 405 | Method not allowed |
| 500 | Internal server error |

---

## Running the Web Server

### Start the Server

```bash
# Default port 8080
auraspeed web

# Custom port
auraspeed web --port 3000
```

### Output:
```
Starting AuraSpeed web server on http://localhost:8080
Press Ctrl+C to stop
```

### Stop the Server

Press `Ctrl+C` in the terminal.

---

## Links

- [GitHub Repository](https://github.com/rkriad585/auraspeed)
- [Web Command Documentation](web-command.md)
- [Getting Started Guide](getting-started.md)
- [API Reference](api-reference.md)
