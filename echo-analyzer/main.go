package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/usecases"
	"github.com/guilehm/echo-vision/echo-analyzer/internal/infra/consumers"
	"github.com/guilehm/echo-vision/echo-analyzer/internal/infra/publishers"
	rabbitmqadapter "github.com/guilehm/echo-vision/echo-analyzer/internal/infra/rabbitmq"
	awsrekognition "github.com/guilehm/echo-vision/echo-analyzer/internal/infra/rekognition"
	"github.com/guilehm/echo-vision/echo-common/pkg/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/rabbitmq"
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

	publishConn, err := rabbitmq.NewRabbitMQClient(
		os.Getenv("RABBITMQ_URL"),
		logger,
		rabbitmq.ConfigWithExchangeName("events"),
		rabbitmq.ConfigWithConfirmMode(),
	)
	if err != nil {
		log.Fatalln("could not create rabbitmq publisher client: ", err)
	}

	publisher, err := publishConn.CreatePublisher()
	if err != nil {
		log.Fatalln("could not create publisher: ", err)
	}
	defer publisher.Close()

	if err := publisher.StartPublisher(context.Background()); err != nil {
		log.Fatalln("could not start publisher", err)
	}

	irs, err := awsrekognition.NewAWSRekognitionAdapter(os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET_NAME"))
	if err != nil {
		log.Fatalln("could not create AWS Rekognition adapter: ", err)
	}

	publisherGroup := publishers.NewPublisherGroup(publisher)
	imageAnalysisUseCase := usecases.NewImageAnalysisUseCase(publisherGroup, irs)
	consumerGroup := consumers.NewConsumerGroup(imageAnalysisUseCase, publisherGroup)

	adapter := rabbitmqadapter.NewRabbitMQAdapter(consumerGroup, publisherGroup)

	err = consumer.Subscribe(
		context.Background(),
		adapter,
	)
	if err != nil {
		log.Fatalln("could not subscribe to queue: ", err)
	}
}
