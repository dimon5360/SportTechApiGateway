package main

import (
	"fmt"
	router "main/core/api"
	"main/core/utils"
	server "main/core/web"
)

/// TODO:
/// 1. write and compile protobuf
/// 2. include grpc
/// 3. transfer queries and process responses
/// 4. connect to kafka

func main() {

	env := utils.Init()
	env.Load("env/app.env")

	fmt.Println("Core service v." + env.Value("VERSION_APP"))

	server.InitServer(router.InitRouter(env.Value("HOST"))).Run()
}
