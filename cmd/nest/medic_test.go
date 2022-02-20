package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/container"
	"github.com/redwebcreation/nest/service"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewMedicCommand(t *testing.T) {
	_ = CommandTest{
		Test: func(console *expect.Console) {
			Err(console.ExpectString("Errors:")).Check(t)
			Err(console.ExpectString("- no errors")).Check(t)
			Err(console.ExpectString("Warnings:")).Check(t)
			Err(console.ExpectString("- no warnings")).Check(t)
		},
		ContextBuilder: []container.Option{
			// As the config is not nil, the container does not try to create it
			container.WithConfig(&service.Locator{}),
			container.WithServicesConfig(&service.Config{}),
		},
		NewCommand: func(ct *container.Container) (*cobra.Command, error) {
			return NewMedicCommand(ct), nil
		},
	}.Run(t)
}
