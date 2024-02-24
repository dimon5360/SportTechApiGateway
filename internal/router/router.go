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

	// static html files
	router.engine.LoadHTMLGlob("../static/html/**/*")

	router.engine.StaticFile("/favicon.ico", "../resources/favicon.ico")
	router.engine.StaticFile("/apple-touch-icon.png", "../resources/apple-touch-icon.png")
	router.engine.StaticFile("/favicon-32x32.png", "../resources/favicon-32x32.png")
	router.engine.Static("/resources", "../resources")

	// static js files
	router.engine.StaticFile("/home.jsx", "../static/js/Home/index.jsx")
	router.engine.StaticFile("/app.jsx", "../static/js/App/App.jsx")
	router.engine.StaticFile("/app/index.js", "../static/js/App/index.js")

	// router.engine.StaticFile("/main.js", "../static/js/Main/main.js") // Test webpack

	router.engine.Use(cors.Default())

	router.setupRouting()

	router.authService = grpc_service.NewAuthService(utils.Env().Value("AUTH_GRPC_HOST"))
	router.profileService = grpc_service.NewProfileService(utils.Env().Value("PROFILE_GRPC_HOST"))

	return router
}

func (r *Router) setupRouting() {

	r.engine.GET("/index", api.Index)
	r.engine.GET("/ping", api.Ping)
	r.engine.GET("/home", api.Home)

	// r.engine.GET("/", api.TestWebpack) // Test webpack

	route := r.engine.Group("/api/v1")
	{
		route.GET("/user/:id", r.GetUser)
		route.POST("/login", r.AuthenticateUser)
		route.POST("/register", r.CreateUser)
		route.POST("/profile", r.CreateProfile)
		route.GET("/profile/:user_id", r.GetProfile)
	}

	r.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
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
