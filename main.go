package main

import (
	"fmt"
	router "main/core/api"
	"main/core/utils"
	server "main/core/web"
)

func main() {

	env := utils.Init()
	env.Load("configs/app.env")

	fmt.Println("Application v." + env.Value("VERSION_APP"))

	server.InitServer(router.InitRouter(env.Value("HOST"))).Run()
}
