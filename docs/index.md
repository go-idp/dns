---
layout: home

hero:
  name: DNS CLI
  text: Simple DNS Client and Server
  tagline: A powerful DNS client and server CLI tool written in Go
  # image:
  #   src: /logo.svg
  #   alt: DNS CLI
  actions:
    - theme: brand
      text: Get Started
      link: /guide/
    - theme: alt
      text: View on GitHub
      link: https://github.com/go-idp/dns

features:
  - icon: 🚀
    title: Fast & Lightweight
    details: Built with Go for high performance and low memory footprint
  - icon: 🔒
    title: Secure Protocols
    details: Support for DoT, DoH, DoQ, and DNSCrypt protocols
  - icon: ⚙️
    title: Flexible Configuration
    details: YAML configuration files with wildcard and regex pattern support
  - icon: 🌐
    title: Multiple Protocols
    details: Plain DNS, DNS-over-TLS, DNS-over-HTTPS, DNS-over-QUIC, and DNSCrypt
  - icon: 📝
    title: Custom Hosts
    details: Custom domain mappings with highest priority, supporting wildcards and regex
  - icon: 🐳
    title: Docker Ready
    details: Pre-built Docker images available for easy deployment
---

## Quick Start

### Installation

```bash
go install github.com/go-idp/dns/cmd/dns@latest
```

### DNS Client Query

```bash
# Query A record
dns client lookup google.com --type A

# Query AAAA record (IPv6)
dns client lookup google.com --type AAAA

# Use DoT server
dns client lookup example.com --server tls://1.1.1.1
```

### DNS Server

```bash
# Start basic DNS server
dns server --port 53

# Start DNS server with DoT support
dns server --port 53 --dot --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem

# Start DNS server with configuration file
dns server --config /path/to/config.yaml
```

Visit the [Guide](/guide/) for detailed usage instructions.

## Features

### Client Features
- ✅ Plain DNS (UDP/TCP)
- ✅ DNS-over-TLS (DoT)
- ✅ DNS-over-HTTPS (DoH)
- ✅ DNS-over-QUIC (DoQ)
- ✅ DNSCrypt

### Server Features
- ✅ Plain DNS (UDP/TCP)
- ✅ DNS-over-TLS (DoT)
- ✅ DNS-over-HTTPS (DoH)
- ✅ DNS-over-QUIC (DoQ)
- ✅ Custom host mappings
- ✅ Wildcard and regex pattern support
- ✅ System hosts file integration
- ✅ Upstream DNS server fallback
- ✅ YAML configuration file support

## License

MIT License - see [LICENSE](https://github.com/go-idp/dns/blob/master/LICENSE) for details.
