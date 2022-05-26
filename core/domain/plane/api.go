package plane

import (
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/domain/token"
	"github.com/vite-cloud/vite/core/static"
	"io"
)

const ApiV1Prefix = "/api/v1"

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(context *gin.Context) {
		username, password, ok := context.Request.BasicAuth()
		if !ok {
			context.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if username != "token" {
			context.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized (accepts: token)",
			})
			return
		}

		tokens, err := resource.List[token.Token](token.Store)
		if err != nil {
			context.AbortWithStatus(500)
			return
		}

		for _, t := range tokens {
			if t.Value == password {
				context.Next()
				return
			}
		}

		context.AbortWithStatusJSON(401, gin.H{
			"error": "unauthorized",
		})
	})

	router.GET(ApiV1Prefix+"/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":    static.Version,
			"commit":     static.Commit,
			"os":         static.OS,
			"go_version": static.GoVersion,
			"built_at":   static.BuiltAt,
		})
	})

	router.GET(ApiV1Prefix+"/config", func(c *gin.Context) {
		conf, err := config.Get()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.JSON(200, conf)
	})

	router.GET(ApiV1Prefix+"/deploy", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		conf, err := config.Get()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		events := make(chan deployment.Event)

		go func() {
			err = deployment.Deploy(events, conf)
			if err != nil {
				events <- deployment.Event{
					ID:   deployment.ErrorEvent,
					Data: err,
				}
			} else {
				events <- deployment.Event{
					ID: deployment.FinishEvent,
				}
			}
		}()

		c.Stream(func(w io.Writer) bool {
			for event := range events {
				if event.ID == deployment.FinishEvent {
					return false
				}

				sse.Encode(w, sse.Event{
					Event: event.ID,
					Data:  event.Data,
				})
			}

			return true
		})

	})

	return router
}
