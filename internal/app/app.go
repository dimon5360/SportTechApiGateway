package app

import (
	"app/main/pkg/utils"
	"fmt"
	"log"
)

type App struct {
	sp *ServiceProvider
}

func New() *App {
	return &App{}
}

func (a *App) Init() error {

	a.sp = NewServiceProvider()
	a.sp.Config()
	a.sp.Init()

	version, err := utils.Env().Value("SERVICE_VERSION")
	if err != nil {
		log.Fatal(err)
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
