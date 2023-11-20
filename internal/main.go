package main

import (
	router "app/main/api"
	"app/main/utils"
	server "app/main/web"
	"fmt"
)

func main() {

	env := utils.Init()
	env.Load("../config/app.env")

	fmt.Println("Core service v." + env.Value("VERSION_APP"))

	server.InitServer(router.InitRouter(env.Value("HOST"))).Run()
}
