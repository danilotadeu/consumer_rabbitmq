package server

import (
	"github.com/consumer_rabbitmq/events/user"
	"github.com/consumer_rabbitmq/rabbitmq"
)

type ServerI interface {
	Start()
}

type Server struct {
	URL string
}

func NewServer(URL string) ServerI {
	return &Server{
		URL: URL,
	}
}

func (s *Server) Start() {
	rabbitMQ := rabbitmq.Connect(s.URL)
	userEvent := user.NewEvent()

	queueUserCreated := rabbitMQ.DeclareAndBind("user_created_service_1", "user_created")
	rabbitMQ.Consume(queueUserCreated, userEvent.UserCreation)

	select {}
}
