package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/medic"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

// runMedicCommand handles the `medic` command.
func runMedicCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.Get()
	if err != nil {
		return err
	}

	diagnosis := medic.Diagnose(conf)

	fmt.Printf("%+v\n", diagnosis)

	return nil
}

// NewMedicCommand creates a new `medic` command.
func NewMedicCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose the configuration",
		RunE:  runMedicCommand,
	}

	return cmd
}
