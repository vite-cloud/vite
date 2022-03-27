package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runDeployCommand(cli *cli.CLI) error {
	conf, err := config.Get()
	if err != nil {
		return err
	}

	events := make(chan deployment.Event)

	go func() {
		err = deployment.Deploy(events, conf.Services)
		if err != nil {
			events <- deployment.Event{
				ID:   deployment.ErrorEvent,
				Data: err,
			}
		} else {
			events <- deployment.Event{
				ID: deployment.FinishEvent,
			}
		}
	}()

	for event := range events {
		if event.ID == deployment.FinishEvent {
			break
		}

		fmt.Fprintf(cli.Out(), "%s: %v\n", event.Label(), event.Data)

		if event.IsError() {
			break
		}
	}

	return nil
}

func NewDeployCommand(cli *cli.CLI) *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "deploy services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployCommand(cli)
		},
	}
}
