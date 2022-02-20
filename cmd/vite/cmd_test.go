package vite

import (
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/container"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

type CommandTest struct {
	Test           func(console *expect.Console)
	NewCommand     func(ct *container.Container) (*cobra.Command, error)
	ContextBuilder []container.Option
	Setup          func(ct *container.Container) []container.Option
}

func (c CommandTest) Run(t *testing.T) *container.Container {
	dir, err := os.MkdirTemp("", "vite")
	assert.NilError(t, err)

	console, _, err := vt10x.NewVT10XConsole()
	assert.NilError(t, err)

	defer console.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		c.Test(console)
	}()

	ct, err := container.New(container.WithConfigHome(dir), container.WithStdio(console.Tty(), console.Tty(), console.Tty()))
	assert.NilError(t, err)

	for _, option := range c.ContextBuilder {
		err = option(ct)
		assert.NilError(t, err)
	}

	if c.Setup != nil {
		for _, opt := range c.Setup(ct) {
			err = opt(ct)
			assert.NilError(t, err)
		}
	}

	cmd, err := c.NewCommand(ct)
	assert.NilError(t, err)
	cmd.SetArgs([]string{})

	err = cmd.Execute()
	assert.NilError(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	console.Tty().Close()
	<-donec

	return ct
}

type ConsoleError struct {
	err error
}

func Err(_ interface{}, err error) ConsoleError {
	return ConsoleError{err}
}

func (c ConsoleError) Check(t *testing.T) {
	// this error is expected??
	// at least it does not change the outcome of the test
	// see https://github.com/creack/pty/issues/21
	if c.err != nil && c.err.Error() == "read /dev/ptmx: input/output error" {
		return
	}

	assert.NilError(t, c.err)
}
