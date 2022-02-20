package cloud

import (
	"github.com/redwebcreation/nest/container"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new `cloud` command.
func NewRootCommand(ct *container.Container) *cobra.Command {
	root := &cobra.Command{
		Use:   "cloud",
		Short: "interact with nest cloud",
	}

	root.AddCommand(
		// login
		NewLoginCommand(ct),
	)

	return root
}
