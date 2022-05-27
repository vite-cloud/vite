package deployment

import (
	"context"
	"github.com/vite-cloud/vite/core/domain/locator"
	"github.com/vite-cloud/vite/core/domain/resource"
	"strconv"
	"sync"
	"time"

	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/runtime"
)

const (
	StartEvent           = "StartEvent"
	StartLayerDeployment = "StartLayerDeployment"
	PullImage            = "PullImage"
	CreateContainer      = "CreateContainer"
	RunHook              = "RunHook"
	StartContainer       = "StartContainer"
	FinishDeployment     = "FinishDeployment"
	ConnectDependency    = "ConnectDependency"
	AcquireSubnet        = "AcquireSubnet"
	CreateNetwork        = "CreateNetwork"
)

const Store = datadir.Store("deployments")

func Deploy(events chan<- Event, locator *locator.Locator) {
	err := deploy(events, locator)
	if err != nil {
		events <- Event{
			ID:   ErrorEvent,
			Data: err,
		}
	} else {
		events <- Event{
			ID: FinishEvent,
		}
	}
}

func deploy(events chan<- Event, locator *locator.Locator) error {
	docker, err := runtime.NewClient()
	if err != nil {
		return err
	}

	depl := Deployment{
		id:      strconv.FormatInt(time.Now().UnixNano(), 10),
		Docker:  docker,
		Bus:     events,
		Locator: locator,
	}
	defer func(depl *Deployment) {
		err = resource.Save[*Deployment](Store, depl, func(d *Deployment) string {
			return d.ID()
		})
		if err != nil {
			events <- Event{
				ID:   ErrorEvent,
				Data: err,
			}
		}
	}(&depl)

	events <- Event{
		ID:   StartEvent,
		Data: depl.ID(),
	}
	conf, err := config.Get(locator)
	if err != nil {
		return err
	}

	layers, err := Layered(conf.Services)
	if err != nil {
		return err
	}

	errored := false

	for i, layer := range layers {
		var wg sync.WaitGroup

		events <- Event{
			ID: StartLayerDeployment,
			Data: struct {
				Current int
				Total   int
			}{i + 1, len(layers)},
		}

		for _, s := range layer {
			wg.Add(1)
			go func(s *config.Service) {
				defer wg.Done()

				err = depl.Deploy(context.Background(), events, s)
				if err != nil {
					events <- Event{
						ID:      ErrorEvent,
						Service: s,
						Data:    err,
					}
					errored = true
					return
				}

				events <- Event{
					ID:      FinishDeployment,
					Service: s,
				}
			}(s)
		}

		wg.Wait()

		if errored {
			break
		}
	}

	return nil
}
