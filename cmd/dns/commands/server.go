package commands

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/go-idp/dns/cmd/dns/config"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/dns"
	"github.com/go-zoox/dns/client"
	"github.com/go-zoox/fs/type/hosts"
	"github.com/go-zoox/logger"
)

// SystemHostsEntry represents an entry in the system hosts file
type SystemHostsEntry struct {
	IP         string
	Domain     string
	IsWildcard bool
	IsRegex    bool
	Regex      *regexp.Regexp
}

// isRegexPattern checks if a string is a valid regex pattern
// Wildcard patterns take priority, so we only check for regex if it doesn't contain *
func isRegexPattern(pattern string) bool {
	// Wildcard takes priority
	if strings.Contains(pattern, "*") {
		return false
	}

	// Try to compile as regex
	if _, err := regexp.Compile(pattern); err != nil {
		return false
	}

	// Check if it contains regex metacharacters (beyond just dots)
	hasRegexMeta := strings.ContainsAny(pattern, "^$+?()[]{}|\\")
	return hasRegexMeta
}

// parseSystemHostsFile parses a system hosts file with support for wildcard and regex patterns
// Uses github.com/go-zoox/fs/type/hosts to parse the file, then enhances entries with pattern matching
// Note: hostsParser.Mapping format is "domain:queryType" -> "IP" (e.g., "frontend:4" -> "10.1.0.169")
func parseSystemHostsFile(filePath string) ([]SystemHostsEntry, error) {
	hostsParser := hosts.New(filePath)
	if err := hostsParser.Load(); err != nil {
		return nil, fmt.Errorf("failed to load hosts file: %w", err)
	}

	domainMap := make(map[string]*SystemHostsEntry) // Use map to merge entries for same domain

	// Iterate through the hosts mapping
	// Note: hostsParser.Mapping format is "domain:queryType" -> "IP" (e.g., "frontend:4" -> "10.1.0.169")
	for key, ip := range hostsParser.Mapping {
		// Parse the key format: "domain:queryType" or just "domain"
		// Extract domain part (remove :4 or :6 suffix)
		domain := key
		if idx := strings.LastIndex(key, ":"); idx > 0 {
			domain = key[:idx]
		}
		domain = strings.TrimSpace(domain)
		if domain == "" {
			continue
		}

		domainLower := strings.ToLower(domain)

		// Get or create entry for this domain
		entry, exists := domainMap[domainLower]
		if !exists {
			// Try to determine if it's a regex pattern by attempting to compile it
			isRegex := isRegexPattern(domain)
			// Check if it contains wildcard (but not if it's already a regex)
			isWildcard := !isRegex && strings.Contains(domain, "*")

			entry = &SystemHostsEntry{
				IP:         ip,
				Domain:     domainLower,
				IsWildcard: isWildcard,
				IsRegex:    isRegex,
			}

			// Compile regex if it's a regex pattern
			if isRegex {
				compiled, err := regexp.Compile(domain)
				if err != nil {
					logger.Warn("Failed to compile regex pattern in hosts file: %s, error: %v", domain, err)
					continue
				}
				entry.Regex = compiled
			}

			domainMap[domainLower] = entry
		}
		// If entry exists, we keep the first IP found (could be enhanced to support multiple IPs)
	}

	// Convert map to slice (deduplicated by domain)
	uniqueEntries := make([]SystemHostsEntry, 0, len(domainMap))
	for _, entry := range domainMap {
		uniqueEntries = append(uniqueEntries, *entry)
	}

	logger.Debug("Parsed %d unique entries from hosts file %s (total mappings: %d)", len(uniqueEntries), filePath, len(hostsParser.Mapping))
	return uniqueEntries, nil
}

