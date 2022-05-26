package deployments

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runShowCommand(cli *cli.CLI, ID string) error {
	dep, err := resource.Get[deployment.Deployment](deployment.Store, ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "%+v", dep)

	return nil
}

func newShowCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [deployment]",
		Short: "show details about a given deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShowCommand(cli, args[0])
		},
	}

	return cmd
}
