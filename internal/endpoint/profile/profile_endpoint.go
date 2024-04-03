package endpoint

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"log"
	"net/http"
	"strconv"

	proto "proto/go"

	"github.com/gin-gonic/gin"
)

type profileEndpoint struct {
	repo repository.ProfileInterface
}

func New(profileRepository repository.ProfileInterface) (endpoint.Profile, error) {
	e := &profileEndpoint{
		repo: profileRepository,
	}

	if err := e.repo.Init(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *profileEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")
	userId, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Read(&proto.GetProfileRequest{
		UserId: userId,
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

func (e *profileEndpoint) Post(c *gin.Context) {

	type createProfileRequest struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		UserId    string `json:"user_id"`
	}

	var req createProfileRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": endpoint.InvalidRequestArgs,
		})
		return
	}

	userId, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Create(&proto.CreateProfileRequest{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		UserId:    userId,
	})

	if err != nil {
		log.Printf("Creation profile failed: %v", err)
		c.Redirect(http.StatusFound, "/profile/"+req.UserId)

		return
	}

	if val, ok := res.(*proto.ProfileResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"profile_id": val.Id,
		})

		c.Redirect(http.StatusFound, "/profile/"+req.UserId)
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
