package medic

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/runtime"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Diagnostic struct {
	Config   *config.Config
	Warnings []Warning
	Errors   []Error
}

type Warning struct {
	Title  string
	Advice string
}

type Error struct {
	Title string
	Error error
}

func Diagnose(config *config.Config) *Diagnostic {
	diagnostic := &Diagnostic{Config: config}

	for _, s := range config.Services {
		diagnostic.diagnoseService(s)
	}

	_, err := deployment.Layered(config.Services)
	diagnostic.ErrorIf(
		err != nil,
		"Circular dependency detected in your services dependencies",
		err,
	)

	diagnostic.ensureDnsRecordsPointToHost()

	return diagnostic
}

func (d *Diagnostic) diagnoseService(service *config.Service) {
	d.ErrorIf(service.Image == "", fmt.Sprintf("Service %s has no image", service.Name), nil)

	re := regexp.MustCompile(`^[a-zA-Z0-9-]+:([a-zA-Z0-9.]+)$`)
	ok := d.ErrorIf(
		!re.MatchString(service.Image),
		fmt.Sprintf("Service %s has an invalid image", service.Name),
		fmt.Errorf("image %s is not in the format <repository>:<tag>", service.Image),
	)
	d.ErrorIf(
		ok && re.FindStringSubmatch(service.Image)[1] == "latest",
		fmt.Sprintf("Service %s uses the `latest` tag, use a specific tag instead.", service.Name),
		nil,
	)

	if service.Registry != nil {
		d.diagnoseRegistry(*service.Registry)
	}

	if service.IsTopLevel && len(service.Hosts) == 0 {
		d.Warnings = append(d.Warnings, Warning{
			Title:  fmt.Sprintf("Service %s has no hosts", service.Name),
			Advice: "Services with no hosts are effectively disabled and will not be deployed.",
		})
	}

	for _, env := range service.Env {
		envKeyRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

		d.ErrorIf(
			!envKeyRegex.MatchString(env),
			fmt.Sprintf("Service %s has an invalid environment variable name: %s", service.Name, env),
			nil,
		)
	}

	for _, host := range service.Hosts {
		d.WarningIf(
			len(host) == 0,
			fmt.Sprintf("Service %s has an empty host", service.Name),
			"An empty host has no effect and will be ignored. Consider removing it.",
		)

		d.ErrorIf(
			host == d.Config.ControlPlane.Host,
			fmt.Sprintf("Service %s has the same host as the control plane", service.Name, host),
			nil,
		)
	}
}

func (d *Diagnostic) ErrorIf(condition bool, message string, err error) bool {
	if condition {
		d.Errors = append(d.Errors, Error{
			Title: message,
			Error: err,
		})
	}

	return !condition
}

func (d *Diagnostic) WarningIf(condition bool, message string, advice string) bool {
	if condition {
		d.Warnings = append(d.Warnings, Warning{
			Title:  message,
			Advice: advice,
		})
	}

	return !condition
}

func (d *Diagnostic) diagnoseRegistry(registry types.AuthConfig) {
	client, err := runtime.NewClient()
	ok := d.ErrorIf(
		err != nil,
		"Failed to create docker client",
		err,
	)
	if !ok {
		return
	}

	err = client.RegistryLogin(context.Background(), registry)
	d.ErrorIf(
		err != nil,
		fmt.Sprintf("Failed to login to registry %s", registry.ServerAddress),
		err,
	)
}

func (d *Diagnostic) ensureDnsRecordsPointToHost() {
	response, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		ok := d.WarningIf(
			strings.Contains(err.Error(), "no such host"),
			"It looks like you're not connected to internet, or AWS is down",
			"Please check your internet connection and try again.",
		)

		d.ErrorIf(
			!ok,
			"Failed to retrieve IP address",
			err,
		)
	}

	defer response.Body.Close()

	raw, err := ioutil.ReadAll(response.Body)
	d.ErrorIf(
		err != nil,
		fmt.Sprintf("Failed to read response body from http://checkip.amazonaws.com"),
		err,
	)

	ip := strings.TrimSpace(string(raw))

	var wg sync.WaitGroup

	for _, service := range d.Config.Services {
		for _, host := range service.Hosts {
			wg.Add(1)

			go func(host string) {
				defer wg.Done()

				hostIPs, err := net.LookupIP(host)
				d.ErrorIf(
					err != nil,
					fmt.Sprintf("Could not resolve DNS records for %s", host),
					err,
				)

				matched := false

				for _, hostIP := range hostIPs {
					if hostIP.String() == ip {
						matched = true
						break
					}
				}

				d.WarningIf(
					!matched,
					fmt.Sprintf("DNS records for %s do not point to your server (%s)", host, ip),
					fmt.Sprintf("Run `dig %s` to check if the DNS record is correct", host),
				)
			}(host)
		}
	}

	wg.Wait()
}
