# DNS-over-QUIC (DoQ)

DNS-over-QUIC (DoQ) provides encrypted DNS queries using QUIC protocol.

## Overview

DoQ encrypts DNS queries between the client and server using QUIC, providing:
- **Privacy**: DNS queries are encrypted
- **Security**: Protection against DNS spoofing
- **Performance**: Lowest latency among encrypted DNS protocols
- **Connection Multiplexing**: Multiple queries over a single connection
- **Improved Error Handling**: Better handling of network errors

## Server Setup

### Generate TLS Certificate

You need a TLS certificate and key file to run a DoQ server.

#### Using OpenSSL

```bash
# Generate private key
openssl genrsa -out key.pem 2048

# Generate certificate
openssl req -new -x509 -key key.pem -out cert.pem -days 365
```

#### Using Let's Encrypt

```bash
# Install certbot
sudo apt-get install certbot

# Get certificate
sudo certbot certonly --standalone -d your-domain.com
```

### Start DoQ Server

#### Using Command Line

```bash
dns server --port 53 --doq --tls-cert cert.pem --tls-key key.pem
```

#### Using Configuration File

```yaml
server:
  port: 53

doq:
  enabled: true
  port: 853
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"
```

Start the server:

```bash
dns server --config config.yaml
```

### Custom DoQ Port

Default DoQ port is 853. You can change it:

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

## Client Usage

### Query Using DoQ

```bash
dns client --domain example.com --server quic://your-doq-server:853
```

### Popular DoQ Servers

- AdGuard: `quic://dns.adguard.com`
- Cloudflare: `quic://cloudflare-dns.com` (if supported)

## Testing

### Test DoQ Server

```bash
# Query your DoQ server
dns client --domain example.com --server quic://localhost:853
```

### Verify Certificate

```bash
openssl s_client -connect localhost:853 -servername localhost
```

## Protocol Details

DoQ uses the QUIC transport protocol (RFC 9250) to send DNS queries:
- **Transport**: QUIC over UDP
- **Port**: 853 (default, as specified in RFC 9250)
- **Encryption**: Built into QUIC protocol
- **Connection Multiplexing**: Multiple streams per connection

## Advantages

1. **Low Latency**: Fastest among encrypted DNS protocols
2. **Connection Multiplexing**: Multiple queries over a single connection
3. **Improved Error Handling**: Better handling of network errors and packet loss
4. **Zero RTT**: Connection establishment with 0-RTT for faster queries
5. **Migration Support**: Connection migration for mobile devices

## Limitations

1. **Firewall**: May be blocked by some firewalls (uses UDP)
2. **NAT Traversal**: May have issues with some NAT configurations
3. **Adoption**: Less widely adopted compared to DoT and DoH
4. **Port Conflicts**: If port 853 is already in use, you need to use a different port

## Security Considerations

1. **Certificate Validation**: Ensure your certificate is valid and trusted
2. **Firewall**: Open port 853 in your firewall if needed
3. **Certificate Renewal**: Set up automatic certificate renewal for production
4. **QUIC Security**: QUIC provides built-in encryption and authentication

## Comparison with Other Protocols

| Feature | DoQ | DoT | DoH |
|---------|-----|-----|-----|
| Port | 853 | 853 | 443 |
| Latency | Very Low | Low | Medium |
| Connection Multiplexing | ✅ | ❌ | ❌ |
| Firewall Friendly | ⚠️ | ⚠️ | ✅ |
| Zero RTT | ✅ | ❌ | ❌ |
| Transport | UDP (QUIC) | TCP (TLS) | TCP (HTTPS) |

## Use Cases

DoQ is ideal for:
- **Performance-Critical Applications**: Where low latency is essential
- **Mobile Networks**: Better handling of network changes
- **High Query Volume**: Connection multiplexing reduces overhead
- **Unstable Networks**: Better error recovery than TCP-based protocols

## Next Steps

- [Server Usage](/guide/server) - Learn more about DNS server
- [Client Usage](/guide/client) - Learn how to use the DNS client
- [Configuration](/guide/configuration) - Learn about configuration options
- [DoH and DoQ Examples](/examples/doh-doq-server) - See practical examples
