package plane

import (
	"github.com/gin-gonic/gin"
	"github.com/vite-cloud/vite/core/static"
)

const ApiV1Prefix = "/api/v1"

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET(ApiV1Prefix+"/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":  static.Version,
			"commit":   static.Commit,
			"os":       static.OS,
			"gov":      static.GoVersion,
			"built_at": static.BuiltAt,
		})
	})

	return router
}
