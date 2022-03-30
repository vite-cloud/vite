package router

import (
	"strings"

	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
)

type Router struct {
	deployment *deployment.Deployment
}

func New(dep *deployment.Deployment) *Router {
	return &Router{deployment: dep}
}

func (r *Router) ServiceFor(host string) *config.Service {
	for _, service := range r.dep.Config.Services {
		for _, h := range service.Hosts {
			if !strings.Contains(h, "*") {
				return service
			}

			panic("not implemented")
		}
	}
}
