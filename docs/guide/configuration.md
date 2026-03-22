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
  # Simple format: single domain to single IP (backward compatible)
  "example.com": "1.2.3.4"

  # Alias target (CNAME-like flattening):
  # if value is a domain (not an IP), it is treated as alias target
  "mysql.ops.ys.idp.internal": "db.tencentcloud.com"

  # Explicit alias format (extension)
  "redis.ops.ys.idp.internal":
    cname: "redis.tencentcloud.com"
  
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

# Optional: in-memory cache for answers that hit upstream (not static hosts IP hits)
# cache:
#   enabled: true
#   positive_ttl: "300s"     # default if omitted when enabled
#   negative_ttl: "60s"       # empty / NXDOMAIN-style answers
#   max_entries: 10000
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

### Alias Target (CNAME-like)

Map a local domain to an upstream domain while still returning final A/AAAA IPs:

```yaml
hosts:
  # Compatible short form: string domain value
  "mysql.ops.ys.idp.internal": "db.tencentcloud.com"

  # Explicit extension form
  "redis.ops.ys.idp.internal":
    cname: "redis.tencentcloud.com"
```

Notes:
- Existing IP mapping behavior is unchanged.
- If a string value is not a valid IP, it is treated as an alias target domain.
- Responses for alias mappings are flattened A/AAAA results (not raw CNAME records).

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

1. **Custom hosts** (from config file) — static IP mappings only
2. **System hosts file** (if enabled) — static IP mappings only
3. **Response cache** (if enabled) — only for names that still need upstream; see below
4. **Custom hosts aliases** — resolve alias target via upstream
5. **System hosts aliases** — resolve alias target via upstream
6. **Upstream DNS servers**

### Response cache

When `cache.enabled: true` or `dns server --cache`:

- After static `hosts` and `/etc/hosts` **direct IP** checks, the server may return a cached answer for the same name and query type (A vs AAAA).
- **Positive cache**: at least one IP was returned; TTL defaults to **300s** unless overridden.
- **Negative cache**: empty or NXDOMAIN-style result; TTL defaults to **60s** unless overridden.
- Cached TTLs are **not** taken from upstream RR TTLs (the upstream client returns only IPs).

CLI flags `--cache-ttl`, `--cache-negative-ttl`, and `--cache-max-entries` have defaults; if you pass them explicitly, they override YAML for those fields. Use `--no-cache` to disable caching even when the config enables it.

## Examples

See `example/conf/server.yaml` for a complete example configuration file.

## Next Steps

- [Server Usage](/guide/server) - Learn how to use the server
- [Examples](/examples/config-file) - See configuration examples
