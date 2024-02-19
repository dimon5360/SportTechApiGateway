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

	var user UserInfo

	info := storage.Redis().Get(ID)

	err := json.Unmarshal(info, &user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := service.GetProfile(&proto.GetProfileRequest{
		UserId: user.Id,
	})

	if err != nil {
		log.Printf("could not get user info: %v", err)
		c.String(http.StatusInternalServerError, "Getting user info failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  res.Username,
		"firstname": res.Firstname,
		"lastname":  res.Lastname,
	})
}

func CreateProfile(service *grpc_service.ProfileService, c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}
