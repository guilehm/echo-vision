package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guilehm/echo-vision/echo-common/rabbitmq"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/usecases"
	bcrypthasher "github.com/guilehm/echo-vision/echo-hub/internal/infra/bcrypt_hasher"
	jwtadapter "github.com/guilehm/echo-vision/echo-hub/internal/infra/jwt"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
)

func main() {
	fmt.Println("hello world")

	conn, err := rabbitmq.NewRabbitMQClient()
	if err != nil {
		log.Fatalln("could not create rabbitmq adapter: ", err)
	}
	defer conn.Close()

	ch, err := conn.CreateChannel()
	if err != nil {
		log.Fatalln("could not create channel: ", err)
	}

	fmt.Println("channel", ch)

	// TODO: use environment variables
	jwtAdapter := jwtadapter.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		1*time.Hour,
		24*time.Hour,
	)
	passwordAdapter := bcrypthasher.NewBcryptAdapter()

	e := postgres.NewEnt(os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_SCHEMA"))
	repo := postgres.NewRepository(e)

	userUseCase := usecases.NewManageUsersUseCase(repo, jwtAdapter, passwordAdapter)
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
