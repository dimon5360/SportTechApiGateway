package api

import (
	"app/main/grpc_service"
	"app/main/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUser(service *grpc_service.AuthService, c *gin.Context) {

	ID := c.Params.ByName("id")

	var user UserInfo

	info := storage.Redis().Get(ID)

	err := json.Unmarshal(info, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestArgs,
		})
		return
	}

	res, err := service.GetUser(&proto.GetUserRequest{
		Id: user.Id,
	})

	if err != nil {
		log.Printf("could not get user info: %v", err)
		c.String(http.StatusInternalServerError, "Getting user info failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": res.Email,
	})
}

func AuthenticateUser(service *grpc_service.AuthService, c *gin.Context) {

	type authUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req authUserRequest
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestArgs,
		})
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

	expireIn := time.Hour * 24

	c.SetCookie("access_token", t, int(time.Now().Add(expireIn).Unix()), "", "", true, false)
	ck, err := c.Cookie("access_token")
	if err != nil {
		log.Println(ck)
	}

	id := uuid.New().String()

	bytes, err := json.Marshal(UserInfo{
		Id:          res.Id,
		AccessToken: t,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	storage.Redis().Store(id, bytes, expireIn)

	c.JSON(http.StatusOK, gin.H{
		"user_id":       id,
		"refresh_token": "default refresh token",
		"access_token":  t,
	})
}

func CreateUser(service *grpc_service.AuthService, c *gin.Context) {

	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestArgs,
		})
		return
	}

	_, err := service.Register(&proto.CreateUserRequst{
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

	// need to transfer frontend to provide redirecting
	// c.Redirect(http.StatusFound, "/api/v1/createprofile") // doesn't work yet
	c.Status(http.StatusOK)
}
