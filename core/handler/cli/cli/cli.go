package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/static"
	"io"
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

	commands []*cobra.Command
}

// Out returns the current output writer.
func (c *CLI) Out() Stdout {
	return c.out
}

// In returns the current input reader.
func (c *CLI) In() Stdin {
	return c.in
}

// Err returns the current error writer.
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
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Version: fmt.Sprintf("%s, build %s", static.Version, static.Commit),
	}

	vite.SetVersionTemplate("Vite version {{.Version}}\n")
	vite.SetHelpCommand(&cobra.Command{
		Use:    "__help",
		Hidden: true,
	})

	vite.AddCommand(c.commands...)

	vite.SetArgs(args)
	vite.SetOut(c.Out())
	vite.SetIn(c.In())
	vite.SetErr(c.Err())

	err := vite.Execute()
	if statusErr, ok := err.(*StatusError); ok {
		fmt.Fprintln(c.Err(), statusErr.Status)
		return statusErr.StatusCode
	} else if err != nil {
		fmt.Fprintln(c.Err(), err)
		return 1
	}

	return 0
}

func (c *CLI) Add(command *cobra.Command) *CLI {
	c.commands = append(c.commands, command)

	return c
}

func New(out Stdout, in Stdin, err io.Writer) *CLI {
	return &CLI{
		out: out,
		in:  in,
		err: err,
	}
}

type StatusError struct {
	Status     string
	StatusCode int
}

func (s *StatusError) Error() string {
	return s.Status
}