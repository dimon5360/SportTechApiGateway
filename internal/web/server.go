package server

import (
	router "app/main/router"
)

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
