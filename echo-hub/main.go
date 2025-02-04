package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
	"github.com/guilehm/echo-vision/echo-common/rabbitmq"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/usecases"
	bcrypthasher "github.com/guilehm/echo-vision/echo-hub/internal/infra/bcrypt_hasher"
	jwtadapter "github.com/guilehm/echo-vision/echo-hub/internal/infra/jwt"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres"
	rabbitmqadapter "github.com/guilehm/echo-vision/echo-hub/internal/infra/rabbitmq"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
)

var logger = logging.NewLogger()

func main() {
	fmt.Println("hello world")

	a, err := filestorage.NewS3Adapter(
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_REGION"),
	)
	if err != nil {
		log.Fatalln("could not create s3 adapter: ", err)
	}

	url, err := a.GeneratePreSignedURL("test.jpeg")
	if err != nil {
		log.Fatalln("could not generate pre-signed URL: ", err)
	}

	fmt.Println("URL", url)

	client, err := rabbitmq.NewRabbitMQClient(
		os.Getenv("RABBITMQ_URL"),
		logger,
		rabbitmq.ConfigConsumerName("echo-hub"),
		rabbitmq.ConfigWithExchangeName("events"),
		rabbitmq.ConfigWithQueueName("echo-hub"),
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
	go func() {
		err = consumer.Subscribe(context.Background(), rabbitmqadapter.NewRabbitMQAdapter())
		if err != nil {
			log.Fatalln("could not subscribe to queue: ", err)
		}
	}()

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

	// go func() {
	// 	for i := 0; i < 500000; i++ {
	// 		go func() {
	// 			o := i
	// 			topic := "event.image_analysis.status_updated.created"
	// 			if o%2 == 0 {
	// 				topic = "whatever"
	// 			}
	//
	// 			err = publisher.Publish(context.Background(), messaging.Message{
	// 				Topic:   topic,
	// 				Payload: []byte("OMG"),
	// 				Headers: map[string]string{},
	// 			})
	// 			if err != nil {
	// 				log.Fatalln("could not publish message OPA: ", err)
	// 			}
	// 		}()
	// 	}
	// }()

	// for i := 0; i < 100000; i++ {
	// 	time.Sleep(1 * time.Second)
	// 	topic := "event.image_analysis.status_updated.*"
	// 	if i%2 == 0 {
	// 		topic = "whatever"
	// 	}
	// 	err = publisher.Publish(context.Background(), rabbitmq.Message{
	// 		Topic:   topic,
	// 		Payload: []byte("OMG"),
	// 		Headers: map[string]string{},
	// 	})
	// 	if err != nil {
	// 		log.Fatalln("could not publish message OPA: ", err)
	// 	}
	// 	fmt.Println("published messages")
	// }

	jwtAdapter := jwtadapter.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		1*time.Hour,
		24*time.Hour,
	)
	passwordAdapter := bcrypthasher.NewBcryptAdapter()

	e := postgres.NewEnt(os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_SCHEMA"))
	repo := postgres.NewRepository(e)

	userUseCase := usecases.NewManageUsersUseCase(repo, jwtAdapter, passwordAdapter)
	eventUseCase := usecases.NewManageEventsUseCase(repo, publisher)

	router := web.NewRouter(userUseCase, eventUseCase, publisher)
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatalln("could not start server: ", err)
	}

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-2"),
	// })
	// if err != nil {
	// 	log.Fatalln("could not create aws session: ", err)
	// }
	// svc := rekognition.New(sess)
	//
	// fileName := "mage.jpeg"
	// f, err := os.Open(fileName)
	// if err != nil {
	// 	log.Fatalln("could not open file: ", err)
	// }
	//
	// reader := bufio.NewReader(f)
	// content, err := io.ReadAll(reader)
	// if err != nil {
	// 	log.Fatalln("could not read file: ", err)
	// }
	//
	// maxLabels := int64(10)
	// minConfidence := float64(70)
	// detectLabelsResult, err := svc.DetectLabels(&rekognition.DetectLabelsInput{
	// 	Image: &rekognition.Image{
	// 		Bytes: content,
	// 	},
	// 	MaxLabels:     &maxLabels,
	// 	MinConfidence: &minConfidence,
	// })
	// if err != nil {
	// 	log.Fatalln("could not detect labels: ", err)
	// }
	// fmt.Println("LABELS", detectLabelsResult)

	select {}
}
