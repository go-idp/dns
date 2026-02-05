# DNS Server Configuration

This directory contains example configuration files for the DNS server.

## Configuration File Format

The DNS server uses YAML format for configuration files. You can specify a configuration file using the `--config` or `-c` flag:

```bash
dns server --config /path/to/config.yaml
```

## Configuration Structure

### Basic Server Settings

```yaml
server:
  host: "0.0.0.0"    # Listen address (default: 0.0.0.0)
  port: 53           # DNS server port (default: 53)
  ttl: 500           # TTL for DNS responses in seconds (default: 500)
```

### DNS-over-TLS (DoT) Configuration

```yaml
dot:
  enabled: false     # Enable DoT server (default: false)
  port: 853         # DoT server port (default: 853)
  tls:
    cert: "/path/to/cert.pem"    # TLS certificate file (required if DoT enabled)
    key: "/path/to/key.pem"      # TLS private key file (required if DoT enabled)
```

### Custom Domain Mappings (Hosts)

Custom domain mappings have the **highest priority** and are checked before upstream DNS servers.

#### Simple Format

```yaml
hosts:
  "example.com": "1.2.3.4"
  "test.local": "192.168.1.100"
```

#### Multiple IPs

```yaml
hosts:
  "www.example.com":
    - "1.2.3.4"
    - "1.2.3.5"
```

#### IPv6 Support

```yaml
hosts:
  "ipv6.example.com": "2001:db8::1"
```

#### Both IPv4 and IPv6

```yaml
hosts:
  "dual.example.com":
    a:     # IPv4 addresses (A records)
      - "1.2.3.4"
      - "1.2.3.5"
    aaaa: # IPv6 addresses (AAAA records)
      - "2001:db8::1"
      - "2001:db8::2"
```

#### Wildcard Patterns

Wildcard patterns use `*` to match any subdomain. The pattern `*.example.com` will match:
- `sub.example.com`
- `www.example.com`
- `api.example.com`
- etc.

```yaml
hosts:
  "*.example.com": "1.2.3.4"
  "*.local.dev": "127.0.0.1"
```

**Note:** Exact matches take priority over wildcard patterns. For example, if you have both `"www.example.com": "2.3.4.5"` and `"*.example.com": "1.2.3.4"`, a query for `www.example.com` will return `2.3.4.5`.

#### Regex Patterns

Regex patterns allow advanced domain matching using regular expressions. Patterns are automatically detected if they contain regex metacharacters and compile successfully.

```yaml
hosts:
  # Match domains starting with "mp-" followed by word characters
  "^mp-\\w+\\.example\\.com$": "1.2.3.4"
  
  # Match domains with specific pattern
  "^api-v\\d+\\.example\\.com$": "1.2.3.4"
```

**Regex Pattern Rules:**
- Patterns are detected automatically if they contain regex metacharacters (`^$+?()[]{}|\`)
- Wildcard patterns (`*`) take priority over regex patterns
- Use standard Go regex syntax
- Patterns must compile successfully to be recognized as regex

**Examples:**
- `^mp-\\w+\\.example\\.com$` - Matches `mp-123.example.com`, `mp-abc.example.com`, etc.
- `^api-v\\d+\\.example\\.com$` - Matches `api-v1.example.com`, `api-v2.example.com`, etc.
- `.*\\.dev$` - Matches any domain ending with `.dev`

### Upstream DNS Servers

Upstream DNS servers are used when custom hosts don't match the query.

```yaml
upstream:
  servers:
    - "114.114.114.114:53"    # Plain DNS
    - "8.8.8.8:53"            # Google DNS
    - "tls://1.1.1.1"         # Cloudflare DoT
    - "https://dns.adguard.com/dns-query"  # DoH
  timeout: "5s"              # Query timeout (default: 5s)
```

### System Hosts File

System hosts file lookup is **enabled by default** and checked after custom hosts but before upstream DNS servers.

```yaml
system_hosts:
  disabled: false             # Disable system hosts file lookup (default: false, i.e., enabled)
  file_path: "/etc/hosts"     # Path to hosts file (default: /etc/hosts)
```

**System Hosts File Format:**

The system hosts file supports the same wildcard and regex patterns as configuration hosts:

```
# Standard format
127.0.0.1 localhost
10.1.0.169 frontend

# Wildcard pattern
1.2.3.4 *.example.com

# Regex pattern
1.2.3.4 ^mp-\w+\.example\.com
```

**Command Line Options:**

```bash
# Disable system hosts file
dns server --disable-system-hosts

# Use custom hosts file path
dns server --system-hosts-file /custom/path/hosts
```

## Priority Order

1. **Custom Hosts** (from `hosts` section) - Highest priority
   - Exact matches are checked first
   - Wildcard patterns are checked next
   - Regex patterns are checked last
2. **System Hosts File** (`/etc/hosts` by default) - Second priority
   - Supports wildcard and regex patterns
   - Enabled by default (can be disabled with `--disable-system-hosts`)
3. **Upstream DNS Servers** (from `upstream.servers`) - Fallback

## Command Line Override

Command line flags **override** configuration file values. For example:

```bash
# Use config file but override port
dns server --config config.yaml --port 5353
```

## Example Files

- `server.yaml` - Complete example with all options
- `test-server.yaml` - Minimal test configuration

## Usage Examples

### Basic Server with Custom Hosts

```yaml
server:
  port: 53

hosts:
  "local.dev": "127.0.0.1"
  "api.local.dev": "127.0.0.1"

upstream:
  servers:
    - "8.8.8.8:53"
```

### Server with DoT Support

```yaml
server:
  port: 53

dot:
  enabled: true
  port: 853
  tls:
    cert: "/etc/ssl/certs/dns.crt"
    key: "/etc/ssl/private/dns.key"

upstream:
  servers:
    - "tls://1.1.1.1"
```
