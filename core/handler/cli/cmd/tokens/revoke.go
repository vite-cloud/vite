package tokens

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runRevokeCommand(cli *cli.CLI) error {
	return nil
}

func NewRevokeCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke a token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCommand(cli)
		},
	}

	return cmd
}
