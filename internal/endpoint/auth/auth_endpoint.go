package endpoint

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	"github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type authEndpoint struct {
	user  repository.Interface
	redis repository.Interface
}

func New(repo ...repository.Interface) endpoint.Interface {
	if len(repo) != 2 {
		log.Fatalln()
		return nil
	}

	e := &authEndpoint{
		user:  repo[0],
		redis: repo[1],
	}

	if err := e.user.Init(); err != nil {
		log.Fatal(err)
	}

	if err := e.redis.Init(); err != nil {
		log.Fatal(err)
	}

	return e
}

func (e *authEndpoint) Get(c *gin.Context) {

	log.Println("unimplemented method")
	c.Status(http.StatusOK)
}

func (e *authEndpoint) Post(c *gin.Context) {

	log.Println("auth user")

	type authUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req authUserRequest
	err := c.Bind(&req)
	if err != nil {
		endpoint.ProcessingFailed(c, err, endpoint.InvalidRequestArgs, http.StatusBadRequest)
		return
	}

	log.Println(req)

	response, err := e.user.Verify(&proto.AuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		endpoint.ProcessingFailed(c, err, "Authentication failed", http.StatusUnauthorized)
		return
	}

	info, ok := response.(*proto.UserInfoResponse)
	if !ok {
		endpoint.ProcessingFailed(c, err, "invalid convert int to string", http.StatusBadRequest)
		return
	}
	log.Printf("auth user %d", info.Id)

	c.AddParam("user_id", strconv.FormatUint(info.Id, 10))
}
