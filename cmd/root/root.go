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
	"auraspeed/internal/logging"
	"auraspeed/internal/ui"

	"github.com/spf13/cobra"
)

var (
	logger *logging.Logger
	cfg    *config.Config

	// Version is the build version, set at build time via ldflags.
	Version = "dev"
	// Commit is the git commit hash, set at build time via ldflags.
	Commit = "none"
	// BuildTime is the build timestamp, set at build time via ldflags.
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "auraspeed",
	Short: "AuraSpeed - Advanced Network Analyzer & System Utility",
	Long: `AuraSpeed is a cross-platform terminal tool for network diagnostics,
system monitoring, and performance optimization.

Features:
  - Speed test with real-time graphs
  - System information panel
  - Network diagnostics
  - Command aliases
  - Configurable profiles

Use 'auraspeed [command] --help' for more information about a command.`,
	SilenceUsage:  true,
	SilenceErrors: false,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel, _ := cmd.Flags().GetString("log-level")
		if err := logging.SetLevel(logLevel); err != nil {
			return fmt.Errorf("invalid log level: %w", err)
		}
		noColor, _ := cmd.Flags().GetBool("no-color")
		logging.SetNoColor(noColor)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}
		return cmd.Help()
	},
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	h, _ := os.UserHomeDir()
	return h
}

func installDir() string {
	dir := filepath.Join(homeDir(), ".config", "neostore", "auraspeed", "bin")
	os.MkdirAll(dir, 0755)
	return dir
}

func selfInstall(version string) int {
	targetDir := installDir()
	binName := "auraspeed"
	if runtime.GOOS == "windows" {
		binName = "auraspeed.exe"
	}
	targetPath := filepath.Join(targetDir, binName)

	fmt.Printf(">>> Installing AuraSpeed %s...\n", version)

	repo := "rkriad585/auraspeed"
	var downloadName string

	switch runtime.GOOS {
	case "windows":
		downloadName = "auraspeed-windows-amd64.exe"
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
		fmt.Fprintf(os.Stderr, "Unsupported platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return 1
	}

	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, version, downloadName)

	fmt.Printf(">>> Downloading %s\n", url)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		fmt.Println("Falling back to copying the current binary...")
		return copySelf(targetPath)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Download failed: HTTP %d\n", resp.StatusCode)
		fmt.Println("Falling back to copying the current binary...")
		return copySelf(targetPath)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return 1
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		out.Close()
		os.Remove(targetPath)
		return 1
	}
	out.Close()

	if written == 0 {
		fmt.Fprintln(os.Stderr, "Downloaded empty file")
		os.Remove(targetPath)
		fmt.Println("Falling back to copying the current binary...")
		return copySelf(targetPath)
	}

	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0755)
	}

	fmt.Printf("OK   Installed to %s (%d bytes)\n", targetPath, written)

	return addToPath(targetDir)
}

func copySelf(targetPath string) int {
	src, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
		return 1
	}

	in, err := os.Open(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening source: %v\n", err)
		return 1
	}
	defer in.Close()

	out, err := os.Create(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating destination: %v\n", err)
		return 1
	}
	defer out.Close()

	written, err := io.Copy(out, in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying file: %v\n", err)
		out.Close()
		os.Remove(targetPath)
		return 1
	}
	out.Close()

	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0755)
	}

	fmt.Printf("OK   Installed to %s (%d bytes)\n", targetPath, written)
	return 0
}

func addToPath(targetDir string) int {
	fmt.Printf("Add %s to your PATH to use the 'auraspeed' command\n", targetDir)
	return 0
}

func Execute() error {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("AuraSpeed %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildTime)
		fmt.Printf("Go: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return nil
	}
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolP("no-color", "", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolP("version", "", false, "Show version information")

	rootCmd.AddCommand(ui.NewTUICommand())
	rootCmd.AddCommand(NewSpeedtestCommand())
	rootCmd.AddCommand(NewInfoCommand())
	rootCmd.AddCommand(NewNetworkCommand())
	rootCmd.AddCommand(NewHistoryCommand())
	rootCmd.AddCommand(NewConfigCommand())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(NewWebCommand())
	rootCmd.AddCommand(NewUpdateCommand())
	rootCmd.AddCommand(NewServersCommand())
	rootCmd.AddCommand(NewInstallCommand())
}
