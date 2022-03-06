package cmd

import (
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
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

	console, _, err := vt10x.NewVT10XConsole()
	assert.NilError(t, err)

	defer func(console *expect.Console) {
		err = console.Close()
		if err != nil {
			panic(err)
		}
	}(console)

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		c.Test(&Expect{
			t:       t,
			console: console,
		})
	}()

	cmd := c.NewCommand(cli.New(console.Tty(), console.Tty(), console.Tty()))
	cmd.SetArgs([]string{})

	err = cmd.Execute()
	assert.NilError(t, err)

	console.Tty().Close()
	<-donec
}
