package nest

import (
	"fmt"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/container"
	"github.com/spf13/cobra"
)

func runVersionCommand(ct *container.Container) error {
	fmt.Fprintf(ct.Out(), "Nest version %s, build %s\n\n", build.Version, build.Commit)
	return nil
}

// NewVersionCommand creates a new `version` command.
func NewVersionCommand(ct *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print nest's version",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runVersionCommand(ct)
		},
	}

	return cmd
}
