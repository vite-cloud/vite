package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/static"
)

func NewVersionCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of vite",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cli.Out(), "Vite version %s, commit %s\n", static.Version, static.Commit)

			return nil
		},
	}

	return cmd
}
