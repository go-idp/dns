# DoH and DoQ Client Examples

This document provides examples for using DNS-over-HTTPS (DoH) and DNS-over-QUIC (DoQ) protocols with the DNS client.

## DNS-over-HTTPS (DoH)

DoH provides DNS queries over HTTPS, offering privacy and security while working through firewalls and proxies.

### Basic DoH Query

```bash
# Query using Cloudflare DoH
dns client --domain google.com --server https://cloudflare-dns.com/dns-query

# Query using AdGuard DoH
dns client --domain example.com --server https://dns.adguard.com/dns-query

# Query using Google DoH
dns client --domain example.com --server https://dns.google/dns-query
```

### DoH with IPv6 Query

```bash
# Query AAAA records using DoH
dns client --domain google.com --type AAAA --server https://cloudflare-dns.com/dns-query
```

### DoH with Custom Timeout

```bash
# Use longer timeout for slow networks
dns client --domain example.com \
  --server https://dns.adguard.com/dns-query \
  --timeout 15s
```

### DoH with Plain Output

```bash
# Get only IP addresses for scripting
dns client --domain example.com \
  --server https://cloudflare-dns.com/dns-query \
  --plain
```

## DNS-over-QUIC (DoQ)

DoQ provides DNS queries over QUIC protocol, offering the lowest latency with encryption.

### Basic DoQ Query

```bash
# Query using AdGuard DoQ
dns client --domain google.com --server quic://dns.adguard.com
```

### DoQ with IPv6 Query

```bash
# Query AAAA records using DoQ
dns client --domain google.com --type AAAA --server quic://dns.adguard.com
```

### DoQ with Custom Timeout

```bash
# Use custom timeout
dns client --domain example.com \
  --server quic://dns.adguard.com \
  --timeout 10s
```

## Using Multiple Protocols

You can specify multiple servers with different protocols. The client will try them in order:

```bash
# Try multiple protocols for redundancy
dns client --domain example.com \
  --server 8.8.8.8 \
  --server tls://1.1.1.1 \
  --server https://cloudflare-dns.com/dns-query \
  --server quic://dns.adguard.com
```

## Environment Variables

You can set default DoH or DoQ servers using environment variables:

```bash
# Set DoH as default
export DNS_SERVER=https://cloudflare-dns.com/dns-query
dns client --domain example.com

# Set DoQ as default
export DNS_SERVER=quic://dns.adguard.com
dns client --domain example.com
```

## Script Examples

### Bash Script with DoH

```bash
#!/bin/bash
# Get IP address using DoH
IP=$(dns client --domain example.com \
  --server https://cloudflare-dns.com/dns-query \
  --plain | head -n 1)

echo "IP address: $IP"
```

### Bash Script with DoQ

```bash
#!/bin/bash
# Get IP address using DoQ
IP=$(dns client --domain example.com \
  --server quic://dns.adguard.com \
  --plain | head -n 1)

echo "IP address: $IP"
```

## Popular DoH Servers

- **Cloudflare**: `https://cloudflare-dns.com/dns-query`
- **AdGuard**: `https://dns.adguard.com/dns-query`
- **Google**: `https://dns.google/dns-query`
- **Quad9**: `https://dns.quad9.net/dns-query`

## Popular DoQ Servers

- **AdGuard**: `quic://dns.adguard.com`

## Protocol Comparison

| Feature | DoH | DoQ |
|---------|-----|-----|
| Encryption | ✅ | ✅ |
| Port | 443 | 853 |
| Latency | Medium | Very Low |
| Firewall Friendly | ✅ | ⚠️ |
| Connection Multiplexing | ❌ | ✅ |

## Next Steps

- [Client Usage Guide](/guide/client) - Learn more about the DNS client
- [Supported Protocols](/guide/client-protocols) - See all supported protocols
- [Quick Start](/guide/quick-start) - Get started quickly
