package network

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// RunPing executes a ping command against the specified host.
// It uses 4 pings and displays the output.
func RunPing(host string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", host)
	} else {
		cmd = exec.Command("ping", "-c", "4", host)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ping failed: %w\n%s", err, string(output))
	}
	fmt.Print(string(output))
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

// RunDNSLookup performs DNS lookup for the given host (forward or reverse).
func RunDNSLookup(host string) error {
	fmt.Printf("DNS Lookup for: %s\n", host)
	fmt.Println(strings.Repeat("-", 50))

	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("DNS lookup failed: %w", err)
	}
	fmt.Println("IP Addresses:")
	for _, ip := range ips {
		fmt.Printf("  %s\n", ip)
	}

	fmt.Println()

	// Try reverse lookup if it looks like an IP
	if net.ParseIP(host) != nil {
		names, err := net.LookupAddr(host)
		if err != nil {
			fmt.Printf("Reverse lookup failed: %v\n", err)
		} else {
			fmt.Println("Hostnames:")
			for _, name := range names {
				fmt.Printf("  %s\n", name)
			}
		}
	}

	return nil
}
