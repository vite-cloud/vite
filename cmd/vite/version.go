package vite

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/build"
	"github.com/vite-cloud/vite/container"
)

func runVersionCommand(ct *container.Container) error {
	fmt.Fprintf(ct.Out(), "Vite version %s, build %s\n\n", build.Version, build.Commit)
	return nil
}

// NewVersionCommand creates a new `version` command.
func NewVersionCommand(ct *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print vite's version",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runVersionCommand(ct)
		},
	}

	return cmd
}
