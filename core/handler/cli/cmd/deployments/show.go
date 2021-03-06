package deployments

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func newShowCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [deployment]",
		Short: "show details about a given deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return manager.ShowCommand(cli, args[0])
		},
	}

	return cmd
}
