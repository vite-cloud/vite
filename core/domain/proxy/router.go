package proxy

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vite-cloud/go-zoup"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	logger     *Logger
	API        *gin.Engine
	config     *config.Config
}

func (r *Router) Proxy(w http.ResponseWriter, req *http.Request) {
	if req.Host == r.config.ControlPlane.Host {
		r.logger.LogR(req, zoup.DebugLevel, "proxy to control plane")
		r.API.ServeHTTP(w, req)
		return
	}

	targetIP, err := r.IPFor(req.Host)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
		r.logger.LogR(req, zoup.ErrorLevel, err.Error())
		return
	}

	if targetIP == "" {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Upstream did not respond."))
		r.logger.LogR(req, zoup.InfoLevel, "host not found")
		return
	}

	httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   targetIP,
	}).ServeHTTP(w, req)

	r.logger.LogR(req, zoup.InfoLevel, "served")
}

func (r *Router) IPFor(host string) (string, error) {
	if ip, ok := r.ips.Load(host); ok {
		return ip.(string), nil
	}

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

	r.ips.Store(host, ins.NetworkSettings.IPAddress)

	return ins.NetworkSettings.IPAddress, nil
}

func (r *Router) serviceFor(host string) (*config.Service, error) {
	for _, service := range r.config.Services {
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
func (r *Router) Accepts(host string) (bool, error) {
	_, err := r.serviceFor(host)
	if err != nil {
		return false, err
	}

	return true, nil
}
