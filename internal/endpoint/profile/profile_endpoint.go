package endpoint

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"log"
	"net/http"
	"strconv"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

type profileEndpoint struct {
	repo repository.Interface
}

func New(repo ...repository.Interface) endpoint.Interface {
	if len(repo) != 1 {
		return nil
	}
	e := &profileEndpoint{
		repo: repo[0],
	}

	if err := e.repo.Init(); err != nil {
		log.Fatal(err)
	}
	return e
}

func (e *profileEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")
	userId, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Get(&proto.GetProfileRequest{
		UserId: userId,
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

	res, err := e.repo.Add(&proto.CreateProfileRequst{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		UserId:    userId,
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
