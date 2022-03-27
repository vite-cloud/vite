package proxy

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"os"
	"os/exec"
)

func runDownCommand(cli *cli.CLI) error {
	state, _, err := proxy.State()
	if err != nil {
		return err
	}

	if state == proxy.Absent {
		fmt.Fprintln(cli.Out(), "Proxy was not running.")
		return nil
	}

	for _, cmd := range proxy.DisableCmds {
		fmt.Fprintf(cli.Out(), "- running %s\n", cmd)
		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%w: %s", err, out)
		}

		if trimmed := bytes.TrimSpace(out); len(trimmed) > 0 {
			fmt.Fprintf(cli.Out(), "  %s\n", trimmed)
		}
	}

	err = os.Remove(proxy.ServiceFile)
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "\nSuccessfully removed service %s.\n", proxy.ServiceFile)

	return nil
}

func NewDownCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "stop proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDownCommand(cli)
		},
	}

	return cmd
}