package cli

import (
	"os"

	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"github.com/vite-cloud/vite/core/handler/cli/cmd"
)

func New() *cli.CLI {
	c := cli.New(os.Stdout, os.Stdin, os.Stderr)

	c.Add(
		cmd.NewVersionCommand(c),
		cmd.NewMedicCommand(c),
		cmd.NewSetupCommand(c),
	)

	return c
}
