package router

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
)

type Router struct {
	deployment *deployment.Deployment
	ips        sync.Map
	mu         sync.Mutex
}

func New(dep *deployment.Deployment) *Router {
	return &Router{deployment: dep, mu: sync.Mutex{}}
}

func (r *Router) IPFor(host string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	service, err := r.serviceFor(host)
	if err != nil {
		return "", err
	}

	id, err := r.deployment.Find("created_containers", service.Name)
	if err != nil {
		return "", err
	}

	ins, err := r.deployment.Docker.ContainerInspect(context.Background(), id.(string))
	if err != nil {
		return "", err
	}

	return ins.NetworkSettings.IPAddress, nil
}

func (r *Router) serviceFor(host string) (*config.Service, error) {
	for _, service := range r.deployment.Config.Services {
		for _, h := range service.Hosts {
			if ok, err := hostMatches(host, h); ok {
				return service, err
			} else if err != nil {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("no service found for host %s", host)
}

func hostMatches(host string, pattern string) (bool, error) {
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", ".*")

	re, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return false, err
	}

	return re.MatchString(host), nil
}
