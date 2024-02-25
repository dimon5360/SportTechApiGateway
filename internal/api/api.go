package api

import (
	"net/http"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

type GetProfileCallback func(req *proto.GetProfileRequest) (*proto.UserProfileResponse, error)
type CreateProfileCallback func(req *proto.CreateProfileRequst) (*proto.UserProfileResponse, error)

type GetUserCallback func(req *proto.GetUserRequest) (*proto.UserInfoResponse, error)
type AuthUserCallback func(req *proto.AuthUserRequest) (*proto.UserInfoResponse, error)
type CreateUserCallback func(req *proto.CreateUserRequst) (*proto.UserInfoResponse, error)

type UserInfo struct {
	Id          uint64
	AccessToken string
}

type ProfileInfo struct {
	ProfileId uint
}

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
	ContentTypeBabel  = "text/babel; charset=utf-8"
)

const InvalidRequestArgs = "Invalid HTTP-Request parameters"

func Index(c *gin.Context) {
	c.Redirect(http.StatusFound, "/home")
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "Home/index.html", gin.H{
		"content": "This is a home page...",
	})
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Hello from server",
	})
}
