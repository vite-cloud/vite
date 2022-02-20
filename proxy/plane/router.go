package plane

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vite-cloud/vite/build"
	"github.com/vite-cloud/vite/docker"
	"github.com/vite-cloud/vite/service"
)

type ControlPlane struct {
	Config          *service.Locator
	ServicesConfig  *service.Config
	ManifestManager *service.ManifestManager
	Docker          *docker.Client
}

func (cp ControlPlane) From(router *gin.Engine) *gin.Engine {
	router.GET("/api/v1/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"software": "vite",
			"version":  build.Version,
			"build":    build.Commit,
		})
	})

	router.GET("/api/v1/server", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"commit":     cp.Config.Commit,
			"branch":     cp.Config.Branch,
			"repository": cp.Config.Repository,
			"remote":     cp.Config.RemoteURL(),
			"provider":   cp.Config.Provider,
			"config":     cp.ServicesConfig,
		})
	})

	router.GET("/api/v1/deploy", func(context *gin.Context) {
		deployment := service.NewDeployment(cp.ServicesConfig, cp.ManifestManager, cp.Docker)

		go func() {
			err := deployment.Start()
			if err != nil {
				deployment.Events <- service.Event{
					Service: nil,
					Value:   service.ErrDeploymentFailed,
				}
			}
		}()

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type")
		context.Header("Content-Type", "text/sseEvent-stream")
		context.Header("Cache-Control", "no-cache")
		context.Header("Connection", "keep-alive")

		var err error

		for e := range deployment.Events {
			if _, ok := e.Value.(error); ok {
				err = e.Value.(error)
				break
			}

			service := "global"

			if e.Service != nil {
				service = e.Service.Name
			}

			data, _ := json.Marshal(sseEvent{
				Kind:    "log",
				Service: service,
				Data:    fmt.Sprintf("%v", e.Value),
			})

			fmt.Fprintf(context.Writer, "data: %s\n\n", data)

			context.Writer.Flush()
		}

		context.Writer.Flush()

		if err != nil {
			fmt.Fprintf(context.Writer, "data: %s\n\n", sseEvent{
				Kind: "error",
				Data: err.Error(),
			})
			context.Writer.Flush()
		}

		if err = cp.ManifestManager.Save(deployment.Manifest); err != nil {
			fmt.Fprintf(context.Writer, "data: %s\n\n", sseEvent{
				Kind: "error",
				Data: fmt.Sprintf("%v", err),
			})
			context.Writer.Flush()
		} else {
			fmt.Fprintf(context.Writer, "data: %s\n\n", sseEvent{
				Kind: "manifest",
				Data: deployment.Manifest,
			})
			context.Writer.Flush()
		}
	})

	return router
}

type sseEvent struct {
	Kind string `json:"kind"`

	Service string `json:"service"`

	Data interface{} `json:"data"`
}

func (e sseEvent) String() string {
	data, _ := json.Marshal(e)

	return string(data)
}
