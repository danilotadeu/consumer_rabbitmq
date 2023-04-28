package main

import (
	"os"

	"github.com/consumer_rabbitmq/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	urlRabbitMQ := os.Getenv("URL_RABBITMQ")
	server := server.NewServer(urlRabbitMQ)

	server.Start()
}
