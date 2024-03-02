package app

import (
	"app/main/internal/service"
	userService "app/main/internal/service/user"
	"app/main/pkg/utils"
	"log"
)

type ServiceProvider struct {
	service service.Interface
}

const serviceEnv = "./config/service.env"

func NewServiceProvider() *ServiceProvider {

	return &ServiceProvider{
		service: getUserService(),
	}
}

func (sp *ServiceProvider) Init() {
	if err := sp.service.Init(); err != nil {
		log.Fatal(err)
	}
}

func (sp *ServiceProvider) Congig() {
	env := utils.Env()

	env.Load(serviceEnv)
	env.Load(env.Value("REDIS_ENV"), env.Value("MONGO_ENV"))
}

func getUserService() service.Interface {
	return userService.NewUserService()
}
