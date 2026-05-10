package root

import (
	"fmt"
	"os"
	"os/exec"
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

func installDir() string {
	dir := filepath.Join(homeDir(), ".config", "neostore", "auraspeed", "bin")
	os.MkdirAll(dir, 0755)
	return dir
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	h, _ := os.UserHomeDir()
	return h
}

func ConfigDir() string {
	return filepath.Join(homeDir(), ".config", "neostore", "auraspeed")
}

func selfUninstall() int {
	configDir := ConfigDir()

	fmt.Println(">>> Uninstalling AuraSpeed...")

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving binary path: %v\n", err)
		return 1
	}

	realPath, err := filepath.EvalSymlinks(exePath)
	if err == nil {
		exePath = realPath
	}

	if runtime.GOOS == "windows" {
		batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak >nul\r\ndel /f /q \"%s\"\r\nif exist \"%s\" (echo ERR Failed to delete binary) else (echo OK   Deleted binary: %s)\r\nrd /s /q \"%s\"\r\ndel /f /q \"%%~f0\"\r\n", exePath, exePath, exePath, configDir)
		batPath := filepath.Join(os.TempDir(), "auraspeed-uninstall.bat")
		if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating uninstall script: %v\n", err)
			fmt.Println("Please manually delete:", exePath)
			return 1
		}
		fmt.Println("OK   Uninstall script created. Binary and config will be deleted shortly.")
		exec.Command("cmd", "/C", "start", "/B", batPath).Start()
	} else {
		if _, err := os.Stat(configDir); err == nil {
			if err := os.RemoveAll(configDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing config directory: %v\n", err)
				return 1
			}
			fmt.Println("OK   Removed config directory:", configDir)
		} else {
			fmt.Println("OK   No config directory found.")
		}

		if err := os.Remove(exePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting binary: %v\n", err)
			fmt.Println("Please manually delete:", exePath)
			return 1
		}
		fmt.Println("OK   Deleted binary:", exePath)
	}

	fmt.Println()
	fmt.Println("To remove AuraSpeed from your PATH, edit your shell rc file")
	fmt.Println("and delete the line containing 'neostore/auraspeed/bin'.")

	fmt.Println("Or run: auraspeed --selfuninstall")

	fmt.Println("Restart your terminal for PATH changes to take effect.")
	return 0
}

func Execute() error {
	for _, arg := range os.Args {
		if arg == "--selfuninstall" || arg == "--uninstall" {
			os.Exit(selfUninstall())
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
