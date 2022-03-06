package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func NewVersionCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of vite",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(cli.Out(), cmd.VersionTemplate())

			return nil
		},
	}

	return cmd
}
