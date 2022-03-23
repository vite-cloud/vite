package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runRunCommand(cli *cli.CLI) error {
	return nil
}

func NewRunCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRunCommand(cli)
		},
	}

	return cmd
}
