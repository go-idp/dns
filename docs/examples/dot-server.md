# DoT Server Example

This example shows how to run a DNS-over-TLS (DoT) server using the CLI.

## Prerequisites

You need a TLS certificate and key file. See [DoT Guide](/guide/dot) for certificate generation.

## Start DoT Server

### Using Command Line

```bash
dns server --port 53 --dot --tls-cert cert.pem --tls-key key.pem
```

### Using Configuration File

Create `config.yaml`:

```yaml
server:
  port: 53

dot:
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

## Custom DoT Port

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

## Testing

```bash
# Query using DoT
dns client --domain example.com --server tls://localhost:853
```

## Next Steps

- [DoT Guide](/guide/dot) - Learn more about DoT
- [Configuration File Example](/examples/config-file) - Learn about configuration files
