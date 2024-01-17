package main

import (
	router "app/main/router"
	"app/main/utils"
	server "app/main/web"
	"fmt"
)

func main() {

	env := utils.Env()
	env.Load("../config/app.env")
	env.Load("../config/api.env")

	fmt.Println("Core service v." + env.Value("VERSION_APP"))

	server.InitServer(router.InitRouter(env.Value("HOST"))).Run()
}
