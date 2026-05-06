package root

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	st "github.com/showwin/speedtest-go/speedtest"

	"auraspeed/internal/info"

	"github.com/spf13/cobra"
)

// NewWebCommand returns the web subcommand.
// It starts a web server to expose AuraSpeed functionality via HTTP.
func NewWebCommand() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start AuraSpeed web server",
		Long:  "Launch a web server that provides HTTP API access to AuraSpeed features.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := startWebServer(port); err != nil {
				return fmt.Errorf("web server failed: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	return cmd
}

func startWebServer(port int) error {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok", "service": "auraspeed"}`)
	})

	// Speed test endpoint
	mux.HandleFunc("/api/speedtest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		client := st.New()
		user, err := client.FetchUserInfo()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch user info: %v", err), http.StatusInternalServerError)
			return
		}

		serverList, err := client.FetchServers()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch servers: %v", err), http.StatusInternalServerError)
			return
		}

		targets, err := serverList.FindServer([]int{})
		if err != nil || len(targets) == 0 {
			http.Error(w, "No servers found", http.StatusInternalServerError)
			return
		}

		s := targets[0]
		if s.Context == nil {
			s.Context = st.New()
		}

		if err := s.DownloadTest(); err != nil {
			http.Error(w, fmt.Sprintf("Download test failed: %v", err), http.StatusInternalServerError)
			return
		}

		if err := s.UploadTest(); err != nil {
			http.Error(w, fmt.Sprintf("Upload test failed: %v", err), http.StatusInternalServerError)
			return
		}

		dlSpeed := float64(s.DLSpeed) / 1000000.0
		ulSpeed := float64(s.ULSpeed) / 1000000.0
		ping := s.Latency.Milliseconds()

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
	fmt.Printf("Starting AuraSpeed web server on http://localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop")
	return http.ListenAndServe(addr, mux)
}
