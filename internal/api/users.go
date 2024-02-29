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
)

func ProcessingFailed(c *gin.Context, err error, message string, status int) {

	log.Println(err.Error())

	c.JSON(status, gin.H{
		"error": InvalidRequestArgs,
	})
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

func GetUser(c *gin.Context) {

	service := grpc_service.AuthServiceInstance()

	ID := c.Params.ByName("id")

	var user UserInfo

	info := storage.Redis().Get(ID)

	err := json.Unmarshal(info, &user)
	if err != nil {
		ProcessingFailed(c, err, InvalidRequestArgs, http.StatusBadRequest)
		return
	}

	res, err := service.GetUser(&proto.GetUserRequest{
		Id: user.Id,
	})

	if err != nil {
		ProcessingFailed(c, err, "Getting user info failed", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": res.Email,
	})
}

func AuthenticateUser(c *gin.Context) {

	service := grpc_service.AuthServiceInstance()

	type authUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req authUserRequest
	err := c.Bind(&req)
	if err != nil {
		ProcessingFailed(c, err, InvalidRequestArgs, http.StatusBadRequest)
		return
	}

	fmt.Printf("%s: %v\n", "Auth user request", req)

	res, err := service.Auth(&proto.AuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ProcessingFailed(c, err, "Authentication failed", http.StatusUnauthorized)
		return
	}

	var id uint64 = res.Id
	token, err := generateToken(id)
	if err != nil {
		ProcessingFailed(c, err, "JWT token signing failed", http.StatusInternalServerError)
		return
	}

	expireIn := time.Hour * 24

	c.SetCookie("access_token", token, int(time.Now().Add(expireIn).Unix()), "", "", true, false)
	ck, err := c.Cookie("access_token")
	if err != nil {
		log.Println(ck)
	}

	bytes, err := json.Marshal(UserInfo{
		Id:          id,
		AccessToken: token,
	})

	if err != nil {
		ProcessingFailed(c, err, "JWT token handling failed", http.StatusInternalServerError)
		return
	}

	storage.Redis().Store(fmt.Sprintf("%d", id), bytes, expireIn)

	if err = VerifyProfile(id); err != nil {
		c.Redirect(http.StatusFound, "/profile/create")
	}

	userInfo := gin.H{
		"user_id":       id,
		"refresh_token": "default refresh token",
		"access_token":  token,
	}

	c.JSON(http.StatusOK, userInfo)
}

func CreateUser(c *gin.Context) {

	service := grpc_service.AuthServiceInstance()

	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req createUserRequest

	if err := c.Bind(&req); err != nil {
		ProcessingFailed(c, err, InvalidRequestArgs, http.StatusBadRequest)
		return
	}

	fmt.Printf("%s: %v\n", "Auth user request", req)

	_, err := service.Register(&proto.CreateUserRequst{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ProcessingFailed(c, err, "User already exists", http.StatusConflict)
		return
	}

	c.Status(http.StatusOK)
}
