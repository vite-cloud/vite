package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func NewProxyCommand(c *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "manage proxy",
	}

	cmd.AddCommand(
		NewRunCommand(c),
		NewLogsCommand(c),
		NewUpCommand(c),
		NewStatusCommand(c),
		NewDownCommand(c),
	)

	return cmd
}
