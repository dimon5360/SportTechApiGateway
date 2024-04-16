package endpoint

import (
	"app/main/internal/endpoint/reportEndpoint"
	"app/main/internal/endpoint/userEndpoint"

	"app/main/internal/repository/reportRepository"
	"app/main/internal/repository/userRepository"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewAuthEndpoint(repo userRepository.Interface) (userEndpoint.Interface, error) {
	if err := repo.Init(); err != nil {
		return nil, err
	}
	return userEndpoint.NewUserEndpoint(repo)
}

func NewReportEndpoint(repo reportRepository.Interface) (reportEndpoint.Interface, error) {
	if err := repo.Init(); err != nil {
		return nil, err
	}
	return reportEndpoint.NewReportEndpoint(repo)
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
