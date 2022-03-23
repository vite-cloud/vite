package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func NewProxyCommand(c *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Proxy commands",
	}

	cmd.AddCommand(
		NewRunCommand(c),
		NewLogsCommand(c),
		NewDisableCommand(c),
		NewEnableCommand(c),
	)

	return cmd
}
