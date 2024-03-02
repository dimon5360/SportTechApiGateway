package app

import (
	"app/main/internal/storage"
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

	fmt.Println("SportTech user API service v." + utils.Env().Value("SERVICE_VERSION"))

	a.sp = NewServiceProvider()
	a.sp.Congig()
	a.sp.Init()

	storage.InitRedis()

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
