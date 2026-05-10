package root

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

// NewUpdateCommand returns the update subcommand
func NewUpdateCommand() *cobra.Command {
	var checkOnly bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update AuraSpeed to the latest version",
		Long:  "Check for and install the latest version of AuraSpeed.",
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
				if checkOnly {
					fmt.Printf("Run 'auraspeed update' to install the latest version.\n")
				} else {
					result := selfInstall(latestVersion)
					if result != 0 {
						return fmt.Errorf("update installation failed")
					}
				}
			case 0:
				fmt.Printf("OK   You are running the latest version: %s\n", currentVersion)
			case 1:
				fmt.Printf("You are running a newer version (%s) than latest release (%s)\n", currentVersion, latestVersion)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates, don't install")

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