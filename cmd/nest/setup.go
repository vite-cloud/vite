package nest

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/redwebcreation/nest/container"
	"regexp"

	"github.com/spf13/cobra"
)

type setupOptions struct {
	UsesFlags  bool
	Provider   string
	Repository string
	Branch     string
	Commit     string
}

func runSetupCommand(ct *container.Container, opts *setupOptions) error {
	oldConfig, err := ct.Config()
	hasConfig := err == nil

	if !opts.UsesFlags {
		// A default value for a select must be one of the options
		var provider = "github"

		if hasConfig {
			provider = oldConfig.Provider
		}

		prompt := &survey.Select{
			Message: "Select your provider:",
			Options: []string{"github", "gitlab", "bitbucket"},
			Default: provider,
		}
		err = survey.AskOne(prompt, &opts.Repository, survey.WithValidator(survey.Required), survey.WithStdio(ct.In(), ct.Out(), ct.Err()))
		if err != nil {
			return err
		}
	}

	if !opts.UsesFlags {
		var repository string

		if hasConfig {
			repository = oldConfig.Repository
		}

		prompt := &survey.Input{
			Message: "Enter your repository:",
			Default: repository,
		}
		err = survey.AskOne(prompt, &opts.Repository, survey.WithValidator(func(ans any) error {
			re := regexp.MustCompile(`[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+`)
			if !re.MatchString(ans.(string)) {
				return fmt.Errorf("repository name must be alphanumeric and can contain hyphens and underscores")
			}
			return nil
		}), survey.WithStdio(ct.In(), ct.Out(), ct.Err()))
		if err != nil {
			return err
		}
	}

	if !opts.UsesFlags {
		var branch string

		if hasConfig {
			branch = oldConfig.Branch
		}

		prompt := &survey.Input{
			Message: "Enter your branch:",
			Default: branch,
		}
		err = survey.AskOne(prompt, &opts.Branch, survey.WithValidator(survey.Required), survey.WithStdio(ct.In(), ct.Out(), ct.Err()))
		if err != nil {
			return err
		}
	}

	config := ct.NewConfig(opts.Provider, opts.Repository, opts.Branch)
	if err = config.Save(); err != nil {
		return err
	}

	err = config.Clone()
	if err != nil {
		return err
	}

	fmt.Fprintln(ct.Out(), "\nYou now need to run `nest use` to specify which version of the oldConfig you want to use.")

	return nil
}

// NewSetupCommand creates a new `setup` command.
func NewSetupCommand(ct *container.Container) *cobra.Command {
	opts := &setupOptions{}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "set up nest",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.UsesFlags = cmd.Flags().NFlag() > 0
			return runSetupCommand(ct, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Provider, "provider", "p", "github", "provider")
	cmd.Flags().StringVarP(&opts.Repository, "repository", "r", "", "repository")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "main", "branch")

	return cmd
}
