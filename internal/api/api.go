package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Id          uint64
	AccessToken string
}

type ProfileInfo struct {
	ProfileId uint
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/user.tmpl",
		gin.H{
			"id":         "1",
			"name":       "Dmitry",
			"created_at": time.Now(),
		})
}
