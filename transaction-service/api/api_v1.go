package api

import (
	"github.com/gin-gonic/gin"
	"transaction-service/api/transactions"
)

type ApiV1 struct {
	transactionHandler transactions.Handler
}

func NewApiV1(transactionHandler transactions.Handler) *ApiV1 {
	return &ApiV1{
		transactionHandler: transactionHandler,
	}
}

func (api *ApiV1) RegisterRoutes(router *gin.Engine) {
	api.transactionHandler.RegisterTransactionRoutes(router)
}