// lookupSystemHosts looks up a domain in system hosts entries with wildcard and regex support
func lookupSystemHosts(entries []SystemHostsEntry, domain string, queryType int) (string, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	domainNoDot := strings.TrimSuffix(domain, ".")

	logger.Debug("Looking up domain: %s (no dot: %s), query type: %d, entries count: %d", domain, domainNoDot, queryType, len(entries))

	// First try exact matches
	for _, entry := range entries {
		if !entry.IsWildcard && !entry.IsRegex {
			logger.Debug("Checking exact match: entry.Domain=%s, domain=%s, domainNoDot=%s", entry.Domain, domain, domainNoDot)
			if entry.Domain == domain || entry.Domain == domainNoDot {
				// Check if IP matches query type
				if queryType == 4 && !config.IsIPv6(entry.IP) {
					logger.Debug("Found exact match: %s -> %s", entry.Domain, entry.IP)
					return entry.IP, nil
				} else if queryType == 6 && config.IsIPv6(entry.IP) {
					logger.Debug("Found exact match: %s -> %s", entry.Domain, entry.IP)
					return entry.IP, nil
				} else {
					logger.Debug("IP type mismatch: queryType=%d, isIPv6=%v, IP=%s", queryType, config.IsIPv6(entry.IP), entry.IP)
				}
			}
		}
	}

	// Then try wildcard and regex patterns
	for _, entry := range entries {
		var matched bool

		if entry.IsRegex && entry.Regex != nil {
			matched = entry.Regex.MatchString(domain) || entry.Regex.MatchString(domainNoDot)
			if matched {
				logger.Debug("Regex match: pattern=%s, domain=%s", entry.Domain, domain)
			}
		} else if entry.IsWildcard {
			matched = config.MatchWildcard(domain, entry.Domain) || config.MatchWildcard(domainNoDot, entry.Domain)
			if matched {
				logger.Debug("Wildcard match: pattern=%s, domain=%s", entry.Domain, domain)
			}
		}

		if matched {
			// Check if IP matches query type
			if queryType == 4 && !config.IsIPv6(entry.IP) {
				return entry.IP, nil
			} else if queryType == 6 && config.IsIPv6(entry.IP) {
				return entry.IP, nil
			}
		}
	}

	logger.Debug("No match found for domain: %s", domain)
	return "", fmt.Errorf("not found")
}

