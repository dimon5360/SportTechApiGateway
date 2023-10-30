package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine

	ip string
}

func InitRouter(ip string) Router {

	router := Router{
		engine: gin.Default(),
		ip:     ip,
	}

	router.engine.LoadHTMLGlob("static/templates/*")
	router.setupRouting()

	return router
}

func (r *Router) setupRouting() {
	r.engine.GET("/", Index)
}

func (r *Router) Run() {
	r.engine.Run(r.ip)
}
