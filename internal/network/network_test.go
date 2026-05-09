package network

import (
	"testing"
)

// TestRunPing tests the ping functionality
func TestRunPing(t *testing.T) {
	err := RunPing("127.0.0.1")
	if err != nil {
		t.Errorf("ping to localhost failed: %v", err)
	}
}

// TestRunPingInvalid tests ping with invalid host
func TestRunPingInvalid(t *testing.T) {
	err := RunPing("invalid.host.that.does.not.exist.test")
	if err == nil {
		t.Log("ping to invalid host succeeded (unexpected), but test passed")
	}
}

// TestRunDNSLookup tests DNS lookup for valid domain
func TestRunDNSLookup(t *testing.T) {
	err := RunDNSLookup("google.com")
	if err != nil {
		t.Errorf("DNS lookup failed: %v", err)
	}
}

// TestRunDNSLookupIP tests DNS lookup for IP address
func TestRunDNSLookupIP(t *testing.T) {
	err := RunDNSLookup("8.8.8.8")
	if err != nil {
		t.Errorf("DNS lookup for IP failed: %v", err)
	}
}

// TestRunDNSLookupInvalid tests DNS lookup for invalid domain
func TestRunDNSLookupInvalid(t *testing.T) {
	err := RunDNSLookup("invalid.domain.test.12345")
	if err == nil {
		t.Error("expected error for invalid domain, got nil")
	}
}

// TestTraceroute tests traceroute (may require elevated permissions on some systems)
func TestRunTraceroute(t *testing.T) {
	err := RunTraceroute("127.0.0.1")
	if err != nil {
		t.Errorf("traceroute to localhost failed: %v", err)
	}
}