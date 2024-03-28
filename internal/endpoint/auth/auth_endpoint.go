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

type authEndpoint struct {
	repo repository.Interface
}

func New(authRepository repository.Interface) (endpoint.Interface, error) {
	e := &authEndpoint{
		repo: authRepository,
	}

	if err := e.repo.Init(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *authEndpoint) Get(c *gin.Context) {

	log.Println("unimplemented method")
	c.Status(http.StatusOK)
}

func (e *authEndpoint) Post(c *gin.Context) {

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

	response, err := e.repo.Verify(&proto.AuthUserRequest{
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
