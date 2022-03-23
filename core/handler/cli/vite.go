package cli

import (
	"github.com/vite-cloud/vite/core/handler/cli/cmd/proxy"
	"os"

	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/handler/cli/cmd"
)

// New returns a new CLI
func New() *cli.CLI {
	c := cli.New(os.Stdout, os.Stdin, os.Stderr)

	c.Add(
		cmd.NewVersionCommand(c),
		cmd.NewMedicCommand(c),
		cmd.NewSetupCommand(c),
		cmd.NewUseCommand(c),
		cmd.NewSelfUpdateCommand(c),
		cmd.NewLogsCommand(c),
		
		proxy.NewProxyCommand(c),
	)

	return c
}
