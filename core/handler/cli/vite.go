package cli

import (
	"github.com/vite-cloud/vite/core/handler/cli/cmd/deployments"
	"github.com/vite-cloud/vite/core/handler/cli/cmd/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cmd/tokens"
	"os"

	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/handler/cli/cmd"
)

// New returns a new CLI
func New() *cli.CLI {
	c := cli.New(os.Stdout, os.Stdin, os.Stderr)

	c.Add(
		cmd.NewDiagnoseCommand(c),
		cmd.NewSetupCommand(c),
		cmd.NewUseCommand(c),
		cmd.NewSelfUpdateCommand(c),
		cmd.NewLogsCommand(c),
		cmd.NewDeployCommand(c),

		proxy.NewProxyCommand(c),

		deployments.NewDeploymentsCommand(c),

		tokens.NewRootCommand(c),
	)

	return c
}
