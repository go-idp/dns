# DNS Client

The `dns client` command has subcommands:

- **`lookup`** — query DNS servers (A/AAAA) using plain DNS, DoT, DoH, DoQ, or DNSCrypt.
- **`stress`** — concurrent plain DNS (UDP/TCP) load test against a **single** server (`host` or `host:port` only).

## Basic Usage (`lookup`)

### Query A Record

```bash
dns client lookup google.com --type A
```

### Query AAAA Record

```bash
dns client lookup google.com --type AAAA
```

## Options

### Domain (positional argument)

Pass the name to resolve as the first argument after `lookup`:

```bash
dns client lookup baidu.com
dns client lookup baidu.com --server 223.5.5.5
```

You can still use **`--domain` / `-d`** instead of a positional argument if you prefer (for example in scripts).

### `--type` / `-t`

The query type. Supported types:
- `A` - IPv4 address (default)
- `AAAA` - IPv6 address

```bash
dns client lookup example.com --type AAAA
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
dns client lookup example.com --server tls://1.1.1.1

# Use multiple servers
dns client lookup example.com --server 8.8.8.8 --server tls://1.1.1.1
```

### `--timeout`

Query timeout. Default: `5s`.

```bash
dns client lookup example.com --timeout 10s
```

### `--plain`

Output only IP addresses, one per line. Useful for scripting.

```bash
dns client lookup google.com --plain
```

## Environment Variables

- `DNS_SERVER` - Default DNS server
- `DNS_TIMEOUT` - Default timeout
- `DNS_PLAIN` - Enable plain output mode

```bash
export DNS_SERVER=tls://1.1.1.1
export DNS_TIMEOUT=10s
dns client lookup example.com
```

## Examples

### Query with DoT

```bash
dns client lookup example.com --server tls://1.1.1.1
```

### Query with DoH

```bash
dns client lookup example.com --server https://dns.adguard.com/dns-query
```

### Query with DoQ

```bash
dns client lookup example.com --server quic://dns.adguard.com
```

### Query with Multiple Protocols

```bash
# Use multiple servers with different protocols
dns client lookup example.com \
  --server 8.8.8.8 \
  --server tls://1.1.1.1 \
  --server https://cloudflare-dns.com/dns-query \
  --server quic://dns.adguard.com
```

### Script-friendly Output

```bash
#!/bin/bash
IP=$(dns client lookup example.com --plain | head -n 1)
echo "IP address: $IP"
```

## Load test (`stress`)

Use this to benchmark a plain DNS listener (for example your `dns server`). It does **not** support `tls://`, `https://`, or `quic://` — only UDP or TCP to `host:port`.

```bash
# 200 workers, 5000 queries over UDP (default)
dns client stress --domain example.com --server 127.0.0.1:5353 --workers 200 --requests 5000

# TCP, custom per-query timeout (`-n` is shorthand for `--requests`)
dns client stress --domain example.com --server 8.8.8.8 --net tcp --timeout 3s --workers 50 -n 500

# Treat NXDOMAIN as success (useful for names you expect not to exist)
dns client stress --domain definitely-missing.example --server 127.0.0.1:53 --accept-nxdomain
```

## Next Steps

- [Supported Protocols](/guide/client-protocols) - Learn about all supported DNS protocols
- [DNS-over-TLS (DoT)](/guide/dot) - Learn about DoT protocol
- [DNS-over-HTTPS (DoH)](/guide/doh) - Learn about DoH protocol
- [DNS-over-QUIC (DoQ)](/guide/doq) - Learn about DoQ protocol
- [Server Usage](/guide/server) - Learn how to run a DNS server
