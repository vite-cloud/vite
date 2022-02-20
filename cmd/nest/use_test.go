package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/container"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewUseCommand(t *testing.T) {
	ct := CommandTest{
		Test: func(console *expect.Console) {
			Err(console.SendLine("")).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		Setup: func(ct *container.Container) []container.Option {
			return []container.Option{
				container.WithConfig(ct.NewConfig("github", "redwebcreation/nest-configs", "empty-config")),
			}
		},
		NewCommand: func(ct *container.Container) (*cobra.Command, error) {
			config, err := ct.Config()
			if err != nil {
				return nil, err
			}

			err = config.Clone()
			if err != nil {
				return nil, err
			}

			return NewUseCommand(ct), nil
		},
	}.Run(t)

	config, err := ct.Config()
	assert.NilError(t, err)

	// see https://github.com/redwebcreation/nest-configs/tree/empty-config
	assert.Equal(t, config.Commit, "3ea941eaf6d2bfcc97480ce5df49bee91d8f09e2")
}
