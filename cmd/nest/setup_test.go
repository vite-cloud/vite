package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/container"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewSetupCommand(t *testing.T) {
	ct := CommandTest{
		Test: func(console *expect.Console) {
			Err(console.ExpectString("Select your provider:")).Check(t)
			Err(console.SendLine("")).Check(t)
			Err(console.ExpectString("Enter your repository:")).Check(t)
			Err(console.SendLine("redwebcreation/nest-configs")).Check(t)
			Err(console.ExpectString("Enter your branch:")).Check(t)
			Err(console.SendLine("empty-config")).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ct *container.Container) (*cobra.Command, error) {
			return NewSetupCommand(ct), nil
		},
	}.Run(t)

	config, err := ct.Config()
	assert.NilError(t, err)

	assert.Equal(t, "redwebcreation/nest-configs", config.Repository)
	assert.Equal(t, "empty-config", config.Branch)
	assert.Equal(t, "github", config.Provider)
}
