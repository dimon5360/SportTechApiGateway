package app

import (
	"fmt"
	"log"
	"os"
)

type App struct {
	sp *ServiceProvider
}

const serviceVersionKey = "SERVICE_VERSION"

func New() *App {
	return &App{
		sp: NewServiceProvider(),
	}
}

func (a *App) Init() error {
	a.sp.Init()

	version := os.Getenv(serviceVersionKey)
	if len(version) == 0 {
		log.Fatal("Service version not found")
	}
	fmt.Println("SportTech user API service v." + version)
	return nil
}

func (a *App) Run() error {

	fmt.Println("service running ...")

	err := a.sp.service.Run()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
