package command

import (
	"fmt"
	"os"

	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func Configure(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	_, disableConfigLocator := cmd.Annotations["config"]
	_, disableMedic := cmd.Annotations["medic"]

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if !disableConfigLocator {
			if _, err := os.Stat(global.ConfigLocatorConfigFile); err != nil {
				return fmt.Errorf("run `nest setup` to setup nest")
			}

			if err := pkg.LoadConfig(); err != nil {
				return err
			}
		}

		if disableMedic {
			return nil
		}

		return pkg.DiagnoseConfiguration().MustPass()
	}
}