package resource

import (
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"sort"
	"time"
)

type resource interface {
	ID() string
	Time() time.Time
}

type Manager[T resource] struct {
	Store datadir.Store
}

func (m Manager[T]) ListCommand(cli *cli.CLI) error {
	// todo: pagination
	deps, err := List[T](m.Store)
	if err != nil {
		return err
	}

	sort.Slice(deps, func(i, j int) bool {
		return (*deps[i]).Time().After((*deps[j]).Time())
	})

	for _, dep := range deps {
		fmt.Fprintf(cli.Out(), "- %s | %s\n", (*dep).ID(), fmtTime((*dep).Time()))
	}

	fmt.Fprintf(cli.Out(), "\n%d found.\n", len(deps))

	return nil
}

func (m Manager[T]) ShowCommand(cli *cli.CLI, ID string) error {
	dep, err := Get[T](m.Store, ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(cli.Out(), "%+v", dep)
	
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
