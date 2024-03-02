package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	profile "app/main/internal/endpoint/profile"
	report "app/main/internal/endpoint/report"
	user "app/main/internal/endpoint/user"

	repository "app/main/internal/repository"
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

func Users(repo repository.Interface) Interface {
	return user.NewUserEndpoint(repo)
}

func Profiles(repo repository.Interface) Interface {
	return profile.NewProfileEndpoint(repo)
}

func Reports(repo repository.Interface) Interface {
	return report.NewReportEndpoint(repo)
}
