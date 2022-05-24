package tokens

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runListCommand(cli *cli.CLI) error {
	return nil
}

func NewListCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list tokens",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cli)
		},
	}

	return cmd
}
