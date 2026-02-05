# DoH and DoQ Server Examples

This document provides examples for running DNS-over-HTTPS (DoH) and DNS-over-QUIC (DoQ) servers using the CLI.

## Prerequisites

You need a TLS certificate and key file for both DoH and DoQ servers. See [DoT Guide](/guide/dot) for certificate generation instructions.

## DNS-over-HTTPS (DoH) Server

DoH provides DNS queries over HTTPS, offering privacy and security while working through firewalls and proxies.

### Start DoH Server

#### Using Command Line

```bash
dns server --port 53 --doh --tls-cert cert.pem --tls-key key.pem
```

#### Using Configuration File

Create `config.yaml`:

```yaml
server:
  port: 53

doh:
  enabled: true
  port: 443
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

hosts:
  "example.com": "1.2.3.4"
```

Start the server:

```bash
dns server --config config.yaml
```

### Custom DoH Port

```bash
dns server --port 53 --doh --doh-port 8443 --tls-cert cert.pem --tls-key key.pem
```

Or in config file:

```yaml
doh:
  enabled: true
  port: 8443  # Custom DoH port
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"
```

### Testing DoH Server

```bash
# Query using DoH
dns client --domain example.com --server https://localhost:443/dns-query
```

## DNS-over-QUIC (DoQ) Server

DoQ provides DNS queries over QUIC protocol, offering the lowest latency with encryption.

### Start DoQ Server

#### Using Command Line

```bash
dns server --port 53 --doq --tls-cert cert.pem --tls-key key.pem
```

#### Using Configuration File

Create `config.yaml`:

```yaml
server:
  port: 53

doq:
  enabled: true
  port: 853
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

hosts:
  "example.com": "1.2.3.4"
```

Start the server:

```bash
dns server --config config.yaml
```

### Custom DoQ Port

```bash
dns server --port 53 --doq --doq-port 8853 --tls-cert cert.pem --tls-key key.pem
```

Or in config file:

```yaml
doq:
  enabled: true
  port: 8853  # Custom DoQ port
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"
```

### Testing DoQ Server

```bash
# Query using DoQ
dns client --domain example.com --server quic://localhost:853
```

## Running Multiple Protocols

You can run DoT, DoH, and DoQ servers simultaneously:

### Using Command Line

```bash
dns server --port 53 \
  --dot --dot-port 853 \
  --doh --doh-port 443 \
  --doq --doq-port 853 \
  --tls-cert cert.pem --tls-key key.pem
```

### Using Configuration File

```yaml
server:
  port: 53

dot:
  enabled: true
  port: 853
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

doh:
  enabled: true
  port: 443
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

doq:
  enabled: true
  port: 853
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

hosts:
  "example.com": "1.2.3.4"
```

## Protocol Comparison

| Feature | DoT | DoH | DoQ |
|---------|-----|-----|-----|
| Encryption | ✅ | ✅ | ✅ |
| Default Port | 853 | 443 | 853 |
| Latency | Low | Medium | Very Low |
| Firewall Friendly | ⚠️ | ✅ | ⚠️ |
| Connection Multiplexing | ❌ | ❌ | ✅ |

## Security Considerations

1. **Certificate Validation**: Ensure your certificate is valid and trusted
2. **Firewall**: Open the appropriate ports in your firewall:
   - DoT: 853 (default)
   - DoH: 443 (default)
   - DoQ: 853 (default)
3. **TLS Configuration**: Use strong TLS certificates and keep them secure

## Next Steps

- [Server Guide](/guide/server) - Learn more about the DNS server
- [DoT Server Example](/examples/dot-server) - Learn about DoT server setup
- [Configuration Guide](/guide/configuration) - Learn about configuration options
