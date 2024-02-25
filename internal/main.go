package main

import (
	router "app/main/router"
	"app/main/storage"
	"app/main/utils"
	server "app/main/web"
	"fmt"
)

const (
	configPath = "/home/dmitry/Projects/SportTechService/SportTechDockerConfig/"
	serviceEnv = "../config/service.env"
	redisEnv   = configPath + "redis.env"
)

func main() {

	utils.Env().Load(serviceEnv, redisEnv)

	fmt.Println("SportTech user API service v." + utils.Env().Value("SERVICE_VERSION"))

	storage.InitRedis()

	server.InitServer(router.InitRouter(utils.Env().Value("SERVICE_HOST"))).Run()
}
