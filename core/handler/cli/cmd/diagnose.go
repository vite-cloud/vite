package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/medic"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

// runDiagnoseCommand handles the `diagnose` command.
func runDiagnoseCommand(cli *cli.CLI) error {
	conf, err := config.Get()
	if err != nil {
		return err
	}

	diagnosis := medic.Diagnose(conf)

	fmt.Fprintln(cli.Out(), "Errors:")

	if len(diagnosis.Errors) == 0 {
		fmt.Fprintln(cli.Out(), "  - no errors")
	} else {
		for _, err := range diagnosis.Errors {
			fmt.Fprintf(cli.Out(), "  - %s\n", err.Title)
			if err.Error != nil {
				fmt.Fprintf(cli.Out(), "    %s\n", err.Error)
			}
		}
	}

	fmt.Fprintln(cli.Out(), "\nWarnings:")

	if len(diagnosis.Warnings) == 0 {
		fmt.Fprintln(cli.Out(), "  - no warnings")
	} else {
		for _, warn := range diagnosis.Warnings {
			fmt.Fprintf(cli.Out(), "  - %s\n", warn.Title)
			if warn.Advice != "" {
				fmt.Fprintf(cli.Out(), "    %s\n", warn.Advice)
			}
		}
	}

	return nil
}

// NewDiagnoseCommand creates a new `diagnose` command.
func NewDiagnoseCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "diagnose the configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiagnoseCommand(cli)
		},
	}

	return cmd
}
