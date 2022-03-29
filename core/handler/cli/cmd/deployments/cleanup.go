package deployments

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runCleanupCommand(cli *cli.CLI) error {
	return nil
}

func newCleanupCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup [deployment]",
		Short: "cleanup a given deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCleanupCommand(cli)
		},
	}

	return cmd
}
