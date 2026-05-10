package root

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	st "github.com/showwin/speedtest-go/speedtest"

	"auraspeed/internal/info"

	"github.com/spf13/cobra"
)

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	mu       sync.Mutex
	tokens   map[string]int
	requests map[string]int
	lastReset time.Time
	rate      int
	burst     int
}

func newRateLimiter(rate, burst int) *rateLimiter {
	return &rateLimiter{
		tokens:   make(map[string]int),
		requests: make(map[string]int),
		lastReset: time.Now(),
		rate:     rate,
		burst:    burst,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastReset) > time.Minute {
		rl.lastReset = now
		rl.requests = make(map[string]int)
	}

	current := rl.requests[key]
	if current >= rl.rate {
		return false
	}
	rl.requests[key] = current + 1
	return true
}

var (
	rl = newRateLimiter(60, 10) // 60 requests per minute, burst of 10
	apiKey = ""
)

// Metrics for monitoring
var (
	metricsRequestsTotal    int64
	metricsSpeedTestsTotal  int64
	metricsSpeedTestsFailed int64
	metricsLastSpeedTest    time.Time
	metricsLastDownload     float64
	metricsLastUpload       float64
	metricsLastPing         int64
	metricsMu              sync.RWMutex
)

// NewWebCommand returns the web subcommand.
// It starts a web server to expose AuraSpeed functionality via HTTP.
func NewWebCommand() *cobra.Command {
	var port int
	var authKey string
	var enableTLS bool
	var certFile string
	var keyFile string

	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start AuraSpeed web server",
		Long:  "Launch a web server that provides HTTP API access to AuraSpeed features.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if authKey != "" {
				apiKey = authKey
			}
			if err := startWebServer(port, enableTLS, certFile, keyFile); err != nil {
				return fmt.Errorf("web server failed: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 59733, "Port to listen on")
	cmd.Flags().StringVar(&authKey, "auth", "", "API key for authentication")
	cmd.Flags().BoolVar(&enableTLS, "tls", false, "Enable TLS/HTTPS")
	cmd.Flags().StringVar(&certFile, "cert", "", "TLS certificate file path")
	cmd.Flags().StringVar(&keyFile, "key", "", "TLS key file path")

	return cmd
}

func startWebServer(port int, enableTLS bool, certFile, keyFile string) error {
	mux := http.NewServeMux()

	// Middleware for auth and rate limiting
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rate limiting
		clientIP := r.RemoteAddr
		if !rl.allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Auth check for API endpoints
		if strings.HasPrefix(r.URL.Path, "/api/") && apiKey != "" {
			providedKey := r.Header.Get("X-API-Key")
			if providedKey != apiKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		mux.ServeHTTP(w, r)
	})

	// Health check (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Deep health check - verify we can access basic info
		_, err := info.GetSystemInfo()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status": "degraded", "service": "auraspeed", "error": "%v"}`, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok", "service": "auraspeed", "version": "%s"}`, Version)
	})

	// Metrics endpoint (Prometheus format)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		metricsMu.RLock()
		defer metricsMu.RUnlock()

		fmt.Fprintf(w, "# HELP auraspeed_requests_total Total number of HTTP requests\n")
		fmt.Fprintf(w, "# TYPE auraspeed_requests_total counter\n")
		fmt.Fprintf(w, "auraspeed_requests_total %d\n", metricsRequestsTotal)

		fmt.Fprintf(w, "# HELP auraspeed_speedtests_total Total number of speed tests performed\n")
		fmt.Fprintf(w, "# TYPE auraspeed_speedtests_total counter\n")
		fmt.Fprintf(w, "auraspeed_speedtests_total %d\n", metricsSpeedTestsTotal)

		fmt.Fprintf(w, "# HELP auraspeed_speedtests_failed Total number of failed speed tests\n")
		fmt.Fprintf(w, "# TYPE auraspeed_speedtests_failed counter\n")
		fmt.Fprintf(w, "auraspeed_speedtests_failed %d\n", metricsSpeedTestsFailed)

		fmt.Fprintf(w, "# HELP auraspeed_last_download Last download speed in Mbps\n")
		fmt.Fprintf(w, "# TYPE auraspeed_last_download gauge\n")
		fmt.Fprintf(w, "auraspeed_last_download %.2f\n", metricsLastDownload)

		fmt.Fprintf(w, "# HELP auraspeed_last_upload Last upload speed in Mbps\n")
		fmt.Fprintf(w, "# TYPE auraspeed_last_upload gauge\n")
		fmt.Fprintf(w, "auraspeed_last_upload %.2f\n", metricsLastUpload)

		fmt.Fprintf(w, "# HELP auraspeed_last_ping Last ping in milliseconds\n")
		fmt.Fprintf(w, "# TYPE auraspeed_last_ping gauge\n")
		fmt.Fprintf(w, "auraspeed_last_ping %d\n", metricsLastPing)

		fmt.Fprintf(w, "# HELP auraspeed_last_test_timestamp Unix timestamp of last speed test\n")
		fmt.Fprintf(w, "# TYPE auraspeed_last_test_timestamp gauge\n")
		if !metricsLastSpeedTest.IsZero() {
			fmt.Fprintf(w, "auraspeed_last_test_timestamp %d\n", metricsLastSpeedTest.Unix())
		}
	})

	// Speed test endpoint
	mux.HandleFunc("/api/speedtest", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&metricsRequestsTotal, 1)

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		client := st.New()

		user, err := client.FetchUserInfo()
		if err != nil {
			atomic.AddInt64(&metricsSpeedTestsFailed, 1)
			http.Error(w, fmt.Sprintf("Failed to fetch user info: %v", err), http.StatusInternalServerError)
			return
		}

		serverList, err := client.FetchServers()
		if err != nil {
			atomic.AddInt64(&metricsSpeedTestsFailed, 1)
			http.Error(w, fmt.Sprintf("Failed to fetch servers: %v", err), http.StatusInternalServerError)
			return
		}

		targets, err := serverList.FindServer([]int{})
		if err != nil || len(targets) == 0 {
			atomic.AddInt64(&metricsSpeedTestsFailed, 1)
			http.Error(w, "No servers found", http.StatusInternalServerError)
			return
		}

		s := targets[0]
		if s.Context == nil {
			s.Context = st.New()
		}

		done := make(chan error, 2)
		go func() {
			done <- s.DownloadTest()
		}()
		go func() {
			done <- s.UploadTest()
		}()

		for i := 0; i < 2; i++ {
			select {
			case <-ctx.Done():
				atomic.AddInt64(&metricsSpeedTestsFailed, 1)
				http.Error(w, "Speed test timed out", http.StatusGatewayTimeout)
				return
			case err := <-done:
				if err != nil {
					atomic.AddInt64(&metricsSpeedTestsFailed, 1)
					http.Error(w, fmt.Sprintf("Test failed: %v", err), http.StatusInternalServerError)
					return
				}
			}
		}

		dlSpeed := float64(s.DLSpeed) / 1000000.0
		ulSpeed := float64(s.ULSpeed) / 1000000.0
		ping := s.Latency.Milliseconds()

		// Update metrics
		atomic.AddInt64(&metricsSpeedTestsTotal, 1)
		metricsMu.Lock()
		metricsLastDownload = dlSpeed
		metricsLastUpload = ulSpeed
		metricsLastPing = ping
		metricsLastSpeedTest = time.Now()
		metricsMu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"download": %.2f, "upload": %.2f, "ping": %d, "isp": "%s", "server": "%s"}`,
			dlSpeed, ulSpeed, ping, user.Isp, s.Name)
	})

	// System info endpoint
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		sysInfo, err := info.GetSystemInfo()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get system info: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, _ := json.Marshal(sysInfo)
		w.Write(jsonData)
	})

	// Simple HTML UI - serve from embedded file
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		htmlFile, err := os.ReadFile("cmd/root/web.html")
		if err != nil {
			http.Error(w, "HTML template not found", http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(htmlFile))
	})

	addr := ":" + strconv.Itoa(port)
	scheme := "http"
	if enableTLS {
		scheme = "https"
	}
	fmt.Printf("Starting AuraSpeed web server on %s://localhost:%d\n", scheme, port)
	if apiKey != "" {
		fmt.Println("API authentication is enabled")
	}
	if enableTLS {
		fmt.Printf("TLS enabled with cert: %s, key: %s\n", certFile, keyFile)
	}
	fmt.Println("Press Ctrl+C to stop")

	server := &http.Server{Addr: addr, Handler: handler}

	// Graceful shutdown
	go func() {
		<-make(chan os.Signal, 1)
		fmt.Println("\nShutting down web server...")
		server.Shutdown(context.Background())
	}()

	if enableTLS {
		return server.ListenAndServeTLS(certFile, keyFile)
	}
	return server.ListenAndServe()
}
