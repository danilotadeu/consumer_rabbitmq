package user

import (
	"encoding/json"
	"fmt"

	"github.com/consumer_rabbitmq/model/user"
	rabbitMQClient "github.com/wagslane/go-rabbitmq"
)

const UserCreatedService = "user_created_service_1"
const UserCreated = "user_created"

type UserEvent interface {
	UserCreation(msg rabbitMQClient.Delivery) rabbitMQClient.Action
}

type userEvent struct{}

func NewEvent() UserEvent {
	return &userEvent{}
}

func (e *userEvent) UserCreation(msg rabbitMQClient.Delivery) rabbitMQClient.Action {
	user := user.User{}
	err := json.Unmarshal(msg.Body, &user)
	if err != nil {
		return rabbitMQClient.NackRequeue
	}

	fmt.Println(user)
	return rabbitMQClient.Ack
}
