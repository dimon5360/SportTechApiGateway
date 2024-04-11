package endpoint

import (
	"app/main/internal/dto"
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	user_id_cookie_key      = "user_id"
	access_token_cookie_key = "access-token"
	refresh_token_key       = "refresh-token"
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
			"error": err.Error(),
		})
		return
	}

	c.SetCookie(user_id_cookie_key, fmt.Sprintf("%d", user.Id), 0, "/", "", true, false)
	c.SetCookie(access_token_cookie_key, user.AccessToken.GetValue(), user.AccessToken.GetAge(), "/", "", true, false)
	c.SetCookie(refresh_token_key, user.RefrestToken.GetValue(), user.RefrestToken.GetAge(), "/", "", true, false)

	if user.ProfileId == 0 {
		c.Redirect(http.StatusFound, endpoint.ApiProfileCreateUrl)
		return
	}
	c.Status(http.StatusOK)
}

func (e *authEndpoint) Register(c *gin.Context) {

	var req dto.RestRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println("http register request:", req)

	err := e.repo.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.Status(http.StatusOK)
}

func (e *authEndpoint) RefreshLogin(c *gin.Context) {

	id, err := c.Cookie(user_id_cookie_key)
	if err != nil {
		c.Redirect(http.StatusFound, endpoint.ApiAuthLoginUrl)
		return
	}

	token, err := c.Cookie(refresh_token_key)
	if err != nil {
		c.Redirect(http.StatusFound, endpoint.ApiAuthLoginUrl)
		return
	}

	user_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	req := dto.RestRefreshTokenRequest{
		Id:           user_id,
		RefreshToken: token,
	}

	log.Println("http refresh token request:", req)

	user, err := e.repo.Refresh(&req)
	if err != nil {
		c.Redirect(http.StatusFound, endpoint.ApiAuthLoginUrl)
		return
	}

	c.SetCookie(user_id_cookie_key, fmt.Sprintf("%d", user.Id), 0, "/", "", true, false)
	c.SetCookie(access_token_cookie_key, user.AccessToken.GetValue(), user.AccessToken.GetAge(), "/", "", true, false)
	c.SetCookie(refresh_token_key, user.RefrestToken.GetValue(), user.RefrestToken.GetAge(), "/", "", true, false)

	c.Status(http.StatusOK)
}
