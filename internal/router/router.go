package router

import (
	"app/main/api"
	"app/main/grpc_service"
	"app/main/utils"
	"net/http"

	"github.com/gin-gonic/gin"

	cors "github.com/rs/cors/wrapper/gin"
)

type Router struct {
	engine *gin.Engine

	authService    *grpc_service.AuthService
	profileService *grpc_service.ProfileService

	ip string
}

func InitRouter(ip string) Router {

	router := Router{
		engine: gin.Default(),
		ip:     ip,
	}

	router.engine.StaticFile("/favicon.ico", "../resources/favicon.ico")
	router.engine.StaticFile("/apple-touch-icon.png", "../resources/apple-touch-icon.png")
	router.engine.StaticFile("/favicon-32x32.png", "../resources/favicon-32x32.png")
	router.engine.Static("/resources", "../resources")

	router.engine.Use(cors.Default())

	router.setupRouting()

	grpc_service.NewAuthService(utils.Env().Value("AUTH_GRPC_HOST"))
	grpc_service.NewProfileService(utils.Env().Value("PROFILE_GRPC_HOST"))

	router.authService = grpc_service.AuthServiceInstance()
	router.profileService = grpc_service.ProfileServiceInstance()

	return router
}

func (r *Router) setupRouting() {

	r.engine.GET("/index", api.Index)
	r.engine.GET("/ping", api.Ping)
	r.engine.GET("/home", api.Home)

	route := r.engine.Group("/api/v1")
	{
		route.GET("/user/:id", api.GetUser)
		route.POST("/login", api.AuthenticateUser)
		route.POST("/register", api.CreateUser)
		route.POST("/create-profile", api.CreateProfile)
		route.GET("/profile/:user_id", api.GetProfile)
	}

	r.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
}

func (r *Router) Run() {
	// env := utils.Env()

	r.engine.Run(r.ip)
	// r.engine.RunTLS(r.ip,
	// 	env.Value("SSL_CERT_PATH"),
	// 	env.Value("SSL_KEY_PATH"))
}
