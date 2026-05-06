package root

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display version, build time, commit hash, and runtime info.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("AuraSpeed %s\n", Version)
			fmt.Printf("Commit: %s\n", Commit)
			fmt.Printf("Built: %s\n", BuildTime)
			fmt.Printf("Go: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}
