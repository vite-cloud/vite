package proxy

import (
	"github.com/redwebcreation/nest/container"
	"github.com/redwebcreation/nest/service"
	"github.com/spf13/cobra"
)

type runOptions struct {
	deployment string
	HTTP       string
	HTTPS      string
	selfSigned bool
}

func runRunCommand(ct *container.Container, opts *runOptions) error {

	var manifest *service.Manifest
	var err error

	if opts.deployment != "" {
		manifest, err = ct.ManifestManager().LoadWithID(opts.deployment)
		if err != nil {
			return err
		}
	} else {
		manifest, err = ct.ManifestManager().Latest()
		if err != nil {
			return err
		}
	}

	config, err := ct.ServicesConfig()
	if err != nil {
		return err
	}

	config.Proxy.HTTP = opts.HTTP
	config.Proxy.HTTPS = opts.HTTPS
	config.Proxy.SelfSigned = opts.selfSigned

	proxy, err := ct.NewProxy(manifest)
	if err != nil {
		return err
	}

	proxy.Run()

	return nil
}

// NewRunCommand creates a new `run` command
func NewRunCommand(ct *container.Container) *cobra.Command {
	opts := &runOptions{}

	cmd := &cobra.Command{
		Use:   "run [deployment]",
		Short: "Starts the proxy",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			servicesConfig, err := ct.ServicesConfig()
			if err != nil {
				return err
			}

			if opts.HTTP == "" {
				opts.HTTP = servicesConfig.Proxy.HTTP
			}

			if opts.HTTPS == "" {
				opts.HTTPS = servicesConfig.Proxy.HTTPS
			}

			if !opts.selfSigned && servicesConfig.Proxy.SelfSigned {
				opts.selfSigned = true
			}

			if len(args) > 0 {
				opts.deployment = args[0]
			}

			return runRunCommand(ct, opts)
		},
	}

	cmd.Flags().StringVar(&opts.HTTP, "http", "", "HTTP port")
	cmd.Flags().StringVar(&opts.HTTPS, "https", "", "HTTPS port")
	cmd.Flags().BoolVarP(&opts.selfSigned, "self-signed", "u", false, "Use a self-signed certificate")

	return cmd
}
