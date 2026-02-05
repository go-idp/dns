# Introduction

DNS CLI is a powerful DNS client and server tool written in Go. It provides both command-line interface for DNS queries and a DNS server with support for custom host mappings and DNS-over-TLS (DoT).

## What is DNS CLI?

DNS CLI is a command-line tool for DNS operations:

- **DNS Client**: Query DNS servers using various protocols (Plain DNS, DoT, DoH, DoQ, DNSCrypt)
- **DNS Server**: Run your own DNS server with custom host mappings, wildcard patterns, and regex support
- **DNS-over-TLS (DoT)**: Secure DNS queries and server support
- **Flexible Configuration**: YAML-based configuration with command-line flag overrides

## Key Features

### For DNS Client
- Support for multiple DNS protocols
- Configurable timeout
- Plain output mode for scripting
- IPv4 and IPv6 support

### For DNS Server
- Custom host mappings with highest priority
- Wildcard pattern matching (`*.example.com`)
- Regex pattern matching (`^mp-\\w+\\.example\\.com$`)
- System hosts file integration
- Upstream DNS server fallback
- DNS-over-TLS (DoT) support
- YAML configuration file support

## Use Cases

- **Development**: Local development with custom domain mappings
- **Testing**: DNS testing and debugging
- **Privacy**: Use secure DNS protocols (DoT, DoH)
- **Custom DNS**: Run your own DNS server with custom rules
- **Network Management**: Manage DNS resolution in your network

## Next Steps

- [Installation](/guide/installation) - Learn how to install DNS CLI
- [Quick Start](/guide/quick-start) - Get started in minutes
- [Client Usage](/guide/client) - Learn how to use the DNS client
- [Server Usage](/guide/server) - Learn how to run a DNS server
