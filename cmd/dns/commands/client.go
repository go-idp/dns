package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-zoox/cli"
)

// normalizeServerAddress normalizes a DNS server address by adding default port if missing
// Supports formats like:
//   - "127.0.0.1:5553" -> "127.0.0.1:5553" (unchanged)
//   - "127.0.0.1" -> "127.0.0.1:53" (adds default port)
//   - "tls://1.1.1.1" -> "tls://1.1.1.1" (protocol prefix, unchanged)
//   - "tls://1.1.1.1:853" -> "tls://1.1.1.1:853" (protocol with port, unchanged)
func normalizeServerAddress(server string) string {
	// Check if it's a protocol-prefixed address (tls://, https://, etc.)
	if strings.Contains(server, "://") {
		// For protocol-prefixed addresses, check if port is already specified
		parts := strings.Split(server, "://")
		if len(parts) != 2 {
			return server
		}
		address := parts[1]

		// Check if address already has a port
		if _, _, err := net.SplitHostPort(address); err == nil {
			// Port already specified
			return server
		}

		// No port specified, but protocol-prefixed addresses usually have default ports
		// For now, return as-is and let the DNS library handle it
		return server
	}

	// For plain addresses, check if port is already specified
	if _, _, err := net.SplitHostPort(server); err == nil {
		// Port already specified
		return server
	}

	// No port specified, add default DNS port 53
	return net.JoinHostPort(server, "53")
}

// NewClientCommand creates the parent `client` command with lookup and stress subcommands.
func NewClientCommand() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "DNS client: resolve names (lookup) or load-test plain DNS (stress)",
		Subcommands: []*cli.Command{
			newClientLookupCommand(),
			newClientStressCommand(),
		},
	}
}

// plainDNSAddressForStress returns host:port for plain UDP/TCP stress; rejects DoT/DoH/DoQ URLs.
func plainDNSAddressForStress(server string) (string, error) {
	n := normalizeServerAddress(strings.TrimSpace(server))
	if strings.Contains(n, "://") {
		return "", fmt.Errorf("stress only supports plain DNS (host or host:port), got %q", server)
	}
	return n, nil
}
