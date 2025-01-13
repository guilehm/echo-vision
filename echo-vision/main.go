package main

import (
	"fmt"
	"log"

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
}
