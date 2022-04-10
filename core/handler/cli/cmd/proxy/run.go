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
	ID       int64
	hasID    bool
	HTTP     string
	HTTPS    string
	Unsecure bool
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

	proxy.Logger.Log(log.DebugLevel, "starting", log.Fields{"http_port": opts.HTTP, "https_port": opts.HTTPS, "secure": !opts.Unsecure})

	proxy.Run(opts.HTTP, opts.HTTPS, opts.Unsecure)

	return nil
}

func newRunCommand(cli *cli.CLI) *cobra.Command {
	opts := &runOpts{}

	cmd := &cobra.Command{
		Use:   "run [id]",
		Short: "run the proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.hasID = len(args) > 0
			if opts.hasID {
				id, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				opts.ID = int64(id)
			}

			conf, err := config.Get()
			if err != nil {
				return err
			}

			if opts.HTTP == "" {
				opts.HTTP = conf.Proxy.HTTP
			}

			if opts.HTTPS == "" {
				opts.HTTPS = conf.Proxy.HTTPS
			}

			return runRunCommand(cli, opts)
		},
	}

	cmd.Flags().StringVar(&opts.HTTP, "http", "", "http port")
	cmd.Flags().StringVar(&opts.HTTPS, "https", "", "https port")
	cmd.Flags().BoolVarP(&opts.Unsecure, "unsecure", "u", false, "use self-signed certificate")

	return cmd
}
