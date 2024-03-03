package middleware

import (
	"app/main/internal/repository"
	"app/main/internal/repository/model"
	tokenRepo "app/main/internal/repository/token"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	userAccessTokenKey = "user%s_access_token"
	redirectLoginUrl   = "/user/login"
)

type Token struct {
	token repository.Interface
}

func redirectToLoginPage(c *gin.Context, message string) {
	log.Println(message)
	c.Redirect(http.StatusFound, redirectLoginUrl)
}

func TokenValidation() func(c *gin.Context) {
	mw := Token{
		token: tokenRepo.NewTokenRepository(),
	}
	return mw.Verify
}

func (s *Token) Verify(c *gin.Context) {

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

	req := model.RedisRequestModel{
		Key:    fmt.Sprintf(userAccessTokenKey, userId),
		Value:  "",
		Expire: 0,
	}

	info, err := s.token.Get(&req)
	if err != nil {
		redirectToLoginPage(c, "not found token")
		return
	}

	if val, ok := info.(string); ok {
		if val != token {
			redirectToLoginPage(c, "expired token")
			return
		}
		c.Next()
		return
	}

	redirectToLoginPage(c, "key not found in redis")
	return
}