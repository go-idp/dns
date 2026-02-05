# Configuration File Example

This example shows how to use a configuration file with the DNS server CLI.

## Configuration File

Create `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 53
  ttl: 500

dot:
  enabled: false

hosts:
  "example.com": "1.2.3.4"
  "www.example.com":
    - "1.2.3.4"
    - "1.2.3.5"
  "*.example.com": "1.2.3.4"

system_hosts:
  disabled: false
  file_path: "/etc/hosts"

upstream:
  servers:
    - "8.8.8.8:53"
    - "tls://1.1.1.1"
  timeout: "5s"
```

## Using Configuration File

Start the server with the configuration file:

```bash
dns server --config config.yaml
```

## Command Line Override

Command line flags override config file values:

```bash
# config.yaml has port: 53, but command line overrides it
dns server --config config.yaml --port 5353
```

## Complete Example

See `example/conf/server.yaml` in the repository for a complete configuration file example.

## Next Steps

- [Configuration Guide](/guide/configuration) - Learn more about configuration options
- [Server Usage](/guide/server) - Learn about server command-line options
