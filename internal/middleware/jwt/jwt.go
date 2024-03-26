package jwt

import (
	"app/main/internal/endpoint"
	"app/main/internal/middleware"
	"app/main/internal/repository"
	"app/main/internal/repository/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	RedisAccessTokenFormat  = "user%v_access_token"
	RedisRefreshTokenFormat = "user%v_refresh_token"
	UserEmailIdKey          = "user%v_id"

	redirectLoginUrl   = "/api/v1/login"
	redirectRefreshUrl = "/api/v1/auth/refresh"

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

func generateTokens(id string) (string, string, error) {
	accessToken, err := generateToken(id, os.Getenv("ACCESS_TOKEN_CRED"), accessTokenExpireIn)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(id, os.Getenv("REFRESH_TOKEN_CRED"), refreshTokenExpireIn)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
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

	log.Println("Generate tokens ...")

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

		accessToken, refreshToken, err := generateTokens(id)
		if err != nil {
			endpoint.ProcessingFailed(c, err, "tokens generating failed", http.StatusInternalServerError)
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
		c.JSON(http.StatusOK, gin.H{
			"access-token":  accessToken,
			"refresh-token": refreshToken,
			"user_id":       id,
		})

		c.Next()
	}
}

func (t *token) Refresh() func(c *gin.Context) {

	log.Println("Refresh tokens ...")

	return func(c *gin.Context) {

		type Token struct {
			ExpiredAt int64  `json:"iat"`
			UserId    string `json:"sub"`
		}

		var token Token

		refreshToken := c.Request.Header.Get("Authorization")
		parts := strings.Split(refreshToken, ".") // split to get payload
		if len(parts) != 3 {
			redirectToLoginPage(c, "refresh token validating failed")
			return
		}

		resp, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(resp, &token)
		if err != nil {
			log.Println(err)
		}

		if time.Now().Unix() > int64(token.ExpiredAt) {
			redirectToLoginPage(c, "refresh token expired")
			return
		}

		err = t.validateRefreshToken(refreshToken, token.UserId)
		if err != nil {
			redirectToLoginPage(c, "refresh token validating failed")
			return
		}

		log.Println("refresh token is valid")

		c.AddParam("user_id", token.UserId)
		c.Next()
	}
}

func (t *token) validateAccessToken(accessToken string, userId string) error {

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
		if string(val) != accessToken {
			return fmt.Errorf("access token doesn't match")
		}
		return nil
	}
	return fmt.Errorf("access token not found")
}

func (t *token) validateRefreshToken(refreshToken string, userId string) error {

	req := model.RedisRequestModel{
		Key:    fmt.Sprintf(RedisRefreshTokenFormat, userId),
		Value:  "",
		Expire: 0,
	}

	info, err := t.repo.Get(&req)
	if err != nil {
		return err
	}

	if val, ok := info.([]byte); ok {
		if string(val) != refreshToken {
			return fmt.Errorf("refresh token doesn't match")
		}
		return nil
	}
	return fmt.Errorf("refreshToken token not found")
}

func (t *token) Validate() func(c *gin.Context) {

	log.Println("Validating token ...")

	return func(c *gin.Context) {

		type Token struct {
			ExpiredAt int64  `json:"iat"`
			UserId    string `json:"sub"`
		}

		var token Token

		accessToken := c.Request.Header.Get("Authorization")
		parts := strings.Split(accessToken, ".") // split to get payload
		if len(parts) != 3 {
			redirectToLoginPage(c, "access token validating failed")
			return
		}

		resp, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(resp, &token)
		if err != nil {
			log.Println(err)
		}

		if time.Now().Unix() > int64(token.ExpiredAt) {
			redirectToRefreshToken(c, "access token expired")
			return
		}

		err = t.validateAccessToken(accessToken, token.UserId)
		if err != nil {
			redirectToLoginPage(c, "access token validating failed")
			return
		}

		log.Println("access token is valid")
		c.Next()
	}
}

func redirectToLoginPage(c *gin.Context, message string) {
	log.Println(message)
	c.Redirect(http.StatusFound, redirectLoginUrl)
}

func redirectToRefreshToken(c *gin.Context, message string) {
	log.Println(message)
	c.Redirect(http.StatusFound, redirectRefreshUrl)
}
