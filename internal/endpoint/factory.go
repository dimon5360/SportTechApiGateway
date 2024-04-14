package endpoint

import (
	"app/main/internal/endpoint/authEndpoint"
	"app/main/internal/endpoint/profileEndpoint"
	"app/main/internal/endpoint/reportEndpoint"

	"app/main/internal/repository/authRepository"
	"app/main/internal/repository/profileRepository"
	"app/main/internal/repository/reportRepository"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewAuthEndpoint(repo authRepository.Interface) (authEndpoint.Interface, error) {
	if err := repo.Init(); err != nil {
		return nil, err
	}
	return authEndpoint.NewAuthEndpoint(repo)
}

func NewProfileEndpoint(repo profileRepository.Interface) (profileEndpoint.Interface, error) {
	if err := repo.Init(); err != nil {
		return nil, err
	}
	return profileEndpoint.NewProfileEndpoint(repo)
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
