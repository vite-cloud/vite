package deployments

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/handler/cli/cmd"
)

func runRollbackCommand(cli *cli.CLI, id string) error {
	dep, err := deployment.Get(id)
	if err != nil {
		return err
	}

	err = dep.Locator.Save()
	if err != nil {
		return err
	}

	return cmd.NewDeployCommand(cli).Execute()
}

func newRollbackCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback [deployment]",
		Short: "rollback a deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRollbackCommand(cli, args[0])
		},
	}

	return cmd
}
