package plane

import "github.com/gin-gonic/gin"

const ApiV1Prefix = "/api/v1"

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET(ApiV1Prefix+"/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": "1.0.0",
		})
	})

	return router
}
