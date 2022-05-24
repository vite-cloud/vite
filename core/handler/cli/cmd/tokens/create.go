package tokens

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/token"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"time"
)

func runCreateCommand(cli *cli.CLI) error {
	tok := token.Token{
		CreatedAt: time.Now(),
		Value:     token.NewWithPrefix("tok"),
	}

	label := survey.Input{
		Message: "Enter a label:",
		Default: "default",
	}
	err := survey.AskOne(&label, &tok.Label, survey.WithStdio(cli.In(), cli.Out(), cli.Err()), survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	err = tok.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "The token %s has been created.\n", tok.Value)

	return nil
}

func NewCreateCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCommand(cli)
		},
	}

	return cmd
}
