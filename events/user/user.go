package user

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/consumer_rabbitmq/model/user"
	"github.com/streadway/amqp"
)

const UserCreatedService = "user_created_service_1"
const UserCreated = "user_created"

type UserEvent interface {
	UserCreation(msg amqp.Delivery) error
}

type userEvent struct {
}

func NewEvent() UserEvent {
	return &userEvent{}
}

func (u *userEvent) UserCreation(msg amqp.Delivery) error {
	user := user.User{}
	err := json.Unmarshal(msg.Body, &user)
	if err != nil {
		return err
	}

	fmt.Println(user)

	err = msg.Ack(false)
	if err != nil {
		log.Printf("error to run ack" + err.Error())
	}
	return nil
}
