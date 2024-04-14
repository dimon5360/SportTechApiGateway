package authEndpoint

import (
	"app/main/internal/dto/constants"
	"app/main/internal/dto/models"
	"app/main/internal/repository/authRepository"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authEndpointInstance struct {
	repo authRepository.Interface
}

func NewAuthEndpoint(authRepository authRepository.Interface) (Interface, error) {
	e := &authEndpointInstance{
		repo: authRepository,
	}
	return e, nil
}

func (e *authEndpointInstance) Login(c *gin.Context) {

	var req models.RestLoginRequest
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

	c.SetCookie(constants.UserIdCookieKey, fmt.Sprintf("%d", user.Id), 0, "/", "", true, false)
	c.SetCookie(constants.AccessTokenCookieKey, user.AccessToken.GetValue(), user.AccessToken.GetAge(), "/", "", true, false)
	c.SetCookie(constants.RefreshTokenCookieKey, user.RefrestToken.GetValue(), user.RefrestToken.GetAge(), "/", "", true, false)

	if user.ProfileId == 0 {
		c.Redirect(http.StatusFound, constants.ApiProfileCreateUrl)
		return
	}
	c.Status(http.StatusOK)
}

func (e *authEndpointInstance) Register(c *gin.Context) {

	var req models.RestRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println("http register request:", req)

	err := e.repo.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (e *authEndpointInstance) RefreshToken(c *gin.Context) {

	id, err := c.Cookie(constants.UserIdCookieKey)
	if err != nil {
		c.Redirect(http.StatusFound, constants.ApiAuthLoginUrl)
		return
	}

	token, err := c.Cookie(constants.RefreshTokenCookieKey)
	if err != nil {
		c.Redirect(http.StatusFound, constants.ApiAuthLoginUrl)
		return
	}

	user_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	req := models.RestRefreshTokenRequest{
		Id:           user_id,
		RefreshToken: token,
	}

	log.Println("http refresh token request:", req)

	newToken, err := e.repo.RefreshToken(&req)
	if err != nil {
		c.Redirect(http.StatusFound, constants.ApiAuthLoginUrl)
		return
	}

	c.SetCookie(constants.UserIdCookieKey, fmt.Sprintf("%d", newToken.UserId), 0, "/", "", true, false)
	c.SetCookie(constants.AccessTokenCookieKey, newToken.AccessToken.GetValue(), newToken.AccessToken.GetAge(), "/", "", true, false)
	c.SetCookie(constants.RefreshTokenCookieKey, newToken.RefrestToken.GetValue(), newToken.RefrestToken.GetAge(), "/", "", true, false)

	c.Status(http.StatusOK)
}
