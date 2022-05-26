package tokens

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/domain/token"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runRevokeCommand(cli *cli.CLI, label string) error {
	tok, err := resource.Get[token.Token](token.Store, label)
	if err != nil {
		return err
	}

	err = resource.Delete[*token.Token](token.Store, tok, func(t *token.Token) string {
		return t.Label
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "The token %s has been revoked.\n", tok.Label)

	return nil
}

func NewRevokeCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [label]",
		Short: "revoke a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRevokeCommand(cli, args[0])
		},
	}

	return cmd
}
