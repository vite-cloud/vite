package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/go-zoup"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	server "github.com/vite-cloud/vite/core/domain/proxy"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"strconv"
)

type runOpts struct {
	deployment *deployment.Deployment
	HTTP       string
	HTTPS      string
	Unsecure   bool
}

func runRunCommand(cli *cli.CLI, opts *runOpts) error {
	proxy, err := server.New(cli.Out(), opts.deployment)
	if err != nil {
		return err
	}

	proxy.Logger.Log(zoup.DebugLevel, "starting", zoup.Fields{"http_port": opts.HTTP, "https_port": opts.HTTPS, "secure": !opts.Unsecure})

	proxy.Run(opts.HTTP, opts.HTTPS, opts.Unsecure)

	return nil
}

func newRunCommand(cli *cli.CLI) *cobra.Command {
	opts := &runOpts{}

	cmd := &cobra.Command{
		Use:   "run [id]",
		Short: "run the proxy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			dep, err := resource.Get[deployment.Deployment](deployment.Store, id)
			if err != nil {
				return err
			}

			conf, err := config.Get(dep.Locator)
			if err != nil {
				return err
			}

			opts.deployment = dep

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
