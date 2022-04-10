package deployments

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"sort"
	"time"
)

func runListCommand(cli *cli.CLI) error {
	// todo: pagination
	deps, err := deployment.List()
	if err != nil {
		return err
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Time().After(deps[j].Time())
	})

	for _, dep := range deps {
		fmt.Fprintf(cli.Out(), "- %d | %s\n", dep.ID, fmtTime(dep.Time()))
	}

	fmt.Fprintf(cli.Out(), "\nFound %d deployments.\n", len(deps))

	return nil
}

func fmtTime(t time.Time) any {
	d := time.Since(t)

	// if less than a minute
	if d < time.Minute {
		secs := int(d.Seconds())
		return fmt.Sprintf("%d %s ago", secs, pluralize(secs, "second"))
	}

	// if less than an hour
	if d < time.Hour {
		mins := int(d.Minutes())
		return fmt.Sprintf("%d %s ago", int(d.Minutes()), pluralize(mins, "minute"))
	}

	// if less than a day
	if d < time.Hour*24 {
		hours := int(d.Hours())
		return fmt.Sprintf("%d %s ago", int(d.Hours()), pluralize(hours, "hour"))
	}

	days := int(d.Hours() / 24)
	return fmt.Sprintf("%d %s ago", days, pluralize(days, "day"))
}

func pluralize(n int, singular string) string {
	if n == 1 {
		return singular
	}
	return singular + "s"
}

func newListCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list deployments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCommand(cli)
		},
	}

	return cmd
}
