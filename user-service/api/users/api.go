package users

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterUserRoutes(router *gin.Engine)
	CreateUser(c *gin.Context)
	GetUserBalance(c *gin.Context)
}
