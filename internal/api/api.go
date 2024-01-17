package api

import (
	"app/main/grpc_service"
	"log"
	"net/http"
	"time"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/user.tmpl",
		gin.H{
			"id":         "1",
			"name":       "Dmitry",
			"created_at": time.Now(),
		})
}

func GetUser(service *grpc_service.AuthService, c *gin.Context) {

	type getUserRequest struct {
		ID uint64 `uri:"id" binding:"required,min=1"`
	}

	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := service.GetUser(&proto.GetUserRequest{
		Id: req.ID,
	})

	if err != nil {
		log.Printf("could not get user info: %v", err)
		c.String(http.StatusInternalServerError, "Getting user info failed")
		return
	}

	c.String(http.StatusOK, res.String())
}

func AuthenticateUser(service *grpc_service.AuthService, c *gin.Context) {

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

	c.JSON(http.StatusOK, gin.H{
		"user_id": res.Id,
		"user":    res.Username,
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

	res, err := service.Register(&proto.CreateUserRequst{
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

	c.JSON(http.StatusOK, gin.H{
		"user_id": res.Id,
		"user":    res.Username,
	})
}
