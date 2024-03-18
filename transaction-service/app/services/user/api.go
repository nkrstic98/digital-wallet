package user

import (
	"transaction-service/config"
)

type Service interface {
	Init(cfg config.Config)
	listenForNewUsers(cfg config.KafkaConfig)
	listenForBalanceRequests(cfg config.NatsConfig)
}
