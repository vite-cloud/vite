package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runMedicCommand(cmd *cobra.Command, args []string) error {

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
