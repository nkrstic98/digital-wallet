package user

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"strconv"
	"time"
	"user-service/app/services/user/repo/user"
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
	db       *gorm.DB
	natsConn *nats.Conn
}

func NewService(userRepository user.Repository, natsConn *nats.Conn) *ServiceImpl {
	return &ServiceImpl{
		userRepo: userRepository,
		natsConn: natsConn,
	}
}

func (service *ServiceImpl) CreateUser(email string) (User, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokerAddresses := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create producer: %s", err))
		return User{}, err
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumer(brokerAddresses, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create consumer: %s", err))
		return User{}, err
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(acknowledgeUserCreatedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create partition consumer: %s", err))
		return User{}, err
	}
	defer partitionConsumer.Close()

	userResponse, err := service.userRepo.Insert(email)
	if err != nil {
		return User{}, err
	}

	marshalledUser, stdErr := json.Marshal(UserCreatedEventRequest{
		UserId:    userResponse.ID,
		Email:     userResponse.Email,
		CreatedAt: userResponse.CreatedAt,
	})
	if stdErr != nil {
		slog.Error(stdErr.Error())
		var deleteErr error
		deleteErr = service.userRepo.Delete(userResponse.ID)
		if deleteErr != nil {
			return User{}, deleteErr
		}
		return User{}, stdErr
	}

	userId := userResponse.ID.String()

	message := &sarama.ProducerMessage{
		Topic: userCreatedTopic,
		Key:   sarama.StringEncoder(userId),
		Value: sarama.StringEncoder(marshalledUser),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to send message: %s", err))
		var deleteErr error
		deleteErr = service.userRepo.Delete(userResponse.ID)
		if deleteErr != nil {
			return User{}, deleteErr
		}
		return User{}, err
	}

	select {
	case msg := <-partitionConsumer.Messages():
		if string(msg.Key) == userId && string(msg.Value) == acknowledgeUserCreatedSuccessMessage {
			slog.Info(fmt.Sprintf("User with email %s created successfully", email))
			return User{
				UserId:    userResponse.ID,
				Email:     userResponse.Email,
				CreatedAt: userResponse.CreatedAt,
			}, nil
		} else if string(msg.Key) == userId && string(msg.Value) == acknowledgeUserCreatedFailedMessage {
			slog.Error(fmt.Sprintf("Failed to create user with email in transaction service %s", email))
			if err = service.userRepo.Delete(userResponse.ID); err != nil {
				return User{}, err
			}
			return User{}, fmt.Errorf("failed to save user to the transactions database")
		}
	case <-partitionConsumer.Errors():
		slog.Error(fmt.Sprintf("Failed to consume message: %s", err))
		if err = service.userRepo.Delete(userResponse.ID); err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf("failed to save user to the transactions database")
	case <-time.After(3 * time.Second):
		slog.Error(fmt.Sprintf("Failed to consume message, request time-out: %s", err))
		if err = service.userRepo.Delete(userResponse.ID); err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf("failed to save user to the transactions database")
	}

	return User{}, fmt.Errorf("user creation failed")
}

func (service *ServiceImpl) GetUserBalance(email string) (float64, error) {
	userResponse, err := service.userRepo.GetByEmail(email)
	if err != nil {
		return 0, err
	}

	msg, err := service.natsConn.Request(natsSubject, []byte(userResponse.ID.String()), 5*time.Second)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending request to NATS: %v", err))
		return 0, err
	}

	data := string(msg.Data)

	if data == natsErrorMessage {
		return 0, fmt.Errorf("error getting user balance")
	}

	balance, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
