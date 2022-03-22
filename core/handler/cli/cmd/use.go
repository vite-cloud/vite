package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/locator"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

func runUseCommand(cli *cli.CLI) error {
	l, err := locator.LoadFromStore()
	if err != nil {
		return err
	}

	commits, err := l.Commits()
	if err != nil {
		return err
	}

	var response string

	err = survey.AskOne(&survey.Select{
		Message: "Select a commit",
		Options: commits.AsOptions(),
	}, &response, survey.WithStdio(cli.In(), cli.Out(), cli.Err()))
	if err != nil {
		return err
	}

	response = response[:40]

	l.Commit = response

	err = l.Save()
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.Out(), "Successfully set commit to", response)

	return nil
}
func NewUseCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use: "use",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUseCommand(cli)
		},
	}
	return cmd
}
