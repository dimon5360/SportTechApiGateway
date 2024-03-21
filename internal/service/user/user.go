package service

import (
	"app/main/internal/endpoint"
	"app/main/internal/service"
	"log"
	"net/http"
	"os"

	"app/main/internal/middleware"

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

	jwt middleware.Token
}

func New(jwt middleware.Token, endpoints ...endpoint.Interface) service.Interface {
	if len(endpoints) != 4 {
		log.Fatal("invalid endpoints number")
		return nil
	}

	s := userService{
		user:    endpoints[0],
		profile: endpoints[1],
		report:  endpoints[2],
		auth:    endpoints[3],

		jwt: jwt,
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

	public := s.engine.Group("/api/v1")
	public.Use(s.jwt.Validate())
	{
		public.GET("/user/get", s.jwt.Validate(), s.user.Get)

		public.GET("/profile/get/:user_id", s.jwt.Validate(), s.profile.Get)
		public.POST("/profile/create", s.jwt.Validate(), s.profile.Post)

		public.POST("/report/get/:user_id", s.jwt.Validate(), s.report.Get)
		public.POST("/report/post", s.jwt.Validate(), s.report.Post)
	}

	private := s.engine.Group("/api/v1")
	{
		private.POST("/user/signup", s.user.Post)
		private.POST("/user/login", s.auth.Post, s.jwt.Generate())
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	return nil
}

func (s *userService) Run() error {

	host := os.Getenv(serviceHostKey)
	if len(host) == 0 {
		log.Fatal("host environment not found")
	}
	key := os.Getenv(sslKetPath)
	if len(key) == 0 {
		log.Fatal("ssl key environment not found")
	}
	cert := os.Getenv(sslCertPath)
	if len(cert) == 0 {
		log.Fatal("ssl cert environment not found")
	}

	err := s.engine.RunTLS(host, cert, key)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
