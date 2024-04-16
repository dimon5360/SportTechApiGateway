package service

import (
	"app/main/internal/dto/constants"
	"app/main/internal/endpoint"
	"app/main/internal/endpoint/reportEndpoint"
	"app/main/internal/endpoint/userEndpoint"
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

	userEndpoint   userEndpoint.Interface
	reportEndpoint reportEndpoint.Interface
}

func New(
	userEndpoint userEndpoint.Interface,
	reportEndpoint reportEndpoint.Interface,
) Interface {

	return &router{
		userEndpoint:   userEndpoint,
		reportEndpoint: reportEndpoint,
	}
}

func (s *router) Init() error {

	s.engine = gin.Default()

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

	s.engine.GET(constants.ApiHomeUrl, endpoint.Index)
	s.engine.GET(constants.ApiPingUrl, endpoint.Ping)

	// token isn't required
	public := s.engine.Group(constants.ApiGroupV1)
	{
		public.POST(constants.ApiAuthLoginUrl, s.userEndpoint.Login)
		public.POST(constants.ApiAuthRegisternUrl, s.userEndpoint.Register)
	}

	// token required
	private := s.engine.Group(constants.ApiGroupV1)
	{
		private.PUT(constants.ApiRefreshTokenUrl, s.userEndpoint.RefreshToken)

		private.GET(constants.ApiProfileGetUrl, s.userEndpoint.GetProfile)
		private.POST(constants.ApiProfileCreateUrl, s.userEndpoint.PostProfile)

		private.GET(constants.ApiReportCreateUrl, s.reportEndpoint.Get)
		private.POST(constants.ApiReportGetUrl, s.reportEndpoint.Post)
	}

	s.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
}
