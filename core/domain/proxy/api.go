package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/locator"
	"github.com/vite-cloud/vite/core/domain/resource"
	"github.com/vite-cloud/vite/core/domain/token"
	"github.com/vite-cloud/vite/core/static"
	"io"
)

const ApiV1Prefix = "/api/v1"

func NewAPI() *gin.Engine {
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
		conf, err := config.GetUsingDefaultLocator()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.JSON(200, conf)
	})

	router.GET(ApiV1Prefix+"/deploy", func(c *gin.Context) {
		loc, err := locator.LoadFromStore()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		events := make(chan deployment.Event)

		go deployment.Deploy(events, loc)

		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if event, ok := <-events; ok {
				if event.Data == nil {
					c.SSEvent(event.ID, "")
				} else {
					c.SSEvent(event.ID, event.Data)
				}
				return true
			}
			return false
		})
	})

	return router
}
