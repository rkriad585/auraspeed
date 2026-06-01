package root

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// NewInstallCommand returns the install subcommand
func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install AuraSpeed as a system service",
		Long:  "Install AuraSpeed as a systemd service (Linux) for automatic startup.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installSystemd()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "systemd",
		Short: "Install systemd service",
		Long:  "Install AuraSpeed as a systemd service for automatic startup on Linux.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installSystemd()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "uninstall",
		Short: "Remove systemd service",
		Long:  "Remove the AuraSpeed systemd service from the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstallSystemd()
		},
	})

	return cmd
}

func installSystemd() error {
	// Check if we're on Linux
	if runtime.GOOS != "linux" {
		fmt.Println("Note: systemd installation is only supported on Linux systems.")
		fmt.Println("For other systems, please use docker-compose or manual startup.")
	}

	// Create systemd directory if it doesn't exist
	systemdDir := "/etc/systemd/system"
	if err := os.MkdirAll(systemdDir, 0755); err != nil {
		return fmt.Errorf("failed to create systemd directory: %w", err)
	}

	// Read the service file template
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	serviceContent, err := os.ReadFile("auraspeed.service")
	if err != nil {
		altPath := filepath.Join(filepath.Dir(execPath), "auraspeed.service")
		serviceContent, err = os.ReadFile(altPath)
		if err != nil {
			return fmt.Errorf("failed to read service file: %w", err)
		}
	}

	servicePath := filepath.Join(systemdDir, "auraspeed.service")
	if err := os.WriteFile(servicePath, serviceContent, 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		fmt.Printf("Warning: Failed to reload systemd: %v\n", err)
		fmt.Println("You may need to run 'sudo systemctl daemon-reload' manually")
	}

	// Enable the service
	if err := exec.Command("systemctl", "enable", "auraspeed.service").Run(); err != nil {
		fmt.Printf("Warning: Failed to enable service: %v\n", err)
		fmt.Println("You may need to run 'sudo systemctl enable auraspeed' manually")
	}

	fmt.Println("AuraSpeed systemd service installed successfully!")
	fmt.Println("")
	fmt.Println("To start the service:")
	fmt.Println("  sudo systemctl start auraspeed")
	fmt.Println("")
	fmt.Println("To check status:")
	fmt.Println("  sudo systemctl status auraspeed")
	fmt.Println("")
	fmt.Println("To view logs:")
	fmt.Println("  journalctl -u auraspeed -f")

	return nil
}

func uninstallSystemd() error {
	// Stop the service
	if err := exec.Command("systemctl", "stop", "auraspeed").Run(); err != nil {
		fmt.Printf("Warning: Failed to stop service: %v\n", err)
	}

	// Disable the service
	if err := exec.Command("systemctl", "disable", "auraspeed").Run(); err != nil {
		fmt.Printf("Warning: Failed to disable service: %v\n", err)
	}

	// Remove the service file
	servicePath := "/etc/systemd/system/auraspeed.service"
	if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		fmt.Printf("Warning: Failed to reload systemd: %v\n", err)
	}

	fmt.Println("AuraSpeed systemd service removed successfully!")
	return nil
}