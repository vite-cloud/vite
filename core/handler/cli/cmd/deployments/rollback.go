package deployments

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/handler/cli/cmd"
	"strconv"
)

func runRollbackCommand(cli *cli.CLI, ID int) error {
	dep, err := deployment.Get(ID)
	if err != nil {
		return err
	}

	err = dep.Config.Locator.Save()
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
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			return runRollbackCommand(cli, id)
		},
	}

	return cmd
}
