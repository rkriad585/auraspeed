package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// PingResult holds the results of a ping operation.
type PingResult struct {
	Host    string `json:"host"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
}

// RunPing executes a ping command against the specified host.
// It uses 4 pings and displays the output.
func RunPing(host string, jsonOutput bool) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", host)
	} else {
		cmd = exec.Command("ping", "-c", "4", host)
	}
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if jsonOutput {
		result := PingResult{
			Host:    host,
			Output:  outputStr,
			Success: err == nil,
		}
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	}

	if err != nil {
		return fmt.Errorf("ping failed: %w\n%s", err, outputStr)
	}
	fmt.Print(outputStr)
	return nil
}

// RunTraceroute executes a traceroute/tracert command to trace the network path.
func RunTraceroute(host string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tracert", host)
	} else {
		cmd = exec.Command("traceroute", host)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("traceroute failed: %w\n%s", err, string(output))
	}
	fmt.Print(string(output))
	return nil
}

// DNSResult holds the results of a DNS lookup.
type DNSResult struct {
	Host       string   `json:"host"`
	Addresses  []string `json:"addresses,omitempty"`
	Hostnames  []string `json:"hostnames,omitempty"`
}

// RunDNSLookup performs DNS lookup for the given host (forward or reverse).
func RunDNSLookup(host string, jsonOutput bool) error {
	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("DNS lookup failed: %w", err)
	}

	var hostnames []string
	if net.ParseIP(host) != nil {
		hostnames, _ = net.LookupAddr(host)
	}

	if jsonOutput {
		result := DNSResult{
			Host:      host,
			Addresses: ips,
			Hostnames: hostnames,
		}
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
		return nil
	}

	fmt.Printf("DNS Lookup for: %s\n", host)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("IP Addresses:")
	for _, ip := range ips {
		fmt.Printf("  %s\n", ip)
	}

	if len(hostnames) > 0 {
		fmt.Println()
		fmt.Println("Hostnames:")
		for _, name := range hostnames {
			fmt.Printf("  %s\n", name)
		}
	}

	return nil
}
