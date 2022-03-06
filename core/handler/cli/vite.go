package cli

import (
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/handler/cli/cmd"
	"os"
)

func New() *cli.CLI {
	c := cli.New(os.Stdout, os.Stdin, os.Stderr)

	c.Add(
		cmd.NewVersionCommand(c),
	)

	return c
}
