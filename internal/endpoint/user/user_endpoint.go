package endpoint

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

type userEndpoint struct {
	user    repository.Interface
	profile repository.Interface
}

func New(repo ...repository.Interface) endpoint.Interface {
	if len(repo) != 2 {
		return nil
	}
	e := &userEndpoint{
		user:    repo[0],
		profile: repo[1],
	}

	if err := e.user.Init(); err != nil {
		log.Fatal(err)
	}
	return e
}

func (e *userEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")
	userId, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.user.Get(&proto.GetUserRequest{
		Id: userId,
	})

	if err != nil {
		log.Printf("Getting user info failed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	isExist, err := e.profile.IsExist(&proto.GetProfileRequest{
		UserId: userId,
	})

	if !isExist {
		c.Redirect(http.StatusFound, "/create-profile")
		return
	}

	if val, ok := res.(*proto.UserInfoResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"email": val.Email,
		})
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}

func (e *userEndpoint) Post(c *gin.Context) {

	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req createUserRequest

	if err := c.Bind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": endpoint.InvalidRequestArgs,
		})
		return
	}

	res, err := e.user.Add(&proto.CreateUserRequst{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Creation user failed: %v", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if val, ok := res.(*proto.UserInfoResponse); ok {
		c.JSON(http.StatusOK, gin.H{
			"user_id": val.Id,
		})
		return
	}

	c.Status(http.StatusInternalServerError)
	log.Println("invalid repository response")
}
