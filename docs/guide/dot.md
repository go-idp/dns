# DNS-over-TLS (DoT)

DNS-over-TLS (DoT) provides encrypted DNS queries using TLS protocol.

## Overview

DoT encrypts DNS queries between the client and server, providing:
- **Privacy**: DNS queries are encrypted
- **Security**: Protection against DNS spoofing
- **Performance**: Low latency compared to DoH

## Server Setup

### Generate TLS Certificate

You need a TLS certificate and key file to run a DoT server.

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

### Start DoT Server

#### Using Command Line

```bash
dns server --port 53 --dot --tls-cert cert.pem --tls-key key.pem
```

#### Using Configuration File

```yaml
server:
  port: 53

dot:
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

### Custom DoT Port

Default DoT port is 853. You can change it:

```bash
dns server --port 53 --dot --dot-port 853 --tls-cert cert.pem --tls-key key.pem
```

Or in config file:

```yaml
dot:
  enabled: true
  port: 853  # Custom DoT port
  tls:
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"
```

## Client Usage

### Query Using DoT

```bash
dns client --domain example.com --server tls://your-dot-server:853
```

### Popular DoT Servers

- Cloudflare: `tls://1.1.1.1`
- Google: `tls://8.8.8.8`
- Quad9: `tls://9.9.9.9`

## Testing

### Test DoT Server

```bash
# Query your DoT server
dns client --domain example.com --server tls://localhost:853
```

### Verify Certificate

```bash
openssl s_client -connect localhost:853 -servername localhost
```

## Security Considerations

1. **Certificate Validation**: Ensure your certificate is valid and trusted
2. **Firewall**: Open port 853 in your firewall if needed
3. **Certificate Renewal**: Set up automatic certificate renewal for production

## Next Steps

- [Server Usage](/guide/server) - Learn more about DNS server
- [Configuration](/guide/configuration) - Learn about configuration options
