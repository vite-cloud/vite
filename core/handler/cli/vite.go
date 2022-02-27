package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/static"
	"io"
	"os"
)

// Stdout provides a minimal interface for writing to stdout.
type Stdout interface {
	io.Writer
	Fd() uintptr
}

// Stdin provides a minimal interface for reading stdin.
type Stdin interface {
	io.Reader
	Fd() uintptr
}

// CLI is the command line interface for Vite.
type CLI struct {
	out Stdout
	in  Stdin
	err io.Writer
}

func (c *CLI) Out() Stdout {
	return c.out
}

func (c *CLI) In() Stdin {
	return c.in
}

func (c *CLI) Err() io.Writer {
	return c.err
}

// Run the vite CLI with the given command line arguments.
// and returns the exit code for the command.
func (c *CLI) Run(args []string) int {
	vite := &cobra.Command{
		Use:           "vite",
		Short:         "Kubernetes alternative for small companies.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       fmt.Sprintf("%s, build %s", static.Version, static.Commit),
	}

	vite.SetHelpCommand(&cobra.Command{
		Use:    "__help",
		Hidden: true,
	})

	vite.SetArgs(args)
	vite.SetOut(c.Out())
	vite.SetIn(c.In())
	vite.SetErr(c.Err())

	return 0
}

func New() *CLI {
	return &CLI{
		out: os.Stdout,
		in:  os.Stdin,
		err: os.Stderr,
	}
}
