package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"auraspeed/cmd/root"
	"auraspeed/internal/config"
	"auraspeed/internal/logging"
	"auraspeed/internal/speedtest"
)

func main() {
	if err := config.Init("auraspeed"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	if err := config.EnsureSensitiveFilePermissions(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	setupSignalHandling()

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	cleanup()
}

func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintln(os.Stderr, "\nShutting down gracefully...")
		cleanup()
		os.Exit(0)
	}()
}

func cleanup() {
	speedtest.StopTUI()
	logging.Get().Sync()
}
