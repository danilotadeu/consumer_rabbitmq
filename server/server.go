package server

import (
	"log"
	"os"
	"os/signal"

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
	rabbitmq := rabbitmq.NewRabbitMQ(s.URL).Connect()
	userEvent := user.NewEvent()

	userCreationConsumer := rabbitmq.Consume(userEvent.UserCreation, "", user.UserCreated, user.UserCreatedService)

	var forever chan bool
	log.Println(" [*] -> listening messages from rabbitMQ")

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	go func() {
		<-gracefulShutdown
		userCreationConsumer.Close()
		_ = rabbitmq.Connection.Close()
		os.Exit(1)
	}()

	<-forever
}
