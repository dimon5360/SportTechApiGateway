package endpoint

import (
	"app/main/internal/dto"
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authEndpoint struct {
	repo repository.AuthInterface
}

func New(authRepository repository.AuthInterface) (endpoint.Auth, error) {
	e := &authEndpoint{
		repo: authRepository,
	}

	if err := e.repo.Init(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *authEndpoint) Login(c *gin.Context) {

	var req dto.RestLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println("http login request:", req)

	user, err := e.repo.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.SetCookie("user_id", fmt.Sprintf("%d", user.Id), 0, "/", "", true, false)
	c.SetCookie("access-token", user.AccessToken.GetValue(), user.AccessToken.GetAge(), "/", "", true, false)
	c.SetCookie("refresh-token", user.RefrestToken.GetValue(), user.RefrestToken.GetAge(), "/", "", true, false)

	if user.ProfileId == 0 {
		c.Redirect(http.StatusFound, "/profile/create")
		return
	}
	c.Status(http.StatusOK)
}

func (e *authEndpoint) Register(c *gin.Context) {

	type registerRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var _ registerRequest
}

func (e *authEndpoint) RefreshLogin(c *gin.Context) {

	type refreshTokenRequest struct {
		Id           uint64 `json:"id"`
		RefreshToken string `json:"refresh-token"`
	}

	var _ refreshTokenRequest
}
