package gin

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGinExample(t *testing.T) {
	server := gin.Default()
	server.GET("/Hello", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
