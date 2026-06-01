package ui

import (
	"auraspeed/internal/logging"
	"auraspeed/internal/speedtest"

	"github.com/spf13/cobra"
)

func NewTUICommand() *cobra.Command {
	var fullscreen bool

	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive TUI",
		Long:  "Open the interactive terminal user interface with full functionality.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.Silence()
			defer logging.Restore()
			if fullscreen {
				return speedtest.RunTUIWithOptions(speedtest.TUIOptions{Fullscreen: true})
			}
			return speedtest.RunTUI()
		},
	}
	cmd.Flags().BoolVarP(&fullscreen, "fullscreen", "f", false, "Run in fullscreen mode")
	return cmd
}
