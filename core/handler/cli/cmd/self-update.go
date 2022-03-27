package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v43/github"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"io"
	"net/http"
	"os"
)

type selfUpdateOptions struct {
	version    string
	hasVersion bool
}

func runSelfUpdate(cli *cli.CLI, opts *selfUpdateOptions) error {
	client := github.NewClient(nil)

	var release *github.RepositoryRelease
	var err error

	if opts.hasVersion {
		release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "vite-cloud", "vite", opts.version)
	} else {
		release, _, err = client.Repositories.GetLatestRelease(context.Background(), "vite-cloud", "vite")
	}

	if err != nil {
		return err
	}

	binary := release.Assets[0]

	if binary.GetState() != "uploaded" {
		return fmt.Errorf("the binary is not available yet, please try again in a few seconds")
	}

	fmt.Fprintf(cli.Out(), "Downloading %s...\n", binary.GetName())

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	response, err := http.Get(binary.GetBrowserDownloadURL())
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(executable+".tmp", body, 0755)
	if err != nil {
		return err
	}

	err = os.Rename(executable+".tmp", executable)
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "Successfully updated to %s\n", release.GetTagName())

	return nil
}

func NewSelfUpdateCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "update vite",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &selfUpdateOptions{hasVersion: len(args) > 0}

			if opts.hasVersion {
				opts.version = args[0]
			}

			return runSelfUpdate(cli, opts)
		},
	}
	return cmd
}
