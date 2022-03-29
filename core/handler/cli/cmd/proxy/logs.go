package proxy

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/log"
	"github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"os"
)

type logsOptions struct {
	stream   bool
	backfill int
}

func runLogsCommand(cli *cli.CLI, opts logsOptions) error {
	dir, err := log.Store.Dir()
	if err != nil {
		return err
	}

	stream, err := log.Tail(dir+"/"+proxy.LogFile, log.TailOptions{
		Stream:   opts.stream,
		Backfill: opts.backfill,
	})
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("proxy has never been started, please start it first")
		}

		return err
	}

	for line := range stream {
		fmt.Fprintln(cli.Out(), line)
	}

	return nil
}

func newLogsCommand(cli *cli.CLI) *cobra.Command {
	opts := logsOptions{}

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "read proxy logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogsCommand(cli, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.stream, "follow", "f", false, "stream logs")
	cmd.Flags().IntVarP(&opts.backfill, "backfill", "n", 10, "number of lines to show")

	return cmd
}
