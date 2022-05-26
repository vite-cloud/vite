package deployments

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

var manager = resource.Manager[*deployment.Deployment]{
	Store: deployment.Store,
}

func newListCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list deployments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return manager.ListCommand(cli)
		},
	}

	return cmd
}
