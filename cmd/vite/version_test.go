package vite

import (
	"fmt"
	"github.com/Netflix/go-expect"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/build"
	"github.com/vite-cloud/vite/container"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	_ = CommandTest{
		Test: func(console *expect.Console) {
			Err(console.ExpectString(fmt.Sprintf("Vite version %s, build %s\n", build.Version, build.Commit))).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ct *container.Container) (*cobra.Command, error) {
			return NewVersionCommand(ct), nil
		},
	}.Run(t)
}
