package profileEndpoint

import (
	"app/main/internal/dto/constants"
	"app/main/internal/repository/profileRepository"
	"log"
	"net/http"
	"strconv"

	proto "proto/go"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	Get(c *gin.Context)
	Post(c *gin.Context)
}

type profileEndpointInstance struct {
	repo profileRepository.Interface
}

func NewProfileEndpoint(profileRepository profileRepository.Interface) (Interface, error) {
	e := &profileEndpointInstance{
		repo: profileRepository,
	}
	return e, nil
}

func (e *profileEndpointInstance) Get(c *gin.Context) {

	userId, err := c.Cookie("user_id")
	if err != nil {
		log.Println("cookie has no user id")
		c.Redirect(http.StatusFound, constants.ApiAuthLoginUrl)
		return
	}

	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Read(&proto.GetProfileRequest{
		Id: id,
	})

	if err != nil {
		log.Printf("could not get profile info: %v", err)
		c.Redirect(http.StatusFound, "/create-profile")
		return
	}

	if val, ok := res.(*proto.ProfileResponse); ok {
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

func (e *profileEndpointInstance) Post(c *gin.Context) {

	type createProfileRequest struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	var req createProfileRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.InvalidRequestArgs,
		})
		return
	}

	userId, err := c.Cookie("user_id")
	if err != nil {
		log.Println("cookie has no user id")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Create(&proto.CreateProfileRequest{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	})

	if err != nil {
		log.Printf("Creation profile failed: %v", err)
		c.Redirect(http.StatusFound, "/profile/"+userId)

		return
	}

	if val, ok := res.(*proto.ProfileResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"profile_id": val.Id,
		})

		c.Redirect(http.StatusFound, "/profile/"+userId)
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
