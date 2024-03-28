package service

import (
	"app/main/internal/endpoint"
	"app/main/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	serviceHostKey = "SERVICE_HOST"
	sslCertPath    = "SSL_CERT_PATH"
	sslKetPath     = "SSL_KEY_PATH"
)

type router struct {
	engine *gin.Engine

	userEndp    endpoint.Interface
	profileEndp endpoint.Interface
	reportEndp  endpoint.Interface
	authEndp    endpoint.Interface
}

func New(
	authEndpoint endpoint.Interface,
	userEndpoint endpoint.Interface,
	profileEndpoint endpoint.Interface,
	reportEndpoint endpoint.Interface,
) service.Interface {
	return &router{
		authEndp:    authEndpoint,
		userEndp:    userEndpoint,
		profileEndp: profileEndpoint,
		reportEndp:  reportEndpoint,
	}
}

func (s *router) Init() error {

	s.engine = gin.Default()

	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())

	s.initStatic()
	s.initEndpoints()

	return nil
}

func (s *router) Run() error {

	host := os.Getenv(serviceHostKey)
	if len(host) == 0 {
		log.Fatal("host environment not found")
	}

	err := s.engine.Run(host)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *router) initStatic() {

	s.engine.Static("/resources", "./resources")
	s.engine.Static("/static", "./static/html")

	s.engine.LoadHTMLGlob("./static/html/*.html")

	s.engine.StaticFile("/favicon.ico", "./resources/favicon.ico")
	s.engine.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	s.engine.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")
}

func (s *router) initEndpoints() {

	s.engine.GET("/", endpoint.Index)
	s.engine.GET("/ping", endpoint.Ping)

	api := s.engine.Group("/api/v1")
	{
		api.GET("/user/:user_id", s.userEndp.Get)
		api.POST("/register", s.userEndp.Post)

		api.GET("/profile/:user_id", s.profileEndp.Get)
		api.POST("/profile/create", s.profileEndp.Post)

		api.POST("/report/:user_id", s.reportEndp.Get)
		api.POST("/report/create", s.reportEndp.Post)

		api.POST("/login", s.authEndp.Post)
		api.POST("/token/refresh", s.authEndp.Get)
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
}
