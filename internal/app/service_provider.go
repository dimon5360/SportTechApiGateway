package app

import (
	"app/main/internal/endpoint"
	authEndpoint "app/main/internal/endpoint/auth"
	profileEndpoint "app/main/internal/endpoint/profile"
	reportEndpoint "app/main/internal/endpoint/report"
	authRepository "app/main/internal/repository/auth"
	profileRepository "app/main/internal/repository/profile"
	reportRepository "app/main/internal/repository/report"
	"app/main/internal/service"
	router "app/main/internal/service/router"
	"app/main/pkg/env"
	"fmt"
	"log"
	"os"
)

type ProviderInterface interface {
	Init() (service.Interface, error)
}

type provider struct {
}

func NewServiceProvider() ProviderInterface {
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
	log.Println("provider initialised")
	return p.initUserService()
}

func (p *provider) initUserService() (service.Interface, error) {

	service := router.New(
		p.getAuthEndpoint(),
		p.getProfileEndpoint(),
		p.getReportEndpoint(),
	)

	if err := service.Init(); err != nil {
		return nil, err
	}

	log.Println("router created")
	return service, nil
}

func (p *provider) getProfileEndpoint() endpoint.Profile {
	endp, err := profileEndpoint.New(profileRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}

func (p *provider) getReportEndpoint() endpoint.Report {
	endp, err := reportEndpoint.New(reportRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}

func (p *provider) getAuthEndpoint() endpoint.Auth {
	endp, err := authEndpoint.New(authRepository.New())
	if err != nil {
		log.Fatal(err)
	}
	return endp
}
