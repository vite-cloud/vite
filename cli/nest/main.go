package main

import (
	"fmt"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/cmd/nest"
	"github.com/redwebcreation/nest/cmd/nest/cloud"
	"github.com/redwebcreation/nest/cmd/nest/proxy"
	"github.com/redwebcreation/nest/container"
	"github.com/redwebcreation/nest/loggy"
	"github.com/spf13/cobra"
	"os"
)

func newNestCommand(ct *container.Container) *cobra.Command {
	cli := &cobra.Command{
		Use:           "nest",
		Short:         "Service orchestrator",
		SilenceUsage:  true,
		SilenceErrors: true,
		Long:          "Nest is a powerful service orchestrator for a single server.",
		Version:       fmt.Sprintf("%s, build %s", build.Version, build.Commit),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ct.Logger().Print(loggy.NewEvent(
				loggy.DebugLevel,
				"command invoked",
				loggy.Fields{
					"tag":     "command.invoke",
					"command": cmd.Name(),
				},
			))
		},
	}

	cli.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	cli.PersistentFlags().StringP("config", "c", ct.Home(), "set the loggy config path")

	cli.AddCommand(
		// version
		nest.NewVersionCommand(ct),

		// setup
		nest.NewSetupCommand(ct),

		// use
		nest.NewUseCommand(ct),

		// medic
		nest.NewMedicCommand(ct),

		// self-update
		nest.NewSelfUpdateCommand(ct),

		// deploy
		nest.NewDeployCommand(ct),

		// proxy commands
		proxy.NewRootCommand(ct),

		// cloud commands
		cloud.NewRootCommand(ct),
	)

	return cli
}

func main() {
	ct, err := container.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = newNestCommand(ct).Execute()
	if err != nil {
		ct.Logger().Print(loggy.NewEvent(loggy.ErrorLevel, err.Error(), nil))

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
