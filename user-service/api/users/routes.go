package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (handler *HandlerImpl) RegisterUserRoutes(router *gin.Engine) {
	router.POST(path("/users"), handler.CreateUser)
	router.GET(path("/users/balance/:email"), handler.GetUserBalance)
}

func path(endpoint string) string {
	return fmt.Sprintf("api/v1/%s", endpoint)
}
