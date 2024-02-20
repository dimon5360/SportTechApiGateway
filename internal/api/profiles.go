package api

import (
	"app/main/grpc_service"
	"app/main/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

func GetProfile(service *grpc_service.ProfileService, c *gin.Context) {

	ID := c.Params.ByName("user_id")

	info := storage.Redis().Get(ID)

	var user UserInfo

	err := json.Unmarshal(info, &user)
	if err != nil {
		c.String(http.StatusInternalServerError, "Getting profile failed")
		return
	}

	res, err := service.GetProfile(&proto.GetProfileRequest{
		UserId: user.Id,
	})

	if err != nil {
		log.Printf("could not get profile info: %v", err)
		c.String(http.StatusNotFound, "Getting profile failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  res.Username,
		"firstname": res.Firstname,
		"lastname":  res.Lastname,
	})
}

func CreateProfile(service *grpc_service.ProfileService, c *gin.Context) {

	type createUserRequest struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		UserId    string `json:"user_id"`
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	info := storage.Redis().Get(req.UserId)

	var user UserInfo
	if err := json.Unmarshal(info, &user); err != nil {
		fmt.Println(err.Error())
		return
	}

	profileReq := proto.CreateProfileRequst{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		UserId:    user.Id,
	}

	_, err := service.CreateProfile(&profileReq)

	if err != nil {
		log.Printf("Creation profile failed: %v", err)
		c.JSON(http.StatusConflict, gin.H{
			"message": "Profile already existsd",
		})
		return
	}

	c.Status(http.StatusOK)
}
