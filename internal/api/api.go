package router

import (
	"context"
	"log"
	"net/http"
	"time"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/user.tmpl",
		gin.H{
			"id":         "1",
			"name":       "Dmitry",
			"created_at": time.Now(),
		})
}

func (r *Router) GetUser(c *gin.Context) {

	type getUserRequest struct {
		ID uint64 `uri:"id" binding:"required,min=1"`
	}

	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := r.grpc.GetUser(ctx, &proto.GetUserRequest{
		Id: req.ID,
	})

	if err != nil {
		log.Printf("could not get user info: %v", err)
		c.String(http.StatusInternalServerError, "Getting user info failed")
		return
	}

	c.String(http.StatusOK, res.String())
}

// test url to auth http://localhost:40401/auth?email=defaultuser@gmail.com&password=defaultuser123
func (r *Router) AuthenticateUser(c *gin.Context) {

	type authUserRequest struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	var req authUserRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log.Print(req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := r.grpc.AuthUser(ctx, &proto.AuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Authentication failed: %v", err)
		c.String(http.StatusUnauthorized, "Authentication failed")
		return
	}

	c.String(http.StatusOK, res.String())
}
