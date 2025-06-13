package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
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
	config *Config

	controller *ChannelController

	wg *sync.WaitGroup
}

// Close implements Consumer.
func (r *RabbitMQConsumer) Close() error {
	logger := r.config.Logger
	logger.Info("closing consumer")

	err := r.controller.ch.Cancel(r.config.ConsumerName.String(), true)
	if err != nil {
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), r.config.GracefulTimeout)
	defer cancel()

	r.waitForShutdown(shutdownCtx)

	err = r.controller.ch.Close()
	if err != nil {
		return err
	}
	logger.Info("consumer closed")
	return nil
}

// Subscribe implements Consumer.
func (r *RabbitMQConsumer) Subscribe(ctx context.Context, handler messaging.Handler) error {
	go r.controller.reconnectLoop()

	for {
		err := r.startConsumers(handler)
		if err != nil {
			return eris.Wrap(err, "calling startConsumers")
		}
		<-r.controller.reconnectCh
		r.controller.logger.Info("restarting consumers after reconnection")
	}
}

func (r *RabbitMQConsumer) startConsumers(handler messaging.Handler) error {
	logger := r.config.Logger
	err := r.declare(handler.Topics())
	if err != nil {
		return err
	}

	msgs, err := r.controller.ch.Consume(
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

	r.wg.Add(r.config.ConcurrentConsumers)
	for range r.config.ConcurrentConsumers {
		go func() {
			r.handler(msgs, handler)
			r.wg.Done()
		}()
	}
	logger.Info("started consumers", slog.Int("concurrency", r.config.ConcurrentConsumers))
	return nil
}

func (r *RabbitMQConsumer) declare(routingKeys []string) error {
	logger := r.config.Logger
	dlxName := r.config.QueueName + "_dlx"
	err := r.deadLetterDeclare(dlxName)
	if err != nil {
		return eris.Wrap(err, "calling deadLetterDeclare")
	}

	err = r.queueDeclare(dlxName)
	if err != nil {
		return eris.Wrap(err, "calling queueDeclare")
	}

	err = r.queueBindDeclare(routingKeys)
	if err != nil {
		return eris.Wrap(err, "calling queueBindDeclare")
	}

	err = r.controller.ch.Qos(
		r.config.PrefetchCount, 0, false,
	)
	if err != nil {
		return eris.Wrap(err, "calling ch.Qos")
	}
	logger.Info("declared queue",
		slog.String("queue", r.config.QueueName.String()),
		slog.String("exchange", r.config.ExchangeName.String()),
		slog.Int("prefetch_count", r.config.PrefetchCount),
	)
	return nil
}

func (r *RabbitMQConsumer) queueDeclare(dlxName QueueName) error {
	err := r.controller.ch.ExchangeDeclare(
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

	_, err = r.controller.ch.QueueDeclare(
		r.config.QueueName.String(),
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type":           "quorum",
			"x-dead-letter-exchange": dlxName.String(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) deadLetterDeclare(dlxName QueueName) error {
	err := r.controller.ch.ExchangeDeclare(
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

	_, err = r.controller.ch.QueueDeclare(
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

	err = r.controller.ch.QueueBind(
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
		err := r.controller.ch.QueueBind(
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

func (r *RabbitMQConsumer) handler(msgs <-chan amqp091.Delivery, handler messaging.Handler) {
	logger := r.config.Logger
	for msg := range msgs {
		message := messaging.Message{
			Payload: msg.Body,
			Topic:   msg.RoutingKey,
			Headers: nil,
		}

		ctx := context.Background()
		res := r.handleWithRecoverer(ctx, handler, message)

		switch res {
		case messaging.Success:
			err := msg.Ack(false)
			if err != nil {
				logger.Error("failed to ack message", slog.String("error", err.Error()))
			}
		case messaging.DeadLetter:
			err := msg.Nack(false, false)
			if err != nil {
				logger.Error("failed to nack message", slog.String("error", err.Error()))
			}
		default:
			err := msg.Nack(false, true)
			if err != nil {
				logger.Error("failed to nack message", slog.String("failed to discard messsage", err.Error()))
			}
		}
	}
}

func (r *RabbitMQConsumer) handleWithRecoverer(
	ctx context.Context,
	handler messaging.Handler,
	msg messaging.Message,
) (res messaging.HandlerResponse) {
	logger := r.config.Logger

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = eris.New(fmt.Sprintf("%v", r))
			}

			err = eris.Wrap(err, "panic")
			logger.Error("panic in handler", slog.String("error", err.Error()))
			res = messaging.DeadLetter
		}
	}()

	return handler.Handle(ctx, msg)
}

func (r *RabbitMQConsumer) waitForShutdown(ctx context.Context) {
	logger := r.config.Logger

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("context cancelled, shutting down consumer")
	case <-done:
		logger.Info("all message processing completed, shutting down consumer")
	}
}
