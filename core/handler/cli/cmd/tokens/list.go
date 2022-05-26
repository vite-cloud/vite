package tokens

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/domain/token"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

var manager = resource.Manager[token.Token]{
	Store: token.Store,
}

func NewListCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list tokens",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return manager.ListCommand(cli)
		},
	}

	return cmd
}
