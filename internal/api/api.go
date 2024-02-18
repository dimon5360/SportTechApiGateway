package api

import (
	"app/main/grpc_service"
	"log"
	"net/http"
	"time"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type UserInfo struct {
	UserId      uint
	AccessToken string
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/user.tmpl",
		gin.H{
			"id":         "1",
			"name":       "Dmitry",
			"created_at": time.Now(),
		})
}

func GetUser(service *grpc_service.AuthService, c *gin.Context) {

	ID, err := uuid.Parse(c.Params.ByName("id"))
	if err != nil {
		log.Printf("could not get request param 'id'")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	info, ok := service.Users[ID]
	if !ok {
		log.Printf("invalid request param 'id'")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	res, err := service.GetUser(&proto.GetUserRequest{
		Id: info.UserId,
	})

	if err != nil {
		log.Printf("could not get user info: %v", err)
		c.String(http.StatusInternalServerError, "Getting user info failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": res.Username,
		"email":    res.Email,
	})
}

func AuthenticateUser(service *grpc_service.AuthService, c *gin.Context) {

	c.SetSameSite(http.SameSiteStrictMode)

	type authUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req authUserRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := service.Auth(&proto.AuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication failed",
		})
		return
	}

	payload := jwt.MapClaims{
		"sub": res.Id,
		"iat": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	var jwtSecretKey = []byte("very-secret-key")

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "JWT token signing failed",
		})
		return
	}

	c.SetCookie("access_token", t, int(time.Now().Add(time.Hour*24).Unix()), "", "", true, false)
	ck, err := c.Cookie("access_token")
	if err != nil {
		log.Println(ck)
	}

	id := uuid.New()

	service.Users[id] = grpc_service.UserInfo{
		UserId:      res.Id,
		AccessToken: t,
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       id,
		"refresh_token": "default refresh token",
		"access_token":  t,
	})
}

func CreateUser(service *grpc_service.AuthService, c *gin.Context) {

	type createUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	_, err := service.Register(&proto.CreateUserRequst{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Creation user failed: %v", err)
		c.JSON(http.StatusConflict, gin.H{
			"message": "User already existsd",
		})
		return
	}

	c.Status(http.StatusOK)
}
