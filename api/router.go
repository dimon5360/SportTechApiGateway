package router

import (
	"app/main/proto"
	"log"

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

	router.engine.SetTrustedProxies([]string{"localhost"})
	router.engine.LoadHTMLGlob("static/templates/*")

	router.engine.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.engine.StaticFile("/site.webmanifest", "./resources/site.webmanifest")
	router.engine.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")

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
	r.engine.GET("/auth", r.AuthenticateUser)
}

func (r *Router) Run() {
	r.engine.Run(r.ip)
}
