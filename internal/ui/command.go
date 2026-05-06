package ui

import (
	"auraspeed/internal/logging"
	"auraspeed/internal/speedtest"

	"github.com/spf13/cobra"
)

func NewTUICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive TUI",
		Long:  "Open the interactive terminal user interface with full functionality.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.Silence()
			defer logging.Restore()
			return speedtest.RunTUI()
		},
	}
	cmd.Flags().BoolP("fullscreen", "f", false, "Run in fullscreen mode")
	return cmd
}
