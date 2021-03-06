package cli

import (
	"fmt"
	"github.com/vite-cloud/go-zoup"
	"github.com/vite-cloud/vite/core/domain/log"
	"io"

	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/static"
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

	vite.SetVersionTemplate("Vite version {{.id}}\n")
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

	command := "vite"
	if len(args) > 0 {
		command = args[0]
	}

	if err == nil {
		log.Log(zoup.InfoLevel, "command ran successfully", zoup.Fields{
			"command": command,
		})

		return 0
	} else if statusErr, ok := err.(*StatusError); ok {
		fmt.Fprintf(c.Err(), "Error: %s\n", statusErr.Status)

		log.Log(zoup.ErrorLevel, "command failed", zoup.Fields{
			"command": command,
			"err":     statusErr.Status,
			"code":    statusErr.StatusCode,
		})

		return statusErr.StatusCode
	} else {
		fmt.Fprintf(c.Err(), "Error: %s\n", err)

		log.Log(zoup.ErrorLevel, "command failed", zoup.Fields{
			"command": command,
			"err":     err,
			"code":    1,
		})

		return 1
	}
}

// Add adds the given commands to the CLI.
func (c *CLI) Add(commands ...*cobra.Command) *CLI {
	c.commands = append(c.commands, commands...)

	return c
}

// New returns a new CLI with the given standard IO.
func New(out Stdout, in Stdin, err io.Writer) *CLI {
	return &CLI{
		out: out,
		in:  in,
		err: err,
	}
}

// StatusError is an error type that contains an exit code.
// It is used to exit with a custom exit code.
type StatusError struct {
	Status     string
	StatusCode int
}

// Error implements the error interface.
func (s *StatusError) Error() string {
	return s.Status
}
