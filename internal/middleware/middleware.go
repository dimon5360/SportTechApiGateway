package middleware

import (
	"github.com/gin-gonic/gin"
)

type Token interface {
	Validate() func(c *gin.Context)
	Refresh() func(c *gin.Context)
	Generate() func(c *gin.Context)
}
