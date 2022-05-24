package tokens

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func NewRootCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "manage tokens",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(NewCreateCommand(cli))
	cmd.AddCommand(NewListCommand(cli))
	cmd.AddCommand(NewRevokeCommand(cli))

	return cmd
}
