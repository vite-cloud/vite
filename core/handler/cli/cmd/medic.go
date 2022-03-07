package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/locator"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runMedicCommand(cmd *cobra.Command, args []string) error {
	l, err := locator.LoadFromStore()
	if err != nil {
		return err
	}

	conf, err := config.Get(l)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", conf)

	return nil
}

func NewMedicCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose the configuration",
		RunE:  runMedicCommand,
	}

	return cmd
}
