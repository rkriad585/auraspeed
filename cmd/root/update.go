package root

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"auraspeed/internal/config"
	"auraspeed/internal/logging"

	"github.com/spf13/cobra"
)

// updateCheckResult represents the response from GitHub releases API
type updateCheckResult struct {
	TagName string `json:"tag_name"`
}

// checkForUpdate checks GitHub for the latest release version
func checkForUpdate() (string, error) {
	const (
		owner = "rkriad585"
		repo  = "auraspeed"
	)

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
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

	var result updateCheckResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.TagName, nil
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
		Short: "Check for updates",
		Long:  "Check if a newer version of AuraSpeed is available.",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentVersion := Version
			if currentVersion == "dev" {
				fmt.Println("Running from source - cannot check for updates.")
				fmt.Println("To check for updates, use a released binary.")
				return nil
			}

			fmt.Println("Checking for updates...")

			latestVersion, err := checkForUpdate()
			if err != nil {
				return fmt.Errorf("failed to check for updates: %w", err)
			}

			comparison := compareVersions(currentVersion, latestVersion)
			switch comparison {
			case -1:
				fmt.Printf("Update available: %s -> %s\n", currentVersion, latestVersion)
				fmt.Println("Download from: https://github.com/rkriad585/auraspeed/releases")
				if !checkOnly {
					fmt.Println("\nTo update, download the new release and replace your binary.")
				}
			case 0:
				fmt.Printf("You are running the latest version: %s\n", currentVersion)
			case 1:
				fmt.Printf("You are running a newer version (%s) than latest release (%s)\n", currentVersion, latestVersion)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates, don't display instructions")

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