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

	userId, err := c.Cookie("user_id")
	if err != nil {
		log.Println("cookie has no user id")
		c.Status(http.StatusInternalServerError)
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

func (e *profileEndpoint) Post(c *gin.Context) {

	type createProfileRequest struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	var req createProfileRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": endpoint.InvalidRequestArgs,
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
