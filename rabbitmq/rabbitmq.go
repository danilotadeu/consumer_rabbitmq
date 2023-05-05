package rabbitmq

import (
	"log"

	rabbitMQClient "github.com/wagslane/go-rabbitmq"
)

type RabbitMQI interface {
	Connect() *RabbitMQ
	Consume(handler func(d rabbitMQClient.Delivery) rabbitMQClient.Action, routingKey string, exchangeName string, queue string) *rabbitMQClient.Consumer
}

type RabbitMQ struct {
	Connection *rabbitMQClient.Conn
	URL        string
}

func NewRabbitMQ(url string) RabbitMQI {
	return &RabbitMQ{
		URL: url,
	}
}

func (r *RabbitMQ) Connect() *RabbitMQ {
	conn, err := rabbitMQClient.NewConn(r.URL, rabbitMQClient.WithConnectionOptionsLogging)
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQ{
		Connection: conn,
	}
}

func (r *RabbitMQ) Consume(
	handler func(d rabbitMQClient.Delivery) rabbitMQClient.Action,
	routingKey string,
	exchangeName string,
	queue string) *rabbitMQClient.Consumer {

	consumer, err := rabbitMQClient.NewConsumer(
		r.Connection,
		handler,
		queue,
		rabbitMQClient.WithConsumerOptionsRoutingKey(routingKey),
		rabbitMQClient.WithConsumerOptionsExchangeName(exchangeName),
		rabbitMQClient.WithConsumerOptionsExchangeKind("fanout"),
		rabbitMQClient.WithConsumerOptionsExchangeDurable,
		rabbitMQClient.WithConsumerOptionsQueueDurable,
		rabbitMQClient.WithConsumerOptionsExchangeDeclare,
		rabbitMQClient.WithConsumerOptionsConcurrency(10),
	)
	if err != nil {
		panic(err)
	}

	return consumer
}
