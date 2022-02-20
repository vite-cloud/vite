package vite

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/container"
	"github.com/vite-cloud/vite/service"
	"io"
)

func runDeployCommand(ct *container.Container) error {
	config, err := ct.ServicesConfig()
	if err != nil {
		return err
	}

	deployment := service.NewDeployment(config, ct.ManifestManager(), ct.DockerClient())

	go func() {
		err = deployment.Start()
		if err != nil {
			deployment.Events <- service.Event{
				Service: nil,
				Value:   service.ErrDeploymentFailed,
			}
		}
	}()

	for event := range deployment.Events {
		if event.Value == service.ErrDeploymentFailed {
			fmt.Fprintln(ct.Out(), "Deployment failed")
			break
		}

		if event.Value == io.EOF {
			break
		}

		if event.Service != nil {
			fmt.Fprintf(ct.Out(), "%s: %v\n", event.Service.Name, event.Value)
		} else {
			fmt.Fprintf(ct.Out(), "global: %v\n", event.Value)
		}
	}

	return ct.ManifestManager().Save(deployment.Manifest)
}

// NewDeployCommand creates a new `deploy` command.
func NewDeployCommand(ct *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployCommand(ct)
		},
	}

	return cmd
}
