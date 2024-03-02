package endpoint

import (
	constants "app/main/internal/endpoint/common"
	"app/main/internal/model"
	"app/main/internal/repository"
	"app/main/internal/storage"
	"encoding/json"
	"log"
	"net/http"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

type profileEndpoint struct {
	repo repository.Interface
}

func NewProfileEndpoint(repo repository.Interface) *profileEndpoint {
	return &profileEndpoint{
		repo: repo,
	}
}

func (e *profileEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")

	info := storage.Redis().Get(ID)

	var user model.UserInfo

	err := json.Unmarshal(info, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	res, err := e.repo.Get(&proto.GetProfileRequest{
		UserId: user.Id,
	})

	if err != nil {
		log.Printf("could not get profile info: %v", err)
		c.Redirect(http.StatusFound, "/profile/create")
		return
	}

	if val, ok := res.(*proto.UserProfileResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"username":  val.Username,
			"firstname": val.Firstname,
			"lastname":  val.Lastname,
		})
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}

func (e *profileEndpoint) Post(c *gin.Context) {

	type createUserRequest struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		UserId    string `json:"user_id"`
	}

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	info := storage.Redis().Get(req.UserId)

	var user model.UserInfo
	if err := json.Unmarshal(info, &user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	res, err := e.repo.Add(&proto.CreateProfileRequst{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		UserId:    user.Id,
	})

	if err != nil {
		log.Printf("Creation profile failed: %v", err)
		c.String(http.StatusConflict, "Profile already exists")
		return
	}

	if val, ok := res.(*proto.UserProfileResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"profile_id": val.Id,
		})

		c.Redirect(http.StatusFound, "/profile/get/"+req.UserId)
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
