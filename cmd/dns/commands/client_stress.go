package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-zoox/cli"
	"github.com/miekg/dns"
)

func newClientStressCommand() *cli.Command {
	return &cli.Command{
		Name:      "stress",
		Aliases:   []string{"bench"},
		Usage:     "Concurrent plain DNS (UDP/TCP) load test against a single server",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Plain DNS server (host or host:port); first value is used. DoT/DoH/DoQ are not supported",
				EnvVars: []string{"DNS_SERVER"},
			},
			&cli.StringFlag{
				Name:    "domain",
				Aliases: []string{"d"},
				Usage:   "QNAME to query",
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Query type (A, AAAA, …)",
				Value:   "A",
			},
			&cli.StringFlag{
				Name:    "net",
				Usage:   "Transport: udp or tcp",
				Value:   "udp",
			},
			&cli.IntFlag{
				Name:    "workers",
				Aliases: []string{"w"},
				Usage:   "Number of concurrent workers (goroutines)",
				Value:   100,
			},
			&cli.IntFlag{
				Name:    "requests",
				Aliases: []string{"n"},
				Usage:   "Total queries to send",
				Value:   1000,
			},
			&cli.StringFlag{
				Name:    "timeout",
				Usage:   "Per-query timeout (e.g. 2s)",
				Value:   "2s",
			},
			&cli.BoolFlag{
				Name:  "accept-nxdomain",
				Usage: "Treat NXDOMAIN responses as success",
			},
		},
		Action: func(ctx *cli.Context) error {
			servers := ctx.StringSlice("server")
			domain := strings.TrimSpace(ctx.String("domain"))
			qtypeStr := strings.ToUpper(strings.TrimSpace(ctx.String("type")))
			netName := strings.ToLower(strings.TrimSpace(ctx.String("net")))
			workers := ctx.Int("workers")
			n := ctx.Int("requests")
			acceptNX := ctx.Bool("accept-nxdomain")

			if domain == "" {
				return fmt.Errorf("domain is required")
			}
			if workers < 1 {
				return fmt.Errorf("workers must be >= 1")
			}
			if n < 1 {
				return fmt.Errorf("n must be >= 1")
			}
			if netName != "udp" && netName != "tcp" {
				return fmt.Errorf("net must be udp or tcp, got %q", netName)
			}

			timeout, err := time.ParseDuration(ctx.String("timeout"))
			if err != nil {
				return fmt.Errorf("invalid timeout: %w", err)
			}

			server := "114.114.114.114"
			if len(servers) > 0 && strings.TrimSpace(servers[0]) != "" {
				server = strings.TrimSpace(servers[0])
			}
			addr, err := plainDNSAddressForStress(server)
			if err != nil {
				return err
			}

			qtype, typeOK := dns.StringToType[qtypeStr]
			if !typeOK {
				return fmt.Errorf("unknown query type %q", qtypeStr)
			}

			dnsClient := &dns.Client{
				Net:     netName,
				Timeout: timeout,
			}

			var ok, fail atomic.Uint64
			var next atomic.Uint64
			start := time.Now()
			var wg sync.WaitGroup

			for w := 0; w < workers; w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						i := next.Add(1)
						if int(i) > n {
							return
						}
						m := new(dns.Msg)
						m.SetQuestion(dns.Fqdn(domain), qtype)
						m.RecursionDesired = true

						rctx, cancel := context.WithTimeout(context.Background(), timeout)
						r, _, err := dnsClient.ExchangeContext(rctx, m, addr)
						cancel()
						if err != nil || r == nil {
							fail.Add(1)
							continue
						}
						switch r.Rcode {
						case dns.RcodeSuccess:
							ok.Add(1)
						case dns.RcodeNameError:
							if acceptNX {
								ok.Add(1)
							} else {
								fail.Add(1)
							}
						default:
							fail.Add(1)
						}
					}
				}()
			}
			wg.Wait()

			elapsed := time.Since(start).Seconds()
			if elapsed < 1e-9 {
				elapsed = 1e-9
			}

			fmt.Fprintf(os.Stdout, "server:   %s (%s)\n", addr, netName)
			fmt.Fprintf(os.Stdout, "domain:   %s %s\n", domain, qtypeStr)
			fmt.Fprintf(os.Stdout, "workers:  %d\n", workers)
			fmt.Fprintf(os.Stdout, "queries:  %d\n", n)
			fmt.Fprintf(os.Stdout, "elapsed:  %.3fs\n", elapsed)
			fmt.Fprintf(os.Stdout, "ok:       %d\n", ok.Load())
			fmt.Fprintf(os.Stdout, "fail:     %d\n", fail.Load())
			fmt.Fprintf(os.Stdout, "qps:      %.0f\n", float64(n)/elapsed)

			if fail.Load() > 0 {
				return fmt.Errorf("%d of %d queries failed", fail.Load(), n)
			}
			return nil
		},
	}
}
