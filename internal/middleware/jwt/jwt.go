package jwt

import (
	"app/main/internal/endpoint"
	"app/main/internal/middleware"
	"app/main/internal/repository"
	"app/main/internal/repository/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	RedisAccessTokenFormat  = "user%v_access_token"
	RedisRefreshTokenFormat = "user%v_refresh_token"
	UserEmailIdKey          = "user%v_id"

	redirectLoginUrl = "/api/v1/login"

	accessTokenExpireIn  = time.Minute * 30    // 30 minutes
	refreshTokenExpireIn = time.Hour * 24 * 30 // 30 days
)

type token struct {
	repo repository.Interface
}

func New(repo repository.Interface) middleware.Token {
	return &token{
		repo: repo,
	}
}

func generateToken(id string, key string, expiredIn time.Duration) (string, error) {

	payload := jwt.MapClaims{
		"sub": id,
		"iat": time.Now().Add(expiredIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return t, nil
}

func generateAccessToken(id string) (string, error) {
	return generateToken(id, os.Getenv("ACCESS_TOKEN_CRED"), accessTokenExpireIn)
}

func generateRefreshToken(id string) (string, error) {
	return generateToken(id, os.Getenv("REFRESH_TOKEN_CRED"), refreshTokenExpireIn)
}

func (t *token) keepToken(key string, value string, expiredIn time.Duration) error {

	obj := model.RedisRequestModel{
		Key:    key,
		Value:  value,
		Expire: expiredIn,
	}
	if _, err := t.repo.Add(&obj); err != nil {
		return err
	}
	return nil
}

func (t *token) keepTokens(id string, accessToken string, refreshToken string) error {
	if err := t.keepToken(fmt.Sprintf(RedisAccessTokenFormat, id), accessToken, accessTokenExpireIn); err != nil {
		return err
	}
	if err := t.keepToken(fmt.Sprintf(RedisRefreshTokenFormat, id), refreshToken, refreshTokenExpireIn); err != nil {
		return err
	}
	return nil
}

func (t *token) Generate() func(c *gin.Context) {
	return func(c *gin.Context) {

		if c.Writer.Status() != http.StatusOK {
			return
		}

		log.Println("generate tokens pair")

		id, isExist := c.Params.Get("user_id")
		if !isExist {
			endpoint.ProcessingFailed(c, nil, "failed parse from interface", http.StatusInternalServerError)
			return
		}

		refreshToken, err := generateRefreshToken(id)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "refresh token generating failed", http.StatusInternalServerError)
			return
		}
		accessToken, err := generateAccessToken(id)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "access token generating failed", http.StatusInternalServerError)
			return
		}

		c.SetCookie("access-token", accessToken, int(time.Now().Add(accessTokenExpireIn).Unix()),
			"/", "", true, false)

		err = t.keepTokens(id, accessToken, refreshToken)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "JWT jwt handling failed", http.StatusInternalServerError)
			return
		}

		c.SetCookie("user_id", id, 0, "/", "", true, false)
		c.Header("Authorization", refreshToken)
		//c.JSON(http.StatusOK, gin.H{
		//	"user_id": id,
		//})

		c.Redirect(http.StatusFound, "/user/"+id)

		c.Next()
	}
}

func (t *token) Refresh() func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("not implemented")
	}
}

func (t *token) validateAccessToken(refreshToken string, userId string) error {

	req := model.RedisRequestModel{
		Key:    fmt.Sprintf(RedisAccessTokenFormat, userId),
		Value:  "",
		Expire: 0,
	}

	info, err := t.repo.Get(&req)
	if err != nil {
		return err
	}

	if val, ok := info.([]byte); ok {
		if string(val) != refreshToken {
			return fmt.Errorf("access token doesn't match")
		}
		return nil
	}
	return fmt.Errorf("access token not found")
}

func (t *token) Validate() func(c *gin.Context) {

	return func(c *gin.Context) {
		var userId, accessToken string

		userId, isExist := c.Params.Get("user_id")

		if !isExist {
			c.Status(http.StatusBadRequest)
			return
		}

		if accessToken = c.Request.Header.Get("Authorization"); len(accessToken) == 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		err := t.validateAccessToken(accessToken, userId)
		if err != nil {
			redirectToLoginPage(c, "access token validating failed")
			return
		}

		c.AddParam("user_id", userId)

		log.Println("access token is valid")
		c.Next()
	}
}

func redirectToLoginPage(c *gin.Context, message string) {
	log.Println(message)
	c.Redirect(http.StatusFound, redirectLoginUrl)
}
