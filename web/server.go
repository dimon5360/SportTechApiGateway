package server

import router "main/core/api"

type Server struct {
	router router.Router
}

func InitServer(router router.Router) *Server {
	var server Server
	server.router = router
	return &server
}

func (server *Server) Run() {
	server.router.Run()
}
