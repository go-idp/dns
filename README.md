# CLI for DNS - Simple DNS Client and Server

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-idp/dns)](https://pkg.go.dev/github.com/go-idp/dns)
[![Build Status](https://github.com/go-idp/dns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-idp/dns/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-idp/dns)](https://goreportcard.com/report/github.com/go-idp/dns)
[![Coverage Status](https://coveralls.io/repos/github/go-idp/dns/badge.svg?branch=master)](https://coveralls.io/github/go-idp/dns?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-idp/dns.svg)](https://github.com/go-idp/dns/issues)
[![Release](https://img.shields.io/github/tag/go-idp/dns.svg?label=Release)](https://github.com/go-idp/dns/tags)

## Installation

### Install CLI Tool

```bash
go install github.com/go-idp/dns/cmd/dns@latest
```

Or build from source:

```bash
git clone https://github.com/go-idp/dns.git
cd dns
go build -o bin/dns ./cmd/dns
```

## CLI Usage

### DNS Client Query
```bash
# Query A record
dns client --domain google.com --type A

# Query AAAA record (IPv6)
dns client --domain google.com --type AAAA

# Use DoT server
dns client --domain example.com --server tls://1.1.1.1

# Use DoH server
dns client --domain example.com --server https://cloudflare-dns.com/dns-query

# Use DoQ server
dns client --domain example.com --server quic://dns.adguard.com

# Use custom timeout
dns client --domain example.com --timeout 10s
```

### DNS Server
```bash
# Start basic DNS server
dns server --port 53

# Start DNS server with DoT support
dns server --port 53 --dot --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem

# Start DNS server with DoH support
dns server --port 53 --doh --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem

# Start DNS server with DoQ support
dns server --port 53 --doq --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem

# Start DNS server with all protocols (DoT, DoH, DoQ)
dns server --port 53 \
  --dot --dot-port 853 \
  --doh --doh-port 443 \
  --doq --doq-port 853 \
  --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem

# Start DNS server with custom upstream
dns server --port 53 --upstream 8.8.8.8:53 --upstream 1.1.1.1:53

# Start DNS server with configuration file
dns server --config /path/to/config.yaml

# Command line flags override config file values
dns server --config /path/to/config.yaml --port 5353
```

### Configuration File

The server supports YAML configuration files for easier management. See `example/conf/server.yaml` for a complete example.

**Configuration File Structure:**

```yaml
# Basic server settings
server:
  host: "0.0.0.0"
  port: 53
  ttl: 500

# DNS-over-TLS (DoT) configuration
dot:
  enabled: false
  port: 853
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

# DNS-over-HTTPS (DoH) configuration
doh:
  enabled: false
  port: 443
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

# DNS-over-QUIC (DoQ) configuration
doq:
  enabled: false
  port: 853
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

# Custom domain mappings (highest priority)
hosts:
  # Simple format: single domain to single IP
  "example.com": "1.2.3.4"
  "www.example.com":
    - "1.2.3.4"
    - "1.2.3.5"
  "dual.example.com":
    a: ["1.2.3.4"]
    aaaa: ["2001:db8::1"]
  
  # Wildcard pattern: matches any subdomain
  "*.example.com": "1.2.3.4"
  
  # Regex pattern: matches domains using regular expressions
  "^mp-\\w+\\.example\\.com$": "1.2.3.4"

# Upstream DNS servers
upstream:
  servers:
    - "114.114.114.114:53"
    - "tls://1.1.1.1"
  timeout: "5s"
```

**Key Features:**
- **Custom Hosts Mapping**: Define custom domain-to-IP mappings with highest priority
- **Multiple IP Support**: Support multiple IPv4 and IPv6 addresses per domain
- **Flexible Format**: Support simple string, list, or structured format
- **Wildcard Patterns**: Use `*.example.com` to match any subdomain
- **Regex Patterns**: Use regular expressions like `^mp-\\w+\\.example\\.com$` for advanced matching
- **System Hosts File**: Support for `/etc/hosts` with wildcard and regex patterns (enabled by default)
- **Priority**: Custom hosts are checked before system hosts and upstream DNS servers
- **Override**: Command line flags override config file values

## Getting Started

### Quick Start

After installation, you can start using the DNS CLI:

```bash
# Query a domain
dns client --domain google.com

# Start a DNS server
dns server --port 53
```

See the [documentation](https://go-idp.github.io/dns/) for more examples and detailed usage.

## Features

### Client
* [x] Plain DNS
	* [x] Plain DNS in UDP
	* [x] Plain DNS in TCP
* [x] DNS-over-TLS (DoT) - Use `tls://` prefix (e.g., `tls://1.1.1.1`)
* [x] DNS-over-HTTPS (DoH)
* [x] DNS-over-QUIC (DoQ)
* [x] DNSCrypt

### Server
* [x] Plain DNS
	* [x] Plain DNS in UDP
	* [x] Plain DNS in TCP
* [x] DNS-over-TLS (DoT)
* [x] DNS-over-HTTPS (DoH)
* [x] DNS-over-QUIC (DoQ)

## Inspired By
* [AdGuardHome](https://github.com/AdguardTeam/AdGuardHome) - Network-wide ads & trackers blocking DNS server.
* [kenshinx/godns](https://github.com/kenshinx/godns) - A fast dns cache server written by go.
* [miekg/dns](https://github.com/miekg/dns) - DNS library in Go.

## Documentation

Full documentation is available at: https://go-idp.github.io/dns/

## License

MIT License - see [LICENSE](./LICENSE) for details.
