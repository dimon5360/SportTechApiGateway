package router

import (
	"log"
	"net/http"

	proto "github.com/dimon5360/SportTechProtos/gen/go/proto"
	"github.com/gin-gonic/gin"

	cors "github.com/rs/cors/wrapper/gin"
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

	router.engine.StaticFile("/favicon.ico", "../resources/favicon.ico")
	router.engine.StaticFile("/apple-touch-icon.png", "../resources/apple-touch-icon.png")
	router.engine.StaticFile("/favicon-32x32.png", "../resources/favicon-32x32.png")
	router.engine.Static("/resources", "../resources")

	router.engine.Use(cors.Default())

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
	r.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello user",
		})
	})
	r.engine.GET("/user/:id", r.GetUser)

	r.engine.POST("/auth", r.AuthenticateUser)
	r.engine.POST("/register", r.CreateUser)
}

var sslkey string = "../../private.key"
var sslcert string = "../../server.crt"

func (r *Router) Run() {
	r.engine.RunTLS(r.ip, sslcert, sslkey)
}
