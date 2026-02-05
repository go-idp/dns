package commands

import (
	"testing"
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
