# DNS-over-HTTPS (DoH)

DNS-over-HTTPS (DoH) provides encrypted DNS queries over HTTPS protocol.

## Overview

DoH encrypts DNS queries between the client and server using HTTPS, providing:
- **Privacy**: DNS queries are encrypted
- **Security**: Protection against DNS spoofing
- **Firewall Friendly**: Works through firewalls and proxies (uses port 443)
- **Compatibility**: Uses standard HTTPS, making it easy to deploy

## Server Setup

### Generate TLS Certificate

You need a TLS certificate and key file to run a DoH server.

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

### Start DoH Server

#### Using Command Line

```bash
dns server --port 53 --doh --tls-cert cert.pem --tls-key key.pem
```

#### Using Configuration File

```yaml
server:
  port: 53

doh:
  enabled: true
  port: 443
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"
```

Start the server:

```bash
dns server --config config.yaml
```

### Custom DoH Port

Default DoH port is 443. You can change it:

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

## Client Usage

### Query Using DoH

```bash
dns client --domain example.com --server https://your-doh-server:443/dns-query
```

### Popular DoH Servers

- Cloudflare: `https://cloudflare-dns.com/dns-query`
- AdGuard: `https://dns.adguard.com/dns-query`
- Google: `https://dns.google/dns-query`
- Quad9: `https://dns.quad9.net/dns-query`

## Testing

### Test DoH Server

```bash
# Query your DoH server
dns client --domain example.com --server https://localhost:443/dns-query
```

### Verify Certificate

```bash
openssl s_client -connect localhost:443 -servername localhost
```

### Test with curl

```bash
# Test DoH endpoint with curl
curl -H "Accept: application/dns-message" \
     -H "Content-Type: application/dns-message" \
     --data-binary @query.bin \
     https://localhost:443/dns-query
```

## Protocol Details

DoH uses the standard HTTPS protocol (RFC 8484) to send DNS queries:
- **Method**: GET or POST
- **Content-Type**: `application/dns-message`
- **Endpoint**: `/dns-query` (default)
- **Port**: 443 (default)

## Advantages

1. **Firewall Friendly**: Uses port 443, same as HTTPS traffic
2. **Proxy Compatible**: Works through HTTP proxies
3. **Standard Protocol**: Uses well-established HTTPS
4. **Privacy**: Encrypted DNS queries

## Limitations

1. **Latency**: Slightly higher latency compared to DoT/DoQ due to HTTP overhead
2. **Connection Overhead**: Each query may require a new HTTPS connection
3. **Port Conflicts**: If port 443 is already in use, you need to use a different port

## Security Considerations

1. **Certificate Validation**: Ensure your certificate is valid and trusted
2. **Firewall**: Open port 443 in your firewall if needed
3. **Certificate Renewal**: Set up automatic certificate renewal for production
4. **HTTPS Security**: Follow HTTPS best practices for secure deployment

## Comparison with Other Protocols

| Feature | DoH | DoT | DoQ |
|---------|-----|-----|-----|
| Port | 443 | 853 | 853 |
| Firewall Friendly | ✅ | ⚠️ | ⚠️ |
| Latency | Medium | Low | Very Low |
| Connection Multiplexing | ❌ | ❌ | ✅ |
| Proxy Support | ✅ | ❌ | ❌ |

## Next Steps

- [Server Usage](/guide/server) - Learn more about DNS server
- [Client Usage](/guide/client) - Learn how to use the DNS client
- [Configuration](/guide/configuration) - Learn about configuration options
- [DoH and DoQ Examples](/examples/doh-doq-server) - See practical examples
