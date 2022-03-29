package deployments

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func NewDeploymentsCommand(c *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployments",
		Short:   "manage deployments",
		Aliases: []string{"ds"},
	}

	cmd.AddCommand(
		newListCommand(c),
		newCleanupCommand(c),
		newRollbackCommand(c),
		newShowCommand(c),
	)

	return cmd
}
