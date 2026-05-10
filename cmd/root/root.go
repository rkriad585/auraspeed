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

func selfUninstall() error {
	fmt.Println(">>> Uninstalling AuraSpeed...")

	binDir := installDir()
	binName := "auraspeed"
	if runtime.GOOS == "windows" {
		binName = "auraspeed.exe"
	}
	binPath := filepath.Join(binDir, binName)

	if _, err := os.Stat(binPath); err == nil {
		if err := os.Remove(binPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove binary: %v\n", err)
		} else {
			fmt.Printf("Removed  %s\n", binPath)
		}
	}

	auraspeedDir := filepath.Join(homeDir(), ".config", "neostore", "auraspeed")
	if _, err := os.Stat(auraspeedDir); err == nil {
		if err := os.RemoveAll(auraspeedDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove app directory: %v\n", err)
		} else {
			fmt.Printf("Removed  %s\n", auraspeedDir)
		}
	}

	configDir := filepath.Join(homeDir(), ".auraspeed")
	if _, err := os.Stat(configDir); err == nil {
		if err := os.RemoveAll(configDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove configuration: %v\n", err)
		} else {
			fmt.Printf("Removed  %s\n", configDir)
		}
	}

	fmt.Println()
	fmt.Println("OK   AuraSpeed has been uninstalled successfully!")

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
