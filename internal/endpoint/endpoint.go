package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	Get(c *gin.Context)
	Post(c *gin.Context)
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": "home page",
	})
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Hello from server",
	})
}
