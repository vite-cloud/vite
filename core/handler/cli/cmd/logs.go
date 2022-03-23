package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/log"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

type logsOptions struct {
	follow bool
	n      int
}

func runLogsCommand(cli *cli.CLI, opts logsOptions) error {
	t, err := log.Tail(opts.follow, opts.n)
	if err != nil {
		return err
	}

	defer t.Cleanup()

	for line := range t.Lines {
		fmt.Fprintln(cli.Out(), line.Text)
	}

	return nil
}

func NewLogsCommand(cli *cli.CLI) *cobra.Command {
	opts := logsOptions{}

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "read proxy logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogsCommand(cli, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.follow, "follow", "f", false, "follow logs")
	cmd.Flags().IntVarP(&opts.n, "lines", "n", 10, "number of lines to show")

	return cmd
}
