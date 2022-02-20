package vite

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/container"
	"github.com/vite-cloud/vite/service"
)

func runMedicCommand(ct *container.Container) error {
	sc, err := ct.ServicesConfig()
	if err != nil {
		return err
	}

	diagnostic := service.DiagnoseConfig(sc)

	fmt.Fprintln(ct.Out())
	fmt.Fprintln(ct.Out(), "Errors:")

	if len(diagnostic.Errors) == 0 {
		fmt.Fprintln(ct.Out(), "  - no errors")
	} else {
		for _, err := range diagnostic.Errors {
			fmt.Fprintf(ct.Out(), "  -  %s\n", err.Title)
			if err.Error != nil {
				fmt.Fprintf(ct.Out(), "    %s\n", err.Error)
			}
		}
	}

	fmt.Fprintln(ct.Out(), "\nWarnings:")

	if len(diagnostic.Warnings) == 0 {
		fmt.Fprintln(ct.Out(), "  - no warnings")
	} else {
		for _, warn := range diagnostic.Warnings {
			fmt.Fprintf(ct.Out(), "  -  %s\n", warn.Title)
			if warn.Advice != "" {
				fmt.Fprintf(ct.Out(), "    %s\n", warn.Advice)
			}
		}
	}

	return nil
}

// NewMedicCommand creates a new `medic` command.
func NewMedicCommand(ct *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose your config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMedicCommand(ct)
		},
	}

	return cmd
}
