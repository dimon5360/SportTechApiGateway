package endpoint

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
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

	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := loginRequest{
		Email:    "admin@test.com",
		Password: "admin1234",
	}

	_, err := e.repo.Login(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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
