# Quick Start

This guide will help you get started with DNS CLI in just a few minutes.

## DNS Client

### Basic Query

Query a domain's A record:

```bash
dns client --domain google.com --type A
```

### Query IPv6

Query a domain's AAAA record (IPv6):

```bash
dns client --domain google.com --type AAAA
```

### Use DoT Server

Query using DNS-over-TLS:

```bash
dns client --domain example.com --server tls://1.1.1.1
```

### Plain Output

Get only IP addresses (useful for scripting):

```bash
dns client --domain google.com --plain
```

## DNS Server

### Start Basic Server

Start a DNS server on port 53:

```bash
dns server --port 53
```

### Start with Custom Upstream

Start a DNS server with custom upstream DNS servers:

```bash
dns server --port 53 --upstream 8.8.8.8:53 --upstream 1.1.1.1:53
```

### Start with Configuration File

Create a configuration file `config.yaml`:

```yaml
server:
  port: 53

hosts:
  "example.com": "1.2.3.4"
  "*.example.com": "1.2.3.4"

upstream:
  servers:
    - "8.8.8.8:53"
```

Start the server:

```bash
dns server --config config.yaml
```

## Next Steps

- [Client Usage](/guide/client) - Learn more about DNS client features
- [Server Usage](/guide/server) - Learn more about DNS server features
- [Configuration](/guide/configuration) - Learn about configuration options
