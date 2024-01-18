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
	apiEnv     = "../config/api.env"
	redisEnv   = configPath + "redis.env"
)

func main() {

	utils.Env().Load(serviceEnv, apiEnv, redisEnv)

	fmt.Println("SportTech core service v." + utils.Env().Value("SERVICE_VERSION"))

	conn := storage.InitRedis()
	conn.TestConnect()

	server.InitServer(router.InitRouter(utils.Env().Value("SERVICE_HOST"))).Run()
}
