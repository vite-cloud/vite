package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"os"
	"os/exec"
	"os/user"
)

type upOptions struct {
	verbose bool
	user    string
}

func runUpCommand(cli *cli.CLI, opts *upOptions) error {
	state, _, err := proxy.State()
	if err != nil {
		return err
	}

	if state != proxy.Absent {
		return fmt.Errorf("proxy is already running (state: %s)", state)
	}

	if opts.user == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}

		opts.user = u.Username
	}

	config, err := proxy.DaemonConfig(opts.user)
	if err != nil {
		return err
	}

	if opts.verbose {
		fmt.Fprintf(cli.Out(), "Installing the below service at %s:\n", proxy.ServiceFile)
		fmt.Fprintf(cli.Out(), "%s\n\n", config)
	}

	if err = os.Remove(proxy.ServiceFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.WriteFile(proxy.ServiceFile, config, 0644)
	if err != nil {
		return err
	}

	for _, cmd := range proxy.EnableCmds {
		fmt.Fprintf(cli.Out(), "- running %s\n", cmd)

		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%w: %s", err, out)
		}

		if trimmed := bytes.TrimSpace(out); len(trimmed) > 0 {
			fmt.Fprintf(cli.Out(), "  %s\n", trimmed)
		}
	}

	fmt.Fprintf(cli.Out(), "\nSuccessfully installed service %s.\n", proxy.ServiceFile)

	return nil
}

func newUpCommand(cli *cli.CLI) *cobra.Command {
	opts := &upOptions{}

	cmd := &cobra.Command{
		Use:               "up",
		Short:             "start proxy",
		PersistentPreRunE: NeedsSystemdAccess,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpCommand(cli, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().StringVarP(&opts.user, "user", "u", "", "user running daemon")

	return cmd
}
