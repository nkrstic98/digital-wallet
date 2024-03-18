//go:build wireinject
// +build wireinject

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
	"transaction-service/api"
	transactions_handler "transaction-service/api/transactions"
	transaction_service "transaction-service/app/services/transaction"
	user_service "transaction-service/app/services/user"
	user_repo "transaction-service/app/services/user/repo/user"
	"transaction-service/config"
	"transaction-service/db"
)

func Build(cfg config.Config) (*gin.Engine, error) {
	panic(wire.Build(
		db.OpenConnection,
		buildUserService,
		buildAPI,
		initializeApp,
	))
}

func buildUserService(db *gorm.DB) *user_service.ServiceImpl {
	panic(wire.Build(
		wire.Bind(new(user_repo.Repository), new(*user_repo.RepositoryImpl)), user_repo.NewRepository,
		user_service.NewService,
	))
}

func buildAPI(db *gorm.DB) *api.ApiV1 {
	panic(wire.Build(
		wire.Bind(new(transaction_service.Service), new(*transaction_service.ServiceImpl)), transaction_service.NewService,
		wire.Bind(new(transactions_handler.Handler), new(*transactions_handler.HandlerImpl)), transactions_handler.NewHandler,
		api.NewApiV1,
	))
}

func initializeApp(api *api.ApiV1, userService *user_service.ServiceImpl, cfg config.Config) *gin.Engine {
	router := gin.Default()
	api.RegisterRoutes(router)

	userService.Init(cfg)

	return router
}
