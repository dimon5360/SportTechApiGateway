package app

import (
	endpoint "app/main/internal/endpoint"
	profileEndpoint "app/main/internal/endpoint/profile"
	reportEndpoint "app/main/internal/endpoint/report"
	userEndpoint "app/main/internal/endpoint/user"
	"app/main/internal/repository"
	profileRepo "app/main/internal/repository/profile"
	reportRepo "app/main/internal/repository/report"
	userRepo "app/main/internal/repository/user"
	"app/main/internal/service"
	userService "app/main/internal/service/user"
	"app/main/pkg/utils"
	"log"
)

const (
	redisEnvKey = "REDIS_ENV"
	mongoEnvKey = "MONGO_ENV"
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

func (sp *ServiceProvider) Config() {
	env := utils.Env()

	env.Load(serviceEnv)
	redisEnv, err := env.Value(redisEnvKey)
	if err != nil {
		log.Fatal(err)
	}

	mongoEnv, err := env.Value(mongoEnvKey)
	if err != nil {
		log.Fatal(err)
	}

	env.Load(redisEnv, mongoEnv)
}

func getUserService() service.Interface {
	return userService.NewUserService(getUserEndpoint(), getProfileEndpoint(), getReportEndpoint())
}

func getUserEndpoint() endpoint.Interface {
	point, err := userEndpoint.NewUserEndpoint(getUserRepository())
	if err != nil {
		log.Fatal(err)
	}
	return point
}
func getProfileEndpoint() endpoint.Interface {
	point, err := profileEndpoint.NewProfileEndpoint(getProfileRepository())
	if err != nil {
		log.Fatal(err)
	}
	return point
}
func getReportEndpoint() endpoint.Interface {
	point, err := reportEndpoint.NewReportEndpoint(getReportRepository())
	if err != nil {
		log.Fatal(err)
	}
	return point
}

func getUserRepository() repository.Interface {
	return userRepo.NewUserRepository()
}

func getProfileRepository() repository.Interface {
	return profileRepo.NewProfileRepository()
}

func getReportRepository() repository.Interface {
	return reportRepo.NewReportRepository()
}
