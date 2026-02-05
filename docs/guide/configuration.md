# Configuration

DNS server supports YAML configuration files for easier management. Command-line flags override config file values.

## Configuration File Structure

```yaml
# Basic server settings
server:
  host: "0.0.0.0"        # Listen address (default: 0.0.0.0)
  port: 53               # DNS server port (default: 53)
  ttl: 500               # TTL for DNS responses in seconds (default: 500)

# DNS-over-TLS (DoT) configuration
dot:
  enabled: false         # Enable DoT server (default: false)
  port: 853             # DoT server port (default: 853)
  tls:
    cert: "/path/to/cert.pem"    # TLS certificate file path
    key: "/path/to/key.pem"      # TLS private key file path

# DNS-over-HTTPS (DoH) configuration
doh:
  enabled: false         # Enable DoH server (default: false)
  port: 443             # DoH server port (default: 443)
  tls:
    cert: "/path/to/cert.pem"    # TLS certificate file path
    key: "/path/to/key.pem"      # TLS private key file path

# DNS-over-QUIC (DoQ) configuration
doq:
  enabled: false         # Enable DoQ server (default: false)
  port: 853             # DoQ server port (default: 853)
  tls:
    cert: "/path/to/cert.pem"    # TLS certificate file path
    key: "/path/to/key.pem"      # TLS private key file path

# Custom domain mappings (highest priority)
hosts:
  # Simple format: single domain to single IP
  "example.com": "1.2.3.4"
  
  # Multiple IPv4 addresses
  "www.example.com":
    - "1.2.3.4"
    - "1.2.3.5"
  
  # Both IPv4 and IPv6 addresses
  "dual.example.com":
    a:     # IPv4 addresses (A records)
      - "1.2.3.4"
      - "1.2.3.5"
    aaaa: # IPv6 addresses (AAAA records)
      - "2001:db8::1"
      - "2001:db8::2"
  
  # Wildcard pattern (matches any subdomain)
  "*.example.com": "1.2.3.4"
  
  # Regex pattern
  "^mp-\\w+\\.example\\.com$": "1.2.3.4"

# System hosts file configuration
system_hosts:
  disabled: false             # Disable system hosts file lookup (default: false)
  file_path: "/etc/hosts"     # Path to hosts file (default: /etc/hosts)

# Upstream DNS servers
upstream:
  servers:
    - "114.114.114.114:53"    # Plain DNS
    - "tls://1.1.1.1"         # Cloudflare DoT
    - "https://dns.adguard.com/dns-query"  # DoH
  timeout: "5s"              # Query timeout (default: 5s)
```

## Host Mappings

### Simple Format

Single domain to single IP:

```yaml
hosts:
  "example.com": "1.2.3.4"
```

### Multiple IPs

Multiple IPv4 addresses:

```yaml
hosts:
  "www.example.com":
    - "1.2.3.4"
    - "1.2.3.5"
```

### Structured Format

Both IPv4 and IPv6:

```yaml
hosts:
  "dual.example.com":
    a:     # IPv4 addresses
      - "1.2.3.4"
      - "1.2.3.5"
    aaaa: # IPv6 addresses
      - "2001:db8::1"
      - "2001:db8::2"
```

### Wildcard Patterns

Match any subdomain:

```yaml
hosts:
  "*.example.com": "1.2.3.4"
```

This matches:
- `test.example.com`
- `api.example.com`
- `www.example.com`
- But NOT `example.com` itself

### Regex Patterns

Advanced pattern matching:

```yaml
hosts:
  "^mp-\\w+\\.example\\.com$": "1.2.3.4"
```

This matches domains like:
- `mp-api.example.com`
- `mp-frontend.example.com`
- But NOT `mp.example.com` (needs word characters after `mp-`)

## Priority Order

DNS resolution follows this priority order:

1. **Custom hosts** (from config file or command line)
2. **System hosts file** (if enabled)
3. **Upstream DNS servers**

## Examples

See `example/conf/server.yaml` for a complete example configuration file.

## Next Steps

- [Server Usage](/guide/server) - Learn how to use the server
- [Examples](/examples/config-file) - See configuration examples
