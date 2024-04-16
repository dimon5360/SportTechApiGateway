package userEndpoint

import (
	"app/main/internal/dto/constants"
	"app/main/internal/dto/models"
	"app/main/internal/repository/userRepository"
	"fmt"
	"log"
	"net/http"
	"strconv"

	proto "proto/go"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)

	GetProfile(c *gin.Context)
	PostProfile(c *gin.Context)
}

type userEndpointInstance struct {
	repo userRepository.Interface
}

func NewUserEndpoint(userRepository userRepository.Interface) (Interface, error) {
	e := &userEndpointInstance{
		repo: userRepository,
	}
	return e, nil
}

func (e *userEndpointInstance) Login(c *gin.Context) {

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

func (e *userEndpointInstance) Register(c *gin.Context) {

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

func (e *userEndpointInstance) RefreshToken(c *gin.Context) {

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

func (e *userEndpointInstance) GetProfile(c *gin.Context) {

	userId, err := c.Cookie("user_id")
	if err != nil {
		log.Println("cookie has no user id")
		c.Redirect(http.StatusFound, constants.ApiAuthLoginUrl)
		return
	}

	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.GetProfile(&proto.GetProfileRequest{
		Id: id,
	})

	if err != nil {
		log.Printf("could not get profile info: %v", err)
		c.Redirect(http.StatusFound, "/create-profile")
		return
	}

	if val, ok := res.(*proto.ProfileResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"username":  val.Username,
			"firstname": val.Firstname,
			"lastname":  val.Lastname,
		})
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}

func (e *userEndpointInstance) PostProfile(c *gin.Context) {

	type createProfileRequest struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	var req createProfileRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	userId, err := c.Cookie("user_id")
	if err != nil {
		log.Println("cookie has no user id")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.CreateProfile(&proto.CreateProfileRequest{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	})

	if err != nil {
		log.Printf("Creation profile failed: %v", err)
		c.Redirect(http.StatusFound, "/profile/"+userId)

		return
	}

	if val, ok := res.(*proto.ProfileResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"profile_id": val.Id,
		})

		c.Redirect(http.StatusFound, "/profile/"+userId)
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
