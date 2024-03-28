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
	repo repository.Interface
}

func New(userRepository repository.Interface) (endpoint.Interface, error) {
	e := &userEndpoint{
		repo: userRepository,
	}

	if err := e.repo.Init(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *userEndpoint) Get(c *gin.Context) {

	ID := c.Params.ByName("user_id")
	userId, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Println("Invalid conversion from string to uint64")
		c.Status(http.StatusInternalServerError)
		return
	}

	res, err := e.repo.Get(&proto.GetUserRequest{
		Id: userId,
	})

	if err != nil {
		log.Printf("Getting user info failed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	isExist, err := e.repo.IsExist(&proto.GetProfileRequest{
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

	res, err := e.repo.Add(&proto.CreateUserRequst{
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
