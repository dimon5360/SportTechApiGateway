package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/user.tmpl",
		gin.H{
			"id":         "1",
			"name":       "Dmitry",
			"created_at": time.Now(),
		})
}
