package service

import "github.com/gin-gonic/gin"

type Interface interface {
	Init() error
	Middleware(mw func(c *gin.Context))
	Run() error
}
