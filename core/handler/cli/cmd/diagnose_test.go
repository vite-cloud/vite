package cmd

import (
	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewMedicCommand(t *testing.T) {
	datadir.UseTestHome(t)

	CommandTest{
		NewCommand: NewDiagnoseCommand,
		ExpectsError: func(t *testing.T, err error) {
			assert.ErrorContains(t, err, "config locator hasn't been configured yet, run `vite setup` first")
		},
	}.Run(t)
}
