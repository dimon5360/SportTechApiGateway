package app

import (
	"app/main/internal/endpoint"
	"app/main/internal/endpoint/reportEndpoint"
	"app/main/internal/endpoint/userEndpoint"
	"app/main/internal/repository"
	"app/main/internal/service"
	"app/main/pkg/env"
	"app/main/pkg/logger"
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

	if err := logger.Init(); err != nil {
		return nil, err
	}

	version := os.Getenv(serviceVersionKey)
	if len(version) == 0 {
		return nil, fmt.Errorf("service version not found")
	}

	log.Println("SportTech API gateway v." + version)
	log.Println("provider initialised")
	return p.initUserService()
}

func (p *provider) initUserService() (service.Interface, error) {

	service := service.New(
		p.getAuthEndpoint(),
		p.getReportEndpoint(),
	)

	if err := service.Init(); err != nil {
		return nil, err
	}

	log.Println("router created")
	return service, nil
}

func (p *provider) getAuthEndpoint() userEndpoint.Interface {
	endp, err := endpoint.NewAuthEndpoint(repository.NewUserRepository())
	if err != nil {
		log.Fatal(err.Error())
	}
	return endp
}

func (p *provider) getReportEndpoint() reportEndpoint.Interface {
	endp, err := endpoint.NewReportEndpoint(repository.NewReportRepository())
	if err != nil {
		log.Fatal(err.Error())
	}
	return endp
}
