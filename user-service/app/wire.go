//go:build wireinject
// +build wireinject

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"user-service/api"
	user_handlers "user-service/api/users"
	user_service "user-service/app/services/user"
	user_repo "user-service/app/services/user/repo/user"
	"user-service/config"
	"user-service/db"
)

func Build(cfg config.Config) (*gin.Engine, func(), error) {
	panic(wire.Build(
		db.OpenConnection,
		provideNATSConnection,
		buildAPI,
		initializeApp,
	))
}

func buildAPI(db *gorm.DB, natsConn *nats.Conn) *api.ApiV1 {
	panic(wire.Build(
		wire.Bind(new(user_repo.Repository), new(*user_repo.RepositoryImpl)), user_repo.NewRepository,
		wire.Bind(new(user_service.Service), new(*user_service.ServiceImpl)), user_service.NewService,
		wire.Bind(new(user_handlers.Handler), new(*user_handlers.HandlerImpl)), user_handlers.NewHandler,
		api.NewApiV1,
	))
}

func provideNATSConnection() (*nats.Conn, func(), error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, nil, err
	}

	return conn, func() {
		conn.Close()
	}, nil
}

func initializeApp(api *api.ApiV1) *gin.Engine {
	router := gin.Default()
	api.RegisterRoutes(router)
	return router
}
