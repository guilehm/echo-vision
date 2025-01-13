package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeType string

func (et ExchangeType) String() string {
	return string(et)
}

const (
	ExchangeTypeDirect  ExchangeType = "direct"
	ExchangeTypeFanout  ExchangeType = "fanout"
	ExchangeTypeTopic   ExchangeType = "topic"
	ExchangeTypeHeaders ExchangeType = "headers"
)

type RabbitMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (*RabbitMQClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{connection: conn, channel: ch}, nil
}

func (rc *RabbitMQClient) Close() error {
	return rc.channel.Close()
}

func (rc *RabbitMQClient) CreateExchange(name ExchangeName, exchangeType ExchangeType, durable bool) error {
	return rc.channel.ExchangeDeclare(name.String(), exchangeType.String(), durable, false, false, false, nil)
}

func (rc *RabbitMQClient) CreateQueue(name QueueName, durable, autodelete bool) error {
	_, err := rc.channel.QueueDeclare(name.String(), durable, autodelete, false, false, nil)
	return err
}

func (rc *RabbitMQClient) CreateBinding(name QueueName, key string, exchange ExchangeName) error {
	return rc.channel.QueueBind(name.String(), key, exchange.String(), false, nil)
}
