# Basic Server Example

This example shows how to run a basic DNS server using the CLI.

## Start Basic Server

```bash
# Start DNS server on port 53
dns server --port 53
```

## Start with Custom Hosts

Create a configuration file `config.yaml`:

```yaml
server:
  port: 53

hosts:
  "example.com": "1.2.3.4"
  "test.example.com": "1.2.3.5"
```

Start the server:

```bash
dns server --config config.yaml
```

## Start with Upstream DNS

```bash
dns server --port 53 --upstream 8.8.8.8:53 --upstream 1.1.1.1:53
```

## Testing

```bash
# Query the server using dig
dig @127.0.0.1 example.com

# Or using DNS CLI
dns client --domain example.com --server 127.0.0.1:53
```

## Next Steps

- [DoT Server Example](/examples/dot-server) - Learn how to run a DoT server
- [Configuration File Example](/examples/config-file) - Learn about configuration files
