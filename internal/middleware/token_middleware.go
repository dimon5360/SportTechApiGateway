package middleware

import (
	"app/main/internal/endpoint"
	repository "app/main/internal/repository"
	"app/main/internal/repository/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	UserAccessTokenKey = "user%v_access_token"
	UserEmailIdKey     = "user_%v__id"
)

const (
	redirectLoginUrl = "/api/v1/user/login"
)

type Token interface {
	Validate() func(c *gin.Context)
	Generate() func(c *gin.Context)
}

type token struct {
	repo repository.Interface
}

func NewJWT(repo repository.Interface) Token {
	return &token{
		repo: repo,
	}
}

func generateToken(id uint64) (string, error) {

	payload := jwt.MapClaims{
		"sub": id,
		"iat": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	var jwtSecretKey = []byte("very-secret-key")

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (t *token) Generate() func(c *gin.Context) {
	return func(c *gin.Context) {

		if c.Writer.Status() != http.StatusOK {
			return
		}

		log.Println("generate new token")

		id, isExist := c.Params.Get("user_id")
		if !isExist {
			endpoint.ProcessingFailed(c, nil, "failed parse from interface", http.StatusInternalServerError)
			return
		}

		res, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "JWT token signing failed", http.StatusInternalServerError)
			return
		}

		token, err := generateToken(res)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "JWT token signing failed", http.StatusInternalServerError)
			return
		}

		expireIn := time.Duration(time.Hour * 24)
		c.SetCookie("access_token", token, int(time.Now().Add(expireIn).Unix()), "", "", true, false)
		ck, err := c.Cookie("access_token")
		if err != nil {
			log.Println(ck)
			endpoint.ProcessingFailed(c, err, "JWT token handling failed", http.StatusInternalServerError)
			return
		}

		log.Println("token access: " + ck)

		tokenObj := model.RedisRequestModel{
			Key:    fmt.Sprintf(UserAccessTokenKey, id),
			Value:  token,
			Expire: expireIn,
		}
		c.SetCookie("user_id", id, 0, "", "", true, false)

		_, err = t.repo.Add(&tokenObj)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "JWT token handling failed", http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":       id,
			"refresh_token": "default refresh token",
			"access_token":  token,
		})

		c.Next()
	}
}

func (t *token) Validate() func(c *gin.Context) {

	return func(c *gin.Context) {
		log.Println("token validation")
		var userId, token string

		params := c.Request.URL.Query()

		if userId = params.Get("id"); len(userId) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user id missing",
			})
			return
		}

		if token = params.Get("access_token"); len(token) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user access token missing",
			})
			return
		}

		log.Println(userId, token)

		req := model.RedisRequestModel{
			Key:    fmt.Sprintf(UserAccessTokenKey, userId),
			Value:  "",
			Expire: 0,
		}

		info, err := t.repo.Get(&req)
		if err != nil {
			redirectToLoginPage(c, "not found token")
			return
		}

		if val, ok := info.([]byte); ok {
			if string(val) != token {
				redirectToLoginPage(c, "expired token")
				return
			}
			c.AddParam("user_id", userId)
			c.Next()
			return
		}

		redirectToLoginPage(c, "key not found in redis")
		return
	}
}

func redirectToLoginPage(c *gin.Context, message string) {
	log.Println(message)
	c.Redirect(http.StatusFound, redirectLoginUrl)
}
