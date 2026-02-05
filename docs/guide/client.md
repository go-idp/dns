# DNS Client

The DNS client allows you to query DNS servers using various protocols.

## Basic Usage

### Query A Record

```bash
dns client --domain google.com --type A
```

### Query AAAA Record

```bash
dns client --domain google.com --type AAAA
```

## Options

### `--domain` / `-d`

The domain name to query.

```bash
dns client --domain example.com
```

### `--type` / `-t`

The query type. Supported types:
- `A` - IPv4 address (default)
- `AAAA` - IPv6 address

```bash
dns client --domain example.com --type AAAA
```

### `--server` / `-s`

DNS server address. Can be specified multiple times. Supports:
- Plain DNS: `8.8.8.8:53` or `8.8.8.8` (default port 53)
- DNS-over-TLS: `tls://1.1.1.1` or `tls://1.1.1.1:853`
- DNS-over-HTTPS: `https://dns.adguard.com/dns-query`
- DNS-over-QUIC: `quic://dns.adguard.com`
- DNSCrypt: `sdns://...`

```bash
# Use DoT server
dns client --domain example.com --server tls://1.1.1.1

# Use multiple servers
dns client --domain example.com --server 8.8.8.8 --server tls://1.1.1.1
```

### `--timeout`

Query timeout. Default: `5s`.

```bash
dns client --domain example.com --timeout 10s
```

### `--plain`

Output only IP addresses, one per line. Useful for scripting.

```bash
dns client --domain google.com --plain
```

## Environment Variables

- `DNS_SERVER` - Default DNS server
- `DNS_TIMEOUT` - Default timeout
- `DNS_PLAIN` - Enable plain output mode

```bash
export DNS_SERVER=tls://1.1.1.1
export DNS_TIMEOUT=10s
dns client --domain example.com
```

## Examples

### Query with DoT

```bash
dns client --domain example.com --server tls://1.1.1.1
```

### Query with DoH

```bash
dns client --domain example.com --server https://dns.adguard.com/dns-query
```

### Query with DoQ

```bash
dns client --domain example.com --server quic://dns.adguard.com
```

### Query with Multiple Protocols

```bash
# Use multiple servers with different protocols
dns client --domain example.com \
  --server 8.8.8.8 \
  --server tls://1.1.1.1 \
  --server https://cloudflare-dns.com/dns-query \
  --server quic://dns.adguard.com
```

### Script-friendly Output

```bash
#!/bin/bash
IP=$(dns client --domain example.com --plain | head -n 1)
echo "IP address: $IP"
```

## Next Steps

- [Supported Protocols](/guide/client-protocols) - Learn about all supported DNS protocols
- [Server Usage](/guide/server) - Learn how to run a DNS server
