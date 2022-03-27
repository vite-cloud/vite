package proxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runStatusCommand(cli *cli.CLI) error {
	status, uptime, err := proxy.State()
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "Status: %s (uptime: %s)\n", status, uptime)

	return nil
}

func NewStatusCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "proxy's status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusCommand(cli)
		},
	}
	return cmd
}
