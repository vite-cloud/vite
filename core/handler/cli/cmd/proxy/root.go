package proxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"syscall"
)

func NewProxyCommand(c *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "manage proxy",
	}

	cmd.AddCommand(
		newRunCommand(c),
		newLogsCommand(c),
		newUpCommand(c),
		newStatusCommand(c),
		newDownCommand(c),
	)

	return cmd
}

func NeedsSystemdAccess(cmd *cobra.Command, args []string) error {
	if err := syscall.Access(proxy.ServiceFile, syscall.O_RDWR); err != nil {
		return fmt.Errorf("you must have elevated privileges to access %s", proxy.ServiceFile)
	}

	return nil
}
