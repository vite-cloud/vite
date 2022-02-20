package nest

import (
	"fmt"
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/container"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	_ = CommandTest{
		Test: func(console *expect.Console) {
			Err(console.ExpectString(fmt.Sprintf("Nest version %s, build %s\n", build.Version, build.Commit))).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ct *container.Container) (*cobra.Command, error) {
			return NewVersionCommand(ct), nil
		},
	}.Run(t)
}
