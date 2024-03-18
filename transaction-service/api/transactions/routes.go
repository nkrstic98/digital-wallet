package transactions

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (handler *HandlerImpl) RegisterTransactionRoutes(router *gin.Engine) {
	router.POST(path("/transactions/add-money"), handler.AddMoney)
	router.POST(path("/transactions/transfer-money"), handler.TransferMoney)
}

func path(endpoint string) string {
	return fmt.Sprintf("api/v1/%s", endpoint)
}
