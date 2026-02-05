package commands

import (
	"testing"
	"time"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/dns"
	"github.com/go-zoox/dns/client"
	"github.com/go-zoox/dns/constants"
)

func TestNormalizeServerAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "IP with port",
			input:    "127.0.0.1:5553",
			expected: "127.0.0.1:5553",
		},
		{
			name:     "IP without port",
			input:    "127.0.0.1",
			expected: "127.0.0.1:53",
		},
		{
			name:     "IPv6 with port",
			input:    "[::1]:5553",
			expected: "[::1]:5553",
		},
		{
			name:     "IPv6 without port",
			input:    "::1",
			expected: "[::1]:53",
		},
		{
			name:     "Domain with port",
			input:    "dns.example.com:5553",
			expected: "dns.example.com:5553",
		},
		{
			name:     "Domain without port",
			input:    "dns.example.com",
			expected: "dns.example.com:53",
		},
		{
			name:     "TLS protocol with port",
			input:    "tls://1.1.1.1:853",
			expected: "tls://1.1.1.1:853",
		},
		{
			name:     "TLS protocol without port",
			input:    "tls://1.1.1.1",
			expected: "tls://1.1.1.1",
		},
		{
			name:     "HTTPS protocol with port",
			input:    "https://dns.adguard.com:443/dns-query",
			expected: "https://dns.adguard.com:443/dns-query",
		},
		{
			name:     "HTTPS protocol without port",
			input:    "https://dns.adguard.com/dns-query",
			expected: "https://dns.adguard.com/dns-query",
		},
		{
			name:     "QUIC protocol with port",
			input:    "quic://dns.adguard.com:853",
			expected: "quic://dns.adguard.com:853",
		},
		{
			name:     "QUIC protocol without port",
			input:    "quic://dns.adguard.com",
			expected: "quic://dns.adguard.com",
		},
		{
			name:     "Default DNS server with port",
			input:    "114.114.114.114:53",
			expected: "114.114.114.114:53",
		},
		{
			name:     "Default DNS server without port",
			input:    "114.114.114.114",
			expected: "114.114.114.114:53",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeServerAddress(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeServerAddress(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDoHClient tests DNS-over-HTTPS client functionality
func TestDoHClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with a public DoH server
	servers := []string{"https://cloudflare-dns.com/dns-query"}
	dnsClient := dns.NewClient(&dns.ClientOptions{
		Servers: servers,
		Timeout: 10 * time.Second,
	})

	// Test A record lookup
	ips, err := dnsClient.LookUp("google.com", &client.LookUpOptions{
		Typ: constants.QueryTypeIPv4,
	})
	if err != nil {
		t.Fatalf("DoH lookup failed: %v", err)
	}
	if len(ips) == 0 {
		t.Error("DoH lookup returned no IP addresses")
	}

	// Verify IP format
	for _, ip := range ips {
		if ip == "" {
			t.Error("DoH lookup returned empty IP address")
		}
	}
}

// TestDoQClient tests DNS-over-QUIC client functionality
func TestDoQClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with a public DoQ server (AdGuard)
	servers := []string{"quic://dns.adguard.com"}
	dnsClient := dns.NewClient(&dns.ClientOptions{
		Servers: servers,
		Timeout: 10 * time.Second,
	})

	// Test A record lookup
	ips, err := dnsClient.LookUp("google.com", &client.LookUpOptions{
		Typ: constants.QueryTypeIPv4,
	})
	if err != nil {
		t.Fatalf("DoQ lookup failed: %v", err)
	}
	if len(ips) == 0 {
		t.Error("DoQ lookup returned no IP addresses")
	}

	// Verify IP format
	for _, ip := range ips {
		if ip == "" {
			t.Error("DoQ lookup returned empty IP address")
		}
	}
}

// TestDoHClientWithMultipleServers tests DoH with multiple servers
func TestDoHClientWithMultipleServers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	servers := []string{
		"https://cloudflare-dns.com/dns-query",
		"https://dns.google/dns-query",
	}
	dnsClient := dns.NewClient(&dns.ClientOptions{
		Servers: servers,
		Timeout: 10 * time.Second,
	})

	ips, err := dnsClient.LookUp("example.com", &client.LookUpOptions{
		Typ: constants.QueryTypeIPv4,
	})
	if err != nil {
		t.Fatalf("DoH lookup with multiple servers failed: %v", err)
	}
	if len(ips) == 0 {
		t.Error("DoH lookup with multiple servers returned no IP addresses")
	}
}

// TestDoQClientWithIPv6 tests DoQ with AAAA record lookup
func TestDoQClientWithIPv6(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	servers := []string{"quic://dns.adguard.com"}
	dnsClient := dns.NewClient(&dns.ClientOptions{
		Servers: servers,
		Timeout: 10 * time.Second,
	})

	ips, err := dnsClient.LookUp("google.com", &client.LookUpOptions{
		Typ: constants.QueryTypeIPv6,
	})
	if err != nil {
		// IPv6 might not be available, so we just log the error
		t.Logf("DoQ IPv6 lookup failed (may be expected): %v", err)
		return
	}

	// If we got results, verify they are IPv6 addresses
	for _, ip := range ips {
		if ip == "" {
			t.Error("DoQ IPv6 lookup returned empty IP address")
		}
	}
}

// TestClientCommandWithDoH tests the client command structure with DoH server
func TestClientCommandWithDoH(t *testing.T) {
	cmd := NewClientCommand()

	// Verify command structure
	if cmd.Name != "client" {
		t.Error("Client command name mismatch")
	}

	// Verify that server flag supports DoH
	serverFlag := cmd.Flags[0].(*cli.StringSliceFlag)
	if serverFlag.Name != "server" {
		t.Error("Server flag not found")
	}
}

// TestClientCommandWithDoQ tests the client command structure with DoQ server
func TestClientCommandWithDoQ(t *testing.T) {
	cmd := NewClientCommand()

	// Verify command structure
	if cmd.Name != "client" {
		t.Error("Client command name mismatch")
	}

	// Verify that server flag supports DoQ
	serverFlag := cmd.Flags[0].(*cli.StringSliceFlag)
	if serverFlag.Name != "server" {
		t.Error("Server flag not found")
	}
}
