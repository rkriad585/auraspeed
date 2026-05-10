package root

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

// getHomeDir returns the user's home directory
func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	h, _ := os.UserHomeDir()
	return h
}

// getInstallDir returns the installation directory for AuraSpeed
func getInstallDir() string {
	dir := filepath.Join(getHomeDir(), ".config", "neostore", "auraspeed", "bin")
	os.MkdirAll(dir, 0755)
	return dir
}

// selfUninstall removes the current binary and its configuration
func selfUninstall() error {
	fmt.Println("Uninstalling AuraSpeed...")

	// Get install directory
	installDir := getInstallDir()

	// Binary name
	binName := "auraspeed"
	if runtime.GOOS == "windows" {
		binName = "auraspeed.exe"
	}

	// Remove binary from install directory
	binPath := filepath.Join(installDir, binName)
	if _, err := os.Stat(binPath); err == nil {
		os.Remove(binPath)
		fmt.Printf("Removed binary: %s\n", binPath)
	}

	// Remove empty bin directory
	binDir := filepath.Join(filepath.Dir(binPath))
	if _, err := os.Stat(binDir); err == nil {
		os.Remove(binDir)
	}

	// Remove auraspeed directory
	auraspeedDir := filepath.Join(getHomeDir(), ".config", "neostore", "auraspeed")
	if _, err := os.Stat(auraspeedDir); err == nil {
		os.RemoveAll(auraspeedDir)
		fmt.Printf("Removed app directory: %s\n", auraspeedDir)
	}

	// Remove configuration directory (~/.auraspeed)
	configDir := filepath.Join(getHomeDir(), ".auraspeed")
	if _, err := os.Stat(configDir); err == nil {
		os.RemoveAll(configDir)
		fmt.Printf("Removed configuration: %s\n", configDir)
	}

	fmt.Println("")
	fmt.Println("AuraSpeed has been uninstalled successfully!")

	return nil
}

func Execute() error {
	// Check for self-uninstall flag first
	for _, arg := range os.Args {
		if arg == "--selfuninstall" || arg == "--uninstall" {
			return selfUninstall()
		}
	}

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
	rootCmd.PersistentFlags().Bool("selfuninstall", false, "Uninstall AuraSpeed from the system")

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
