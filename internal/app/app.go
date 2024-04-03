package app

import (
	"app/main/internal/service"
	"log"
)

type AppInterface interface {
	Init() error
	Run() error
}

type app struct {
	provider ProviderInterface
	service  service.Interface
}

const serviceVersionKey = "SERVICE_VERSION"

func New() AppInterface {
	return &app{
		provider: NewServiceProvider(),
	}
}

func (a *app) Init() error {

	service, err := a.provider.Init()
	if err != nil {
		return err
	}

	a.service = service
	return nil
}

func (a *app) Run() error {

	log.Println("service running ...")
	return a.service.Run()
}
