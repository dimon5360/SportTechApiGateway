package router

import (
	"app/main/proto"
	"context"
	"log"
	"net/http"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := r.grpc.GetUser(ctx, &proto.GetUserRequest{
		Id: "1",
	})

	if err != nil {
		log.Fatalf("could not get drink: %v", err)
		c.String(http.StatusInternalServerError, "Getting bar failed")
	}

	// TODO: serialize to JSON
	c.String(http.StatusOK, res.String())
}
