package service

import (
	"app/main/internal/endpoint"
	"app/main/internal/service"
	"app/main/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ service.Interface = (*userService)(nil)

const (
	serviveHostKey = "SERVICE_HOST"
	sslCertPath    = "SSL_CERT_PATH"
	sslKetPath     = "SSL_KEY_PATH"
)

type userService struct {
	engine *gin.Engine

	user    endpoint.Interface
	profile endpoint.Interface
	report  endpoint.Interface
}

func NewUserService(
	user endpoint.Interface,
	profile endpoint.Interface,
	report endpoint.Interface,
) service.Interface {
	return &userService{
		user:    user,
		profile: profile,
		report:  report,
	}
}

func (s *userService) Init() error {

	s.engine = gin.Default()

	s.engine.Static("/resources", "./resources")
	s.engine.Static("/static", "./static/html")

	s.engine.LoadHTMLGlob("./static/html/*.html")

	s.engine.StaticFile("/favicon.ico", "./resources/favicon.ico")
	s.engine.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	s.engine.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")

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

func (s *userService) Run() error {
	env := utils.Env()

	host, err := env.Value(serviveHostKey)
	if err != nil {
		log.Fatal(err)
	}
	cert, err := env.Value(sslCertPath)
	if err != nil {
		log.Fatal(err)
	}
	key, err := env.Value(sslKetPath)
	if err != nil {
		log.Fatal(err)
	}

	err = s.engine.RunTLS(host, cert, key)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
