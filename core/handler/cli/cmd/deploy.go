package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/locator"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runDeployCommand(cli *cli.CLI) error {
	loc, err := locator.LoadFromStore()
	if err != nil {
		return err
	}

	events := make(chan deployment.Event)

	go deployment.Deploy(events, loc)

	for event := range events {
		if event.ID == deployment.FinishEvent {
			break
		}

		fmt.Fprintf(cli.Out(), "%s(%s): %v\n", event.Label(), event.ID, event.Data)

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
