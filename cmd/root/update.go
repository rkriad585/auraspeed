package root

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"auraspeed/internal/config"
	"auraspeed/internal/logging"

	"github.com/spf13/cobra"
)

// checkForUpdate fetches the latest version from the GitHub .version file
func checkForUpdate() (string, error) {
	versionURL := "https://raw.githubusercontent.com/rkriad585/auraspeed/main/.version"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, versionURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "AuraSpeed-UpdateChecker")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("update check returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	version := strings.TrimSpace(string(body))
	if version == "" {
		return "", fmt.Errorf("empty version from .version file")
	}

	return version, nil
}

// compareVersions compares two version strings (v1.0.0 format)
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Remove 'v' prefix if present
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		p1 := 0
		p2 := 0

		if i < len(v1Parts) {
			fmt.Sscanf(v1Parts[i], "%d", &p1)
		}
		if i < len(v2Parts) {
			fmt.Sscanf(v2Parts[i], "%d", &p2)
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}
	return 0
}

// downloadUpdate downloads the binary for the given version and installs it
func downloadUpdate(version string) error {
	targetDir := installDir()
	binName := "auraspeed"
	if runtime.GOOS == "windows" {
		binName = "auraspeed.exe"
	}
	targetPath := filepath.Join(targetDir, binName)

	repo := "rkriad585/auraspeed"
	var downloadName string

	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "arm64":
			downloadName = "auraspeed-windows-arm64.exe"
		default:
			downloadName = "auraspeed-windows-amd64.exe"
		}
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			downloadName = "auraspeed-darwin-arm64"
		default:
			downloadName = "auraspeed-darwin-amd64"
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			downloadName = "auraspeed-linux-arm64"
		default:
			downloadName = "auraspeed-linux-amd64"
		}
	default:
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, version, downloadName)

	fmt.Printf(">>> Downloading %s\n", url)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		os.Remove(targetPath)
		return fmt.Errorf("failed to write file: %w", err)
	}
	out.Close()

	if written == 0 {
		os.Remove(targetPath)
		return fmt.Errorf("downloaded empty file")
	}

	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0755)
	}

	fmt.Printf("OK   Installed to %s (%d bytes)\n", targetPath, written)
	return nil
}

// NewUpdateCommand returns the update subcommand
func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update AuraSpeed to the latest version",
		Long:  "Check for and download the latest version of AuraSpeed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentVersion := Version
			if currentVersion == "dev" {
				fmt.Println("Running from source - cannot update.")
				return nil
			}

			fmt.Println(">>> Checking for updates...")

			latestVersion, err := checkForUpdate()
			if err != nil {
				return fmt.Errorf("failed to check for updates: %w", err)
			}

			comparison := compareVersions(currentVersion, latestVersion)
			switch comparison {
			case -1:
				fmt.Printf("Update available: %s -> %s\n", currentVersion, latestVersion)
				if err := downloadUpdate(latestVersion); err != nil {
					return fmt.Errorf("update failed: %w", err)
				}
				fmt.Println("OK   AuraSpeed has been updated! Restart the application to use the new version.")
			case 0:
				fmt.Printf("OK   You are running the latest version: %s\n", currentVersion)
			case 1:
				fmt.Printf("You are running a newer version (%s) than latest release (%s)\n", currentVersion, latestVersion)
			}

			return nil
		},
	}

	return cmd
}

// NewAutoUpdateCheck runs automatic update check if enabled in config
func NewAutoUpdateCheck() {
	cfg := config.Get()
	if !cfg.Global.AutoUpdate {
		return
	}

	// Only check if not running dev version
	if Version == "dev" {
		return
	}

	go func() {
		time.Sleep(5 * time.Second) // Wait a bit before checking

		latestVersion, err := checkForUpdate()
		if err != nil {
			return // Silent failure for auto-update check
		}

		comparison := compareVersions(Version, latestVersion)
		if comparison < 0 {
			logger := logging.Get()
			logger.Info("Update available: " + Version + " -> " + latestVersion)
		}
	}()
}