package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/dns"
	"github.com/go-zoox/dns/client"
	"github.com/go-zoox/dns/constants"
	ucli "github.com/urfave/cli/v2"
)

// errLookupHelp signals that lookup usage should be printed (see parseLookupArgv).
var errLookupHelp = errors.New("lookup help")

type lookupParsed struct {
	domain   string
	servers  []string
	qtype    string
	timeout  string
	plain    bool
}

// parseLookupArgv parses interspersed flags and positional domain (urfave/cli stops
// flag parsing at the first non-flag token, so we parse ourselves when SkipFlagParsing is set).
func parseLookupArgv(argv []string) (*lookupParsed, error) {
	o := &lookupParsed{
		qtype:   "A",
		timeout: "5s",
	}
	var positional []string

	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		if arg == "--" {
			positional = append(positional, argv[i+1:]...)
			break
		}
		if arg == "-h" || arg == "--help" {
			return nil, errLookupHelp
		}
		if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
			continue
		}

		long, shortVal, eqOK := strings.Cut(arg, "=")
		var flagName string
		var inlineVal string
		if eqOK {
			flagName = long
			inlineVal = shortVal
		} else {
			flagName = arg
		}

		needVal := func() (string, error) {
			if inlineVal != "" {
				v := inlineVal
				inlineVal = ""
				return v, nil
			}
			i++
			if i >= len(argv) {
				return "", fmt.Errorf("flag %s requires a value", flagName)
			}
			if strings.HasPrefix(argv[i], "-") {
				return "", fmt.Errorf("flag %s requires a value", flagName)
			}
			return argv[i], nil
		}

		switch flagName {
		case "--server", "-s":
			v, err := needVal()
			if err != nil {
				return nil, err
			}
			o.servers = append(o.servers, v)
		case "--domain", "-d":
			v, err := needVal()
			if err != nil {
				return nil, err
			}
			o.domain = strings.TrimSpace(v)
		case "--type", "-t":
			v, err := needVal()
			if err != nil {
				return nil, err
			}
			o.qtype = strings.TrimSpace(v)
		case "--timeout":
			v, err := needVal()
			if err != nil {
				return nil, err
			}
			o.timeout = strings.TrimSpace(v)
		case "--plain":
			if inlineVal != "" {
				return nil, fmt.Errorf("invalid use of --plain")
			}
			o.plain = true
		default:
			return nil, fmt.Errorf("unknown flag %q", arg)
		}
	}

	if o.domain == "" && len(positional) > 0 {
		o.domain = strings.TrimSpace(positional[0])
		positional = positional[1:]
	}
	if len(positional) > 0 {
		return nil, fmt.Errorf("unexpected arguments after domain: %q", strings.Join(positional, " "))
	}

	applyLookupEnv(o)
	return o, nil
}

func applyLookupEnv(o *lookupParsed) {
	if len(o.servers) == 0 {
		if ev := strings.TrimSpace(os.Getenv("DNS_SERVER")); ev != "" {
			for _, p := range strings.Split(ev, ",") {
				p = strings.TrimSpace(p)
				if p != "" {
					o.servers = append(o.servers, p)
				}
			}
		}
	}
	if os.Getenv("DNS_TIMEOUT") != "" && o.timeout == "5s" {
		if v := strings.TrimSpace(os.Getenv("DNS_TIMEOUT")); v != "" {
			o.timeout = v
		}
	}
	if envPlain := strings.TrimSpace(strings.ToLower(os.Getenv("DNS_PLAIN"))); envPlain != "" {
		switch envPlain {
		case "1", "true", "yes", "on":
			o.plain = true
		}
	}
}

func newClientLookupCommand() *cli.Command {
	return &cli.Command{
		Name:            "lookup",
		Aliases:         []string{"l"},
		Usage:           "Query a DNS server for A or AAAA records",
		ArgsUsage:       "<domain>",
		SkipFlagParsing: true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "DNS server address (supports plain DNS, DoT, DoH, etc.)",
				EnvVars: []string{"DNS_SERVER"},
			},
			&cli.StringFlag{
				Name:    "domain",
				Aliases: []string{"d"},
				Usage:   "Domain to query (optional if given as the first argument)",
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Query type (A, AAAA)",
				Value:   "A",
			},
			&cli.StringFlag{
				Name:    "timeout",
				Usage:   "Timeout for DNS query (e.g., 5s, 10s)",
				Value:   "5s",
				EnvVars: []string{"DNS_TIMEOUT"},
			},
			&cli.BoolFlag{
				Name:    "plain",
				Usage:   "Output only IP addresses, one per line",
				EnvVars: []string{"DNS_PLAIN"},
			},
		},
		Action: func(ctx *cli.Context) error {
			opts, err := parseLookupArgv(ctx.Args().Slice())
			if errors.Is(err, errLookupHelp) {
				lineage := ctx.Lineage()
				if len(lineage) >= 2 {
					return ucli.ShowCommandHelp(lineage[1], "lookup")
				}
				return ucli.ShowCommandHelp(ctx, "lookup")
			}
			if err != nil {
				return err
			}

			domain := opts.domain
			if domain == "" {
				return fmt.Errorf("domain is required (e.g. dns client lookup example.com)")
			}

			timeout, err := time.ParseDuration(opts.timeout)
			if err != nil {
				return fmt.Errorf("invalid timeout format: %v", err)
			}

			servers := opts.servers
			if len(servers) == 0 {
				servers = []string{"114.114.114.114:53"}
			}

			queryType := strings.ToUpper(opts.qtype)
			normalizedServers := make([]string, len(servers))
			for i, server := range servers {
				normalizedServers[i] = normalizeServerAddress(server)
			}

			dnsClient := dns.NewClient(&dns.ClientOptions{
				Servers: normalizedServers,
				Timeout: timeout,
			})

			var typ int
			switch queryType {
			case "A":
				typ = constants.QueryTypeIPv4
			case "AAAA":
				typ = constants.QueryTypeIPv6
			default:
				return fmt.Errorf("unsupported query type: %s (supported: A, AAAA)", queryType)
			}

			ips, err := dnsClient.LookUp(domain, &client.LookUpOptions{
				Typ: typ,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if len(ips) == 0 {
				if !opts.plain {
					fmt.Printf("No %s records found for %s\n", queryType, domain)
				}
				return nil
			}

			if opts.plain {
				for _, ip := range ips {
					fmt.Println(ip)
				}
			} else {
				fmt.Printf("%s records for %s:\n", queryType, domain)
				for _, ip := range ips {
					fmt.Printf("  %s\n", ip)
				}
			}

			return nil
		},
	}
}
