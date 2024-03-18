package api

import (
	"github.com/gin-gonic/gin"
	"user-service/api/users"
)

type ApiV1 struct {
	userHandler users.Handler
}

func NewApiV1(userHandler users.Handler) *ApiV1 {
	return &ApiV1{
		userHandler: userHandler,
	}
}

func (api *ApiV1) RegisterRoutes(router *gin.Engine) {
	api.userHandler.RegisterUserRoutes(router)
}
