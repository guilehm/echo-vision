package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/guilehm/echo-vision/internal/infra/rabbitmq"
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

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})
	svc := rekognition.New(sess)

	fileName := "mage.jpeg"
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("could not open file: ", err)
	}

	reader := bufio.NewReader(f)
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln("could not read file: ", err)
	}

	maxLabels := int64(10)
	minConfidence := float64(70)
	detectLabelsResult, err := svc.DetectLabels(&rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: content,
		},
		MaxLabels:     &maxLabels,
		MinConfidence: &minConfidence,
	})
	if err != nil {
		log.Fatalln("could not detect labels: ", err)
	}
	fmt.Println("LABELS", detectLabelsResult)
}