// NewServerCommand creates a new server command
func NewServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start a DNS server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to configuration file (YAML)",
				EnvVars: []string{"DNS_CONFIG"},
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "DNS server port (UDP/TCP)",
				Value:   53,
				EnvVars: []string{"DNS_PORT"},
			},
			&cli.StringFlag{
				Name:    "host",
				Usage:   "DNS server host",
				Value:   "0.0.0.0",
				EnvVars: []string{"DNS_HOST"},
			},
			&cli.UintFlag{
				Name:    "ttl",
				Usage:   "TTL for DNS responses (in seconds)",
				Value:   500,
				EnvVars: []string{"DNS_TTL"},
			},
			&cli.BoolFlag{
				Name:    "dot",
				Usage:   "Enable DNS-over-TLS (DoT)",
				EnvVars: []string{"DNS_DOT"},
			},
			&cli.IntFlag{
				Name:    "dot-port",
				Usage:   "DoT server port",
				Value:   853,
				EnvVars: []string{"DNS_DOT_PORT"},
			},
			&cli.StringFlag{
				Name:    "tls-cert",
				Usage:   "TLS certificate file path (required for DoT)",
				EnvVars: []string{"DNS_TLS_CERT"},
			},
			&cli.StringFlag{
				Name:    "tls-key",
				Usage:   "TLS private key file path (required for DoT)",
				EnvVars: []string{"DNS_TLS_KEY"},
			},
			&cli.StringSliceFlag{
				Name:    "upstream",
				Aliases: []string{"u"},
				Usage:   "Upstream DNS servers",
				EnvVars: []string{"DNS_UPSTREAM"},
			},
			&cli.BoolFlag{
				Name:    "disable-system-hosts",
				Usage:   "Disable system hosts file lookup (enabled by default)",
				EnvVars: []string{"DNS_DISABLE_SYSTEM_HOSTS"},
			},
			&cli.StringFlag{
				Name:    "system-hosts-file",
				Usage:   "Path to system hosts file (default: /etc/hosts)",
				Value:   "/etc/hosts",
				EnvVars: []string{"DNS_SYSTEM_HOSTS_FILE"},
			},
		},
		Action: func(ctx *cli.Context) error {
			var cfg *config.Config
			var err error

			// Load configuration file if provided
			configPath := ctx.String("config")
			if configPath != "" {
				cfg, err = config.LoadConfig(configPath)
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}
				logger.Info("Loaded configuration from %s", configPath)
			}

			// Get values from config file or command line flags (CLI flags override config)
			port := ctx.Int("port")
			host := ctx.String("host")
			ttl := ctx.Uint("ttl")
			enableDoT := ctx.Bool("dot")
			dotPort := ctx.Int("dot-port")
			tlsCert := ctx.String("tls-cert")
			tlsKey := ctx.String("tls-key")
			upstreams := ctx.StringSlice("upstream")
			disableSystemHosts := ctx.Bool("disable-system-hosts")
			systemHostsFile := ctx.String("system-hosts-file")

			// Merge config file values if config was loaded
			if cfg != nil {
				if port == 53 && cfg.Server.Port != 0 {
					port = cfg.Server.Port
				}
				if host == "0.0.0.0" && cfg.Server.Host != "" {
					host = cfg.Server.Host
				}
				if ttl == 500 && cfg.Server.TTL != 0 {
					ttl = uint(cfg.Server.TTL)
				}
				if !enableDoT && cfg.DoT.Enabled {
					enableDoT = cfg.DoT.Enabled
				}
				if dotPort == 853 && cfg.DoT.Port != 0 {
					dotPort = cfg.DoT.Port
				}
				if tlsCert == "" && cfg.DoT.TLS.Cert != "" {
					tlsCert = cfg.DoT.TLS.Cert
				}
				if tlsKey == "" && cfg.DoT.TLS.Key != "" {
					tlsKey = cfg.DoT.TLS.Key
				}
				if len(upstreams) == 0 && len(cfg.Upstream.Servers) > 0 {
					upstreams = cfg.Upstream.Servers
				}
				// Merge system hosts config (CLI flags override config)
				// System hosts is enabled by default, unless explicitly disabled
				if !disableSystemHosts && cfg.SystemHosts.Disabled {
					disableSystemHosts = cfg.SystemHosts.Disabled
				}
				if systemHostsFile == "/etc/hosts" && cfg.SystemHosts.FilePath != "" {
					systemHostsFile = cfg.SystemHosts.FilePath
				}
			}

			// Default upstream if still empty
			if len(upstreams) == 0 {
				upstreams = []string{"114.114.114.114:53"}
			}

			// Parse upstream timeout
			upstreamTimeout := 5 * time.Second
			if cfg != nil && cfg.Upstream.Timeout != "" {
				if parsed, err := time.ParseDuration(cfg.Upstream.Timeout); err == nil {
					upstreamTimeout = parsed
				}
			}

			// Validate DoT configuration
			if enableDoT {
				if tlsCert == "" || tlsKey == "" {
					return fmt.Errorf("TLS certificate and key are required when DoT is enabled (use --tls-cert and --tls-key or config file)")
				}
			}

			// Create upstream client
			upstreamClient := dns.NewClient(&dns.ClientOptions{
				Servers: upstreams,
				Timeout: upstreamTimeout,
			})

			// Create server
			serverOptions := &dns.ServerOptions{
				Port:      port,
				Host:      host,
				TTL:       uint32(ttl),
				EnableDoT: enableDoT,
			}

			if enableDoT {
				serverOptions.DoTPort = dotPort
				serverOptions.TLSCertFile = tlsCert
				serverOptions.TLSKeyFile = tlsKey
			}

			server := dns.NewServer(serverOptions)

			// Initialize system hosts file parser (enabled by default unless disabled)
			var systemHostsEntries []SystemHostsEntry
			if !disableSystemHosts {
				// Use custom parser (supports wildcard and regex)
				entries, err := parseSystemHostsFile(systemHostsFile)
				if err != nil {
					logger.Warn("Failed to load system hosts file %s: %v", systemHostsFile, err)
					systemHostsEntries = []SystemHostsEntry{}
				} else {
					systemHostsEntries = entries
					logger.Info("Loaded system hosts file: %s (with wildcard/regex support, %d entries)", systemHostsFile, len(entries))
					// Debug: log first few entries
					for i, entry := range entries {
						if i < 5 {
							logger.Debug("System hosts entry %d: %s -> %s (wildcard: %v, regex: %v)", i, entry.Domain, entry.IP, entry.IsWildcard, entry.IsRegex)
						}
					}
				}
			}

			// Set up handler with three-tier priority:
			// 1. Configuration hosts (highest priority)
			// 2. System /etc/hosts file (if enabled)
			// 3. Upstream DNS servers
			server.Handle(func(hostname string, typ int) ([]string, error) {
				// Map query type to string for better logging
				queryType := "A"
				if typ == 6 {
					queryType = "AAAA"
				}
				logger.Debug("DNS query received: %s (type: %s, code: %d)", hostname, queryType, typ)

				// Priority 1: Check configuration hosts mapping
				if cfg != nil {
					ips, err := cfg.LookupHost(hostname, typ)
					if err == nil && len(ips) > 0 {
						logger.Info("[channel: config.hosts] Resolved %s (%s) from config hosts -> %v", hostname, queryType, ips)
						return ips, nil
					}
					logger.Debug("No match found in config hosts for %s (%s)", hostname, queryType)
				} else {
					logger.Debug("Config hosts not available, skipping priority 1")
				}

				// Priority 2: Check system hosts file (if enabled)
				if len(systemHostsEntries) > 0 {
					logger.Debug("Checking system hosts for %s (%s), total entries: %d", hostname, queryType, len(systemHostsEntries))
					ip, err := lookupSystemHosts(systemHostsEntries, hostname, typ)
					if err == nil && ip != "" {
						logger.Info("[channel: system.hosts] Resolved %s (%s) from system hosts -> %v", hostname, queryType, []string{ip})
						return []string{ip}, nil
					}
					logger.Debug("No match found in system hosts for %s (%s), error: %v", hostname, queryType, err)
				} else {
					logger.Debug("System hosts not enabled or empty, skipping priority 2")
				}

				// Priority 3: Fallback to upstream DNS servers
				logger.Debug("Querying upstream DNS servers for %s (%s)", hostname, queryType)
				ips, err := upstreamClient.LookUp(hostname, &client.LookUpOptions{
					Typ: typ,
				})
				if err != nil {
					logger.Error("Failed to resolve %s (%s) from upstream: %v", hostname, queryType, err)
					return nil, err
				}

				if len(ips) > 0 {
					logger.Info("[channel: upstream] Resolved %s (%s) from upstream -> %v", hostname, queryType, ips)
				} else {
					logger.Warn("No results found for %s (%s) from upstream", hostname, queryType)
				}
				return ips, nil
			})

			// Handle graceful shutdown
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-sigChan
				logger.Info("Shutting down DNS server...")
				os.Exit(0)
			}()

			// Start server
			if enableDoT {
				logger.Info("Starting DNS server on %s:%d (UDP/TCP) and DoT on %s:%d", host, port, host, dotPort)
			} else {
				logger.Info("Starting DNS server on %s:%d (UDP/TCP)", host, port)
			}

			server.Serve()

			return nil
		},
	}
}
