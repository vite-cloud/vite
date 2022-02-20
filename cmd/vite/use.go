package vite

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/container"
)

type useOptions struct {
	commit string
}

func runUseCommand(ct *container.Container, opts *useOptions) error {
	config, err := ct.Config()
	if err != nil {
		return err
	}

	err = config.Pull()
	if err != nil {
		return err
	}

	commits, err := config.Git.ListCommits(config.StorePath(), config.Branch)
	if err != nil {
		return err
	}

	fmt.Fprintf(ct.Out(), "Inspecting %d commits...\n", len(commits))

	if opts.commit == "" {
		prompt := survey.Select{
			Message: "Select a commit to use",
			Options: commits.Hashes(),
		}
		err = survey.AskOne(&prompt, &opts.commit, survey.WithStdio(ct.In(), ct.Out(), ct.Err()))
		if err != nil {
			return err
		}
	}

	if len(opts.commit) != 40 {
		return fmt.Errorf("invalid commit hash (must be full): %s", opts.commit)
	}

	err = config.LoadCommit(opts.commit)
	if err != nil {
		return err
	}

	fmt.Fprintf(ct.Out(), "Updated the config. Now using %s.\n", opts.commit[:7])

	return nil
}

// NewUseCommand creates a new `use` command
func NewUseCommand(ct *container.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [commit]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "use a specific commit for the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &useOptions{
				commit: "",
			}

			if len(args) > 0 {
				opts.commit = args[0]
			}

			return runUseCommand(ct, opts)
		},
	}

	return cmd
}
