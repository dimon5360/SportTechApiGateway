package app

import (
	"app/main/internal/endpoint"
	authEndpoint "app/main/internal/endpoint/auth"
	profileEndpoint "app/main/internal/endpoint/profile"
	reportEndpoint "app/main/internal/endpoint/report"
	userEndpoint "app/main/internal/endpoint/user"
	authRepository "app/main/internal/repository/auth"
	profileRepository "app/main/internal/repository/profile"
	reportRepository "app/main/internal/repository/report"
	userRepository "app/main/internal/repository/user"
	"app/main/internal/service"
	router "app/main/internal/service/router"
	"app/main/pkg/env"
	"fmt"
	"log"
	"os"
)

type IServiceProvider interface {
	Init() (service.Interface, error)
}

type provider struct {
	service service.Interface
}

func NewServiceProvider() IServiceProvider {
	return &provider{}
}

func (p *provider) Init() (service.Interface, error) {

	if err := env.Init(); err != nil {
		return nil, err
	}

	version := os.Getenv(serviceVersionKey)
	if len(version) == 0 {
		return nil, fmt.Errorf("service version not found")
	}

	fmt.Println("SportTech API gateway v." + version)
	return p.initUserService()
}

func (p *provider) Run() error {
	return p.service.Run()
}

func (p *provider) initUserService() (service.Interface, error) {
	return router.New(
		p.getAuthEndpoint(),
		p.getUserEndpoint(),
		p.getProfileEndpoint(),
		p.getReportEndpoint(),
	), nil
}

func (p *provider) getUserEndpoint() endpoint.Interface {
	endp, err := userEndpoint.New(userRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}

func (p *provider) getProfileEndpoint() endpoint.Interface {
	endp, err := profileEndpoint.New(profileRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}

func (p *provider) getReportEndpoint() endpoint.Interface {
	endp, err := reportEndpoint.New(reportRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}

func (p *provider) getAuthEndpoint() endpoint.Interface {
	endp, err := authEndpoint.New(authRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}
