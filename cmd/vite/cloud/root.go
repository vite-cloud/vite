package cloud

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/container"
)

// NewRootCommand creates a new `cloud` command.
func NewRootCommand(ct *container.Container) *cobra.Command {
	root := &cobra.Command{
		Use:   "cloud",
		Short: "interact with vite cloud",
	}

	root.AddCommand(
		// login
		NewLoginCommand(ct),
	)

	return root
}
