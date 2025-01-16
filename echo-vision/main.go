package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/guilehm/echo-vision/internal/app/usecases"
	"github.com/guilehm/echo-vision/internal/infra/postgres"
	"github.com/guilehm/echo-vision/internal/infra/rabbitmq"
	"github.com/guilehm/echo-vision/internal/infra/web"
)

func main() {
	fmt.Println("hello world")

	conn, err := rabbitmq.NewRabbitMQAdapter()
	if err != nil {
		log.Fatalln("could not create rabbitmq adapter: ", err)
	}
	defer conn.Close()

	ch, err := conn.CreateChannel()
	if err != nil {
		log.Fatalln("could not create channel: ", err)
	}

	fmt.Println("channel", ch)

	e := postgres.NewEnt(os.Getenv("DATABASE_URL"))

	repo := postgres.NewRepository(e)
	userUseCase := usecases.NewManageUsersUseCase(repo)
	eventUseCase := usecases.NewManageEventsUseCase(repo)
	router := web.NewRouter(userUseCase, eventUseCase)

	err = http.ListenAndServe(":8080", router)
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
