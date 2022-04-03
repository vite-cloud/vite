package proxy

import (
	"context"
	"github.com/vite-cloud/vite/core/domain/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Shutdown(ctx context.Context) error
}

type Keeper struct {
	Name    string
	Server  Server
	Timeout time.Duration
}

type Finisher struct {
	Keepers []*Keeper
	logger  *Logger
}

func (f Finisher) Wait() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	for _, keeper := range f.Keepers {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), keeper.Timeout)
			defer cancel()

			if err := keeper.Server.Shutdown(ctx); err != nil {
				f.logger.Log(log.ErrorLevel, "graceful shutdown failed", log.Fields{
					"name": keeper.Name,
					"err":  err,
				})
			} else {
				f.logger.Log(log.InfoLevel, "graceful shutdown", log.Fields{
					"name": keeper.Name,
				})
			}
		}()
	}
}
