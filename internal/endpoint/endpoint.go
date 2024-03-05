package endpoint

import (
	"log"
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

func ProcessingFailed(c *gin.Context, err error, message string, status int) {

	if err != nil {
		log.Println(err.Error())
	}

	c.JSON(status, gin.H{
		"error": message,
	})
}
