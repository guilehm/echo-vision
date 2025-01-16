package rabbitmq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rotisserie/eris"
)

type ConsumerName string

func (cn ConsumerName) String() string {
	return string(cn)
}

type Topic string

func (t Topic) String() string {
	return string(t)
}

type RabbitMQConsumer struct {
	config Config

	ch          *amqp091.Channel
	consumeFunc func(Message)
}

// Publish implements Publisher.
func (r *RabbitMQConsumer) Publish(ctx context.Context, topic string, message Message) error {
	panic("unimplemented")
}

// Close implements Consumer.
func (r *RabbitMQConsumer) Close() error {
	panic("unimplemented")
}

// Subscribe implements Consumer.
func (r *RabbitMQConsumer) Subscribe(ctx context.Context, topic string, handler func(msg Message) error) error {
	panic("unimplemented")
}

func (r *RabbitMQConsumer) startConsumers(handler Handler) error {
	err := r.declare(handler.Topics())
	if err != nil {
		return err
	}

	msgs, err := r.ch.Consume(
		r.config.QueueName.String(),
		r.config.ConsumerName.String(),
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < r.config.ConcurrentConsumers; i++ {
		go func() {
			r.handler(msgs, handler)
		}()
	}
	// TODO: add a logger
	fmt.Printf("proessing messages in %s go routines\n", r.config.ConcurrentConsumers)
}

func (r *RabbitMQConsumer) declare(routingKeys []string) error {
	dlxName := r.config.QueueName + "_dlx"
	err := r.deadLetterDeclare(dlxName)
	if err != nil {
		return err
	}

	err = r.queueDeclare(dlxName)
	if err != nil {
		return err
	}

	err = r.queueBindDeclare(routingKeys)
	if err != nil {
		return err
	}

	// TODO: log qos settings
	err = r.ch.Qos(
		r.config.PrefetchCount, 0, false,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) queueDeclare(dlxName QueueName) error {
	err := r.ch.ExchangeDeclare(
		r.config.ExchangeName.String(),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = r.ch.QueueDeclare(
		r.config.QueueName.String(),
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type":           "quorum",
			"x-dead-letter-exchange": dlxName,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) deadLetterDeclare(dlxName QueueName) error {
	err := r.ch.ExchangeDeclare(
		dlxName.String(),
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = r.ch.QueueDeclare(
		dlxName.String(),
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type": "quorum",
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		dlxName.String(),
		"",
		dlxName.String(),
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQConsumer) queueBindDeclare(routingKeys []string) error {
	for _, routingKey := range routingKeys {
		err := r.ch.QueueBind(
			r.config.QueueName.String(),
			routingKey,
			r.config.ExchangeName.String(),
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (r *RabbitMQConsumer) handler(msgs <-chan amqp091.Delivery, handler Handler) {
// 	for d := range msgs {
// 		message := Message{
// 			Body: d.Body,
// 		}
//
// 		err := handler.Handle(message)
// 		if err != nil {
// }

func (r *RabbitMQConsumer) handleWithRecoverer(
	ctx context.Context,
	handler Handler,
	topic Topic,
	msg Message,
) (res HandlerResponse) {
	// TODO: log messages

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = eris.New(fmt.Sprintf("%v", r))
			}

			err = eris.Wrap(err, "panic")
			// TODO: log error

			res = DeadLetter
		}
	}()

	return handler.Handle(ctx, topic, msg)
}
