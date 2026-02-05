# Supported DNS Protocols

DNS CLI client supports multiple DNS protocols for secure and flexible DNS queries.

## Plain DNS

The standard DNS protocol over UDP or TCP.

### Usage

```bash
# Default port 53
dns client --domain example.com --server 8.8.8.8

# Custom port
dns client --domain example.com --server 8.8.8.8:5353
```

## DNS-over-TLS (DoT)

DNS queries encrypted using TLS. Default port: 853.

### Usage

```bash
# Cloudflare DoT
dns client --domain example.com --server tls://1.1.1.1

# Custom port
dns client --domain example.com --server tls://1.1.1.1:853
```

### Popular DoT Servers

- Cloudflare: `tls://1.1.1.1`
- Google: `tls://8.8.8.8`
- Quad9: `tls://9.9.9.9`

## DNS-over-HTTPS (DoH)

DNS queries over HTTPS. Provides privacy and security.

### Usage

```bash
# AdGuard DoH
dns client --domain example.com --server https://dns.adguard.com/dns-query

# Cloudflare DoH
dns client --domain example.com --server https://cloudflare-dns.com/dns-query
```

### Popular DoH Servers

- AdGuard: `https://dns.adguard.com/dns-query`
- Cloudflare: `https://cloudflare-dns.com/dns-query`
- Google: `https://dns.google/dns-query`

## DNS-over-QUIC (DoQ)

DNS queries over QUIC protocol. Provides low latency and security. DoQ uses the QUIC transport protocol, which offers connection multiplexing, improved error handling, and reduced latency compared to DoT.

### Usage

```bash
# AdGuard DoQ
dns client --domain example.com --server quic://dns.adguard.com

# DoQ with custom port
dns client --domain example.com --server quic://dns.adguard.com:853
```

### Popular DoQ Servers

- AdGuard: `quic://dns.adguard.com`
- Cloudflare: `quic://cloudflare-dns.com` (if supported)

## DNSCrypt

DNS encryption protocol using DNSCrypt.

### Usage

```bash
dns client --domain example.com --server sdns://AQcAAAAAAAAAAAAQ1syXTkwLjE3OC4xNzIuMTc4XQ
```

## Protocol Comparison

| Protocol | Encryption | Port | Latency | Privacy |
|----------|-----------|------|---------|---------|
| Plain DNS | ❌ | 53 | Low | ❌ |
| DoT | ✅ | 853 | Low | ✅ |
| DoH | ✅ | 443 | Medium | ✅ |
| DoQ | ✅ | 853 | Very Low | ✅ |
| DNSCrypt | ✅ | 443/5353 | Low | ✅ |

## Choosing a Protocol

- **Plain DNS**: Fastest, but no encryption. Use for local networks.
- **DoT**: Good balance of security and performance. Recommended for most use cases.
- **DoH**: Works through firewalls and proxies. Good for restricted networks.
- **DoQ**: Lowest latency with encryption. Best for performance-critical applications.
- **DNSCrypt**: Alternative encryption protocol. Good compatibility.

## Next Steps

- [Client Usage](/guide/client) - Learn more about using the DNS client
- [DNS-over-TLS (DoT)](/guide/dot) - Detailed DoT guide
- [DNS-over-HTTPS (DoH)](/guide/doh) - Detailed DoH guide
- [DNS-over-QUIC (DoQ)](/guide/doq) - Detailed DoQ guide
- [Server Usage](/guide/server) - Learn how to run a DNS server
