package cloud

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/cloud"
	"github.com/vite-cloud/vite/container"
)

type loginOptions struct {
	id          string
	accessToken string
}

func runLoginCommand(ct *container.Container, opts *loginOptions) error {
	err := ct.SetCloudCredentials(opts.id, opts.accessToken)
	if err != nil {
		return err
	}

	client, err := ct.CloudClient()
	if err != nil {
		return err
	}

	err = client.Ping()
	if err == cloud.ErrResourceNotFound {
		fmt.Fprintln(ct.Out(), "Invalid token.")
	} else if err != nil {
		return err
	} else {
		fmt.Fprintln(ct.Out(), "Successfully logged in.")
	}

	return nil
}

// NewLoginCommand creates a new `login` command.
func NewLoginCommand(ct *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to vite cloud",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args[0]) != 45 {
				return fmt.Errorf("invalid token")
			}

			return runLoginCommand(ct, &loginOptions{
				id:          args[0][:22],
				accessToken: args[0][23:],
			})
		},
	}

	return cmd
}
