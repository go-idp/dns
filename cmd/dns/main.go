package main

import (
	"github.com/go-idp/dns"
	"github.com/go-idp/dns/cmd/dns/commands"
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "dns",
		Usage:   "A simple and powerful DNS client and server",
		Version: dns.Version,
	})

	app.Register("client", commands.NewClientCommand())
	app.Register("server", commands.NewServerCommand())

	app.Run()
}
