package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the DNS server configuration
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	DoT         DoTConfig         `yaml:"dot"`
	Hosts       HostsConfig       `yaml:"hosts"`
	SystemHosts SystemHostsConfig `yaml:"system_hosts"`
	Upstream    UpstreamConfig    `yaml:"upstream"`
}

// ServerConfig represents basic server settings
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	TTL  uint32 `yaml:"ttl"`
}

// DoTConfig represents DNS-over-TLS configuration
type DoTConfig struct {
	Enabled bool      `yaml:"enabled"`
	Port    int       `yaml:"port"`
	TLS     TLSConfig `yaml:"tls"`
}

// TLSConfig represents TLS certificate configuration
type TLSConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// HostsConfig represents custom domain to IP mappings
// Supports multiple formats:
//   - Simple: "example.com": "1.2.3.4"
//   - Multiple IPs: "example.com": ["1.2.3.4", "1.2.3.5"]
//   - With type: "example.com": {"a": ["1.2.3.4"], "aaaa": ["2001:db8::1"]}
type HostsConfig map[string]interface{}

// HostMapping represents a parsed host mapping
type HostMapping struct {
	Domain     string
	IPv4       []string
	IPv6       []string
	IsWildcard bool           // true if domain contains wildcard (*)
	IsRegex    bool           // true if domain is a regex pattern (starts with ^)
	Regex      *regexp.Regexp // compiled regex pattern if IsRegex is true
}

// SystemHostsConfig represents system hosts file configuration
type SystemHostsConfig struct {
	Disabled bool   `yaml:"disabled"`
	FilePath string `yaml:"file_path"`
}

