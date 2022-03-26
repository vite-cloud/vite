package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runDisableCommand(cli *cli.CLI) error {
	return nil
}

func NewDisableCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisableCommand(cli)
		},
	}

	return cmd
}
