package transactions

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterTransactionRoutes(router *gin.Engine)
	AddMoney(c *gin.Context)
	TransferMoney(c *gin.Context)
}