// UpstreamConfig represents upstream DNS servers configuration
type UpstreamConfig struct {
	Servers []string `yaml:"servers"`
	Timeout string   `yaml:"timeout"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 53
	}
	if config.Server.TTL == 0 {
		config.Server.TTL = 500
	}
	if config.DoT.Port == 0 {
		config.DoT.Port = 853
	}
	if config.Upstream.Timeout == "" {
		config.Upstream.Timeout = "5s"
	}
	if len(config.Upstream.Servers) == 0 {
		config.Upstream.Servers = []string{"114.114.114.114:53"}
	}

	// Set default system hosts file path if not disabled and not specified
	if !config.SystemHosts.Disabled && config.SystemHosts.FilePath == "" {
		config.SystemHosts.FilePath = "/etc/hosts"
	}

	return &config, nil
}

// ParseHosts parses the hosts configuration into a map of domain to IP mappings
func (c *Config) ParseHosts() (map[string]*HostMapping, error) {
	hosts := make(map[string]*HostMapping)

	for domain, value := range c.Hosts {
		domain = strings.TrimSpace(domain)
		domainLower := strings.ToLower(domain)

		// Check if it contains wildcard first (wildcard takes priority)
		isWildcard := strings.Contains(domain, "*")

		// Try to determine if it's a regex pattern (only if not a wildcard)
		var isRegex bool
		var compiledRegex *regexp.Regexp
		if !isWildcard {
			// Try to compile as regex
			if compiled, err := regexp.Compile(domain); err == nil {
				// Check if it's actually a regex (not just a plain domain)
				// A plain domain like "example.com" contains dots which are regex metacharacters
				// but we want to treat it as a literal. We consider it a regex if it contains
				// regex metacharacters beyond just dots
				hasRegexMeta := strings.ContainsAny(domain, "^$+?()[]{}|\\")
				// Also check if it has escaped characters or other regex features
				hasEscape := strings.Contains(domain, "\\")
				if hasRegexMeta || hasEscape {
					isRegex = true
					compiledRegex = compiled
				}
			}
		}

		mapping := &HostMapping{
			Domain:     domainLower,
			IPv4:       []string{},
			IPv6:       []string{},
			IsWildcard: isWildcard,
			IsRegex:    isRegex,
			Regex:      compiledRegex,
		}

		switch v := value.(type) {
		case string:
			// Simple format: "example.com": "1.2.3.4"
			ip := strings.TrimSpace(v)
			if IsIPv6(ip) {
				mapping.IPv6 = append(mapping.IPv6, ip)
			} else {
				mapping.IPv4 = append(mapping.IPv4, ip)
			}

		case []interface{}:
			// Multiple IPs: "example.com": ["1.2.3.4", "1.2.3.5"]
			for _, item := range v {
				ip := strings.TrimSpace(fmt.Sprintf("%v", item))
				if IsIPv6(ip) {
					mapping.IPv6 = append(mapping.IPv6, ip)
				} else {
					mapping.IPv4 = append(mapping.IPv4, ip)
				}
			}

		case map[string]interface{}:
			// Structured format: "example.com": {"a": [...], "aaaa": [...]}
			if aList, ok := v["a"].([]interface{}); ok {
				for _, item := range aList {
					ip := strings.TrimSpace(fmt.Sprintf("%v", item))
					mapping.IPv4 = append(mapping.IPv4, ip)
				}
			}
			if aaaaList, ok := v["aaaa"].([]interface{}); ok {
				for _, item := range aaaaList {
					ip := strings.TrimSpace(fmt.Sprintf("%v", item))
					mapping.IPv6 = append(mapping.IPv6, ip)
				}
			}
			// Also support single string values
			if aStr, ok := v["a"].(string); ok {
				mapping.IPv4 = append(mapping.IPv4, strings.TrimSpace(aStr))
			}
			if aaaaStr, ok := v["aaaa"].(string); ok {
				mapping.IPv6 = append(mapping.IPv6, strings.TrimSpace(aaaaStr))
			}
		}

		if len(mapping.IPv4) > 0 || len(mapping.IPv6) > 0 {
			hosts[domainLower] = mapping
		}
	}

	return hosts, nil
}

// MatchWildcard checks if a domain matches a wildcard pattern
// This is exported so it can be used by the server command
func MatchWildcard(domain, pattern string) bool {
	// Convert wildcard pattern to regex
	// *.example.com -> ^.*\.example\.com$
	// *.*.example.com -> ^.*\..*\.example\.com$
	regexPattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\*", ".*") + "$"
	matched, _ := regexp.MatchString(regexPattern, domain)
	return matched
}

// IsIPv6 checks if an IP address is IPv6
// This is exported so it can be used by the server command
func IsIPv6(ip string) bool {
	return strings.Contains(ip, ":")
}

// LookupHost looks up a domain in the hosts configuration
func (c *Config) LookupHost(domain string, queryType int) ([]string, error) {
	hosts, err := c.ParseHosts()
	if err != nil {
		return nil, err
	}

	domain = strings.ToLower(strings.TrimSpace(domain))
	domainNoDot := strings.TrimSuffix(domain, ".")

	// Try exact match first
	if mapping, ok := hosts[domain]; ok && !mapping.IsWildcard && !mapping.IsRegex {
		if queryType == 4 { // A record
			if len(mapping.IPv4) > 0 {
				return mapping.IPv4, nil
			}
		} else if queryType == 6 { // AAAA record
			if len(mapping.IPv6) > 0 {
				return mapping.IPv6, nil
			}
		}
	}

	// Try with trailing dot removed
	if mapping, ok := hosts[domainNoDot]; ok && !mapping.IsWildcard && !mapping.IsRegex {
		if queryType == 4 { // A record
			if len(mapping.IPv4) > 0 {
				return mapping.IPv4, nil
			}
		} else if queryType == 6 { // AAAA record
			if len(mapping.IPv6) > 0 {
				return mapping.IPv6, nil
			}
		}
	}

	// Try wildcard and regex patterns
	for pattern, mapping := range hosts {
		var matched bool

		if mapping.IsRegex && mapping.Regex != nil {
			// Match against regex pattern
			matched = mapping.Regex.MatchString(domain) || mapping.Regex.MatchString(domainNoDot)
		} else if mapping.IsWildcard {
			// Match against wildcard pattern
			matched = MatchWildcard(domain, pattern) || MatchWildcard(domainNoDot, pattern)
		}

		if matched {
			if queryType == 4 { // A record
				if len(mapping.IPv4) > 0 {
					return mapping.IPv4, nil
				}
			} else if queryType == 6 { // AAAA record
				if len(mapping.IPv6) > 0 {
					return mapping.IPv6, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("not found in hosts")
}
