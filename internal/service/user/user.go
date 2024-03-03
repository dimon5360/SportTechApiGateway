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
	serviceHostKey = "SERVICE_HOST"
	sslCertPath    = "SSL_CERT_PATH"
	sslKetPath     = "SSL_KEY_PATH"
)

type userService struct {
	engine *gin.Engine

	user    endpoint.Interface
	profile endpoint.Interface
	report  endpoint.Interface
	auth    endpoint.Interface
}

func New(endpoints ...endpoint.Interface) service.Interface {
	if len(endpoints) != 4 {
		log.Fatal("invalid endpoints number")
		return nil
	}

	s := userService{
		user:    endpoints[0],
		profile: endpoints[1],
		report:  endpoints[2],
		auth:    endpoints[3],
	}

	s.engine = gin.Default()
	return &s
}

func (s *userService) Init() error {

	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())

	if err := s.initStatic(); err != nil {
		log.Fatal(err)
	}

	if err := s.initEndpoints(); err != nil {
		log.Fatal(err)
	}
	return nil
}
func (s *userService) Middleware(mw func(c *gin.Context)) {
	s.engine.Use(mw)
}

func (s *userService) initStatic() error {

	s.engine.Static("/resources", "./resources")
	s.engine.Static("/static", "./static/html")

	s.engine.LoadHTMLGlob("./static/html/*.html")

	s.engine.StaticFile("/favicon.ico", "./resources/favicon.ico")
	s.engine.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	s.engine.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")

	return nil
}

func (s *userService) initEndpoints() error {

	s.engine.GET("/", endpoint.Index)
	s.engine.GET("/ping", endpoint.Ping)

	route := s.engine.Group("/api/v1")
	{
		route.GET("/user/get/:id", s.user.Get)
		route.POST("/user/signup", s.user.Post)

		route.GET("/profile/get/:user_id", s.profile.Get)
		route.POST("/profile/create", s.profile.Post)

		route.POST("/report/get/:user_id", s.report.Get)
		route.POST("/report/post", s.report.Post)

		route.POST("/user/login", s.auth.Post)
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	return nil
}

func (s *userService) Run() error {
	env := utils.Env()

	host, err := env.Value(serviceHostKey)
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
