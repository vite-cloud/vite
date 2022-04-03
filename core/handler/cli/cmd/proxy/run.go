package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/log"
	server "github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"strconv"
)

type runOpts struct {
	ID    int
	hasID bool
}

func runRunCommand(cli *cli.CLI, opts *runOpts) error {
	if !opts.hasID {
		id, err := deployment.Latest()
		if err != nil {
			return err
		}

		opts.ID = id
	}

	depl, err := deployment.Get(opts.ID)
	if err != nil {
		return err
	}

	proxy, err := server.New(cli.Out(), depl)
	if err != nil {
		return err
	}

	conf, err := config.Get()
	if err != nil {
		return err
	}

	proxy.Logger.Log(log.DebugLevel, "starting", log.Fields{"http_port": conf.Proxy.HTTP, "https_port": conf.Proxy.HTTPS})

	proxy.Run(conf.Proxy.HTTP, conf.Proxy.HTTPS)

	return nil
}

func newRunCommand(cli *cli.CLI) *cobra.Command {
	opts := &runOpts{}

	cmd := &cobra.Command{
		Use:    "run [id]",
		Short:  "run the proxy",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.hasID = len(args) > 0
			if opts.hasID {
				id, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				opts.ID = id
			}

			return runRunCommand(cli, opts)
		},
	}

	return cmd
}
