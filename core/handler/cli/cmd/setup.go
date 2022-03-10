package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/locator"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
)

type wr struct {
	orig terminal.FileWriter
}

func (w *wr) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))

	return w.orig.Write(p)
}

func (w *wr) Fd() uintptr {
	return w.orig.Fd()
}

func runSetupCommand(cli *cli.CLI) error {
	fmt.Fprintln(cli.Out(), "Welcome to Vite!")
	var qs = []*survey.Question{
		{
			Name: "provider",
			Prompt: &survey.Select{
				Message: "Select your provider:",
				Options: []string{"github", "gitlab", "bitbucket"},
			},
			Validate: survey.Required,
		},
		{
			Name: "protocol",
			Prompt: &survey.Select{
				Message: "Select your protocol:",
				Options: []string{"ssh", "https", "auto"},
				Default: "ssh",
			},
			Validate: survey.Required,
		},
		{
			Name: "repository",
			Prompt: &survey.Input{
				Message: "Enter your repository:",
			},
			//Validate: func(ans interface{}) error {
			//	fmt.Println("INVALID")
			//	re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")
			//	if !re.MatchString(ans.(string)) {
			//		return fmt.Errorf("repository must be in format: username/repository")
			//	}
			//	return nil
			//},
		},
		//	{
		//		Name: "branch",
		//		Prompt: &survey.Input{
		//			Message: "Enter your branch:",
		//			Default: "main",
		//		},
		//		Validate: survey.Required,
		//	},
		//	{
		//		Name: "path",
		//		Prompt: &survey.Input{
		//			Message: "Enter a sub-path (optional):",
		//			Default: "",
		//		},
		//	},
	}

	answers := struct {
		Provider   string
		Protocol   string
		Repository string
		Branch     string
		Path       string
	}{}

	err := survey.Ask(qs, &answers, survey.WithStdio(cli.In(), &wr{cli.Out()}, cli.Err()))
	if err != nil {
		return err
	}

	l := locator.Locator{
		Provider:   locator.Provider(answers.Provider),
		Protocol:   answers.Protocol,
		Repository: answers.Repository,
		Branch:     answers.Branch,
		Path:       answers.Path,
	}

	err = l.Save()
	if err != nil {
		return err
	}

	//fmt.Fprintln(cli.Out(), "\nSetup successfully. You may now run `vite use` to select a commit to use.")
	//
	//return nil
	return nil
}

func NewSetupCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup vite",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetupCommand(cli)
		},
	}

	return cmd
}
