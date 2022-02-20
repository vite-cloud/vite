package proxy

import (
	"github.com/redwebcreation/nest/container"
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the proxy root command
func NewRootCommand(ct *container.Container) *cobra.Command {
	root := &cobra.Command{
		Use:   "proxy",
		Short: "manage the proxy",
	}

	root.AddCommand(
		// run
		NewRunCommand(ct),
	)

	return root
}
