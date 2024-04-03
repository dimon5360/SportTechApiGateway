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
	serviceHostKey = "GATEWAY_SERVICE_HOST"
	sslCertPath    = "SSL_CERT_PATH"
	sslKetPath     = "SSL_KEY_PATH"
)

type router struct {
	engine *gin.Engine

	authEndpoint    endpoint.Auth
	profileEndpoint endpoint.Profile
	reportEndpoint  endpoint.Report
}

func New(
	authEndpoint endpoint.Auth,
	profileEndpoint endpoint.Profile,
	reportEndpoint endpoint.Report,
) service.Interface {

	return &router{
		authEndpoint:    authEndpoint,
		profileEndpoint: profileEndpoint,
		reportEndpoint:  reportEndpoint,
	}
}

func (s *router) Init() error {

	s.engine = gin.Default()

	log.Println(s.engine)

	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())

	s.initStatic()
	s.initEndpoints()

	log.Println("gin initialized")

	return nil
}

func (s *router) Run() error {

	host := os.Getenv(serviceHostKey)
	if len(host) == 0 {
		log.Fatal("host environment not found")
	}

	log.Printf("listening host %s ...", host)
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
		api.POST("/login", s.authEndpoint.Login)
		api.POST("/token-refresh", s.authEndpoint.RefreshLogin)
		api.POST("/register", s.authEndpoint.Register)

		api.GET("/profile/:user_id", s.profileEndpoint.Get)
		api.POST("/profile/create", s.profileEndpoint.Post)

		api.POST("/report/:user_id", s.reportEndpoint.Get)
		api.POST("/report/create", s.reportEndpoint.Post)
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
}
