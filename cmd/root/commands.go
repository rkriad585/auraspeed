package root

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"auraspeed/internal/config"
	"auraspeed/internal/info"
	"auraspeed/internal/network"
	"auraspeed/internal/speedtest"

	st "github.com/showwin/speedtest-go/speedtest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSpeedtestCommand returns the speedtest subcommand.
// It runs a network speed test and outputs results to terminal.
func NewSpeedtestCommand() *cobra.Command {
	var serverID string
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "speedtest",
		Short: "Run network speed test",
		Long:  "Perform a comprehensive network speed test with download, upload, and latency measurements.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !jsonOutput {
				fmt.Println("Running speed test...")
			}

			client := st.New()

			user, err := client.FetchUserInfo()
			if err != nil {
				return fmt.Errorf("failed to fetch user info: %w", err)
			}

			var serverList st.Servers
			if cachedServers, cacheHit := speedtest.GetCachedServers(); cacheHit {
				if !jsonOutput {
					fmt.Println("Using cached server list...")
				}
				serverList = cachedServers
			} else {
				var fetchErr error
				for attempt := 1; attempt <= 3; attempt++ {
					serverList, fetchErr = client.FetchServers()
					if fetchErr == nil && len(serverList) > 0 {
						break
					}
					if attempt < 3 {
						if !jsonOutput {
							fmt.Printf("Retrying... (attempt %d/3)\n", attempt+1)
						}
						time.Sleep(2 * time.Second)
					}
				}
				if fetchErr != nil {
					return fmt.Errorf("failed to fetch servers after 3 attempts: %w", fetchErr)
				}
				if len(serverList) == 0 {
					return fmt.Errorf("no servers found after 3 attempts")
				}
				if err := speedtest.SaveServersToCache(serverList); err != nil {
					if !jsonOutput {
						fmt.Printf("Warning: Failed to save server cache: %v\n", err)
					}
				}
			}

			var targets []*st.Server
			if serverID != "" {
				sid, err := strconv.Atoi(serverID)
				if err != nil {
					return fmt.Errorf("invalid server ID: %w", err)
				}
				targets, err = serverList.FindServer([]int{sid})
			} else {
				targets, err = serverList.FindServer([]int{})
			}
			if err != nil {
				return fmt.Errorf("failed to find server: %w", err)
			}
			if len(targets) == 0 {
				return fmt.Errorf("no servers found")
			}

			s := targets[0]

			if err := s.DownloadTest(); err != nil {
				return fmt.Errorf("download test failed: %w", err)
			}

			if err := s.UploadTest(); err != nil {
				return fmt.Errorf("upload test failed: %w", err)
			}

			dlMbps := float64(s.DLSpeed) / 1000000.0
			ulMbps := float64(s.ULSpeed) / 1000000.0
			pingMs := s.Latency.Milliseconds()
			isp := user.Isp

			if jsonOutput {
				result := struct {
					Download float64 `json:"download"`
					Upload   float64 `json:"upload"`
					Ping     int64   `json:"ping"`
					ISP      string  `json:"isp"`
					Server   string  `json:"server"`
				}{
					Download: dlMbps,
					Upload:   ulMbps,
					Ping:     pingMs,
					ISP:      isp,
					Server:   s.Name,
				}
				jsonData, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			} else {
				fmt.Printf("Server: %s\n", s.Name)
				fmt.Printf("ISP:    %s\n", isp)
				fmt.Println("\nResults:")
				fmt.Println("--------")
				fmt.Printf("Download: %.2f Mbps\n", dlMbps)
				fmt.Printf("Upload:   %.2f Mbps\n", ulMbps)
				fmt.Printf("Ping:     %d ms\n", pingMs)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&serverID, "server-id", "", "Specify server ID to use for speed test")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")

	return cmd
}

// NewInfoCommand returns the info subcommand.
// It displays system information (OS, CPU, memory, disk).
func NewInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display system information",
		Long:  "Show system information including CPU, memory, disk, and network stats.",
		RunE: func(cmd *cobra.Command, args []string) error {
			sysInfo, err := info.GetSystemInfo()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting system info: %v\n", err)
				return err
			}

			fmt.Println("System Information")
			fmt.Println("------------------")
			fmt.Printf("OS:       %s\n", sysInfo.OS)
			fmt.Printf("CPU:      %s\n", sysInfo.CPU)
			fmt.Printf("Memory:   %s\n", sysInfo.Memory)
			fmt.Printf("Disk:     %s\n", sysInfo.Disk)
			fmt.Printf("Hostname: %s\n", sysInfo.Hostname)

			return nil
		},
	}
	return cmd
}

// NewNetworkCommand returns the network subcommand.
// It provides network diagnostic tools (ping, traceroute, dns).
func NewNetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Network diagnostics",
		Long:  "Network diagnostic tools including ping, traceroute, and DNS lookup.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	pingCmd := &cobra.Command{
		Use:   "ping <host>",
		Short: "Ping a host",
		Long:  "Send ICMP echo requests to a host to test connectivity.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network.RunPing(args[0])
		},
	}

	tracerouteCmd := &cobra.Command{
		Use:   "traceroute <host>",
		Short: "Trace route to host",
		Long:  "Trace the network path to a host showing all hops.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network.RunTraceroute(args[0])
		},
	}

	dnsCmd := &cobra.Command{
		Use:   "dns <host>",
		Short: "DNS lookup",
		Long:  "Perform DNS lookup for a hostname or IP address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network.RunDNSLookup(args[0])
		},
	}

	cmd.AddCommand(pingCmd)
	cmd.AddCommand(tracerouteCmd)
	cmd.AddCommand(dnsCmd)

	return cmd
}

