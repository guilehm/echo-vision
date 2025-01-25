package rabbitmq

import (
	"log/slog"
	"time"
)

type Config struct {
	URL          string
	ExchangeName ExchangeName
	QueueName    QueueName
	ConsumerName ConsumerName

	ConcurrentConsumers int
	PrefetchCount       int

	Logger         *slog.Logger
	PublishTimeout time.Duration
	IsConfirmMode  bool

	ReconnSleepDuration time.Duration
	GracefulTimeout     time.Duration

	Opts []ConfigOpt
}
type ConfigOpt func(c *Config)

func ConfigConcurrentConsumers(concurrency int) func(c *Config) {
	return func(c *Config) {
		c.ConcurrentConsumers = concurrency
		c.PrefetchCount = concurrency
	}
}

func ConfigWithLogger(logger *slog.Logger) func(c *Config) {
	return func(c *Config) {
		c.Logger = logger
	}
}

func ConfigWithExchangeName(exchangeName ExchangeName) func(c *Config) {
	return func(c *Config) {
		c.ExchangeName = exchangeName
	}
}

func ConfigWithQueueName(queueName QueueName) func(c *Config) {
	return func(c *Config) {
		c.QueueName = queueName
	}
}

func ConfigConsumerName(consumerName ConsumerName) func(c *Config) {
	return func(c *Config) {
		c.ConsumerName = consumerName
	}
}

func ConfigWithPublishTimeout(timeout time.Duration) func(c *Config) {
	return func(c *Config) {
		c.PublishTimeout = timeout
	}
}

func ConfigWithPrefetchCount(prefetchCount int) func(c *Config) {
	return func(c *Config) {
		c.PrefetchCount = prefetchCount
	}
}

func ConfigWithConfirmMode() func(c *Config) {
	return func(c *Config) {
		c.IsConfirmMode = true
	}
}

func ConfigWithReconnSleepDuration(duration time.Duration) func(c *Config) {
	return func(c *Config) {
		c.ReconnSleepDuration = duration
	}
}

func ConfigWithGracefulTimeout(timeout time.Duration) func(c *Config) {
	return func(c *Config) {
		c.GracefulTimeout = timeout
	}
}

func newRabbitMQConfig(url string, logger *slog.Logger, opts ...ConfigOpt) *Config {
	// set default values
	c := &Config{
		ExchangeName:        "events",
		QueueName:           "",
		ConsumerName:        "",
		URL:                 url,
		ConcurrentConsumers: 1,
		PrefetchCount:       1,
		Logger:              logger,
		PublishTimeout:      10 * time.Second,
		ReconnSleepDuration: 5 * time.Second,
		GracefulTimeout:     30 * time.Second,
		Opts:                opts,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
