package app

import (
	"app/main/internal/endpoint"
	authEndpoint "app/main/internal/endpoint/auth"
	profileEndpoint "app/main/internal/endpoint/profile"
	reportEndpoint "app/main/internal/endpoint/report"
	userEndpoint "app/main/internal/endpoint/user"
	"app/main/internal/middleware"
	"app/main/internal/repository"
	profileRepository "app/main/internal/repository/profile"
	reportRepository "app/main/internal/repository/report"
	tokenRepository "app/main/internal/repository/token"
	userRepository "app/main/internal/repository/user"
	"app/main/internal/service"
	userSerivce "app/main/internal/service/user"
	"app/main/pkg/utils"
	"log"
)

const (
	redisEnvKey = "REDIS_ENV"
	mongoEnvKey = "MONGO_ENV"
)

type ServiceProvider struct {
	service service.Interface

	eUser    endpoint.Interface
	eProfile endpoint.Interface
	eReport  endpoint.Interface
	eAuth    endpoint.Interface

	rUser    repository.Interface
	rProfile repository.Interface
	rReport  repository.Interface
	rToken   repository.Interface
}

const serviceEnv = "./config/service.env"

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
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

func (sp *ServiceProvider) Init() {
	sp.initUserService()
}

func (sp *ServiceProvider) initUserService() {
	sp.service = userSerivce.New(sp.getUserEndpoint(),
		sp.getProfileEndpoint(),
		sp.getReportEndpoint(),
		sp.getAuthEndpoint())

	sp.service.Middleware(middleware.TokenValidation(sp.rToken))
	if err := sp.service.Init(); err != nil {

		log.Fatal(err)
	}
}

func (sp *ServiceProvider) getUserEndpoint() endpoint.Interface {
	sp.eUser = userEndpoint.New(sp.userRepository())
	if sp.eUser == nil {
		log.Fatal("Failed endpoint creation")
	}
	return sp.eUser
}

func (sp *ServiceProvider) getProfileEndpoint() endpoint.Interface {
	if sp.eProfile == nil {
		sp.eProfile = profileEndpoint.New(sp.profileRepository())
	}
	return sp.eProfile
}

func (sp *ServiceProvider) getReportEndpoint() endpoint.Interface {
	if sp.eReport == nil {
		sp.eReport = reportEndpoint.New(sp.reportRepository())
	}
	return sp.eReport
}

func (sp *ServiceProvider) getAuthEndpoint() endpoint.Interface {
	if sp.eAuth == nil {
		sp.eAuth = authEndpoint.New(sp.userRepository(), sp.tokenRepository())
	}
	return sp.eAuth
}

func (sp *ServiceProvider) userRepository() repository.Interface {
	if sp.rUser == nil {
		sp.rUser = userRepository.New()
	}
	return sp.rUser
}

func (sp *ServiceProvider) profileRepository() repository.Interface {
	if sp.rProfile == nil {
		sp.rProfile = profileRepository.New()
	}
	return sp.rProfile
}

func (sp *ServiceProvider) reportRepository() repository.Interface {
	if sp.rReport == nil {
		sp.rReport = reportRepository.New()
	}
	return sp.rReport
}

func (sp *ServiceProvider) tokenRepository() repository.Interface {
	if sp.rToken == nil {
		sp.rToken = tokenRepository.New()
	}
	return sp.rToken
}
