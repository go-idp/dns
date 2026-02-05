# Installation

## Prerequisites

- Go 1.25.6 or later
- For DNS-over-TLS server: TLS certificate and key files

## Install as CLI Tool

### Using go install

```bash
go install github.com/go-idp/dns/cmd/dns@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/go-idp/dns.git
cd dns

# Build the binary
go build -o bin/dns ./cmd/dns

# Or install globally
go install ./cmd/dns
```

### Verify Installation

```bash
dns --version
```

## Docker

Docker images are available on [Docker Hub](https://hub.docker.com/r/goidp/dns).

```bash
docker pull goidp/dns:latest
```

## Next Steps

- [Quick Start](/guide/quick-start) - Get started with your first DNS query
- [Client Usage](/guide/client) - Learn how to use the DNS client
- [Server Usage](/guide/server) - Learn how to run a DNS server
