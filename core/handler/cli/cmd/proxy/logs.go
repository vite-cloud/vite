package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

type logsOptions struct {
	follow bool
}

var follow bool

func runLogsCommand(cli *cli.CLI, opts logsOptions) error {
	return nil
}

func NewLogsCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [-f] [-n <size>]",
		Short: "read proxy logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := logsOptions{
				follow: follow,
			}

			return runLogsCommand(cli, opts)
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow logs")

	return cmd
}
