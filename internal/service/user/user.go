package service

import (
	"app/main/internal/endpoint"
	"app/main/internal/repository"
	def "app/main/internal/service"
	"app/main/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ def.Interface = (*userService)(nil)

type userService struct {
	engine *gin.Engine

	user    endpoint.Interface
	profile endpoint.Interface
	report  endpoint.Interface
}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) Init() error {

	s.engine = gin.Default()

	s.engine.Static("/resources", "./resources")
	s.engine.Static("/static", "./static/html")

	s.engine.LoadHTMLGlob("./static/html/*.html")

	s.engine.StaticFile("/favicon.ico", "./resources/favicon.ico")
	s.engine.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	s.engine.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")

	// s.engine.StaticFile("/index.html", "./static/html/index.html")

	s.user = endpoint.Users(repository.Users())
	s.profile = endpoint.Profiles(repository.Profiles())
	s.report = endpoint.Reports(repository.Reports())

	s.engine.GET("/", endpoint.Index)
	// s.engine.GET("/index", endpoint.Index)
	s.engine.GET("/ping", endpoint.Ping)

	route := s.engine.Group("/api/v1")
	{
		route.GET("/user/get/:id", s.user.Get)
		// route.POST("/user/login", api.AuthenticateUser)
		route.POST("/user/signup", s.user.Post)

		route.GET("/profile/get/:user_id", s.profile.Get)
		route.POST("/profile/create", s.profile.Post)

		route.POST("/report/get/:user_id", s.report.Get)
		route.POST("/report/post", s.report.Post)
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	return nil
}

func (r *userService) Run() error {
	env := utils.Env()

	err := r.engine.RunTLS(
		utils.Env().Value("SERVICE_HOST"),
		env.Value("SSL_CERT_PATH"),
		env.Value("SSL_KEY_PATH"),
	)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
