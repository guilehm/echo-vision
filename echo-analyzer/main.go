package main

import (
	"context"
	"fmt"
	"log"
	"os"

	rabbitmqadapter "github.com/guilehm/echo-vision/echo-analyzer/internal/infra/rabbitmq"
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/rabbitmq"
)

var logger = logging.NewLogger()

func main() {
	fmt.Println("hello world from echo-analyzer")

	client, err := rabbitmq.NewRabbitMQClient(
		os.Getenv("RABBITMQ_URL"),
		logger,
		rabbitmq.ConfigConsumerName("echo-analyzer"),
		rabbitmq.ConfigWithExchangeName("events"),
		rabbitmq.ConfigWithQueueName("echo-analyzer"),
		rabbitmq.ConfigConcurrentConsumers(5),
	)
	if err != nil {
		log.Fatalln("could not create rabbitmq consumer client: ", err)
	}

	consumer, err := client.CreateConsumer()
	if err != nil {
		log.Fatalln("could not create consumer: ", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe(context.Background(), rabbitmqadapter.NewRabbitMQAdapter())
	if err != nil {
		log.Fatalln("could not subscribe to queue: ", err)
	}
}
