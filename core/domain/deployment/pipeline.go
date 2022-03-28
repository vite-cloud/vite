package deployment

import (
	"context"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/runtime"
	"strconv"
	"sync"
	"time"
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

func Deploy(events chan<- Event, services map[string]*config.Service) error {
	docker, err := runtime.NewClient()
	if err != nil {
		return err
	}

	depl := Deployment{
		ID:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Docker: docker,
		Bus:    events,
	}

	events <- Event{
		ID:   StartEvent,
		Data: depl.ID,
	}

	layers, err := Layered(services)
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
