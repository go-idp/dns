# AGENTS.md

This document provides guidance for AI agents working on this codebase.

## Project Overview

This is a DNS client and server CLI tool written in Go. It provides both a command-line interface for DNS queries and a DNS server with support for custom host mappings and DNS-over-TLS (DoT).

## Project Structure

```
.
├── cmd/dns/              # Main CLI application
│   ├── commands/         # CLI command implementations
│   │   ├── client.go          # DNS client parent command + shared helpers
│   │   ├── client_lookup.go   # client lookup subcommand
│   │   ├── client_stress.go   # client stress subcommand (plain DNS load test)
│   │   └── server.go          # DNS server command
│   ├── config/           # Configuration management
│   │   ├── config.go     # Config struct and parsing
│   │   └── config_test.go # Configuration tests
│   └── main.go           # Application entry point
├── example/              # Example code and configurations
│   ├── conf/             # Example configuration files
│   └── main.go           # Example usage
├── .github/workflows/    # CI/CD workflows
├── Dockerfile            # Docker container definition
├── go.mod                # Go module dependencies
└── README.md             # Project documentation
```

## Key Components

### DNS Client (`cmd/dns/commands/client*.go`)
- **`lookup`**: A and AAAA queries; multiple DNS server types (plain DNS, DoT, DoH, DoQ, DNSCrypt); configurable timeout; `--plain` for scripting
- **`stress`**: concurrent plain DNS (UDP/TCP) load test via `github.com/miekg/dns` against one `host:port`

### DNS Server (`cmd/dns/commands/server.go`)
- Plain DNS server (UDP/TCP)
- DNS-over-TLS (DoT) support
- Custom host mappings with highest priority
- Upstream DNS server fallback
- YAML configuration file support
- Command-line flags override config file values

### Configuration (`cmd/dns/config/config.go`)
- YAML-based configuration
- Supports multiple host mapping formats:
  - Simple: `"example.com": "1.2.3.4"`
  - Multiple IPs: `"example.com": ["1.2.3.4", "1.2.3.5"]`
  - Structured: `"example.com": {a: [...], aaaa: [...]}`
- IPv4 and IPv6 support
- Case-insensitive domain matching

## Development Guidelines

### Commit Message Convention
This project follows the [Google Commit Message Convention](https://google.github.io/eng-practices/review/developers.html#commit-messages):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

**Examples:**
- `feat(server): add DoH support`
- `fix(config): handle empty hosts section`
- `docs(readme): update installation instructions`
- `test(config): add tests for IPv6 parsing`

### Code Style
- Follow Go standard formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small
- Handle errors explicitly

### CLI flags (`github.com/go-zoox/cli` / `urfave/cli/v2`)
When defining `cli.StringFlag`, `cli.IntFlag`, etc.:
- **`Name`** — canonical long flag (e.g. `requests`, `workers`, `domain`). This is what `ctx.Int("requests")` / `ctx.String("domain")` use.
- **`Aliases`** — short or alternate names (e.g. `n`, `w`, `d`). Users type `-n` / `-w` / `-d`.

Do **not** put the short letter in `Name` and the long word in `Aliases` (that inverts the library’s convention and breaks help text). Example: total query count should be `Name: "requests", Aliases: []string{"n"}`, not the reverse.

For **`client lookup`**, the desired UX is `lookup <domain> --flags…`. The stdlib flag parser stops at the first non-flag token, so flags after the domain would not apply. The lookup subcommand sets **`SkipFlagParsing: true`** and parses `ctx.Args()` manually (`parseLookupArgv`) so `lookup baidu.com --server 223.5.5.5` works.

### Testing
- Write tests for all configuration parsing logic
- Test edge cases (empty configs, invalid formats, etc.)
- Use table-driven tests where appropriate
- Test files should be named `*_test.go`

### Dependencies
- Main dependencies are in `go.mod`
- Uses `github.com/go-zoox/*` packages for CLI, DNS, and logging
- Uses `gopkg.in/yaml.v3` for YAML parsing

## Common Tasks

### Adding a New Feature
1. Create a feature branch
2. Implement the feature with tests
3. Update documentation (README.md if needed)
4. Commit using Google Commit Message Convention
5. Ensure CI passes

### Modifying Configuration
- Configuration changes should maintain backward compatibility when possible
- Update `config_test.go` with new test cases
- Update example configuration files in `example/conf/`
- Update README.md if configuration format changes

### Adding CLI Flags
- Add flags to the appropriate command in `cmd/dns/commands/`
- Support environment variable overrides using `EnvVars`
- Document in README.md
- Ensure flags can override config file values

## CI/CD

The project uses GitHub Actions for:
- Continuous Integration (`.github/workflows/ci.yml`)
- Docker builds (`.github/workflows/docker.yml`)
- Releases (`.github/workflows/release.yml`)

## Version Management

Version is defined in `version.go` and can be set during build or release process.

## License

MIT License - see LICENSE file for details.