// NewHistoryCommand returns the history subcommand.
// It views and manages speed test history.
func NewHistoryCommand() *cobra.Command {
	var limit int
	var clearHistory bool
	var exportFile string
	var saveOverride bool

	cmd := &cobra.Command{
		Use:   "history",
		Short: "View test history",
		Long:  "View and manage speed test history and results.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Get()
			if !cfg.UI.SaveHistory && !saveOverride {
				fmt.Println("Note: History saving is disabled. Results are not being saved.")
				fmt.Println("To enable, run: auraspeed config set ui.savehistory true")
				fmt.Println("Use --save flag to temporarily enable saving for this session.")
			}
			if saveOverride {
				cfg.UI.SaveHistory = true
				fmt.Println("History saving temporarily enabled for this session.")
			}

			if clearHistory {
				historyPath := config.GetHistoryFile()
				return clearHistoryFile(historyPath)
			}

			results, err := speedtest.LoadHistory()
			if err != nil {
				return fmt.Errorf("failed to load history: %w", err)
			}

			if len(results) == 0 {
				fmt.Println("No test history found.")
				return nil
			}

			if exportFile != "" {
				return exportHistory(results, exportFile)
			}

			if limit > 0 && limit < len(results) {
				results = results[len(results)-limit:]
			}

			displayHistoryTable(results)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "Show last N results")
	cmd.Flags().BoolVar(&clearHistory, "clear", false, "Clear all history")
	cmd.Flags().StringVar(&exportFile, "export", "", "Export history to file")
	cmd.Flags().BoolVar(&saveOverride, "save", false, "Temporarily enable history saving")

	return cmd
}

func displayHistoryTable(results []speedtest.TestResult) {
	fmt.Println("Test History")
	fmt.Println("------------")

	fmt.Printf("%-10s %-12s %-12s %-8s %s\n", "Time", "Download", "Upload", "Ping", "ISP")
	fmt.Println("------------------------------------------------------------------------")

	for _, r := range results {
		fmt.Printf("%-10s %-12.2f %-12.2f %-8d %s\n", r.Timestamp, r.Download, r.Upload, r.Ping, r.ISP)
	}
}

func clearHistoryFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear history: %w", err)
	}
	fmt.Println("History cleared successfully.")
	return nil
}

func exportHistory(results []speedtest.TestResult, filename string) error {
	ext := filepath.Ext(filename)
	if ext == ".json" {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create export file: %w", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(results); err != nil {
			return fmt.Errorf("failed to write JSON: %w", err)
		}
	} else {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create export file: %w", err)
		}
		defer file.Close()

		fmt.Fprintf(file, "Time\tDownload\tUpload\tPing\tISP\n")
		for _, r := range results {
			fmt.Fprintf(file, "%s\t%.2f\t%.2f\t%d\t%s\n", r.Timestamp, r.Download, r.Upload, r.Ping, r.ISP)
		}
	}

	fmt.Printf("History exported to %s\n", filename)
	return nil
}

func runConfigView(filterSection string) error {
	settings := viper.AllSettings()
	var lines []string

	for section, val := range settings {
		if filterSection != "" && section != filterSection {
			continue
		}
		switch v := val.(type) {
		case map[string]interface{}:
			for key, subVal := range v {
				lines = append(lines, fmt.Sprintf("%s.%s: %v", section, key, subVal))
			}
		default:
			lines = append(lines, fmt.Sprintf("%s: %v", section, val))
		}
	}

	sort.Strings(lines)
	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}

func newConfigViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [section]",
		Short: "View current configuration",
		Long:  "Display all current configuration settings. Optionally specify a section (global, speedtest, ui).",
		RunE: func(cmd *cobra.Command, args []string) error {
			section := ""
			if len(args) > 0 {
				section = args[0]
			}
			return runConfigView(section)
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a specific configuration key to a value.
Example: auraspeed config set speedtest.timeout 60`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			current := viper.Get(key)
			switch current.(type) {
			case int:
				var intVal int
				if _, err := fmt.Sscanf(value, "%d", &intVal); err == nil {
					viper.Set(key, intVal)
				} else {
					viper.Set(key, value)
				}
			case bool:
				switch value {
				case "true":
					viper.Set(key, true)
				case "false":
					viper.Set(key, false)
				default:
					viper.Set(key, value)
				}
			case float64:
				var floatVal float64
				if _, err := fmt.Sscanf(value, "%f", &floatVal); err == nil {
					viper.Set(key, floatVal)
				} else {
					viper.Set(key, value)
				}
			default:
				viper.Set(key, value)
			}

			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("Set %s to %v\n", key, viper.Get(key))
			return nil
		},
	}
	return cmd
}

func newConfigResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  "Reset all configuration settings to their default values.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := config.GetConfigFile()

			if err := os.Remove(configFile); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to delete config file: %w", err)
			}

			viper.Reset()

			if err := config.Init("auraspeed"); err != nil {
				return fmt.Errorf("failed to re-initialize config: %w", err)
			}

			fmt.Println("Configuration reset to defaults.")
			return nil
		},
	}
}

// NewConfigCommand returns the config subcommand.
// It views, sets, and resets configuration values.
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  "View, edit, or reset AuraSpeed configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			section := ""
			if len(args) > 0 {
				section = args[0]
			}
			return runConfigView(section)
		},
	}

	cmd.AddCommand(newConfigViewCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigResetCmd())

	return cmd
}
