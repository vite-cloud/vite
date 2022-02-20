package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/build"
	"github.com/vite-cloud/vite/cmd/vite"
	"github.com/vite-cloud/vite/cmd/vite/cloud"
	"github.com/vite-cloud/vite/cmd/vite/proxy"
	"github.com/vite-cloud/vite/container"
	"github.com/vite-cloud/vite/loggy"
	"os"
)

func newViteCommand(ct *container.Container) *cobra.Command {
	cli := &cobra.Command{
		Use:           "vite",
		Short:         "Service orchestrator",
		SilenceUsage:  true,
		SilenceErrors: true,
		Long:          "Vite is a powerful service orchestrator for a single server.",
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
		vite.NewVersionCommand(ct),

		// setup
		vite.NewSetupCommand(ct),

		// use
		vite.NewUseCommand(ct),

		// medic
		vite.NewMedicCommand(ct),

		// self-update
		vite.NewSelfUpdateCommand(ct),

		// deploy
		vite.NewDeployCommand(ct),

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

	err = newViteCommand(ct).Execute()
	if err != nil {
		ct.Logger().Print(loggy.NewEvent(loggy.ErrorLevel, err.Error(), nil))

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
