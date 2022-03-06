package cmd

import (
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/kr/pty"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

type Expect struct {
	console *expect.Console
	t       *testing.T
}

func (e *Expect) String(s string) *Expect {
	_, err := e.console.ExpectString(s)
	e.check(err)

	return e
}

func (e *Expect) check(err error) {
	if err != nil && err.Error() == "read /dev/ptmx: input/output error" {
		return
	}

	if err != nil {
		e.t.Fatal(err)
	}
}

func (e *Expect) EOF() *Expect {
	_, err := e.console.ExpectEOF()
	e.check(err)
	return e
}

type CommandTest struct {
	Test       func(console *Expect)
	NewCommand func(cli *cli.CLI) *cobra.Command
}

func (c CommandTest) Run(t *testing.T) {
	dir, err := os.MkdirTemp("", "vite-home")
	assert.NilError(t, err)

	datadir.SetHomeDir(dir)

	console, err := newConsole()
	assert.NilError(t, err)

	defer func(console *expect.Console) {
		err = console.Close()
		assert.NilError(t, err)
	}(console)

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.Test(&Expect{
			console: console,
			t:       t,
		})
	}()

	cmd := c.NewCommand(cli.New(console.Tty(), console.Tty(), console.Tty()))
	cmd.SetArgs(nil)

	err = cmd.Execute()
	assert.NilError(t, err)

	err = console.Tty().Close()
	assert.NilError(t, err)
	<-donec
}

func newConsole() (*expect.Console, error) {
	ptm, pts, err := pty.Open()
	if err != nil {
		return nil, err
	}

	term := vt10x.New(vt10x.WithWriter(ptm))

	c, err := expect.NewConsole(expect.WithStdin(ptm), expect.WithStdout(term), expect.WithCloser(pts, ptm))
	if err != nil {
		return nil, err
	}

	return c, nil
}
