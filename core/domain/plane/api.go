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
	//router.Use(func(context *gin.Context) {
	//	username, password, ok := context.Request.BasicAuth()
	//	if !ok {
	//		context.AbortWithStatusJSON(401, gin.H{
	//			"error": "unauthorized",
	//		})
	//		return
	//	}
	//
	//	if username != "token" {
	//		context.AbortWithStatusJSON(401, gin.H{
	//			"error": "unauthorized (accepts: token)",
	//		})
	//		return
	//	}
	//
	//})

	router.GET(ApiV1Prefix+"/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":    static.Version,
			"commit":     static.Commit,
			"os":         static.OS,
			"go_version": static.GoVersion,
			"built_at":   static.BuiltAt,
		})
	})

	return router
}
