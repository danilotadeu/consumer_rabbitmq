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
	rabbitMQ := rabbitmq.Connect(s.URL)
	userEvent := user.NewEvent()

	var forever chan bool

	queueUserCreated := rabbitMQ.DeclareAndBind(user.UserCreated, user.UserCreatedService)
	rabbitMQ.Consume(queueUserCreated, userEvent.UserCreation)

	log.Println(" [*] - listening messages from rabbitMQ")

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	go func() {
		<-gracefulShutdown
		_ = rabbitMQ.Connection.Close()
		_ = rabbitMQ.Channel.Close()
		os.Exit(1)
	}()

	<-forever
}
