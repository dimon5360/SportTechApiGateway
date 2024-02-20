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

	router.engine.LoadHTMLGlob("../static/templates/*")

	router.engine.StaticFile("/favicxon.ico", "../resources/favicon.ico")
	router.engine.StaticFile("/apple-touch-icon.png", "../resources/apple-touch-icon.png")
	router.engine.StaticFile("/favicon-32x32.png", "../resources/favicon-32x32.png")
	router.engine.Static("/resources", "../resources")

	router.engine.Use(cors.Default())

	router.setupRouting()

	router.authService = grpc_service.NewAuthService(utils.Env().Value("AUTH_GRPC_HOST"))
	router.profileService = grpc_service.NewProfileService(utils.Env().Value("PROFILE_GRPC_HOST"))

	return router
}

func (r *Router) setupRouting() {
	env := utils.Env()

	// auth users service
	r.engine.GET(env.Value("API_V1_INDEX"), api.Index)
	r.engine.GET(env.Value("API_V1_GET_USER_BY_ID"), r.GetUser)

	r.engine.POST(env.Value("API_V1_LOGIN"), r.AuthenticateUser)
	r.engine.POST(env.Value("API_V1_REGISTER"), r.CreateUser)

	// profiles service

	r.engine.POST(env.Value("API_V1_CREATE_PROFILE"), r.CreateProfile)
	r.engine.GET(env.Value("API_V1_GET_PROFILE"), r.GetProfile)

	// test api with grpc
	r.engine.GET(env.Value("API_V1_PING"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping message hello from server",
		})
	})
}

func (r *Router) Run() {
	env := utils.Env()

	r.engine.RunTLS(r.ip,
		env.Value("SSL_CERT_PATH"),
		env.Value("SSL_KEY_PATH"))
}

func (r *Router) GetUser(c *gin.Context) {
	api.GetUser(r.authService, c)
}

func (r *Router) AuthenticateUser(c *gin.Context) {
	api.AuthenticateUser(r.authService, c)
}

func (r *Router) CreateUser(c *gin.Context) {
	api.CreateUser(r.authService, c)
}

func (r *Router) CreateProfile(c *gin.Context) {
	api.CreateProfile(r.profileService, c)
}

func (r *Router) GetProfile(c *gin.Context) {
	api.GetProfile(r.profileService, c)
}
