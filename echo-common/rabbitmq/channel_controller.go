package rabbitmq

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// ChannelController handles pausing and resuming of a RabbitMQ channel.
type ChannelController struct {
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	logger *slog.Logger
	config *Config

	mu          sync.Mutex
	paused      atomic.Bool
	pauseCh     chan struct{}
	reconnectCh chan struct{}
}

// NewChannelController initializes a new ChannelController.
func NewChannelController(
	logger *slog.Logger,
	conn *amqp091.Connection,
	ch *amqp091.Channel,
	config *Config,
) *ChannelController {
	pauseCh := make(chan struct{}, 1)
	controller := &ChannelController{
		conn:   conn,
		ch:     ch,
		logger: logger,
		config: config,

		mu:          sync.Mutex{},
		paused:      atomic.Bool{},
		pauseCh:     pauseCh,
		reconnectCh: make(chan struct{}),
	}
	return controller
}

func (cc *ChannelController) Pause() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if !cc.paused.Load() {
		cc.paused.Store(true)
		cc.pauseCh <- struct{}{}
		cc.logger.Warn("channel paused")
	}
}

func (cc *ChannelController) Resume() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if cc.paused.Load() {
		cc.paused.Store(false)
		select {
		case <-cc.pauseCh:
		default:
		}
		cc.logger.Info("channel resumed")
	}
}

// WaitForResume blocks until the channel is resumed.
func (cc *ChannelController) WaitForResume() {
	for cc.paused.Load() {
		cc.logger.Info("waiting for channel to resume")
		time.Sleep(5 * time.Second)
	}
}

func (cc *ChannelController) reconnectLoop() {
	logger := cc.logger
	logger.Info("reconnect loop started")

	closeCh := cc.ch.NotifyClose(make(chan *amqp091.Error))
	cancelCh := cc.ch.NotifyCancel(make(chan string))

	select {
	case err := <-closeCh:
		if err != nil {
			cc.reconnect()
		}
		if err == nil {
			logger.Info("channel closed gracefully")
		}

	case cancelMsg := <-cancelCh:
		logger.Warn("channel cancelled, reconnecting", slog.String("reason", cancelMsg))
		cc.reconnect()
	}
}

func (cc *ChannelController) reconnect() {
	logger := cc.logger
	for {
		// TODO: implement exponential backoff
		time.Sleep(cc.config.ReconnSleepDuration)
		cc.mu.Lock()

		logger.Info("attempting to reconnect to RabbitMQ")

		conn, err := createConnection(cc.config)
		if err != nil {
			logger.Error("could not create connection", slog.String("error", err.Error()))
			cc.mu.Unlock()
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			logger.Error("could not create channel", slog.String("error", err.Error()))
			cc.mu.Unlock()
			continue
		}
		if cc.config.IsConfirmMode {
			err = ch.Confirm(false)
			if err != nil {
				logger.Error("could not set channel to confirm mode", slog.String("error", err.Error()))
				cc.mu.Unlock()
				continue
			}
		}

		err = cc.ch.Close()
		if err != nil {
			logger.Error("could not close channel", slog.String("error", err.Error()))
		}
		err = cc.conn.Close()
		if err != nil {
			logger.Error("could not close connection", slog.String("error", err.Error()))
		}

		logger.Info("reconnected to RabbitMQ, resuming channel")
		cc.ch = ch
		cc.mu.Unlock()
		go cc.reconnectLoop()
		cc.Resume()
		cc.reconnectCh <- struct{}{}
		return
	}
}
