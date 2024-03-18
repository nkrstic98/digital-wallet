package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"transaction-service/app/services/transaction"
)

type HandlerImpl struct {
	userService transaction.Service
}

func NewHandler(userService transaction.Service) *HandlerImpl {
	return &HandlerImpl{
		userService: userService,
	}
}

func (handler *HandlerImpl) AddMoney(c *gin.Context) {
	var addMoneyRequest AddMoneyRequest
	if err := c.BindJSON(&addMoneyRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := handler.userService.AddMoney(transaction.AddMoneyInput{
		UserID: addMoneyRequest.UserID,
		Amount: addMoneyRequest.Amount,
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"data": AddMoneyResponse{Balance: balance},
	})
}

func (handler *HandlerImpl) TransferMoney(c *gin.Context) {
	var transferMoneyRequest TransferMoneyRequest
	if err := c.BindJSON(&transferMoneyRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.userService.TransferMoney(transaction.TransferMoneyInput{
		FromUserID: transferMoneyRequest.FromUserID,
		ToUserID:   transferMoneyRequest.ToUserID,
		Amount:     transferMoneyRequest.Amount,
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{})

}
