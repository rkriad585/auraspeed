# Web Server Mode

AuraSpeed includes a built-in web server that provides a browser-based interface and a JSON API for remote monitoring.

## Starting the Server

Run the following command to start the server:
```bash
auraspeed web --port 8080
```

## Web Interface

Once started, navigate to `http://localhost:8080`. The web interface provides:
- **Live Speed Test:** Run tests directly from the browser.
- **System Stats:** View CPU, RAM, and Disk usage of the host machine.
- **Dark/Light Mode:** Toggle themes for comfortable viewing.
- **Local History:** The web UI maintains its own history in the browser's local storage.

## API Endpoints

The web server also acts as a REST API:

### `GET /api/info`
Returns system information.
```json
{
  "os": "linux 5.15.0",
  "cpu": "AMD Ryzen 7 (16 cores)",
  "memory": "32 GB total, 8 GB used (25%)",
  "disk": "500 GB total, 120 GB used (24%)",
  "hostname": "node-01"
}
```

### `GET /api/speedtest`
Triggers a speed test and returns the results.
*Note: This is a blocking call and may take 10-20 seconds.*

### `GET /health`
Simple health check endpoint for monitoring tools or Docker health checks.

## Security Note
The web server is intended for local network use. If exposing to the internet, it is recommended to use a reverse proxy (like Nginx) with Authentication and TLS.

written by Neorwc