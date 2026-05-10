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

// selfUninstall removes the current binary and its configuration
func selfUninstall() error {
	fmt.Println("Uninstalling AuraSpeed...")

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Define possible binary locations
	locations := []string{
		filepath.Join(homeDir, ".config", "neostore", "auraspeed", "bin", "auraspeed"),
		filepath.Join(homeDir, ".local", "bin", "auraspeed"),
		filepath.Join(homeDir, "bin", "auraspeed"),
		"/usr/local/bin/auraspeed",
		"/usr/bin/auraspeed",
	}

	if runtime.GOOS == "windows" {
		for i := range locations {
			locations[i] += ".exe"
		}
	}

	removedAny := false

	// Remove binaries
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			if err := os.Remove(loc); err == nil {
				fmt.Printf("Removed binary: %s\n", loc)
				removedAny = true
			}
		}
	}

	// Remove configuration directory (~/.auraspeed)
	configDir := filepath.Join(homeDir, ".auraspeed")
	if _, err := os.Stat(configDir); err == nil {
		os.RemoveAll(configDir)
		fmt.Printf("Removed configuration: %s\n", configDir)
	}

	// Remove auraspeed app directory (~/.config/neostore/auraspeed)
	auraspeedBinDir := filepath.Join(homeDir, ".config", "neostore", "auraspeed")
	if _, err := os.Stat(auraspeedBinDir); err == nil {
		// Remove all files in bin directory first
		binPath := filepath.Join(auraspeedBinDir, "bin")
		if _, err := os.Stat(binPath); err == nil {
			entries, _ := os.ReadDir(binPath)
			for _, entry := range entries {
				filePath := filepath.Join(binPath, entry.Name())
				os.Remove(filePath)
			}
			os.Remove(binPath)
		}
		// Remove auraspeed directory
		os.RemoveAll(auraspeedBinDir)
		fmt.Printf("Removed app directory: %s\n", auraspeedBinDir)
	}

	if !removedAny {
		fmt.Println("No AuraSpeed binary found in standard locations.")
	}

	fmt.Println("")
	fmt.Println("AuraSpeed has been uninstalled successfully!")
	fmt.Println("Note: If the binary was running, it may not have been removed.")

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
