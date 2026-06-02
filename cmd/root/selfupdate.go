package root

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"auraspeed/internal/config"

	"github.com/spf13/cobra"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	fmt.Printf("\r>>> Downloading... %s complete", bytesToMB(wc.Total))
}

func bytesToMB(bytes uint64) string {
	return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
}

func getBinaryName() string {
	bin := "auraspeed-" + runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	return bin
}

func getDownloadURL(version string) string {
	return fmt.Sprintf(
		"https://github.com/rkriad585/auraspeed/releases/download/%s/%s",
		version, getBinaryName(),
	)
}

func downloadBinary(version, proxyURL string) (string, error) {
	url := getDownloadURL(version)

	client := &http.Client{Timeout: 120 * time.Second}

	if proxyURL != "" {
		os.Setenv("HTTP_PROXY", proxyURL)
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "auraspeed-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	counter := &WriteCounter{}
	_, err = io.Copy(tmpFile, io.TeeReader(resp.Body, counter))
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write binary: %w", err)
	}

	tmpFile.Close()
	fmt.Println()

	if runtime.GOOS != "windows" {
		os.Chmod(tmpFile.Name(), 0755)
	}

	return tmpFile.Name(), nil
}

func replaceBinary(tmpPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	if runtime.GOOS == "windows" {
		oldPath := exePath + ".old"
		os.Remove(oldPath)

		if err := os.Rename(exePath, oldPath); err != nil {
			return fmt.Errorf("failed to rename current binary: %w", err)
		}

		err = os.Rename(tmpPath, exePath)
		if err != nil {
			os.Rename(oldPath, exePath)
			return fmt.Errorf("failed to replace binary, restored original: %w", err)
		}

		os.Remove(oldPath)
	} else {
		if err := os.Rename(tmpPath, exePath); err != nil {
			return fmt.Errorf("failed to replace binary: %w", err)
		}
	}

	return nil
}

func selfUpdateRun(cmd *cobra.Command, args []string) error {
	currentVersion := Version
	if currentVersion == "dev" {
		fmt.Println("Developer version detected. Update skipped.")
		return nil
	}

	proxyFlag, _ := cmd.Flags().GetString("proxy")
	proxyURL := proxyFlag
	if proxyURL == "" {
		proxyURL = config.Get().Network.Proxy
	}

	fmt.Printf(">>> Current version: %s\n", currentVersion)

	fmt.Println(">>> Checking for updates...")
	latestVersion, err := checkForUpdate()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	fmt.Printf(">>> Latest version: %s\n", latestVersion)

	comparison := compareVersions(currentVersion, latestVersion)
	if comparison >= 0 {
		fmt.Println("OK   You are already up to date!")
		return nil
	}

	fmt.Printf(">>> Downloading %s...\n", latestVersion)

	tmpPath, err := downloadBinary(latestVersion, proxyURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpPath)

	fmt.Println(">>> Installing update...")
	if err := replaceBinary(tmpPath); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Println("OK   AuraSpeed has been updated! Please restart the application.")
	return nil
}

func NewSelfUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update AuraSpeed to the latest version",
		Long:  "Check for updates, download, and replace the running AuraSpeed binary.",
		RunE:  selfUpdateRun,
	}

	cmd.Flags().StringP("proxy", "p", "", "Proxy URL for download (e.g., http://proxy:port)")
	return cmd
}
