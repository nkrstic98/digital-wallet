package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/app/services/user"
)

type HandlerImpl struct {
	userService user.Service
}

func NewHandler(userService user.Service) *HandlerImpl {
	return &HandlerImpl{
		userService: userService,
	}
}

func (handler *HandlerImpl) CreateUser(c *gin.Context) {
	var createRequest CreateRequest
	if err := c.BindJSON(&createRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := handler.userService.CreateUser(createRequest.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"data": CreateResponse{
			UserId:    newUser.UserId,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
		},
	})
}

func (handler *HandlerImpl) GetUserBalance(c *gin.Context) {
	userEmail := c.Param("email")

	balance, err := handler.userService.GetUserBalance(userEmail)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"data": GetUserBalanceResponse{
			Email:   userEmail,
			Balance: balance,
		},
	})
}
