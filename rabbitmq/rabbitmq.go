package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQI interface {
	DeclareAndBind(queueName, exchange string) amqp.Queue
	Consume(queue amqp.Queue, handler func(amqp.Delivery) error)
}

type rabbitMQ struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func Connect(URL string) *rabbitMQ {
	conn, err := amqp.Dial(URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		panic(err)
	}

	return &rabbitMQ{
		Connection: conn,
		Channel:    ch,
	}
}

func (r *rabbitMQ) DeclareAndBind(queueName, exchange string) amqp.Queue {
	q, err := r.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
		panic(err)
	}

	err = r.Channel.QueueBind(
		q.Name, "", exchange, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
		panic(err)
	}

	return q
}

func (r *rabbitMQ) Consume(queue amqp.Queue, handler func(amqp.Delivery) error) {
	msgs, err := r.Channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for i := 0; i < 10; i++ {
		go consumeMessages(msgs, handler)
	}
}

func consumeMessages(msgs <-chan amqp.Delivery, handler func(amqp.Delivery) error) {
	for msg := range msgs {
		err := handler(msg)
		if err != nil {
			fmt.Println("error in handler, you can put this message in DLQ " + err.Error())
		}
	}
}
