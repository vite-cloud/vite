package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runEnableCommand(cli *cli.CLI) error {
	return nil
}

func NewEnableCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRunCommand(cli)
		},
	}

	return cmd
}
