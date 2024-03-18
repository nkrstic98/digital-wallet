package user

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
	"transaction-service/app/services/user/repo/user"
	"transaction-service/config"
)

const (
	natsSubject      = "user.balance"
	natsErrorMessage = "user.balance-error"

	userCreatedTopic = "user-created-topic"

	acknowledgeUserCreatedTopic          = "acknowledge-user-created-topic"
	acknowledgeUserCreatedSuccessMessage = "acknowledge-user-created-success"
	acknowledgeUserCreatedFailedMessage  = "acknowledge-user-created-failed"
)

type ServiceImpl struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) *ServiceImpl {
	return &ServiceImpl{
		userRepo: userRepo,
	}
}

func (service *ServiceImpl) Init(cfg config.Config) {
	go service.listenForNewUsers(cfg.Kafka)
	go service.listenForBalanceRequests(cfg.Nats)
}

func (service *ServiceImpl) listenForNewUsers(cfg config.KafkaConfig) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokerAddresses := []string{fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)}

	consumer, err := sarama.NewConsumer(brokerAddresses, config)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create consumer: %s", err))
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(userCreatedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create partition consumer: %s", err))
	}
	defer partitionConsumer.Close()

	producer, err := sarama.NewSyncProducer(brokerAddresses, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create producer: %s", err))
	}
	defer producer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			return
		case msg := <-partitionConsumer.Messages():
			var receivedUser UserCreatedEventRequest
			if stdErr := json.Unmarshal(msg.Value, &receivedUser); stdErr != nil {
				slog.Error(stdErr.Error())
				if _, _, err = producer.SendMessage(&sarama.ProducerMessage{
					Topic: acknowledgeUserCreatedTopic,
					Key:   sarama.StringEncoder(msg.Key),
					Value: sarama.StringEncoder(acknowledgeUserCreatedFailedMessage),
				}); err != nil {
					slog.Error(fmt.Sprintf("Error sending message: %v", err))
				}
			}

			if _, err = service.userRepo.Insert(receivedUser.UserId, receivedUser.CreatedAt); err != nil {
				if _, _, err = producer.SendMessage(&sarama.ProducerMessage{
					Topic: acknowledgeUserCreatedTopic,
					Key:   sarama.StringEncoder(msg.Key),
					Value: sarama.StringEncoder(acknowledgeUserCreatedFailedMessage),
				}); err != nil {
					slog.Error(fmt.Sprintf("Error sending message: %v", err))
				}
			}

			if _, _, err = producer.SendMessage(&sarama.ProducerMessage{
				Topic: acknowledgeUserCreatedTopic,
				Key:   sarama.StringEncoder(msg.Key),
				Value: sarama.StringEncoder(acknowledgeUserCreatedSuccessMessage),
			}); err != nil {
				slog.Error(fmt.Sprintf("Error sending message: %v", err))
				err = service.userRepo.Delete(receivedUser.UserId)
			}
		}
	}
}

func (service *ServiceImpl) listenForBalanceRequests(cfg config.NatsConfig) {
	natsConn, err := nats.Connect(cfg.URL)
	if err != nil {
		return
	}
	defer natsConn.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	if _, err = natsConn.Subscribe(natsSubject, func(msg *nats.Msg) {
		var data string

		userResponse, err := service.userRepo.Get(string(msg.Data))
		if err != nil {
			data = natsErrorMessage
		} else {
			data = fmt.Sprintf("%f", userResponse.Balance)
		}

		err = natsConn.Publish(msg.Reply, []byte(data))
		if err != nil {
			slog.Error(fmt.Sprintf("Error publishing response: %v", err))
			return
		}
	}); err != nil {
		return
	}

	select {
	case <-signals:
		return
	}
}
