package plane

import (
	"github.com/gin-gonic/gin"
)

func Auth(id string, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		if username != id || password != token {
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}
