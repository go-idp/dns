# DNS Server

The DNS server allows you to run your own DNS server with custom host mappings and upstream DNS fallback.

## Basic Usage

### Start Basic Server

Start a DNS server on the default port 53:

```bash
dns server --port 53
```

### Start with Custom Port

```bash
dns server --port 5353
```

## Options

### `--port`

DNS server port. Default: 53.

```bash
dns server --port 5353
```

### `--host`

Listen address. Default: `0.0.0.0`.

```bash
dns server --host 127.0.0.1 --port 53
```

### `--upstream`

Upstream DNS server. Can be specified multiple times. Supports:
- Plain DNS: `8.8.8.8:53`
- DNS-over-TLS: `tls://1.1.1.1`
- DNS-over-HTTPS: `https://dns.adguard.com/dns-query`

```bash
dns server --port 53 --upstream 8.8.8.8:53 --upstream tls://1.1.1.1
```

### `--config`

Path to configuration file. See [Configuration](/guide/configuration) for details.

```bash
dns server --config /path/to/config.yaml
```

### `--dot`

Enable DNS-over-TLS (DoT) server.

```bash
dns server --port 53 --dot --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem
```

### `--tls-cert` / `--tls-key`

TLS certificate and key files for DoT, DoH, and DoQ servers.

```bash
dns server --dot --tls-cert cert.pem --tls-key key.pem
```

### `--dot-port`

DoT server port. Default: 853.

```bash
dns server --dot --dot-port 853 --tls-cert cert.pem --tls-key key.pem
```

### `--doh`

Enable DNS-over-HTTPS (DoH) server.

```bash
dns server --port 53 --doh --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem
```

### `--doh-port`

DoH server port. Default: 443.

```bash
dns server --doh --doh-port 443 --tls-cert cert.pem --tls-key key.pem
```

### `--doq`

Enable DNS-over-QUIC (DoQ) server.

```bash
dns server --port 53 --doq --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem
```

### `--doq-port`

DoQ server port. Default: 853.

```bash
dns server --doq --doq-port 853 --tls-cert cert.pem --tls-key key.pem
```

### `--ttl`

TTL for DNS responses in seconds. Default: 500.

```bash
dns server --port 53 --ttl 300
```

### `--system-hosts-disabled`

Disable system hosts file lookup.

```bash
dns server --port 53 --system-hosts-disabled
```

### `--system-hosts-file`

Path to system hosts file. Default: `/etc/hosts`.

```bash
dns server --port 53 --system-hosts-file /custom/hosts
```

## Command Line Flags Override Config File

Command line flags take precedence over configuration file values:

```bash
# config.yaml has port: 53, but command line overrides it
dns server --config config.yaml --port 5353
```

## Examples

### Basic Server with Upstream

```bash
dns server --port 53 --upstream 8.8.8.8:53
```

### Server with DoT

```bash
dns server --port 53 --dot --tls-cert cert.pem --tls-key key.pem
```

### Server with DoH

```bash
dns server --port 53 --doh --tls-cert cert.pem --tls-key key.pem
```

### Server with DoQ

```bash
dns server --port 53 --doq --tls-cert cert.pem --tls-key key.pem
```

### Server with Multiple Protocols

```bash
# Enable all protocols (DoT, DoH, DoQ)
dns server --port 53 \
  --dot --dot-port 853 \
  --doh --doh-port 443 \
  --doq --doq-port 853 \
  --tls-cert cert.pem --tls-key key.pem
```

### Server with Configuration File

```bash
dns server --config config.yaml
```

## Next Steps

- [Configuration](/guide/configuration) - Learn about configuration options
- [DNS-over-TLS](/guide/dot) - Learn about DoT server setup
- [Examples](/examples/) - See more examples
- [DoH and DoQ Client Examples](/examples/doh-doq-client) - Learn about DoH and DoQ client usage
