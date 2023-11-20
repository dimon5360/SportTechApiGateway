package router

import (
	"log"

	proto "github.com/dimon5360/SportTechProtos/gen/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Router struct {
	engine *gin.Engine

	grpc proto.AuthUsersServiceClient
	ip   string
}

func InitRouter(ip string) Router {

	router := Router{
		engine: gin.Default(),
		ip:     ip,
	}

	router.engine.LoadHTMLGlob("../static/templates/*")
	router.setupRouting()

	conn, err := grpc.Dial("localhost:40402",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	router.grpc = proto.NewAuthUsersServiceClient(conn)

	return router
}

func (r *Router) setupRouting() {
	r.engine.GET("/", Index)
	r.engine.GET("/user/:id", r.GetUser)
}

func (r *Router) Run() {
	r.engine.Run(r.ip)
}
